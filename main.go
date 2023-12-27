package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	photoDir := "test_assets/_R000372.DNG"

	//Open Image
	file, err := os.Open(photoPath)
	if err != nil {
		fmt.Println("Error processing image file", err)
		return
	}
	//Close Image
	defer file.Close()

	//Get Exif Data
	exifData, err := exif.Decode(file)
	if err != nil {
		fmt.Println("Error decoding Exif data from image", err)
		return
	}

	// Extract camera model and date taken
	modelTag, _ := exifData.Get(exif.Model)
	dateTakenTag, _ := exifData.Get(exif.DateTime)

	model, _ := modelTag.StringVal()
	dateTaken, _ := dateTakenTag.StringVal()

	// Parse date taken into year and quarter
	time, err := time.Parse("2006:01:02 15:04:05", dateTaken)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

	year, quarter := getYearQuarter(time.Month())

	// Create directory structure
	baseDir := filepath.Join("/", "/Users/brandonbaker/Pictures/") // Replace with your base directory
	cameraDir := filepath.Join(baseDir, model)
	yearDir := filepath.Join(cameraDir, fmt.Sprintf("%d", year))
	quarterDir := filepath.Join(yearDir, fmt.Sprintf("Q%d", quarter))

	// Create directories if they don't exist
	err = os.MkdirAll(quarterDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return
	}

	// Copy or move the photo to the new directory
	newFilePath := filepath.Join(quarterDir, filepath.Base(photoPath))
	err = copyFile(photoPath, newFilePath) // Implement a function to copy or move the file
	if err != nil {
		fmt.Println("Error copying/moving file:", err)
		return
	}

	fmt.Println("File moved to:", newFilePath)
}

func getYearQuarter(month time.Month) (year, quarter int) {
	switch {
	case month >= time.January && month <= time.March:
		return time.Now().Year(), 1
	case month >= time.April && month <= time.June:
		return time.Now().Year(), 2
	case month >= time.July && month <= time.September:
		return time.Now().Year(), 3
	default:
		return time.Now().Year(), 4
	}
}

func copyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}
	return nil
}

