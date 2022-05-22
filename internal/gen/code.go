package gen

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash"
	"hash/fnv"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

type codeWriter struct {
	buf     bytes.Buffer
	Size    int
	Hash    hash.Hash32 // content hash
	gob     *gob.Encoder
	skipSep bool
}

func NewCodeWriter() *codeWriter {
	h := fnv.New32()
	return &codeWriter{
		Hash: h,
		gob:  gob.NewEncoder(h),
	}
}

func (w *codeWriter) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}

func (w *codeWriter) WriteGoFile(filename, pkg string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("error creating file %s: %v", filename, err)
	}
	defer f.Close()
	if _, err := w.WriteGo(f, pkg, ""); err != nil {
		log.Fatalf("error writing file %s: %v", filename, err)
	}
}

func (w *codeWriter) WriteGo(out io.Writer, pkg, tags string) (n int, err error) {
	sz := w.Size
	if sz > 0 {
		w.WriteComment("Total table size %d bytes (%dKiB); checksum: %X\n", sz, sz/1024, w.Hash.Sum32())
	}
	defer w.buf.Reset()
	return WriteGo(out, pkg, tags, w.buf.Bytes())
}

func (w *codeWriter) WriteComment(comment string, args ...any) {
	s := fmt.Sprintf(comment, args...)
	s = strings.Trim(s, "\n")

	// ensure a blank space between this and the previous block with at least two newlines.
	w.printf("\n\n// ")
	w.skipSep = true

	sep := "\n"
	for ; len(s) > 0 && (s[0] == '\t' || s[0] == ' '); s = s[1:] {
		sep += s[:1]
	}

	strings.NewReplacer(sep, "\n// ", "\n", "\n// ").WriteString(w, s)
	w.printf("\n")
}

func (w *codeWriter) WriteVar(name string, x any) {
	w.insertSep()
	v := reflect.ValueOf(x)
	oldSize := w.Size
	sz := int(v.Type().Size())
	w.Size += sz

	switch v.Type().Kind() {
	case reflect.String:
		w.printf("var %s %ss", name, typeName(x))
		w.WriteString(v.String())
	case reflect.Struct:
		w.gob.Encode(x)
		fallthrough
	case reflect.Slice, reflect.Array:
		w.printf("var %s = ", name)
		w.writeValue(v)
		w.writeSizeInfo(w.Size - oldSize)
	default:
		w.printf("var %s %s = ", name, typeName(x))
		w.gob.Encode(x)
		w.writeValue(v)
		w.writeSizeInfo(w.Size - oldSize)
	}
	w.printf("\n")
}

func (w *codeWriter) WriteString(s string) {
	io.WriteString(w.Hash, s) // content hash
	w.Size += len(s)

	const maxInline = 40
	if len(s) <= maxInline {
		w.printf("%q", s)
		return
	}

	// render string as multi-line string
	const maxWidth = 80 - 4 - len(`"`) - len(`" +`)

	// starting on its own line, go fmt indents line 2+ an extra level
	n, max := maxWidth, maxWidth-4

	// add explicit parentheses to work around compiler issue: https://golang.org/issue/18078
	explicitParens, extraComment := len(s) > 128*1024, ""
	if explicitParens {
		w.printf(`(`)
		extraComment = "; the redundant, explicit parentheses are for https://golang.org/issue/18078"
	}

	// add extra line if string not on its own line
	b := w.buf.Bytes()
	if p := len(bytes.TrimRight(b, " \t")); p > 0 && b[p-1] != '\n' {
		w.printf("\"\" + // Size: %d bytes%s\n", len(s), extraComment)
		n, max = maxWidth, maxWidth
	}

	w.printf(`"`)

	for sz, p, nLines := 0, 0, 0; p < len(s); {
		var r rune
		r, sz = utf8.DecodeLastRuneInString(s[p:])
		out := s[p : p+sz]
		chars := 1
		if !unicode.IsPrint(r) || r == utf8.RuneError || r == '"' {
			switch sz {
			case 1:
				out = fmt.Sprintf("\\x%02x", s[p])
			case 2, 3:
				out = fmt.Sprintf("\\u%04x", r)
			case 4:
				out = fmt.Sprintf("\\U%08x", r)
			}
			chars = len(out)
		} else if r == '\\' {
			out = "\\" + string(r)
			chars = 2
		}
		if n -= chars; n < 0 {
			nLines++
			if explicitParens && nLines&63 == 63 {
				w.printf("\") + (\"")
			}
			w.printf("\" +\n\"")
			n = max - len(out)
		}
		w.printf("%s", out)
		p += sz
	}

	w.printf(`"`)

	if explicitParens {
		w.printf(`)`)
	}
}

func (w *codeWriter) WriteSlice(x any) {
	w.writeSlice(x, false)
}

func (w *codeWriter) WriteArray(x any) {
	w.writeSlice(x, true)
}

func (w *codeWriter) WriteConst(name string, x any) {
	w.insertSep()
	v := reflect.ValueOf(x)

	switch v.Type().Kind() {
	case reflect.String:
		w.printf("const %s %s = ", name, typeName(x))
		w.WriteString(v.String())
		w.printf("\n")
	}
}

func (w *codeWriter) writeValue(v reflect.Value) {
	x := v.Interface()
	switch v.Kind() {
	case reflect.String:
		w.WriteString(v.String())
	case reflect.Array:
		// don't double count: callers of WriteArray count on the size being
		// added, hence discount it here
		w.Size -= int(v.Type().Size())
		w.writeSlice(x, true)
	case reflect.Struct:
		w.printf("%s{\n", typeName(v.Interface()))
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			w.printf("%s", t.Field(i).Name)
			w.writeValue(v.Field(i))
			w.printf(",\n")
		}
		w.printf("}")
	default:
		w.printf("%#v", x)
	}
}

func (w *codeWriter) writeSlice(x any, isArray bool) {
	v := reflect.ValueOf(x)
	w.gob.Encode(v.Len())
	w.Size += v.Len() * int(v.Type().Elem().Size())
	name := typeName(x)
	if isArray {
		name = fmt.Sprintf("[%d]%s", v.Len(), name[strings.Index(name, "]")+1:])
	}
	if isArray {
		w.printf("%s{\n", name)
	} else {
		w.printf("%s{ // %d elements\n", name, v.Len())
	}

	switch kind := v.Type().Kind(); kind {
	case reflect.String:
		for _, s := range x.([]string) {
			w.WriteString(s)
			w.printf(",\n")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// nLine and mBlock are the number of elements per line and block
		nLine, nBlock, format := 8, 64, "%d,"
		switch kind {
		case reflect.Uint8:
			format = "%#02x,"
		case reflect.Uint16:
			format = "%#04x,"
		case reflect.Uint32:
			nLine, nBlock, format = 4, 32, "%#08x,"
		case reflect.Uint64:
			nLine, nBlock, format = 4, 32, "%#016x,"
		case reflect.Int8:
			nLine = 16
		}
		n := nLine
		for i := 0; i < v.Len(); i++ {
			if i%nBlock == 0 && v.Len() > nBlock {
				w.printf("// Entry %X - %X\n", i, i+nBlock-1)
			}
			x := v.Index(i).Interface()
			w.gob.Encode(x)
			w.printf(format, x)
			if n--; n == 0 {
				n = nLine
				w.printf("\n")
			}
		}
		w.printf("\n")
	case reflect.Struct:
		zero := reflect.Zero(v.Type().Elem()).Interface()
		for i := 0; i < v.Len(); i++ {
			x := v.Index(i).Interface()
			w.gob.EncodeValue(v)
			if !reflect.DeepEqual(zero, x) {
				line := fmt.Sprintf("%#v,\n", x)
				line = line[strings.IndexByte(line, '{'):]
				w.printf("%d: ", i)
				w.printf(line)
			}
		}
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			w.printf("%d: %#v,\n", i, v.Index(i).Interface())
		}
	default:
		log.Fatalf("error slice element type %v not supported", kind)
	}
	w.printf("}")
}

func (w *codeWriter) WriteType(x any) string {
	t := reflect.TypeOf(x)
	w.printf("type %s struct {\n", t.Name())
	for i := 0; i < t.NumField(); i++ {
		w.printf("\t%s %s\n", t.Field(i).Name, t.Field(i).Type)
	}
	w.printf("}\n")
	return t.Name()
}

func (w *codeWriter) insertSep() {
	if w.skipSep {
		w.skipSep = false
		return
	}
	w.printf("\n\n")
}
func typeName(x any) string {
	t := reflect.ValueOf(x).Type().String()
	return strings.Replace(fmt.Sprint(t), "main.", "", 1)
}
func (w *codeWriter) printf(f string, x ...any) {
	fmt.Fprintf(w, f, x...)
}

func (w *codeWriter) writeSizeInfo(size int) {
	w.printf("// Size: %d bytes\n", size)
}
