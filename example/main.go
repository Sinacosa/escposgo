package main

import (
	escposgo "escpos-go"
	"fmt"
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

	// Print some message from the user
	var feedback = escposgo.NewDocument()

	feedback.Add(escposgo.NewTitle("New feedback!"))
	feedback.Add(escposgo.NewBody("Hello, this is a simple mesage to tell you your app is amazing. Thank you so much for making it free.").NewLine())
	feedback.Add(escposgo.NewBody("Mathieu."))

	err = printer.PrintDocument(feedback)
	if err != nil {
		fmt.Println(err)
	}
}
