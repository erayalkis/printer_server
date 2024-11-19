package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kenshaw/escpos"
)

func printText(c *gin.Context, p *escpos.Escpos) {
	p.Init()

	p.SetSmooth(1)
	p.SetFont("A")

	p.Write("Hello, World!\n")
	p.SetFont("B")
	p.Write("Hello, World! In a new font!")

	p.Formfeed()

	p.End()
}
