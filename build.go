//go:build ignore
// +build ignore

// A small helper you can run with:  go run build.go
// It discovers all main packages in the current module and builds them into ./bin.

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type target struct {
	dir        string // absolute directory of the package
	importPath string // full import path
	exeName    string // output executable name (without extension)
}

func main() {
	// Determine working directory and output folder
	rootDir, err := os.Getwd()
	if err != nil {
		fatalf("getting current directory: %v", err)
	}

	binDir := filepath.Join(rootDir, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		fatalf("creating bin directory: %v", err)
	}

	// Discover main packages in the module
	mains, err := discoverMainPackages()
	if err != nil {
		fatalf("discovering main packages: %v", err)
	}

	// Optionally include extra build targets from a .build file (one import path per line)
	if extras, err := readBuildFile(".build"); err == nil {
		mains = mergeTargets(mains, extras)
	}

	if len(mains) == 0 {
		fatalf("no main packages found to build")
	}

	// Build each target
	// Determine output extension based on target OS: prefer GOOS env if set, else host runtime
	targetOS := os.Getenv("GOOS")
	if targetOS == "" {
		targetOS = runtime.GOOS
	}
	ext := ""
	if targetOS == "windows" {
		ext = ".exe"
	}

	for _, t := range mains {
		outPath := filepath.Join(binDir, t.exeName+ext)
		fmt.Printf("Building %s -> %s\n", t.importPath, outPath)
		if err := runCmd("go", "build", "-o", outPath, t.importPath); err != nil {
			fatalf("building %s: %v", t.importPath, err)
		}
	}

	fmt.Println("Build completed successfully!")
	fmt.Printf("Executables are in: %s\n", binDir)
}

func discoverMainPackages() ([]target, error) {
	// Use `go list` to enumerate packages and filter to Name==main
	// Format: Dir|Name|ImportPath
	out, err := exec.Command("go", "list", "-f", "{{.Dir}}|{{.Name}}|{{.ImportPath}}", "./...").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("go list failed: %v\n%s", err, string(out))
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var res []target
	seen := map[string]bool{}
	rootDir, _ := os.Getwd()

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) != 3 {
			continue
		}
		dir, name, importPath := parts[0], parts[1], parts[2]
		if name != "main" {
			continue
		}
		// Determine executable name from directory base
		exe := filepath.Base(dir)
		// If this is the module root (contains go.mod and equals cwd), prefer folder name
		if sameDir(dir, rootDir) {
			exe = filepath.Base(rootDir)
		}
		key := dir + "::" + importPath
		if !seen[key] {
			seen[key] = true
			res = append(res, target{dir: dir, importPath: importPath, exeName: exe})
		}
	}

	// Sort output by executable name for stable ordering
	sort.Slice(res, func(i, j int) bool { return res[i].exeName < res[j].exeName })
	return res, nil
}

func readBuildFile(path string) ([]target, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var res []target
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Treat the line as an import path or relative directory
		imp := line
		// If it looks like a relative path, keep it as-is for `go build`
		exe := filepath.Base(strings.TrimRight(line, "/"))
		res = append(res, target{importPath: imp, exeName: exe})
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, errors.New("no targets")
	}
	return res, nil
}

func mergeTargets(a, b []target) []target {
	seen := map[string]bool{}
	var out []target
	for _, t := range a {
		k := t.importPath + "::" + t.exeName
		if !seen[k] {
			seen[k] = true
			out = append(out, t)
		}
	}
	for _, t := range b {
		k := t.importPath + "::" + t.exeName
		if !seen[k] {
			seen[k] = true
			out = append(out, t)
		}
	}
	return out
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

func sameDir(a, b string) bool {
	ap, _ := filepath.Abs(a)
	bp, _ := filepath.Abs(b)
	// Normalize path separators for safety
	ap = filepath.Clean(ap)
	bp = filepath.Clean(bp)
	return ap == bp
}

func fatalf(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	// Ensure trailing newline
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	// Prefix to stand out
	var buf bytes.Buffer
	buf.WriteString("error: ")
	buf.WriteString(msg)
	fmt.Fprint(os.Stderr, buf.String())
	os.Exit(1)
}
