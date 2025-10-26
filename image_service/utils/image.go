package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func SaveImageToStorage(imageData []byte, imageName string) (string, error) {
	log.Println("Saving image to storage...")
	storagePath := os.Getenv("IMAGE_STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage/images"
	}

	//Add create directory if not exists
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		log.Printf("Storage directory does not exist, creating: %s", storagePath)
		err := os.MkdirAll(storagePath, os.ModePerm)
		if err != nil {
			log.Printf("Error creating storage directory: %v", err)
			return "", err
		}
	}

	// Save image to storage folder and add path YMD
	today := time.Now().Format("2006/01/02")
	dateStoragePath := filepath.Join(storagePath, today)

	// Create date directory if not exists
	if _, err := os.Stat(dateStoragePath); os.IsNotExist(err) {
		log.Printf("Date storage directory does not exist, creating: %s", dateStoragePath)
		err := os.MkdirAll(dateStoragePath, os.ModePerm)
		if err != nil {
			log.Printf("Error creating date storage directory: %v", err)
			return "", err
		}
	}

	imagePath := filepath.Join(dateStoragePath, imageName)

	//Check if file already exists, if yes add timestamp to filename
	if _, err := os.Stat(imagePath); err == nil {
		timestamp := time.Now().Unix()
		ext := filepath.Ext(imageName)
		nameOnly := imageName[0 : len(imageName)-len(ext)]
		imageName = fmt.Sprintf("%s_%d%s", nameOnly, timestamp, ext)
		imagePath = filepath.Join(dateStoragePath, imageName)
	}
	//Write file
	err := os.WriteFile(imagePath, imageData, 0644)
	if err != nil {
		log.Printf("Error saving image to storage: %v", err)
		return "", err
	}
	log.Printf("Image saved successfully: %s", imagePath)

	//Image url path
	imageEndpoint := os.Getenv("APP_IMAGE_URL")
	if imageEndpoint == "" {
		imageEndpoint = "http://localhost:8081"
	}
	fullImagePath := imageEndpoint + "/" + filepath.ToSlash(imagePath)
	return fullImagePath, nil
}
