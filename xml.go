package asciidocxml

import (
	"encoding/xml"
)

const XMLNamespace = "http://asciidoc.org/ns"

// Document is the root element
type Document struct {
	XMLName xml.Name `xml:"http://asciidoc.org/ns asciidoc"`
	DocType string   `xml:"doctype,attr,omitempty"`
	Header  *Header  `xml:"header,omitempty"`
	Content Content  `xml:"content"`
}

// Header contains document metadata
type Header struct {
	Title     string       `xml:"title"`
	Authors   []Author     `xml:"author,omitempty"`
	Revision  *Revision   `xml:"revision,omitempty"`
	Attributes []Attribute `xml:"attribute,omitempty"`
}

// Author represents document author
type Author struct {
	Name  string `xml:"name"`
	Email string `xml:"email,omitempty"`
}

// Revision represents document revision info
type Revision struct {
	Number string `xml:"number,omitempty"`
	Date   string `xml:"date,omitempty"`
	Remark string `xml:"remark,omitempty"`
}

// Attribute represents a document attribute
type Attribute struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// Content is a container for various block elements
type Content struct {
	Items []ContentItem `xml:",any"`
}

// ContentItem represents any content block element
type ContentItem struct {
	XMLName      xml.Name
	Section      *Section      `xml:"section,omitempty"`
	Paragraph    *Paragraph    `xml:"paragraph,omitempty"`
	CodeBlock    *CodeBlock    `xml:"codeblock,omitempty"`
	LiteralBlock *LiteralBlock `xml:"literalblock,omitempty"`
	ListingBlock *ListingBlock `xml:"listingblock,omitempty"`
	Example      *Example      `xml:"example,omitempty"`
	Sidebar      *Sidebar      `xml:"sidebar,omitempty"`
	Quote        *Quote        `xml:"quote,omitempty"`
	Verse        *Verse        `xml:"verse,omitempty"`
	Table        *Table        `xml:"table,omitempty"`
	List         *List         `xml:"list,omitempty"`
	Image        *Image        `xml:"image,omitempty"`
	Video        *Video        `xml:"video,omitempty"`
	Audio        *Audio        `xml:"audio,omitempty"`
	PageBreak    *PageBreak    `xml:"pagebreak,omitempty"`
	ThematicBreak *ThematicBreak `xml:"thematicbreak,omitempty"`
	Admonition   *Admonition   `xml:"admonition,omitempty"`
	Passthrough  *Passthrough  `xml:"passthrough,omitempty"`
}

// Section represents a document section
type Section struct {
	Level    int     `xml:"level,attr"`
	ID       string  `xml:"id,attr,omitempty"`
	Role     string  `xml:"role,attr,omitempty"`
	Numbered *bool   `xml:"numbered,attr,omitempty"`
	Title    InlineContent `xml:"title"`
	Content  Content `xml:"content,omitempty"`
}

// Paragraph represents a paragraph block
type Paragraph struct {
	Role  string        `xml:"role,attr,omitempty"`
	ID    string        `xml:"id,attr,omitempty"`
	Items []InlineItem  `xml:",any"`
}

// CodeBlock represents a source code block
type CodeBlock struct {
	Language     string `xml:"language,attr,omitempty"`
	Source       string `xml:"source,attr,omitempty"`
	LineNums     *bool  `xml:"linenums,attr,omitempty"`
	FirstLineNum *int   `xml:"firstlinenum,attr,omitempty"`
	Highlight    string `xml:"highlight,attr,omitempty"`
	Indent       *int   `xml:"indent,attr,omitempty"`
	Role         string `xml:"role,attr,omitempty"`
	ID           string `xml:"id,attr,omitempty"`
	Title        string `xml:"title,attr,omitempty"`
	Content      string `xml:",chardata"`
}

// LiteralBlock represents a literal/preformatted block
type LiteralBlock struct {
	Style string `xml:"style,attr,omitempty"`
	Role  string `xml:"role,attr,omitempty"`
	ID    string `xml:"id,attr,omitempty"`
	Content string `xml:",chardata"`
}

// ListingBlock represents a listing block
type ListingBlock struct {
	Language string `xml:"language,attr,omitempty"`
	Source   string `xml:"source,attr,omitempty"`
	LineNums *bool  `xml:"linenums,attr,omitempty"`
	Role     string `xml:"role,attr,omitempty"`
	ID       string `xml:"id,attr,omitempty"`
	Title    string `xml:"title,attr,omitempty"`
	Content  string `xml:",chardata"`
}

// Example represents an example block
type Example struct {
	Role    string  `xml:"role,attr,omitempty"`
	ID      string  `xml:"id,attr,omitempty"`
	Title   string  `xml:"title,omitempty"`
	Content Content `xml:"content"`
}

// Sidebar represents a sidebar block
type Sidebar struct {
	Role    string  `xml:"role,attr,omitempty"`
	ID      string  `xml:"id,attr,omitempty"`
	Title   string  `xml:"title,omitempty"`
	Content Content `xml:"content"`
}

// Quote represents a quote block
type Quote struct {
	Role       string        `xml:"role,attr,omitempty"`
	ID         string        `xml:"id,attr,omitempty"`
	Attribution InlineContent `xml:"attribution,omitempty"`
	Citation   string        `xml:"citation,omitempty"`
	Content    Content       `xml:"content"`
}

// Verse represents a verse block
type Verse struct {
	Role       string        `xml:"role,attr,omitempty"`
	ID         string        `xml:"id,attr,omitempty"`
	Attribution InlineContent `xml:"attribution,omitempty"`
	Citation   string        `xml:"citation,omitempty"`
	Content    Content       `xml:"content"`
}

// Table represents a table
type Table struct {
	Frame   string    `xml:"frame,attr,omitempty"`
	Grid    string    `xml:"grid,attr,omitempty"`
	Stripes string    `xml:"stripes,attr,omitempty"`
	Role    string    `xml:"role,attr,omitempty"`
	ID      string    `xml:"id,attr,omitempty"`
	Cols    string    `xml:"cols,attr,omitempty"`
	Title   string    `xml:"title,omitempty"`
	Header  *TableRow `xml:"header,omitempty"`
	Rows    []TableRow `xml:"row"`
}

// TableRow represents a table row
type TableRow struct {
	Cells []TableCell `xml:"cell"`
}

// TableCell represents a table cell
type TableCell struct {
	ColSpan *int         `xml:"colspan,attr,omitempty"`
	RowSpan *int         `xml:"rowspan,attr,omitempty"`
	Align   string       `xml:"align,attr,omitempty"`
	VAlign  string       `xml:"valign,attr,omitempty"`
	Style   string       `xml:"style,attr,omitempty"`
	Role    string       `xml:"role,attr,omitempty"`
	Items   []InlineItem `xml:",any"`
}

// List represents a list (ordered, unordered, labeled, or callout)
type List struct {
	Style  string     `xml:"style,attr"`
	Marker string     `xml:"marker,attr,omitempty"`
	Start  *int       `xml:"start,attr,omitempty"`
	Role   string     `xml:"role,attr,omitempty"`
	ID     string     `xml:"id,attr,omitempty"`
	Items  []ListItem `xml:"item"`
}

// ListItem represents a list item
type ListItem struct {
	Marker string        `xml:"marker,attr,omitempty"`
	Term   InlineContent `xml:"term,omitempty"`
	Items  []ListItemContentItem `xml:",any"`
}

// ListItemContentItem represents content within a list item
type ListItemContentItem struct {
	XMLName      xml.Name
	InlineItems  []InlineItem `xml:",any"`
	List         *List        `xml:"list,omitempty"`
	Paragraph    *Paragraph   `xml:"paragraph,omitempty"`
	CodeBlock    *CodeBlock   `xml:"codeblock,omitempty"`
	LiteralBlock *LiteralBlock `xml:"literalblock,omitempty"`
	Table        *Table       `xml:"table,omitempty"`
}

// Image represents an image
type Image struct {
	Src    string `xml:"src,attr"`
	Alt    string `xml:"alt,attr,omitempty"`
	Title  string `xml:"title,attr,omitempty"`
	Width  string `xml:"width,attr,omitempty"`
	Height string `xml:"height,attr,omitempty"`
	Link   string `xml:"link,attr,omitempty"`
	Role   string `xml:"role,attr,omitempty"`
	ID     string `xml:"id,attr,omitempty"`
}

// Video represents a video
type Video struct {
	Src      string `xml:"src,attr"`
	Poster   string `xml:"poster,attr,omitempty"`
	Width    string `xml:"width,attr,omitempty"`
	Height   string `xml:"height,attr,omitempty"`
	Autoplay *bool  `xml:"autoplay,attr,omitempty"`
	Loop     *bool  `xml:"loop,attr,omitempty"`
	Controls *bool  `xml:"controls,attr,omitempty"`
	Role     string `xml:"role,attr,omitempty"`
	ID       string `xml:"id,attr,omitempty"`
}

// Audio represents an audio element
type Audio struct {
	Src      string `xml:"src,attr"`
	Autoplay *bool  `xml:"autoplay,attr,omitempty"`
	Loop     *bool  `xml:"loop,attr,omitempty"`
	Controls *bool  `xml:"controls,attr,omitempty"`
	Role     string `xml:"role,attr,omitempty"`
	ID       string `xml:"id,attr,omitempty"`
}

// PageBreak represents a page break
type PageBreak struct{}

// ThematicBreak represents a thematic break (horizontal rule)
type ThematicBreak struct{}

// Admonition represents an admonition block
type Admonition struct {
	Type    string  `xml:"type,attr"`
	Role    string  `xml:"role,attr,omitempty"`
	ID      string  `xml:"id,attr,omitempty"`
	Title   string  `xml:"title,omitempty"`
	Content Content `xml:"content"`
}

// Passthrough represents a passthrough block
type Passthrough struct {
	Subs    string `xml:"subs,attr,omitempty"`
	Content string `xml:",chardata"`
}

// InlineContent represents inline content (mixed content)
type InlineContent struct {
	Items []InlineItem `xml:",any"`
}

// InlineItem represents any inline element
type InlineItem struct {
	XMLName     xml.Name
	Text        string      `xml:"text,omitempty"`
	Strong      *InlineContent `xml:"strong,omitempty"`
	Emphasis    *InlineContent `xml:"emphasis,omitempty"`
	Monospace   *InlineContent `xml:"monospace,omitempty"`
	Superscript *InlineContent `xml:"superscript,omitempty"`
	Subscript   *InlineContent `xml:"subscript,omitempty"`
	Mark        *InlineContent `xml:"mark,omitempty"`
	Link        *Link          `xml:"link,omitempty"`
	XRef        *XRef          `xml:"xref,omitempty"`
	Image       *Image         `xml:"image,omitempty"`
	Kbd         *InlineContent `xml:"kbd,omitempty"`
	Button      *InlineContent `xml:"button,omitempty"`
	Menu        *Menu          `xml:"menu,omitempty"`
	Attribute   *AttributeRef  `xml:"attribute,omitempty"`
	CharData    string         `xml:",chardata"`
}

// Link represents a hyperlink
type Link struct {
	Href   string        `xml:"href,attr"`
	Title  string        `xml:"title,attr,omitempty"`
	Role   string        `xml:"role,attr,omitempty"`
	Window string        `xml:"window,attr,omitempty"`
	Items  []InlineItem  `xml:",any"`
}

// XRef represents a cross-reference
type XRef struct {
	RefID string        `xml:"refid,attr"`
	Path  string        `xml:"path,attr,omitempty"`
	Items []InlineItem  `xml:",any"`
}

// Menu represents a menu path
type Menu struct {
	MenuItems []string `xml:"menuitem"`
}

// AttributeRef represents an attribute reference
type AttributeRef struct {
	Name string `xml:"name,attr"`
}

