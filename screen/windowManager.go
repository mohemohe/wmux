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
	}
)

func NewWindowManager(screen tcell.Screen) *WindowManager {
	wm := &WindowManager{
		screen: screen,
		windows: make([]*Window, 0),
	}
	return wm
}

func (this *WindowManager) CreateWindow() {
	w := NewWindow(this.screen)
	this.windows = append(this.windows, w)
	w.Open(true)
}

func (this *WindowManager) CloseWindow() {
	// TODO: impl
}

func (this *WindowManager) ForceRender() {
	for _, w := range this.windows {
		w.ForceRender()
	}
}

func (this *WindowManager) OnLeftMouseDown(x int, y int) {
	this.isLeftMouseDown = true
	this.changeActiveWindow(x, y)
}

func (this *WindowManager) OnMouseMove(x int, y int) {
	if this.isLeftMouseDown {

	}
}

func (this *WindowManager) OnLeftMouseUp(x int, y int) {
	this.isLeftMouseDown = false
}

func (this *WindowManager) changeActiveWindow(x int, y int) bool {
	changed := false

	newWindows := make([]*Window, len(this.windows))
	for i, w := range this.windows {
		if w.TryClick(x, y) {
			if i == 0 {
				newWindows[i] = w
			} else {
				changed = true
				newWindows, newWindows[0] = append(newWindows[0:1], newWindows[0:i]...), w
				newWindows = append(newWindows, this.windows[i:]...)
			}
			w.Active(true)
			break
		} else {
			newWindows[i] = w
		}
	}
	this.windows = newWindows

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