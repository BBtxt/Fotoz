package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("Usage: <program> <path to directory>")
		return
	}

	photoDir := os.Args[1]
	absPhotoDir, err := filepath.Abs(photoDir)
	if err != nil {
		log.Println("Error resolving directory path:", err)
		return
	}

	files, err := os.ReadDir(absPhotoDir)
	if err != nil {
		log.Println("Error reading directory:", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}

		filePath := filepath.Join(absPhotoDir, file.Name())
		processPhoto(filePath) // Process each photo
	}
}

func processPhoto(photoPath string) {
	file, err := os.Open(photoPath)
	if err != nil {
		log.Printf("Error opening file %s: %v", photoPath, err)
		return
	}
	defer file.Close()

	exifData, err := exif.Decode(file)
	if err != nil {
		log.Printf("Error decoding Exif data from image %s: %v", photoPath, err)
		return
	}

	modelTag, _ := exifData.Get(exif.Model)
	dateTakenTag, _ := exifData.Get(exif.DateTime)

	model, _ := modelTag.StringVal()
	dateTaken, _ := dateTakenTag.StringVal()

	time, err := time.Parse("2006:01:02 15:04:05", dateTaken)
	if err != nil {
		log.Printf("Error parsing date from image %s: %v", photoPath, err)
		return
	}

	year, quarter := getYearQuarter(time.Month())

	baseDir := filepath.Join("/", "/Users/brandonbaker/Pictures/")
	cameraDir := filepath.Join(baseDir, model)
	yearDir := filepath.Join(cameraDir, fmt.Sprintf("%d", year))
	quarterDir := filepath.Join(yearDir, fmt.Sprintf("Q%d", quarter))

	err = os.MkdirAll(quarterDir, os.ModePerm)
	if err != nil {
		log.Printf("Error creating directories for image %s: %v", photoPath, err)
		return
	}

	newFilePath := filepath.Join(quarterDir, filepath.Base(photoPath))
	err = copyFile(photoPath, newFilePath)
	if err != nil {
		log.Printf("Error copying/moving file %s: %v", photoPath, err)
		return
	}

	fmt.Printf("File %s moved to: %s\n", photoPath, newFilePath)
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
