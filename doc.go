package dox2go

type PageSize int32

const (
	PS_A1 PageSize = iota
	PS_A2
	PS_A3
	PS_A4
	PS_A5
)

type PageOrientation int32

const (
	PO_Landscape PageOrientation = iota
	PO_Portrait
)

type Document interface {
	CreatePage(pu Unit, ps PageSize, po PageOrientation) Page

	CreateFont(name string, fs FontStyle, size float64) Font

	Close() error
}

type Page interface {
	Surface() Surface
}
