package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	inputFile := "cdn/json/cities.json"
	data, err := os.ReadFile(inputFile)
	if err != nil {
		panic(fmt.Errorf("failed to read file %s: %v", inputFile, err))
	}

	var items []map[string]interface{}
	err = json.Unmarshal(data, &items)
	if err != nil {
		panic(fmt.Errorf("invalid JSON structure: %v", err))
	}

	if len(items) == 0 {
		panic("JSON array is empty")
	}

	outputFile := "cdn/json/output.csv"
	file, err := os.Create(outputFile)
	if err != nil {
		panic(fmt.Errorf("failed to create CSV file: %v", err))
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{}
	for key := range items[0] {
		headers = append(headers, key)
	}
	_ = writer.Write(headers)

	for _, obj := range items {
		row := []string{}
		for _, h := range headers {
			value := ""
			if obj[h] != nil {
				value = fmt.Sprintf("%v", obj[h])
			}
			row = append(row, value)
		}
		_ = writer.Write(row)
	}

	fmt.Printf("CSV file created successfully: %s\n", outputFile)
}
