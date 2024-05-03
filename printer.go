package escposgo

import (
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/google/gousb"
)

// minimal implementatio of a printer
type Printer struct {
	VendorID     gousb.ID
	ProductID    gousb.ID
	Device       *gousb.Device
	writer       io.Writer
	ctx          *gousb.Context
	doneCallback func()
	intf         *gousb.Interface
}

// this is how we close the printer
func (p *Printer) Close() {
	p.Device.Close()
	p.ctx.Close()
	p.doneCallback()
}

// NewPrinter returns a new printer for a given vendor and product productID
// It returns an error if the printer is not found or if the endpoint is not available
func NewPrinter(vendorID, productID gousb.ID) (*Printer, error) {

	// Initialize a new Context.
	ctx := gousb.NewContext()

	// Open any device with a given VID/PID using a convenience function.
	dev, err := ctx.OpenDeviceWithVIDPID(vendorID, productID)
	if err != nil {
		return nil, err
	}

	// Claim the default interface using a convenience function.
	// The default interface is always #0 alt #0 in the currently active
	// config.
	intf, done, err := dev.DefaultInterface()
	if err != nil {
		return nil, err
	}

	// Open an OUT endpoint.
	ep, err := intf.OutEndpoint(3)
	if err != nil {
		return nil, err
	}

	printer := &Printer{
		VendorID:     vendorID,
		ProductID:    productID,
		Device:       dev,
		writer:       ep,
		ctx:          ctx,
		doneCallback: done,
		intf:         intf,
	}

	return printer, nil
}

func (p *Printer) Print(data []byte) error {
	expected := len(data)
	n, err := p.writer.Write(data)
	if err != nil {
		slog.Error("could not write to the printer", "err", err)
		return err
	}

	if expected != n {
		return errors.New(fmt.Sprintf("expected %d bytes, actuall %d bytes", expected, n))
	}
	slog.Info("data written to the writer", "len", n)

	return nil
}

func (p *Printer) SimplePrint(message string) {

	RESET := []byte{ESC, byte(0x40)}
	resetCommand := NewCommand(RESET)
	initCommand := NewCommand([]byte{ESC, byte('@')})
	lineFeedCommand := NewCommand([]byte{0x0A})
	textCommand := NewCommand([]byte(message + " \n\n"))
	eotCommand := NewCommand([]byte{ESC, EOT})

	commands := []Command{
		resetCommand,
		initCommand,
		textCommand,
		lineFeedCommand,
		eotCommand,
	}

	p.Print(mergeCommands(commands))
}

func (p *Printer) PrintDocument(doc *Document) error {
	return p.Print(doc.ToBytes())
}
