package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/otiai10/copy"
	"golang.org/x/xerrors"
)

func main() {
	if err := run(os.Args[1:]); err != nil && err != flag.ErrHelp {
		log.Println(err)
		os.Exit(1)
	}
}

func isSkipDir(fi os.FileInfo) bool {
	if !fi.IsDir() {
		return false
	}
	n := fi.Name()
	switch n {
	case "", "testdata":
		return true
	}
	switch n[0] {
	case '.', '_':
		return true
	}
	return false
}

func run(args []string) error {
	bs, err := exec.Command("ghq", "root").Output()
	if err != nil {
		return xerrors.Errorf("failed to execute `ghq root`: %w", err)
	}
	goroot := filepath.Join(strings.TrimSpace(string(bs)), "github.com/golang/go/src")

	if err := place(goroot, "cmd/go/internal/modfile"); err != nil {
		return err
	}

	wd, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	return filepath.Walk(wd, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if isSkipDir(fi) {
			return filepath.SkipDir
		}
		if !strings.HasSuffix(p, ".go") {
			return nil
		}
		return rewriteGo(p)
	})
}

func place(goroot, dir string) error {
	target := filepath.Join(goroot, dir)
	dest := destDir(target)
	if err := copy.Copy(target, dest); err != nil {
		return err
	}
	bs, err := exec.Command("go", "list", "-f", `{{join .Imports "\x00"}}`, target).Output()
	if err != nil {
		return xerrors.Errorf("failed to execute `go list`: %w", err)
	}
	deps := strings.Split(string(bs), "\x00")
	for _, pkg := range deps {
		if strings.Contains(pkg, "internal/") {
			if err := place(goroot, pkg); err != nil {
				return err
			}
		}
	}
	return nil
}

func destDir(dir string) string {
	if strings.HasSuffix(dir, "modfile") {
		return "."
	}
	return filepath.Join("internal", filepath.Base(dir))
}

var rewriteReg = regexp.MustCompile(`"[^"]*\binternal/[^"]+"`)

func rewriteGo(f string) error {
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	str := rewriteReg.ReplaceAllStringFunc(string(bs), func(match string) string {
		return fmt.Sprintf(`"github.com/Songmu/modfile/internal/%s`, path.Base(match))
	})
	str = fmt.Sprintf(`// Code generated by _tools/place.go. DO NOT EDIT.
// This is copied from Go source code and modified by above script.
// Here is the original copyright notice:

%s`, str)

	return ioutil.WriteFile(f, []byte(str), 0644)
}
