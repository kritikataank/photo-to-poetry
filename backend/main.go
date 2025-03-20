package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Structs for API requests & responses
type ImageUploadRequest struct {
	Image string `json:"image"` // Base64-encoded image
}

type CaptionResponse struct {
	Caption string `json:"caption"`
}

type ConvertRequest struct {
	Caption string `json:"caption"`
}

type PoemResponse struct {
	Poem string `json:"poem"`
}

var captions = make(map[string]string) // Store captions in memory

func main() {
	r := gin.Default()

	// ✅ CORS configuration
	r.Use(cors.Default())

	// ✅ Test route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Go backend is running!"})
	})

	// ✅ API Routes
	r.POST("/upload", handleImageUpload)      // Upload image + generate caption
	r.GET("/image/:image_name", serveImage)   // Fetch image
	r.GET("/caption/:image_name", getCaption) // Fetch caption
	r.POST("/convert", convertTextToPoetry)   // Convert caption to poetry

	// ✅ Start backend server
	log.Println("🚀 Backend running on http://localhost:8080")
	r.Run(":8080")
}

// Handle image upload + captioning
func handleImageUpload(c *gin.Context) {
	var request ImageUploadRequest

	// Parse JSON request body
	if err := c.BindJSON(&request); err != nil {
		log.Println("❌ Error parsing JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Save image and continue processing
	imageName, imagePath := saveImage(request.Image)
	if imagePath == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// Generate caption
	caption, err := getImageCaption(imagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate caption"})
		return
	}

	// Store caption in memory
	captions[imageName] = caption

	// ✅ Send image details to frontend
	c.JSON(http.StatusOK, gin.H{
		"image_name":  imageName,
		"image_url":   fmt.Sprintf("http://localhost:8080/image/%s", imageName),
		"caption_url": fmt.Sprintf("http://localhost:8080/caption/%s", imageName),
	})
}

// Serve saved image
func serveImage(c *gin.Context) {
	imageName := c.Param("image_name")
	imagePath := filepath.Join("uploads", imageName)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	c.File(imagePath)
}

// Fetch generated caption
func getCaption(c *gin.Context) {
	imageName := c.Param("image_name")

	if caption, exists := captions[imageName]; exists {
		c.JSON(http.StatusOK, gin.H{"caption": caption})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Caption not found"})
	}
}

// Convert caption to poetry
func convertTextToPoetry(c *gin.Context) {
	var request ConvertRequest

	if err := c.BindJSON(&request); err != nil {
		log.Println("❌ Error parsing JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Run Python script for text conversion
	cmd := exec.Command("python", "convert.py")
	cmd.Stdin = strings.NewReader(fmt.Sprintf(`{"caption": "%s"}`, request.Caption))

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println("❌ Error running Python script:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate poetry"})
		return
	}

	// Parse Python response
	var response PoemResponse
	err = json.Unmarshal(out.Bytes(), &response)
	if err != nil {
		log.Println("❌ Error parsing poetry response:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid poetry response"})
		return
	}

	// Return the generated poem
	c.JSON(http.StatusOK, gin.H{"poem": response.Poem})
}

// Save base64 image as a file
func saveImage(base64Data string) (string, string) {
	// Remove Base64 prefix
	if strings.Contains(base64Data, ",") {
		parts := strings.Split(base64Data, ",")
		base64Data = parts[1]
	}

	// Decode Base64
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		log.Println("❌ Error decoding image:", err)
		return "", ""
	}

	// Ensure "uploads" directory exists
	err = os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		log.Println("❌ Error creating uploads directory:", err)
		return "", ""
	}

	// Generate unique filename
	imageName := fmt.Sprintf("image_%d.jpg", time.Now().Unix())
	imagePath := filepath.Join("uploads", imageName)

	// Save image file
	err = os.WriteFile(imagePath, data, 0644)
	if err != nil {
		log.Println("❌ Error saving image:", err)
		return "", ""
	}

	log.Println("✅ Image saved successfully:", imagePath)
	return imageName, imagePath
}

// Call Python script to generate caption
func getImageCaption(imagePath string) (string, error) {
	cmd := exec.Command("python", "caption.py", imagePath)

	// Capture output
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println("❌ Error running Python script:", err)
		return "", err
	}

	// Parse response
	var response CaptionResponse
	err = json.Unmarshal(out.Bytes(), &response)
	if err != nil {
		log.Println("❌ Error parsing caption response:", err)
		return "", err
	}

	return response.Caption, nil
}
