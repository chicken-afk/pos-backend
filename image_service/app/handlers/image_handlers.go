package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	services "pos/image_service/app/service"
	"pos/image_service/pb"
	"strings"
)

type ImageHandler struct {
	pb.UnimplementedImageServiceServer
	imageService services.ImageService
}

func NewImageHandler(imageService services.ImageService) *ImageHandler {
	return &ImageHandler{
		imageService: imageService,
	}
}

func (h *ImageHandler) UploadImage(ctx context.Context, req *pb.UploadImageRequest) (*pb.UploadImageResponse, error) {
	imageBase64 := req.GetImageData() // string base64
	imageName := req.GetImageName()

	// Remove data URL prefix if present (e.g., "data:image/jpeg;base64,")
	if strings.Contains(imageBase64, ",") {
		parts := strings.SplitN(imageBase64, ",", 2)
		imageBase64 = parts[1]
	}

	// Decode base64 string to bytes
	imageBytes, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		log.Printf("Error decoding base64 image data: %v", err)
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Upload to image service
	imageURL, err := h.imageService.UploadImage("", imageBytes, imageName)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}

	return &pb.UploadImageResponse{ImageUrl: imageURL}, nil
}
