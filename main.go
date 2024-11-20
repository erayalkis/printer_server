package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kenshaw/escpos"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("Cannot load env! This should never happen(?)")
	}

	printerPath, exists := os.LookupEnv("POS_PATH")

	if !exists {
		panic("Cannot start server without a printer path present")
	}

	f, err := os.OpenFile(printerPath, os.O_RDWR, 0)

	if err != nil {
		panic(fmt.Errorf("error while opening printer path %s", err))
	}

	defer f.Close()

	p := escpos.New(f)

	r := gin.Default()

	v1 := r.Group("/v1")

	v1.POST("/text", func(c *gin.Context) {
		var body TextPrintPayload

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		printText(p, body)

		c.JSON(http.StatusOK, gin.H{
			"message": "Request processed successfully",
			"data":    body,
		})
	})

	v1.POST("/image", func(c *gin.Context) {
		file, err := c.FormFile("image")

		if err != nil {
			c.JSON(400, gin.H{"error": "Image is required"})
			return
		}

		src, err := file.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": "Could not open uploaded image"})
			return
		}
		defer src.Close()

		// Preprocess the image
		bitmap, err := preprocessImage(src, 384) // Use 384px width for a 58mm printer
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to process image"})
			return
		}

		// Print the image
		p.Init()
		_, err = p.WriteRaw(bitmap)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to send image to printer"})
			return
		}

		// Cut the paper and finish
		p.Formfeed()
		p.Cut()
		p.End()

		// Success response
		c.JSON(200, gin.H{"message": "Image printed successfully"})

	})

	r.Run()
}
