package main

import (
	"fmt"

	"github.com/kenshaw/escpos"
)

func printText(p *escpos.Escpos, body TextPrintPayload) {
	p.Init()

	p.SetSmooth(1)
	p.SetFont("A")

	parsedBody := fmt.Sprintf("%s\n", body.Text)
	p.Write(parsedBody)

	p.Formfeed()

	p.End()
}
