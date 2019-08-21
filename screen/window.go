package screen

import (
	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

type (
	Rect struct {
		X int
		Y int
	}
	Window struct {
		screen tcell.Screen
		title string
		open bool
		active bool
		movable bool
		resizable bool
		origin Rect
		size Rect
	}
)

var (
	activeTitleBarStyle = tcell.StyleDefault.Background(tcell.ColorTeal).Foreground(tcell.ColorWhite).Bold(true)
	inactiveTitleBarStyle = tcell.StyleDefault.Background(tcell.ColorDarkSeaGreen).Foreground(tcell.ColorWhite).Bold(false)
	bodyStyle = tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack).Bold(false)
)

func NewWindow(screen tcell.Screen) *Window {
	w := &Window{
		screen: screen,
		title: "ウィンドウ",
		open: false,
		movable: true,
		resizable: true,
		origin: Rect{0, 0},
		size: Rect{80, 24},
	}
	return w
}

func (this *Window) SetTitle(title string) {
	if this.title != title {
		this.title = title
		this.render()
	}
}

func (this *Window) Open(isOpen bool) {
	this.open = isOpen
	this.active = true
	this.render()
}

func (this *Window) Active(isActive bool) {
	this.active = isActive
	this.render()
}

func (this *Window) Move(dx int, dy int) {

}

func (this *Window) TryClick(x int, y int) bool {
	clicked := this.open && this.origin.X <= x && x < this.origin.X + this.size.X && this.origin.Y <= y && y < this.origin.Y + this.size.Y
	this.Active(clicked)
	return clicked
}

func (this *Window) IsClickTitleBar(x int, y int) bool {
	return this.open && this.origin.X <= x && x < this.origin.X + this.size.X && this.origin.Y == y
}

func (this *Window) IsClickBody(x int, y int) bool {
	return this.open && this.origin.X <= x && x < this.origin.X + this.size.X && this.origin.Y + 1 <= y && y < this.origin.Y + this.size.Y
}

func (this *Window) render() {
	var titleBarStyle tcell.Style
	if this.active {
		titleBarStyle = activeTitleBarStyle
	} else {
		titleBarStyle = inactiveTitleBarStyle
	}

	for x := this.origin.X; x <= this.origin.X + this.size.X; x++ {
		this.screen.SetContent(x, this.origin.Y, ' ', nil, titleBarStyle)
	}
	this.drawString(this.origin.X, this.origin.Y, this.title, titleBarStyle)

	for x := this.origin.X; x <= this.origin.X + this.size.X; x++ {
		for y := this.origin.Y + 1; y < this.origin.Y + this.size.Y; y++ {
			this.screen.SetContent(x, y, ' ', nil, bodyStyle)
		}
	}
}

func (this *Window) ForceRender() {
	if this.open {
		this.render()
	}
}

// REF: https://github.com/gdamore/tcell/blob/master/_demos/mouse.go#L34
func (this *Window) drawString(x int, y int, str string, style tcell.Style) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		this.screen.SetContent(x, y, c, comb, style)
		x += w
	}
}