package blackfriday

import (
	"bytes"
	"github.com/davidoram/flower"
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

// block-level callbacks
func (options *Flower) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	options.interpreter.EvaluateCode(text)
	options.renderer.BlockCode(out, text, lang)
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

func (options *Flower) CodeSpan(out *bytes.Buffer, text []byte) {
	options.interpreter.EvaluateCode(text)
	options.renderer.CodeSpan(out, text)
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
	options.renderer.Entity(out, options.interpreter.SummaryReport())
	options.renderer.DocumentFooter(out)
}
