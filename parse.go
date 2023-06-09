package hclutil

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"golang.org/x/term"
)

type DiagnosticsWriter struct {
	once        sync.Once
	diagsWriter hcl.DiagnosticWriter
	diagsOutput io.Writer
	width       uint
	color       bool
	files       map[string]*hcl.File
}

func (w *DiagnosticsWriter) Files() []string {
	files := make([]string, 0, len(w.files))
	for k := range w.files {
		files = append(files, k)
	}
	return files
}

func (w *DiagnosticsWriter) SetOutput(output io.Writer) {
	w.once = sync.Once{}
	width := uint(400)
	color := false
	if output == nil {
		w.width = width
		w.color = color
		w.diagsOutput = io.Discard
		return
	}
	w.diagsOutput = output
	if f, ok := output.(*os.File); ok {
		fd := int(f.Fd())
		color = term.IsTerminal(fd)
		if w, _, err := term.GetSize(int(f.Fd())); err == nil {
			width = uint(w)
		}
	}
	w.width = width
	w.color = color
}

func (w *DiagnosticsWriter) SetColor(color bool) {
	w.once = sync.Once{}
	w.color = color
}

func (w *DiagnosticsWriter) SetWidth(width uint) {
	w.once = sync.Once{}
	w.width = width
}

func (w *DiagnosticsWriter) Output() io.Writer {
	return w.diagsOutput
}

func (w *DiagnosticsWriter) Color() bool {
	return w.color
}

func (w *DiagnosticsWriter) Width() uint {
	return w.width
}

func (w *DiagnosticsWriter) WriteDiagnostics(diags hcl.Diagnostics) error {
	w.once.Do(func() {
		w.diagsWriter = hcl.NewDiagnosticTextWriter(w.diagsOutput, w.files, w.width, w.color)
	})
	w.diagsWriter.WriteDiagnostics(diags)
	if diags.HasErrors() {
		return errors.New("diagnostics had errors, see above for details")
	}
	return nil
}

func newDiagnosticsWriter(files map[string]*hcl.File) *DiagnosticsWriter {
	w := &DiagnosticsWriter{
		files: files,
	}
	w.SetOutput(os.Stderr)
	return w
}

func Parse(p string) (hcl.Body, *DiagnosticsWriter, hcl.Diagnostics) {
	parser := hclparse.NewParser()
	stat, err := os.Stat(p)
	if err != nil {
		return nil, newDiagnosticsWriter(parser.Files()), hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Parse failed",
			Detail:   err.Error(),
		}}
	}
	if !stat.IsDir() {
		switch filepath.Ext(p) {
		case ".hcl":
			file, diags := parser.ParseHCLFile(p)
			return file.Body, newDiagnosticsWriter(parser.Files()), diags
		case ".json":
			file, diags := parser.ParseJSONFile(p)
			return file.Body, newDiagnosticsWriter(parser.Files()), diags
		default:
			return nil, newDiagnosticsWriter(parser.Files()), hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Unsupported file extension",
				Detail:   "Only .hcl and .json are supported",
			}}
		}
	}
	entires, err := os.ReadDir(p)
	if err != nil {
		return nil, newDiagnosticsWriter(parser.Files()), hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to read directory",
			Detail:   err.Error(),
		}}
	}
	var diags hcl.Diagnostics
	var files []*hcl.File
	for _, entry := range entires {
		if entry.IsDir() {
			continue
		}
		// if ext is .hcl execute parser.ParseHCLFile
		ext := filepath.Ext(entry.Name())
		if ext == ".hcl" {
			file, d := parser.ParseHCLFile(filepath.Join(p, entry.Name()))
			files = append(files, file)
			diags = append(diags, d...)
			continue
		}
		// if entity.Name() is *.hcl.json execute parser.ParseJSONFile
		if ext != ".json" {
			continue
		}
		baseName := filepath.Base(entry.Name())
		fileNameWithoutExt := baseName[:len(baseName)-len(ext)]
		if filepath.Ext(fileNameWithoutExt) == ".hcl" {
			file, d := parser.ParseJSONFile(filepath.Join(p, entry.Name()))
			files = append(files, file)
			diags = append(diags, d...)
		}
	}
	return hcl.MergeFiles(files), newDiagnosticsWriter(parser.Files()), diags
}
