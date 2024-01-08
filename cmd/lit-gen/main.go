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
		root      = flag.String("root", "", "Root Directory")
		outputDir = flag.String("o", "", "Output Directory")
		pkg       = flag.String("pkg", os.Getenv("GOPACKAGE"), "go package")
	)
	flag.Parse()

	if *root == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		*root = wd
	} else {
		abs, err := filepath.Abs(*root)
		if err != nil {
			log.Fatal(err)
		}

		*root = abs
	}

	if *outputDir == "" {
		log.Fatal("-o (output directory) is required")
	} else {
		abs, err := filepath.Abs(*outputDir)
		if err != nil {
			log.Fatal(err)
		}

		*outputDir = abs
	}

	log.Printf("[attempting to generate go from *.go.md files in=%s out=%s]", *root, *outputDir)

	err := walk(*root, generate(*root, *outputDir, *pkg))
	if err != nil {
		log.Fatal(err)
	}
}

func walk(root string, genF func(string) error) error {
	return filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("path=%s err=%s", path, err)
			return nil
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
		log.Printf("-> %s", in)

		f, err := os.Open(in)
		if err != nil {
			return err
		}
		defer f.Close()

		dest := path.Join(destination, strings.TrimSuffix(strings.TrimPrefix(in, root), ".md"))
		log.Printf("  :: destination=%s", dest)

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

		_, err = w.WriteString(PREFIX + "\npackage " + pkg + "\n")
		if err != nil {
			return err
		}

		for scanner.Scan() {
			line := scanner.Text()
			switch line {
			case "```":
				if shouldWrite {
					shouldWrite = false
					_, err = w.WriteString("\n")
				}
			case "```go":
				shouldWrite = true
			default:
				if shouldWrite {
					_, err = w.WriteString(line + "\n")
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("scanning: %w", err)
		}

		log.Printf("  :: Formatting")
		formatted, err := format.Source(w.Bytes())
		if err != nil {
			return fmt.Errorf("formatting error: %w", err)
		}

		_, err = out.Write(formatted)
		if err != nil {
			return fmt.Errorf("writing file: %w", err)
		}

		log.Printf("  :: DONE")
		return nil
	}
}

const PREFIX = `// CODE GENERATED BY lit-gen; DO NOT EDIT.`
