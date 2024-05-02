package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"

	"github.com/google/gousb"
)

type Command struct {
	data []byte
}

func NewCommand(b []byte) Command {
	return Command{
		data: b,
	}
}

func (cmd Command) toByteArray() []byte {
	return cmd.data
}

func sendCommand(cmd Command, writer io.Writer) error {
	return writeData(cmd.toByteArray(), writer)
}

func mergeCommands(commands []Command) []byte {
	var data []byte

	for _, c := range commands {
		data = append(data, c.toByteArray()...)
	}

	return data
}

func writeData(data []byte, writer io.Writer) error {
	expected := len(data)
	n, err := writer.Write(data)
	if err != nil {
		slog.Error("could not write to the printer", "err", err)
		return err
	}

	fmt.Println("EXPECTED", expected)

	if expected != n {
		return errors.New(fmt.Sprintf("expected %d bytes, actuall %d bytes", expected, n))
	}
	slog.Info("data written to the writer", "len", n)

	return nil
}

// ESC 0x1B
const ESC = byte(0x1b)
const EOT = byte(0x04)

func main() {
	// Initialize a new Context.
	ctx := gousb.NewContext()
	defer ctx.Close()

	//var printer, err = escposgo.NewPrinter(0x0456, 0x0808, ctx)
	//	if err != nil {
	//		log.Fatalf("could not open the device")
	//	}
	//	defer printer.Close()
	//	printer.Print([]byte{ESC, '@'})

	// Open any device with a given VID/PID using a convenience function.
	dev, err := ctx.OpenDeviceWithVIDPID(0x0456, 0x0808)
	if err != nil {
		log.Fatalf("Could not open a device: %v", err)
	}
	defer dev.Close()

	// Claim the default interface using a convenience function.
	// The default interface is always #0 alt #0 in the currently active
	// config.
	intf, done, err := dev.DefaultInterface()
	if err != nil {
		log.Fatalf("%s.DefaultInterface(): %v", dev, err)
	}
	defer done()

	// Open an OUT endpoint.
	ep, err := intf.OutEndpoint(3)
	if err != nil {
		log.Fatalf("%s.OutEndpoint(7): %v", intf, err)
	}

	// Generate some data to write.
	RESET := []byte{ESC, byte(0x40)}
	resetCommand := NewCommand(RESET)
	initCommand := NewCommand([]byte{ESC, byte('@')})
	// lineFeedCommand := NewCommand([]byte{LF})
	textCommand := NewCommand([]byte("coucou mes loulous. Let's see how it is behaving... \n\n"))
	eotCommand := NewCommand([]byte{ESC, EOT})
	alignLeft := NewCommand([]byte{ESC, byte('a'), byte(0)})
	alignCenter := NewCommand([]byte{ESC, byte('a'), byte(1)})
	alignRight := NewCommand([]byte{ESC, byte('a'), byte(2)})

	commands := []Command{
		resetCommand,
		initCommand,
		alignRight,
		textCommand,
		alignCenter,
		textCommand,
		alignLeft,
		textCommand,
		eotCommand}

	ep.Write(mergeCommands(commands))
	_ = commands
	// printer.Print(mergeCommands(commands))
}
