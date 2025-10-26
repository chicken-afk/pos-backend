package services

import (
	"os"
	"pos/image_service/utils"
)

type ImageService interface {
	UploadImage(userID string, imageData []byte, imageName string) (string, error)
	// GetImageURL(imageID string) (string, error)
}

type imageService struct {
	appURL string
}

func NewImageService() ImageService {
	appURL := os.Getenv("APP_URL")
	return &imageService{
		appURL: appURL,
	}
}

func (s *imageService) UploadImage(userID string, imageData []byte, imageName string) (string, error) {

	//Save image to storage folder dont use repository for this
	imageURL, err := utils.SaveImageToStorage(imageData, imageName)

	if err != nil {
		return "", err
	}
	return imageURL, nil
}

// func (s *imageService) GetImageURL(imageID string) (string, error) {
// 	exists, err := s.imageRepo.ImageExists(imageID)
// 	if err != nil {
// 		return "", err
// 	}
// 	if !exists {
// 		return "", utils.ErrImageNotFound
// 	}
// 	imageURL := s.appURL + "/images/" + imageID
// 	return imageURL, nil
// }
