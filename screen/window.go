package screen

import (
	"github.com/creack/pty"
	"github.com/gdamore/tcell"
	"github.com/mattn/go-libvterm"
	"github.com/mattn/go-runewidth"
	"io"
	"os"
	"os/exec"
	"time"
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
		vt *vterm.VTerm
		pty *os.File
		requestRender func()
	}
)

var (
	activeTitleBarStyle = tcell.StyleDefault.Background(tcell.ColorTeal).Foreground(tcell.ColorWhite).Bold(true)
	inactiveTitleBarStyle = tcell.StyleDefault.Background(tcell.ColorDarkSeaGreen).Foreground(tcell.ColorWhite).Bold(false)
	bodyStyle = tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack).Bold(false)
)

func NewWindow(screen tcell.Screen, requestRender func()) *Window {
	w := &Window{
		screen: screen,
		title: "うんこ",
		open: false,
		movable: true,
		resizable: true,
		origin: Rect{0, 0},
		size: Rect{80, 24},
		vt: vterm.New(23, 80),
		requestRender: requestRender,
	}
	w.vt.SetUTF8(true)
	w.vt.ObtainScreen().Reset(true)

	c := exec.Command("zsh", "-l")
	ptmx, err := pty.Start(c)
	if err != nil {
		_, _ = w.vt.Write([]byte("\033[31mTTY error\033[0m"))
	} else {
		w.pty = ptmx
	}

	_ = w.vt.ObtainScreen().Flush()

	return w
}

func (this *Window) SetTitle(title string) {
	if this.title != title {
		this.title = title
	}
}

func (this *Window) Open(isOpen bool) {
	this.open = isOpen
	this.Active(true)

	go func() {
		for {
			buff := make([]byte, 1024)
			_, err := this.pty.Read(buff)
			if err != nil{
				if err == io.EOF {
					time.Sleep(time.Millisecond)
					continue
				}
			}
			_, err = this.vt.Write(buff)
			this.vt.ObtainScreen().Flush()
			this.requestRender()
		}
	}()

}

func (this *Window) Close() {
	this.open = false
	this.active = false
	_ = this.vt.Close()

	if this.pty != nil {
		_ = this.pty.Close()
	}
}

func (this *Window) Active(isActive bool) {
	this.active = isActive
}

func (this *Window) Move(dx int, dy int) {
	this.origin.X += dx
	this.origin.Y += dy
}

func (this *Window) TryClick(x int, y int) bool {
	clicked := this.open && this.origin.X <= x && x < this.origin.X + this.size.X && this.origin.Y <= y && y < this.origin.Y + this.size.Y
	this.Active(clicked)
	return clicked
}

func (this *Window) IsClickTitleBar(x int, y int) bool {
	return this.open && this.origin.X <= x && x <= this.origin.X + this.size.X && this.origin.Y == y
}

func (this *Window) IsClickBody(x int, y int) bool {
	return this.open && this.origin.X <= x && x <= this.origin.X + this.size.X && this.origin.Y + 1 <= y && y < this.origin.Y + this.size.Y
}

func (this *Window) Input(b []byte) {
	_ , _ = this.pty.Write(b)
}

func (this *Window) render() {
	var titleBarStyle tcell.Style
	if this.active {
		titleBarStyle = activeTitleBarStyle
	} else {
		titleBarStyle = inactiveTitleBarStyle
	}

	for x := this.origin.X; x < this.origin.X + this.size.X; x++ {
		this.screen.SetContent(x, this.origin.Y, ' ', nil, titleBarStyle)
	}
	this.drawString(this.origin.X, this.origin.Y, this.title, titleBarStyle)

	for x := this.origin.X; x < this.origin.X + this.size.X; x++ {
		for y := this.origin.Y + 1; y < this.origin.Y + this.size.Y; y++ {
			this.screen.SetContent(x, y, ' ', nil, bodyStyle)
		}
	}
	vtH, vtW := this.vt.Size()
	vtS := this.vt.ObtainScreen()

	for y := 0; y < vtH; y++ {
		for x := 0; x < vtW; x++ {
			cell, err := vtS.GetCellAt(y, x)
			if err != nil {
				continue
			}
			runes := cell.Chars()
			br, bg, bb, _ := cell.Bg().RGBA()
			fr, fg, fb, _ := cell.Fg().RGBA()
			style := bodyStyle.Background(tcell.NewRGBColor(int32(br), int32(bg), int32(bb))).Foreground(tcell.NewRGBColor(int32(fr), int32(fg), int32(fb)))
			this.screen.SetContent(x + this.origin.X, y + this.origin.Y + 1, runes[0], runes[1:], style)
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