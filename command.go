package escposgo

const ESC = byte(0x1b)
const EOT = byte(0x04)
const GS = byte(0x1D)

// Reset command
var ResetCommand = NewCommand([]byte{ESC, byte(0x40)})

// Init command
var InitCommand = NewCommand([]byte{ESC, byte('@')})

// line feed command
var LineFeedCommand = NewCommand([]byte{0x0A})

// Justify left command
var AlignLeft = NewCommand([]byte{ESC, byte('a'), byte(0)})

// Justify center command
var AlignCenter = NewCommand([]byte{ESC, byte('a'), byte(1)})

// Justify right command
var AlignRight = NewCommand([]byte{ESC, byte('a'), byte(2)})

// End of transmission command
var EOTCommand = NewCommand([]byte{ESC, EOT})

// var textCommand = NewCommand([]byte(message + " \n\n"))
var eotCommand = NewCommand([]byte{ESC, EOT})

// Select Title font size
var SelectTitleFontSizeCommand = NewCommand([]byte{ESC, byte('!'), 0x20})

// Select body font size
var SelectBodyFontSizeCommand = NewCommand([]byte{ESC, byte('!'), 0x00})

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

func MergeCommands(commands []Command) []byte {
	var data []byte

	for _, c := range commands {
		data = append(data, c.toByteArray()...)
	}

	return data
}
