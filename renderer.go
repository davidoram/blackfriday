package blackfriday

import (
	"bytes"
	//"fmt"
	"github.com/davidoram/blackfriday/flower"
	"strings"
)

// Flower is a type that implements the Renderer interface
//
// Do not create this directly, instead use the FlowerRenderer function.
type Flower struct {
	renderer    Renderer                    // the renderer that we are wrapping,
	interpreter *flower.StandardInterpreter // parses code blocks for flower directives
}

// WrappedRenderer creates and configures an Renderer object, which
// satisfies the Renderer interface,  intercepts calls to parse code
// and interprets any embedded flower commands, returning the output
// as markdown which is passed on to the wrapped
// Renderer which performs the end user rendering
func WrappedRenderer(wrapped_renderer Renderer) Renderer {

	return &Flower{
		renderer:    wrapped_renderer,
		interpreter: flower.NewInterpreter(),
	}
}

// Surround HTML code with tags that can be used to identify and style the flower command contained within
func (options *Flower) CommandTagStart(out *bytes.Buffer, command flower.Command) {
	options.renderer.CommandTagStart(out, command)
}

func (options *Flower) CommandTagEnd(out *bytes.Buffer, command flower.Command) {
	options.renderer.CommandTagEnd(out, command)
}


// block-level callbacks
func (options *Flower) BlockCodeStart(out *bytes.Buffer, text []byte, lang string) {
	options.renderer.BlockCodeStart(out, text, lang)
}

func (options *Flower) BlockCodeBody(out *bytes.Buffer, text []byte, lang string) {
	lines := strings.Split(string(text[:]), "\n")
	for _, line := range lines {
		command := options.interpreter.EvaluateCode(line)
		options.renderer.CommandTagStart(out, command)
		options.renderer.BlockCodeBody(out, []byte(line), lang)
		options.renderer.CommandTagEnd(out, command)
	}
}

func (options *Flower) BlockCodeEnd(out *bytes.Buffer, text []byte, lang string) {
	options.renderer.BlockCodeEnd(out, text, lang)
}

func (options *Flower) BlockQuote(out *bytes.Buffer, text []byte) {
	options.renderer.BlockQuote(out, text)
}
func (options *Flower) BlockHtml(out *bytes.Buffer, text []byte) {
	options.renderer.BlockHtml(out, text)
}

func (options *Flower) Header(out *bytes.Buffer, text func() bool, level int) {
	options.renderer.Header(out, text, level)
}

func (options *Flower) HRule(out *bytes.Buffer) {
	options.renderer.HRule(out)
}

func (options *Flower) List(out *bytes.Buffer, text func() bool, flags int) {
	options.renderer.List(out, text, flags)
}

func (options *Flower) ListItem(out *bytes.Buffer, text []byte, flags int) {
	options.renderer.ListItem(out, text, flags)
}

func (options *Flower) Paragraph(out *bytes.Buffer, text func() bool) {
	options.renderer.Paragraph(out, text)
}

func (options *Flower) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
	options.renderer.Table(out, header, body, columnData)
}

func (options *Flower) TableRow(out *bytes.Buffer, text []byte) {
	options.renderer.TableRow(out, text)
}

func (options *Flower) TableCell(out *bytes.Buffer, text []byte, flags int) {
	options.renderer.TableCell(out, text, flags)
}

func (options *Flower) Footnotes(out *bytes.Buffer, text func() bool) {
	options.renderer.Footnotes(out, text)
}

func (options *Flower) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
	options.renderer.FootnoteItem(out, name, text, flags)
}

// Span-level callbacks
func (options *Flower) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	options.renderer.AutoLink(out, link, kind)
}

func (options *Flower) CodeSpanStart(out *bytes.Buffer, text []byte) {
	options.renderer.CodeSpanStart(out, text)
}

func (options *Flower) CodeSpanBody(out *bytes.Buffer, text []byte) {
	lines := strings.Split(string(text[:]), "\n")
	for _, line := range lines {
		command := options.interpreter.EvaluateCode(line)
		options.renderer.CommandTagStart(out, command)
		options.renderer.CodeSpanBody(out, []byte(line))
		options.renderer.CommandTagEnd(out, command)
	}
}

func (options *Flower) CodeSpanEnd(out *bytes.Buffer, text []byte) {
	options.renderer.CodeSpanStart(out, text)
}

func (options *Flower) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	options.renderer.DoubleEmphasis(out, text)
}
func (options *Flower) Emphasis(out *bytes.Buffer, text []byte) {
	options.renderer.Emphasis(out, text)
}
func (options *Flower) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	options.renderer.Image(out, link, title, alt)
}
func (options *Flower) LineBreak(out *bytes.Buffer) {
	options.renderer.LineBreak(out)
}
func (options *Flower) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	options.renderer.Link(out, link, title, content)
}
func (options *Flower) RawHtmlTag(out *bytes.Buffer, tag []byte) {
	options.renderer.RawHtmlTag(out, tag)
}
func (options *Flower) TripleEmphasis(out *bytes.Buffer, text []byte) {
	options.renderer.TripleEmphasis(out, text)
}
func (options *Flower) StrikeThrough(out *bytes.Buffer, text []byte) {
	options.renderer.StrikeThrough(out, text)
}
func (options *Flower) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
	options.renderer.FootnoteRef(out, ref, id)
}

// Low-level callbacks
func (options *Flower) Entity(out *bytes.Buffer, entity []byte) {
	options.renderer.Entity(out, entity)
}
func (options *Flower) NormalText(out *bytes.Buffer, text []byte) {
	options.renderer.NormalText(out, text)
}

// Header and footer
func (options *Flower) DocumentHeader(out *bytes.Buffer) {
	options.renderer.DocumentHeader(out)
}
func (options *Flower) DocumentFooter(out *bytes.Buffer) {
	options.renderer.DocumentFooter(out)
}
