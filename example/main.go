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
	defer printer.Close()

	document := escposgo.NewDocument()
	document.Add(escposgo.NewTitle("New purchase!"))
	document.Add(escposgo.NewBody("This is my second body").NewLine())
	document.Add(escposgo.NewBody("And this is the third one...").NewLine())
	err = printer.PrintDocument(document)
	if err != nil {
		slog.Info("error occured while printing", "err", err)
	}
}
