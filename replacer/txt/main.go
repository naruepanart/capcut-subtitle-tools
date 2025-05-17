package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type TextReplacer struct {
	replacements map[string]string
}

func NewTextReplacer(replacements map[string]string) *TextReplacer {
	return &TextReplacer{replacements: replacements}
}

func (tr *TextReplacer) Replace(text string) string {
	if len(tr.replacements) == 0 {
		return text
	}

	var oldnew []string
	for old, new := range tr.replacements {
		oldnew = append(oldnew, old, new)
	}
	return strings.NewReplacer(oldnew...).Replace(text)
}

func loadConfig(filename string) (map[string]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("config read error: %w", err)
	}

	var config map[string]string
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("config parse error: %w", err)
	}
	return config, nil
}

func processFile(inputFile string, replacer *TextReplacer) error {
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("read error: %w", err)
	}

	result := replacer.Replace(string(content))
	if err := os.WriteFile(inputFile, []byte(result), 0644); err != nil {
		return fmt.Errorf("write error: %w", err)
	}
	return nil
}

func main() {
	// Load configuration
	replacements, err := loadConfig("text.json")
	if err != nil {
		fmt.Printf("Configuration error: %v\n", err)
		os.Exit(1)
	}

	// Initialize replacer
	replacer := NewTextReplacer(replacements)

	// Process file
	if err := processFile("a.txt", replacer); err != nil {
		fmt.Printf("Processing error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Text replacement completed successfully")
}
