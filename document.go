package escposgo

const PAGE_WIDTH = 32

type Doc interface {
	ToBytes() []byte
}

type Document struct {
	Items []Doc
}

func NewDocument() *Document {
	return &Document{}
}

func (d *Document) Add(item Doc) {
	d.Items = append(d.Items, item)
}

func (d *Document) ToBytes() []byte {
	var b = []byte{}
	b = append(b, ResetCommand.toByteArray()...)
	b = append(b, InitCommand.toByteArray()...)
	for _, item := range d.Items {
		b = append(b, item.ToBytes()...)
	}
	b = append(b, LineFeedCommand.toByteArray()...)
	b = append(b, LineFeedCommand.toByteArray()...)
	b = append(b, LineFeedCommand.toByteArray()...)
	b = append(b, EOTCommand.toByteArray()...)
	return b
}

type Title struct {
	content string
}

func (t *Title) ToBytes() []byte {
	bytes := []byte{}
	bytes = append(bytes, SelectTitleFontSizeCommand.toByteArray()...)
	bytes = append(bytes, []byte(t.content)...)
	bytes = append(bytes, []byte("\n\n")...)
	return bytes
}

func NewTitle(content string) *Title {
	return &Title{
		content: content,
	}
}

type Body struct {
	content string
	newLine bool
}

func NewBody(content string) *Body {
	return &Body{
		content: content,
	}
}

func (b *Body) NewLine() *Body {
	b.newLine = true
	return b
}

func (b *Body) ToBytes() []byte {
	bytes := []byte{}
	bytes = append(bytes, SelectBodyFontSizeCommand.toByteArray()...)
	bytes = append(bytes, []byte(b.content)...)

	if b.newLine {
		bytes = append(bytes, []byte("\n")...)
	}
	return bytes
}

type Table struct {
	Header Row
	Rows   []Row
	Footer Row
}

func (t *Table) ToBytes() []byte {
	// encode the title data
	// measure the spaces
	bytes := []byte{}
	bytes = append(bytes, t.Header.ToBytes()...)
	// add the separator
	bytes = append(bytes, NewSeparator('=').ToBytes()...)
	for _, row := range t.Rows {
		bytes = append(bytes, row.ToBytes()...)
	}
	bytes = append(bytes, NewSeparator('=').ToBytes()...)
	bytes = append(bytes, t.Footer.ToBytes()...)
	return bytes
}

func NewTable(header Row, rows []Row, footer Row) *Table {
	return &Table{
		Header: header,
		Rows:   rows,
		Footer: footer,
	}
}

type Row struct {
	Left  string
	Right string
}

func (r Row) ToBytes() []byte {
	numSpaces := PAGE_WIDTH - (len(r.Left) + len(r.Right))
	bytes := []byte{}
	bytes = append(bytes, []byte(r.Left)...)
	for i := 0; i < numSpaces; i++ {
		bytes = append(bytes, byte(' '))
	}
	bytes = append(bytes, []byte(r.Right)...)
	return bytes
}

func NewRow(left, right string) Row {
	return Row{
		Left:  left,
		Right: right,
	}
}

type Separator struct {
	rune rune
}

func NewSeparator(rune rune) *Separator {
	return &Separator{
		rune: rune,
	}
}

func (s *Separator) ToBytes() []byte {
	bytes := []byte{byte('\n')}
	for i := 0; i < PAGE_WIDTH; i++ {
		bytes = append(bytes, byte(s.rune))
	}
	bytes = append(bytes, byte('\n'))
	return bytes
}
