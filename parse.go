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

// DiagnosticWriter は hcl.DiagnosticWriter のラッパーです。
// 通常のDiagnosticWriterに加えて以下の機能を追加します。
//
//	Output がターミナルであるのかどうかを検出し、色と幅を自動的に設定します。
//	Parse済みのファイル情報を保持します。
//
// DiagnosticWriter is a wrapper for hcl.DiagnosticWriter.
// In addition to the normal DiagnosticWriter, it adds the following features.
//
//	Detects whether Output is a terminal and automatically sets color and width.
//	Holds parsed file information.
type DiagnosticsWriter struct {
	once        sync.Once
	diagsWriter hcl.DiagnosticWriter
	diagsOutput io.Writer
	width       uint
	color       bool
	files       map[string]*hcl.File
}

// Files は Parse済みのファイルのPath名を返します。
// Files returns the Path name of the parsed file.
func (w *DiagnosticsWriter) Files() []string {
	files := make([]string, 0, len(w.files))
	for k := range w.files {
		files = append(files, k)
	}
	return files
}

// SetOutput は出力先を設定します。
// SetOutput sets the output destination.
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

// SetColor は色を設定します。
// SetColor sets the color.
func (w *DiagnosticsWriter) SetColor(color bool) {
	w.once = sync.Once{}
	w.color = color
}

// SetWidth は幅を設定します。
// SetWidth sets the width.
func (w *DiagnosticsWriter) SetWidth(width uint) {
	w.once = sync.Once{}
	w.width = width
}

// Output は出力先を返します。
// Output returns the output destination.
func (w *DiagnosticsWriter) Output() io.Writer {
	return w.diagsOutput
}

// Color をつけるかどうかを返します。
// Returns whether to add Color.
func (w *DiagnosticsWriter) Color() bool {
	return w.color
}

// Width は幅を返します。
// Width returns the width.
func (w *DiagnosticsWriter) Width() uint {
	return w.width
}

// WriteDiagnostics は診断情報を出力します。
// WriteDiagnostics outputs diagnostic information.
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

// Parse は与えられたPathをHCLとして解析します。
// Parse parses the given Path as HCL.
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
