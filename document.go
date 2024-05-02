package escposgo

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
	// append small font command
	// append justify right
	bytes := []byte{}
	bytes = append(bytes, SelectBodyFontSizeCommand.toByteArray()...)
	bytes = append(bytes, []byte(b.content)...)

	if b.newLine {
		bytes = append(bytes, []byte("\n")...)
	}
	return bytes
}
