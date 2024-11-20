package main

type TextPrintPayload struct {
	Text string `json:"text"`
}

type TicketPrintPayload struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	Due      string `json:"due"`
	Assigner string `json:"assigner"`
	Link     string `json:"link"`
}
