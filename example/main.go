package main

import (
	escposgo "escpos-go"
	"log"
	"log/slog"
)

func main() {
	// get the printer
	var printer, err = escposgo.NewPrinter(0x0456, 0x0808)
	if err != nil {
		log.Fatalf("could not open the device")
	}
	// defer the closing of the printer
	defer printer.Close()

	// Create a new document to print
	document := escposgo.NewDocument()
	document.Add(escposgo.NewTitle("New purchase!"))
	document.Add(escposgo.NewBody("This is my second body").NewLine())
	document.Add(escposgo.NewBody("And this is the third one...").NewLine())
	rows := []escposgo.Row{
		{Left: "Monthly Sub", Right: "199.99"},
		{Left: "Monthly Sub", Right: "49.99"},
		{Left: "Monthly Sub", Right: "99.99"},
	}
	table := escposgo.NewTable(escposgo.NewRow("Name", "Price"), rows, escposgo.NewRow("Total: ", "349,97"))
	document.Add(table)

	// print the document
	err = printer.PrintDocument(document)
	if err != nil {
		slog.Info("error occured while printing", "err", err)
	}
}
