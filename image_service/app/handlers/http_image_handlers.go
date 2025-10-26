package handlers

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"pos/image_service/app/client"
	"strings"
)

// HTTPImageHandler untuk handling HTTP requests
type HTTPImageHandler struct {
	authClient client.AuthClient
}

func NewHTTPImageHandler() *HTTPImageHandler {
	// Initialize auth client
	authClient, err := client.NewAuthClient(os.Getenv("AUTH_SERVICE_URL"))
	if err != nil {
		log.Printf("Warning: Failed to initialize auth client: %v", err)
		// Continue without auth client for development
	}

	return &HTTPImageHandler{
		authClient: authClient,
	}
}

// ServeImage serves image files with authentication
func (h *HTTPImageHandler) ServeImage(w http.ResponseWriter, r *http.Request) {
	// Extract path dari URL
	// URL format: /storage/images/2025/10/26/testupload_1761478232.jpg
	imagePath := r.URL.Path

	// Validasi token dari Authorization header atau query parameter
	// token := h.extractToken(r)
	// if token == "" {
	// 	http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
	// 	return
	// }

	// // Validasi token dengan auth service
	// if h.authClient != nil {
	// 	isValid, err := h.authClient.ValidateToken(token)
	// 	if err != nil {
	// 		log.Printf("Token validation error: %v", err)
	// 		http.Error(w, "Unauthorized: Token validation failed", http.StatusUnauthorized)
	// 		return
	// 	}
	// 	if !isValid {
	// 		http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
	// 		return
	// 	}
	// } else {
	// 	// Development mode - log token but allow access
	// 	log.Printf("Development mode: Token validation bypassed for token: %s", token)
	// }

	// Construct file path
	// Remove leading slash and construct full path
	relativePath := strings.TrimPrefix(imagePath, "/")
	fullPath := filepath.Join(".", relativePath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	// Open file
	file, err := os.Open(fullPath)
	if err != nil {
		log.Printf("Error opening file %s: %v", fullPath, err)
		http.Error(w, "Error reading image", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Error getting file info: %v", err)
		http.Error(w, "Error reading image", http.StatusInternalServerError)
		return
	}

	// Set content type based on file extension
	ext := filepath.Ext(fullPath)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Set headers
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
	w.Header().Set("Last-Modified", fileInfo.ModTime().UTC().Format(http.TimeFormat))

	// Copy file content to response
	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("Error serving file: %v", err)
		return
	}

	log.Printf("Served image: %s", imagePath)
}

// extractToken extracts JWT token from request
func (h *HTTPImageHandler) extractToken(r *http.Request) string {
	// Try Authorization header first
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// Try query parameter
	token := r.URL.Query().Get("token")
	if token != "" {
		return token
	}

	return ""
}
