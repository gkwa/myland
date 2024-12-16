package template

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update .golden files")

type testCase struct {
	name   string
	input  string
	golden string
}

func readAndTrim(t *testing.T, path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	return strings.TrimSpace(string(content))
}

func TestDelimiterEscaper(t *testing.T) {
	tests := []testCase{
		{
			name:   "simple_template",
			input:  "../testdata/simple_template.input",
			golden: "../testdata/simple_template.golden",
		},
		{
			name:   "simple_template2",
			input:  "../testdata/simple_template2.input",
			golden: "../testdata/simple_template2.golden",
		},
		{
			name:   "already_escaped",
			input:  "../testdata/already_escaped.input",
			golden: "../testdata/already_escaped.golden",
		},
		{
			name:   "mixed_content",
			input:  "../testdata/mixed_content.input",
			golden: "../testdata/mixed_content.golden",
		},
	}

	escaper := NewDelimiterEscaper()
	runTests(t, escaper, tests)
}

func TestFileHandler(t *testing.T) {
	tests := []testCase{
		{
			name:   "file_processing",
			input:  "../testdata/simple_template.input",
			golden: "../testdata/file_processing.golden",
		},
	}

	escaper := NewDelimiterEscaper()
	handler := NewFileHandler(escaper)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputContent := readAndTrim(t, tt.input)

			var buf bytes.Buffer
			reader := strings.NewReader(inputContent)

			if err := handler.Process(reader, &buf); err != nil {
				t.Fatalf("Process failed: %v", err)
			}

			result := strings.TrimSpace(buf.String())

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
		})
	}
}

func runTests(t *testing.T, escaper Escaper, tests []testCase) {
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
