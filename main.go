package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	// Check command line arguments
	if len(os.Args) < 3 {
		log.Println("Usage: <program> <path to directory>")
		return
	}
	// this is the directory containing the photos, the first argument on the script
	photoDir := os.Args[1]

	// this is the directory that the photos will be moved to, the second argument on the script
	baseDir, err := filepath.Abs(os.Args[2])
	if err != nil {
		log.Println("Error resolving base directory path:", err)
		return
	}
	// this is the directory that the photos will be moved to, the second argument on the script
	absPhotoDir, err := filepath.Abs(photoDir)
	if err != nil {
		log.Println("Error resolving directory path:", err)
		return
	}
	// this is the directory that the photos will be moved to, the second argument on the script
	files, err := os.ReadDir(absPhotoDir)
	if err != nil {
		log.Println("Error reading directory:", err)
		return
	}
	// Process each photo
	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}

		filePath := filepath.Join(absPhotoDir, file.Name())
		processPhoto(filePath, absPhotoDir, baseDir) // Process each photo
	}
}

func processPhoto(photoPath string, absPhotoDir string, baseDir string) {
	// processPhoto organizes a photo by extracting its Exif data, creating a directory structure based on the camera model, year, and quarter, and moving related files to the appropriate directory.
	//
	// Parameters:
	// - photoPath: the path to the photo file.
	// - absPhotoDir: the absolute path to the directory containing the photo file.
	// - baseDir: the base directory where the organized photos will be stored.
	//
	// Returns: None.

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

	// newDirName := "OrganizedPhotos"
	cameraDir := filepath.Join(baseDir, model)
	yearDir := filepath.Join(cameraDir, fmt.Sprintf("%d", year))
	quarterDir := filepath.Join(yearDir, fmt.Sprintf("Q%d", quarter))

	err = os.MkdirAll(quarterDir, os.ModePerm)
	if err != nil {
		log.Printf("Error creating directories for image %s: %v", photoPath, err)
		return
	}

	baseName := filepath.Base(photoPath)
	nameWithoutExt := baseName[:len(baseName)-len(filepath.Ext(baseName))]

	// Move all related files
	files, err := os.ReadDir(absPhotoDir)
	if err != nil {
		log.Println("Error reading directory to find matching files:", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasPrefix(file.Name(), nameWithoutExt) {
			srcFilePath := filepath.Join(absPhotoDir, file.Name())
			destFilePath := filepath.Join(quarterDir, file.Name())
			err = moveFile(srcFilePath, destFilePath)
			if err != nil {
				log.Printf("Error moving file %s to %s: %v", srcFilePath, destFilePath, err)
			} else {
				fmt.Printf("File %s moved to: %s\n", srcFilePath, destFilePath)
			}
		}
	}
}

func getYearQuarter(month time.Month) (year, quarter int) {
	// getYearQuarter returns the year and quarter corresponding to a given month.
	//
	// month: the month for which to determine the year and quarter.
	// year: the year corresponding to the given month.
	// quarter: the quarter corresponding to the given month.

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

func moveFile(src, dest string) error {
	// moveFile moves a file from the source directory to the destination directory.
	//
	// Parameters:
	// - src: the source directory of the file.
	// - dest: the destination directory where the file will be moved.
	//
	// Return type:
	// - error: returns an error if the file cannot be moved.
	err := os.Rename(src, dest)
	if err != nil {
		return err
	}
	return nil
}
