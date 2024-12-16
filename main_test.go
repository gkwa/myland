package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gkwa/myland/template"
)

var update = flag.Bool("update", false, "update .golden files")

func readAndTrim(t *testing.T, path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	return strings.TrimSpace(string(content))
}

func TestTemplateEscaper(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		golden string
	}{
		{
			name:   "simple_template",
			input:  "testdata/simple_template.input",
			golden: "testdata/simple_template.golden",
		},
		{
			name:   "already_escaped",
			input:  "testdata/already_escaped.input",
			golden: "testdata/already_escaped.golden",
		},
		{
			name:   "mixed_content",
			input:  "testdata/mixed_content.input",
			golden: "testdata/mixed_content.golden",
		},
	}

	if err := os.MkdirAll("testdata", 0o755); err != nil {
		t.Fatalf("Failed to create testdata directory: %v", err)
	}

	flag.Parse()

	escaper := template.NewDelimiterEscaper()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := readAndTrim(t, tt.input)
			result := escaper.Escape(input)

			if *update {
				err := os.WriteFile(tt.golden, []byte(result+"\n"), 0o644)
				if err != nil {
					t.Fatalf("Failed to update golden file: %v", err)
				}
				return
			}

			expected := readAndTrim(t, tt.golden)
			result = strings.TrimSpace(result)

			if expected != result {
				t.Errorf("Result does not match golden file.\nExpected:\n%s\nGot:\n%s", expected, result)
			}
		})
	}
}

func TestProcessFile(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		golden string
	}{
		{
			name:   "file_processing",
			input:  "testdata/simple_template.input",
			golden: "testdata/file_processing.golden",
		},
	}

	if err := os.MkdirAll("testdata", 0o755); err != nil {
		t.Fatalf("Failed to create testdata directory: %v", err)
	}

	escaper := template.NewDelimiterEscaper()
	handler := template.NewFileHandler(escaper)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpOutput := filepath.Join("testdata", "tmp_output.txt")
			inFile, err := os.Open(tt.input)
			if err != nil {
				t.Fatalf("Failed to open input file: %v", err)
			}
			defer inFile.Close()

			outFile, err := os.Create(tmpOutput)
			if err != nil {
				t.Fatalf("Failed to create output file: %v", err)
			}
			defer outFile.Close()

			if err := handler.Process(inFile, outFile); err != nil {
				t.Fatalf("Process failed: %v", err)
			}

			result := readAndTrim(t, tmpOutput)

			if *update {
				err := os.WriteFile(tt.golden, []byte(result+"\n"), 0o644)
				if err != nil {
					t.Fatalf("Failed to update golden file: %v", err)
				}
				return
			}

			expected := readAndTrim(t, tt.golden)

			if expected != result {
				t.Errorf("Result does not match golden file.\nExpected:\n%s\nGot:\n%s", expected, result)
			}

			os.Remove(tmpOutput)
		})
	}
}
