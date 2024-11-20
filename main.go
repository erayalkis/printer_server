package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hennedo/escpos"
	"github.com/joho/godotenv"
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
		p.Write(body.Text)
		p.LineFeed()

		p.Print()

		c.JSON(http.StatusOK, gin.H{
			"message": "Request processed successfully",
		})
	})

	v1.POST("/qr", func(c *gin.Context) {
		var body TextPrintPayload

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		p.QRCode(body.Text, false, 5, escpos.QRCodeErrorCorrectionLevelH)
		p.LineFeed()

		p.Print()

		c.JSON(http.StatusOK, gin.H{
			"message": "Request processed successfully",
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

		img, imgFormat, err := image.Decode(src)

		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "Could not decode image"})
			return
		}

		log.Print("Loaded image, format: ", imgFormat)

		_, err = p.PrintImage(img)

		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "Could not print image"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Request processed successfully",
		})
	})

	v1.POST("/ticket", func(c *gin.Context) {
		var body TicketPrintPayload

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		file, err := os.Open("./cadenza_c.png")

		if err != nil {
			c.JSON(500, gin.H{"error": "Could not open uploaded image"})
			return
		}

		defer file.Close()

		img, imgFormat, err := image.Decode(file)

		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "Could not decode image"})
			return
		}

		log.Print("Loaded image, format: ", imgFormat)

		p.Justify(1)
		_, err = p.PrintImage(img)

		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "Could not print image"})
			return
		}

		p.Size(2, 2)
		p.Bold(true).Write("===CADENZABOX===")

		p.Size(1, 1)
		p.Write("\n")
		p.Write("\n")

		p.Bold(true).Underline(2).Write("=============TITLE=============\n")
		p.Bold(false)
		p.Underline(0).Write(body.Title)
		p.Write("\n")
		p.LineFeed()

		p.Bold(true).Underline(2).Write("==========DESCRIPTION==========\n")
		p.Bold(false)
		p.Underline(0).Write(body.Body)
		p.Write("\n")
		p.LineFeed()

		p.Bold(true).Underline(2).Write("===========DUE DATE============\n")
		p.Bold(false)
		p.Underline(0).Write(body.Due)
		p.Write("\n")
		p.LineFeed()

		p.Bold(true).Underline(2).Write("===========ASSIGNER============\n")
		p.Bold(false)
		p.Underline(0).Write(body.Assigner)
		p.Write("\n")
		p.LineFeed()

		p.QRCode(body.Link, false, 5, escpos.QRCodeErrorCorrectionLevelH)

		p.LineFeed()
		p.LineFeed()

		p.Print()

		c.JSON(http.StatusOK, gin.H{
			"message": "Request processed successfully",
		})
	})

	v1.POST("/tickets", func(c *gin.Context) {
		var body MultiTicketPrintPaylod

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		file, err := os.Open("./cadenza_c.png")

		if err != nil {
			c.JSON(500, gin.H{"error": "Could not open uploaded image"})
			return
		}

		defer file.Close()

		img, imgFormat, err := image.Decode(file)

		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "Could not decode image"})
			return
		}

		log.Print("Loaded image, format: ", imgFormat)

		p.Justify(1)
		_, err = p.PrintImage(img)

		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "Could not print image"})
			return
		}

		p.Size(2, 2)
		p.Bold(true).Write("===CADENZABOX===")

		p.Size(1, 1)
		p.Write("\n")
		p.Write("\n")

		for _, ticket := range body.Tickets {
			p.Bold(true).Underline(2).Write("=============TITLE=============\n")
			p.Bold(false)
			p.Underline(0).Write(ticket.Title)
			p.Write("\n")
			p.LineFeed()

			p.Bold(true).Underline(2).Write("==========DESCRIPTION==========\n")
			p.Bold(false)
			p.Underline(0).Write(ticket.Body)
			p.Write("\n")
			p.LineFeed()

			p.Bold(true).Underline(2).Write("===========DUE DATE============\n")
			p.Bold(false)
			p.Underline(0).Write(ticket.Due)
			p.Write("\n")
			p.LineFeed()

			p.Bold(true).Underline(2).Write("===========ASSIGNER============\n")
			p.Bold(false)
			p.Underline(0).Write(ticket.Assigner)
			p.Write("\n")
			p.LineFeed()

			p.Bold(true).Underline(2).Write("===============================\n")
			p.Bold(false)
			p.Write("\n")
			p.LineFeed()
		}

		p.QRCode(body.Link, false, 5, escpos.QRCodeErrorCorrectionLevelH)

		p.LineFeed()
		p.LineFeed()

		p.Print()

		c.JSON(http.StatusOK, gin.H{
			"message": "Request processed successfully",
		})
	})

	r.Run()
}
