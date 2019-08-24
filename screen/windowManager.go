package screen

import (
	"github.com/gdamore/tcell"
)

type (
	WindowManager struct {
		screen tcell.Screen

		// 0に近いほど手前に表示
		windows []*Window

		isLeftMouseDown bool
		previousMouseLocation Rect
	}
)

func NewWindowManager(screen tcell.Screen) *WindowManager {
	wm := &WindowManager{
		screen: screen,
		windows: make([]*Window, 0),
	}

	return wm
}

func (this *WindowManager) Dispose() {
	for _, w := range this.windows {
		w.Close()
	}
}

func (this *WindowManager) CreateWindow() *Window {
	for _, w := range this.windows {
		w.Active(false)
	}

	w := NewWindow(this.screen, RequestCallback{
		render: func() {
			this.ForceRender()
			this.ForceUpdate()
		},
		close: func(window *Window) {
			this.CloseWindow(window)
		},
	})
	this.windows = append([]*Window{w}, this.windows...)
	w.Open(true)

	return w
}

func (this *WindowManager) CloseWindow(window *Window) {
	for i, w := range this.windows {
		if w == window {
			this.windows = append(this.windows[:i], this.windows[i+1:]...)
			break
		}
	}
	this.ForceRender()
	this.ForceUpdate()
}

func (this *WindowManager) ForceRender() {
	this.screen.Clear()
	for i := len(this.windows) - 1; i >= 0; i-- {
		this.windows[i].ForceRender()
	}
	this.ForceUpdate()
}

func (this *WindowManager) ForceUpdate() {
	_ = this.screen.PostEvent(nil) // HACK: なんかいい感じに更新してくれる
}

func (this *WindowManager) OnLeftMouseDown(x int, y int) {
	if !this.isLeftMouseDown {
		this.changeActiveWindow(x, y)
		this.previousMouseLocation.X = x
		this.previousMouseLocation.Y = y
		this.isLeftMouseDown = true
	}
}

func (this *WindowManager) OnMouseMove(x int, y int) {
	if this.isLeftMouseDown {
		dx := x - this.previousMouseLocation.X
		dy := y - this.previousMouseLocation.Y

		w := this.activeWindow()
		if w != nil {
			if w.IsClickTitleBar(this.previousMouseLocation.X, this.previousMouseLocation.Y) {
				w.Move(dx, dy)
			}
			this.ForceRender()
		}

		this.previousMouseLocation.X = x
		this.previousMouseLocation.Y = y
	}
}

func (this *WindowManager) OnLeftMouseUp() {
	this.isLeftMouseDown = false
	this.previousMouseLocation.X = -1
	this.previousMouseLocation.Y = -1
}

func (this *WindowManager) OnKeyDown(key *tcell.EventKey) {
	w := this.activeWindow()
	if w != nil {
		if key.Modifiers() == tcell.ModNone {
			switch key.Key() {
			case tcell.KeyEnter:
				fallthrough
			case tcell.KeyBackspace:
				fallthrough
			case tcell.KeyBackspace2:
				fallthrough
			case tcell.KeyRune:
				b := []byte(string(key.Rune()))
				w.Input(b)
				return
			}
		}

		// TODO: #9
	}
}

func (this *WindowManager) changeActiveWindow(x int, y int) bool {
	changed := false

	newWindows := make([]*Window, len(this.windows))
	for i, w := range this.windows {
		if w.TryClick(x, y) {
			if i == 0 {
				newWindows = this.windows
				break
			}
			changed = true
			newWindows = append([]*Window{w}, newWindows[0:i]...)
			newWindows = append(newWindows, this.windows[i+1:]...)
			w.Active(true)
			break
		} else {
			newWindows[i] = w
		}
	}
	this.windows = newWindows
	this.ForceRender()
	return changed
}

func (this *WindowManager) activeWindow() *Window {
	if len(this.windows) == 0 {
		return nil
	}
	if this.windows[0].active {
		return this.windows[0]
	}
	return nil
}