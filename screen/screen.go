package screen

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"os"
)

func Start() {
	encoding.Register()
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer func() {
		err := recover()
		screen.Fini()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	if e := screen.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorGray).Foreground(tcell.ColorWhite))
	screen.EnableMouse()
	screen.Clear()

	w, h := screen.Size()
	ecnt := 0

	wm := NewWindowManager(screen)
	wm.CreateWindow()
	wm.CreateWindow().Move(5, 3)
	wm.CreateWindow().Move(10, 6)

	if err != nil {
		panic(err)
	}

	for {
		screen.Show()
		ev := screen.PollEvent()
		st := tcell.StyleDefault.Background(tcell.ColorRed)
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
			}
			wm.OnKeyDown(ev)
		case *tcell.EventMouse:
			x, y := ev.Position()
			button := ev.Buttons()
			if button&tcell.WheelUp != 0 {
			}
			if button&tcell.WheelDown != 0 {
			}
			button &= tcell.ButtonMask(0xff)

			switch ev.Buttons() {
			case tcell.ButtonNone:
				wm.OnLeftMouseUp()
			case tcell.Button1:
				wm.OnLeftMouseDown(x, y)
			}
			if button != tcell.ButtonNone {
				wm.OnMouseMove(x, y)
			}
		default:
			wm.OnLeftMouseUp()
		}
	}
}
