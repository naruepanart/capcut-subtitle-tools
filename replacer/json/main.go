package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config map[string]string

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("JSON file processed successfully")
}

func run() error {
	cfg, err := load("capcut-subtitle-replacer.json")
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	target, err := path("file-path.txt")
	if err != nil {
		return fmt.Errorf("path error: %w", err)
	}

	return process(target, cfg)
}

func load(name string) (Config, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func path(name string) (string, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}

	target := filepath.Clean(strings.TrimSpace(string(data)))
	if target == "" || strings.Contains(target, "..") {
		return "", fmt.Errorf("invalid file path")
	}

	return target, nil
}

func process(path string, cfg Config) error {
	r := replacer(cfg)

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	data, err := jsonProcess(content, r)
	if err != nil {
		return err
	}

	return write(path, data, info.Mode())
}

func replacer(cfg Config) *strings.Replacer {
	pairs := make([]string, 0, len(cfg)*2)
	for k, v := range cfg {
		pairs = append(pairs, k, v)
	}
	return strings.NewReplacer(pairs...)
}

func jsonProcess(content []byte, r *strings.Replacer) ([]byte, error) {
	var data interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, err
	}

	processed := value(data, r)
	return json.Marshal(processed)
}

func value(v interface{}, r *strings.Replacer) interface{} {
	switch x := v.(type) {
	case string:
		return r.Replace(x)
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, val := range x {
			result[k] = value(val, r)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(x))
		for i, val := range x {
			result[i] = value(val, r)
		}
		return result
	default:
		return v
	}
}

func write(path string, data []byte, mode os.FileMode) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, mode); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
