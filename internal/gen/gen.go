package gen

import (
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
)

const (
	header = `// Generated by github.com/paulwerner/gomoney/gen. DO NOT EDIT.`
)

func WriteGoFile(filename, pkg string, b []byte) {
	w, err := os.Create(filename)
	if err != nil {
		log.Fatalf("error creating file %s: %v", filename, err)
	}
	defer w.Close()

	if _, err = WriteGo(w, pkg, "", b); err != nil {
		log.Fatalf("error writing file %s: %v", filename, err)
	}
}

func WriteGo(w io.Writer, pkg, tags string, b []byte) (n int, err error) {
	src := []byte(header)
	if tags != "" {
		src = append(src, fmt.Sprintf("// +build %s\n\n", tags)...)
	}
	src = append(src, fmt.Sprintf("package %s\n\n", pkg)...)
	src = append(src, b...)
	formatted, err := format.Source(src)
	if err != nil {
		// write generated code despite error allowing interpretation
		n, _ = w.Write(src)
		return n, err
	}
	return w.Write(formatted)
}