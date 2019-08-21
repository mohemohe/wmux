package screen

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/mattn/go-runewidth"
	"os"
)

func emitStr(screen tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		screen.SetContent(x, y, c, comb, style)
		x += w
	}
}

func drawWindow(screen tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, r rune) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	if runewidth.IsEastAsian() {
		if x1 % 2 != 0 {
			x1++
		}
		if x2 % 4 != 0 {
			x2++
		}
		// space := []rune{' '}
		combH := []rune{tcell.RuneHLine}
		combV := []rune{tcell.RuneVLine}
		combUR := []rune{tcell.RuneURCorner}
		combLR := []rune{tcell.RuneLRCorner}
		for col := x1; col <= x2; col += runewidth.RuneWidth(tcell.RuneHLine) {
			screen.SetContent(col, y1, tcell.RuneHLine, combH, style)
			screen.SetContent(col, y2, tcell.RuneHLine, combH, style)
		}
		for row := y1 + 1; row < y2; row++ {
			screen.SetContent(x1, row, tcell.RuneVLine, nil, style)
			screen.SetContent(x2, row, 0, combV, style)
		}
		if y1 != y2 && x1 != x2 {
			// Only add corners if we need to
			screen.SetContent(x1, y1, tcell.RuneULCorner, combH, style)
			screen.SetContent(x2, y1, tcell.RuneHLine, combUR, style)
			screen.SetContent(x1, y2, tcell.RuneLLCorner, combH, style)
			screen.SetContent(x2, y2, tcell.RuneHLine, combLR, style)
		}

		// for row := y1 + 1; row < y2; row++ {
		// 	for col := x1 + 1; col < x2; col++ {
		// 		screen.SetContent(col, row, ' ', space, style)
		// 	}
		// }
	} else {
		for col := x1; col <= x2; col+=runewidth.RuneWidth(tcell.RuneHLine) {
			screen.SetContent(col, y1, tcell.RuneHLine, nil, style)
			screen.SetContent(col, y2, tcell.RuneHLine, nil, style)
		}
		for row := y1 + 1; row < y2; row++ {
			screen.SetContent(x1, row, tcell.RuneVLine, nil, style)
			screen.SetContent(x2, row, tcell.RuneVLine, nil, style)
		}
		if y1 != y2 && x1 != x2 {
			// Only add corners if we need to
			screen.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
			screen.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
			screen.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
			screen.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
		}
	}

	for row := y1 + 1; row < y2; row++ {
		for col := x1 + 1; col < x2; col++ {
			screen.SetContent(col, row, r, nil, style)
		}
	}
}

func Start() {
	encoding.Register()
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if e := screen.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	screen.EnableMouse()
	screen.Clear()

	// mx, my := -1, -1
	ox, oy := -1, -1
	// bx, by := -1, -1
	w, h := screen.Size()
	ecnt := 0

	wm := NewWindowManager(screen)
	wm.CreateWindow()

	for {
		screen.Show()
		ev := screen.PollEvent()
		st := tcell.StyleDefault.Background(tcell.ColorRed)
		// up := tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorBlack)
		w, h = screen.Size()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
			screen.Clear()
			screen.SetContent(w-1, h-1, 'R', nil, st)
			wm.ForceRender()
		case *tcell.EventKey:
			screen.SetContent(w-2, h-2, ev.Rune(), nil, st)
			screen.SetContent(w-1, h-1, 'K', nil, st)
			if ev.Key() == tcell.KeyEscape {
				ecnt++
				if ecnt > 1 {
					screen.Fini()
					os.Exit(0)
				}
			} else if ev.Key() == tcell.KeyCtrlL {
				screen.Sync()
			} else {
				ecnt = 0
				if ev.Rune() == 'C' || ev.Rune() == 'c' {
					screen.Clear()
				}
			}
		case *tcell.EventMouse:
			x, y := ev.Position()
			button := ev.Buttons()
			for i := uint(0); i < 8; i++ {
				if int(button)&(1<<i) != 0 {
				}
			}
			if button&tcell.WheelUp != 0 {
			}
			if button&tcell.WheelDown != 0 {
			}
			button &= tcell.ButtonMask(0xff)

			if button != tcell.ButtonNone && ox < 0 {
				ox, oy = x, y
			}
			switch ev.Buttons() {
			case tcell.ButtonNone:
				if ox >= 0 && oy >= 0 {
					// bg := tcell.Color((lchar - '0') * 2)
					// drawWindow(screen, ox, oy, x, y,
					// 	up.Background(bg),
					// 	lchar)
					ox, oy = -1, -1
					// bx, by = -1, -1
				}
			case tcell.Button1:
				wm.OnLeftMouseDown(x, y)
			}
			if button != tcell.ButtonNone {
				// bx, by = x, y
			}
			// lchar = ch
			// screen.SetContent(w-1, h-1, 'M', nil, st)
			// mx, my = x, y
		default:
			// screen.SetContent(w-1, h-1, 'X', nil, st)
		}
	}
}
