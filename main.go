package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"slices"
)

func main() {

	GitPull()

	// Open the CSV file
	file, err := os.Open("vocab.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	//read function
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//open previousFile CSV file
	prevFile, err := os.Open("previousFile.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer prevFile.Close()

	prevReader := csv.NewReader(prevFile)

	previousRecords, err := prevReader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	records = slices.DeleteFunc(records, func(row []string) bool {
		for _, prevRow := range previousRecords {
			if reflect.DeepEqual(row, prevRow) {

				return true
			}
		}
		return false
	})

	//print
	for _, row := range records {
		for _, col := range row {
			fmt.Printf("%s\t", col)
		}
		fmt.Println()
	}

	var outputText string

	for _, row := range records {
		if string(row[3]) != "x" {
			// the structure must be this: Finnish|English|Part of Speech|Inflections|Example|Notes
			outputText += fmt.Sprintf("%s|%s|%s|%s|%s|%s", row[0], row[6], row[4], row[5], row[7], row[8])
		}
		outputText += "\n"
	}

	//print out
	outputErr := os.WriteFile("output.csv", []byte(outputText), 0644)
	if err != nil {
		log.Fatalln("Error writing file:", outputErr)
	}

	//copy vocab to previousFile
	copyErr := copyFile("vocab.csv", "previousFile.csv")
	if err != nil {
		fmt.Println("Error:", copyErr)
		return
	}
	fmt.Println("File copied successfully!")
}

// copy function i stole
func copyFile(src, dst string) error {
	// 1. Open the source file for reading
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source: %w", err)
	}
	defer source.Close() // Ensure the file is closed later

	// 2. Create or truncate the destination file
	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination: %w", err)
	}
	defer destination.Close()

	// 3. Stream data from source to destination
	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("failed to copy data: %w", err)
	}

	// 4. Commit file contents to stable storage
	return destination.Sync()
}
