package core

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gkwa/myland/template"
	"github.com/go-logr/logr"
)

type Processor struct {
	logger  logr.Logger
	handler *template.FileHandler
}

func NewProcessor(logger logr.Logger) *Processor {
	escaper := template.NewDelimiterEscaper()
	handler := template.NewFileHandler(escaper)

	return &Processor{
		logger:  logger,
		handler: handler,
	}
}

func (p *Processor) ProcessFiles(paths []string) error {
	for _, inputPath := range paths {
		if err := p.processFile(inputPath); err != nil {
			absPath, absErr := filepath.Abs(inputPath)
			if absErr != nil {
				return fmt.Errorf("processing %s: %w", inputPath, err)
			}
			return fmt.Errorf("processing %s: %w", absPath, err)
		}
	}
	return nil
}

func (p *Processor) processFile(inputPath string) error {
	original, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	var output bytes.Buffer
	if err := p.handler.Process(bytes.NewReader(original), &output); err != nil {
		return fmt.Errorf("processing content: %w", err)
	}

	processed := output.Bytes()

	originalSum := calculateChecksum(original)
	processedSum := calculateChecksum(processed)

	if originalSum == processedSum {
		p.logger.V(1).Info("file unchanged, skipping", "path", inputPath)
		return nil
	}

	if err := os.WriteFile(inputPath, processed, 0o644); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	p.logger.V(1).Info("updated file", "path", inputPath)
	return nil
}
