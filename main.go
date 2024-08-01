package main

import (
	"encoding/base64"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// Secret Define the structure of the YAML file
type Secret struct {
	APIVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Metadata   map[string]interface{} `yaml:"metadata"`
	Data       map[string]interface{} `yaml:"data"`
}

func main() {
	// Directory containing YAML files
	dir := "."

	// List all .yaml files in the directory
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		log.Fatalf("Failed to list YAML files: %v", err)
	}

	for _, file := range files {
		fmt.Printf("Processing file: %s\n", file)

		// Read the YAML file
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Failed to read file %s: %v\n", file, err)
			continue
		}

		// Parse the YAML data
		var secret Secret
		err = yaml.Unmarshal(data, &secret)
		if err != nil {
			fmt.Printf("Failed to parse YAML in file %s: %v\n", file, err)
			continue
		}

		// Validate the expected structure
		if secret.APIVersion == "" || secret.Kind == "" || secret.Data == nil {
			fmt.Printf("Invalid structure in file %s: skipping\n", file)
			continue
		}

		// Encode values to base64
		for key, value := range secret.Data {
			var encoded string
			switch v := value.(type) {
			case int:
				encoded = base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(v)))
			case bool:
				encoded = base64.StdEncoding.EncodeToString([]byte(strconv.FormatBool(v)))
			case string:
				encoded = base64.StdEncoding.EncodeToString([]byte(v))
			default:
				encoded = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", v)))
			}

			secret.Data[key] = encoded
		}

		// Marshal the new YAML data
		outputData, err := yaml.Marshal(&secret)
		if err != nil {
			fmt.Printf("Failed to marshal YAML for file %s: %v\n", file, err)
			continue
		}

		// Create output file name
		outputFile := file[:len(file)-len(filepath.Ext(file))] + "_encoded.yaml"

		// Write the output YAML file
		err = os.WriteFile(outputFile, outputData, 0644)
		if err != nil {
			fmt.Printf("Failed to write output file %s: %v\n", outputFile, err)
			continue
		}

		fmt.Printf("Encoded YAML saved to %s\n", outputFile)
	}
}
