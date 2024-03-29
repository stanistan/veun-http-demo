package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	var (
		// rootPath directory specifies where the project root is. If not specified
		// we will find it by looking for the go.mod file from the current working
		// directory.
		rootPath = flag.String("root", "", "Root Directory")

		// docsPath is the relative path from the root to the docs. defaults to docs.
		docsPath = flag.String("docs", "docs", "relative path from root to docs")

		// pkgName is the package we are going to build. Defaults to GOPACKAGE env var.
		pkgName = flag.String("pkg", os.Getenv("GOPACKAGE"), "go package")
	)

	flag.Parse()

	// current working directory is required.
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	rp := cwd
	for *rootPath == "" {
		if rp == "/" {
			log.Fatal("root volume")
		}
		if _, err := os.Stat(path.Join(rp, "go.mod")); err == nil {
			*rootPath = rp
			break
		} else {
			rp = filepath.Dir(rp)
		}
	}

	outputDir := cwd
	rel, err := filepath.Rel(*rootPath, cwd)
	if err != nil {
		log.Fatal(err)
	}

	docRoot := filepath.Join(*rootPath, *docsPath, rel)

	fmt.Printf("[lit-gen] :: source=%s dest=%s\n", docRoot, outputDir)

	if err := walk(docRoot, generate(docRoot, outputDir, *pkgName)); err != nil {
		log.Fatal(err)
	}
}

func walk(root string, genF func(string) error) error {
	files, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go.md") {
			if err := genF(filepath.Join(root, file.Name())); err != nil {
				return err
			}
		}
	}

	return nil

	return filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".go.md") && !info.IsDir() {
			if err := genF(path); err != nil {
				return err
			}
		}

		return nil
	})
}

func generate(root, destination, pkg string) func(string) error {
	return func(in string) error {

		f, err := os.Open(in)
		if err != nil {
			return err
		}
		defer f.Close()

		fmt.Printf("[lit-gen] -> %s\n", filepath.Base(in))

		dest := filepath.Join(
			destination,
			strings.TrimSuffix(strings.TrimPrefix(in, root), ".go.md"),
		) + ".generated.go"

		out, err := os.Create(dest)
		if err != nil {
			return fmt.Errorf("opening destination file: %w", err)
		}
		defer out.Close()

		var (
			scanner = bufio.NewScanner(f)
			w       bytes.Buffer
		)

		var (
			shouldWrite bool
		)

		_, _ = w.WriteString(PREFIX + "\n")
		_, _ = w.WriteString("package " + pkg + "\n")

		for scanner.Scan() {
			line := scanner.Text()
			switch line {
			case "```":
				if shouldWrite {
					shouldWrite = false
					_, _ = w.WriteString("\n")
				}
			case "```go":
				shouldWrite = true
			default:
				if shouldWrite {
					_, _ = w.WriteString(line + "\n")
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("scanning: %w", err)
		}

		formatted, err := format.Source(w.Bytes())
		if err != nil {
			return fmt.Errorf("%s formatting error: %w", in, err)
		}

		_, err = out.Write(formatted)
		if err != nil {
			return fmt.Errorf("writing file: %w", err)
		}

		fmt.Printf("[lit-gen] -> DONE\n")
		return nil
	}
}

const PREFIX = `// Code Generated by github.com/stanistan/veun-http-demo/cmd/lit-gen; DO NOT EDIT.`
