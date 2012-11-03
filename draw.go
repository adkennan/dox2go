package dox2go

type Point struct {
	X float64
	Y float64
}

func (p *Point) ChangeUnit(from Unit, to Unit) Point {
	return Point{
		ConvertUnit(p.X, from, to),
		ConvertUnit(p.Y, from, to)}
}

type Color struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

func RGB(r uint8, g uint8, b uint8) Color {
	return Color{r, g, b, 255}
}

func RGBA(r uint8, g uint8, b uint8, a uint8) Color {
	return Color{r, g, b, a}
}

type LineCapStyle int32

const (
	LC_ButtCap   LineCapStyle = 0
	LC_RoundCap               = 1
	LC_SquareCap              = 2
)

type LineJoinStyle int32

const (
	LJ_MitreJoin LineJoinStyle = iota
	LJ_RoundJoin
	LJ_BevelJoin
)

type FontStyle int32

const (
	FS_Regular FontStyle = 0
	FS_Bold    FontStyle = 1
	FS_Italic  FontStyle = 2
)

type Font interface {
	Id() int
	Style() FontStyle
	Size() float64
}

type pathCmdType int32

const (
	moveCmdType pathCmdType = iota
	lineCmdType
	curveCmdType
	rectCmdType
	closeCmdType
)

type pathCmdPoints [3]Point

type pathCmd struct {
	t pathCmdType
	p pathCmdPoints
}

type Path struct {
	elements []pathCmd
}

func NewPath() *Path {
	p := new(Path)
	p.elements = make([]pathCmd, 0, 8)
	return p
}

func (path *Path) Move(to Point) {
	path.elements = append(path.elements, pathCmd{moveCmdType, pathCmdPoints{to}})
}

func (path *Path) Line(to Point) {
	path.elements = append(path.elements, pathCmd{lineCmdType, pathCmdPoints{to}})
}

func (path *Path) Curve(control1 Point, control2 Point, to Point) {
	path.elements = append(path.elements, pathCmd{curveCmdType, pathCmdPoints{control1, control2, to}})
}

func (path *Path) Rect(from Point, to Point) {
	path.elements = append(path.elements, pathCmd{rectCmdType, pathCmdPoints{from, to}})
}

func (path *Path) Close() {
	path.elements = append(path.elements, pathCmd{closeCmdType, pathCmdPoints{}})
}

type Surface interface {
	PushState()
	PopState()

	Rotate(byRadians float64)
	Skew(xRadians float64, yRadians float64)
	Translate(byDistance Point)
	Scale(xScale float64, yScale float64)

	Fg(color Color)
	Bg(color Color)
	LineWidth(width float64)
	LineCap(capStyle LineCapStyle)
	LineJoin(joinStyle LineJoinStyle)
	LinePattern(pattern []float64, phase float64)
	Text(f Font, p Point, text string)

	Stroke(path *Path)
	Fill(path *Path)
}
