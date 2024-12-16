package template

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Delimiter struct {
	Open  string
	Close string
}

type Escaper interface {
	Escape(input string) string
}

type FileProcessor interface {
	Process(reader io.Reader, writer io.Writer) error
}

type DelimiterEscaper struct {
	escapedDelims   Delimiter
	savePattern     string
	unescapePattern *regexp.Regexp
	templatePattern *regexp.Regexp
}

func NewDelimiterEscaper() *DelimiterEscaper {
	return &DelimiterEscaper{
		escapedDelims: Delimiter{
			Open:  `{{"{{"}}`,
			Close: `{{"}}"}}`,
		},
		savePattern:     "SAVED_SEQUENCE_%d",
		unescapePattern: regexp.MustCompile(`{{"\{\{"}}[^}]*{{"}}"}}`),
		templatePattern: regexp.MustCompile(`{{\s*([^}]+?)\s*}}`),
	}
}

func (e *DelimiterEscaper) Escape(input string) string {
	saved := make(map[string]string)
	counter := 0

	input = e.unescapePattern.ReplaceAllStringFunc(input, func(match string) string {
		placeholder := fmt.Sprintf(e.savePattern, counter)
		saved[placeholder] = match
		counter++
		return placeholder
	})

	input = e.templatePattern.ReplaceAllString(input, e.escapedDelims.Open+` $1 `+e.escapedDelims.Close)

	for placeholder, original := range saved {
		input = strings.ReplaceAll(input, placeholder, original)
	}

	return input
}

type FileHandler struct {
	escaper Escaper
}

func NewFileHandler(escaper Escaper) *FileHandler {
	return &FileHandler{
		escaper: escaper,
	}
}

func (f *FileHandler) Process(reader io.Reader, writer io.Writer) error {
	input, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	processed := f.escaper.Escape(string(input))

	_, err = writer.Write([]byte(processed))
	if err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	return nil
}
