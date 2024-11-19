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
		printText(c, p)
		c.JSON(http.StatusOK, gin.H{
			"message": "done",
		})
	})

	r.Run()
}
