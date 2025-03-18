package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// Request struct to parse JSON data
type ImageUploadRequest struct {
	Image string `json:"image"` // Base64-encoded image
}

func main() {
	r := gin.Default()

	// ✅ CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// ✅ Test route
	r.GET("/", func(c *gin.Context) {
		log.Println("📡 GET / - Backend is running")
		c.JSON(http.StatusOK, gin.H{"message": "Go backend is running!"})
	})

	// ✅ API to receive image from frontend
	r.POST("/upload", handleImageUpload)

	// ✅ New route to show success message after upload
	r.GET("/uploaded", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully! You have been redirected here."})
	})

	// ✅ Start server
	log.Println("🚀 Backend running on http://localhost:8080")
	r.Run(":8080")
}

// ✅ Image upload handler with redirection
func handleImageUpload(c *gin.Context) {
	var request ImageUploadRequest

	// Log incoming request
	log.Println("🔄 Receiving image upload request...")

	// Parse JSON request body
	if err := c.BindJSON(&request); err != nil {
		log.Println("❌ Error parsing JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// ✅ Log success
	log.Println("✅ Image received successfully!")

	// ✅ Redirect user to "/uploaded"
	c.Redirect(http.StatusSeeOther, "/uploaded")
}