package screen

import (
	"github.com/gdamore/tcell"
	"math"
	"strconv"
)

type (
	WindowManager struct {
		screen tcell.Screen

		// 0に近いほど手前に表示
		windows []*Window

		tasks   []*Window

		isLeftMouseDown bool
		previousMouseLocation Rect
	}
)

const (
	buttonSize = 3
)

var (
	taskBarStyle = tcell.StyleDefault.Background(tcell.ColorDarkSeaGreen).Foreground(tcell.ColorWhite)
	activeTaskButtonStyle = tcell.StyleDefault.Background(tcell.ColorTeal).Foreground(tcell.ColorWhite).Bold(true)
)

func NewWindowManager(screen tcell.Screen) *WindowManager {
	wm := &WindowManager{
		screen: screen,
		windows: make([]*Window, 0),
		tasks: make([]*Window, 0),
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

	lastTask := w
	if len(this.tasks) != 0 {
		lastTask = this.tasks[len(this.tasks) - 1]
	}
	lastTaskName, _ := strconv.Atoi(lastTask.GetTitle())
	w.SetTitle(strconv.Itoa(lastTaskName + 1))
	w.Open(true)

	this.tasks = append(this.tasks, w)

	return w
}

func (this *WindowManager) CloseWindow(window *Window) {
	for i, w := range this.windows {
		if w == window {
			this.windows = append(this.windows[:i], this.windows[i+1:]...)
			break
		}
	}
	for i, w := range this.tasks {
		if w == window {
			this.tasks = append(this.tasks[:i], this.tasks[i+1:]...)
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
	this.renderTaskBar()
	this.ForceUpdate()
}

func (this *WindowManager) ForceUpdate() {
	_ = this.screen.PostEvent(nil) // HACK: なんかいい感じに更新してくれる
}

func (this *WindowManager) OnLeftMouseDown(x int, y int) {
	if !this.isLeftMouseDown {
		if changed := this.changeActiveTask(x, y); !changed {
			this.changeActiveWindow(x, y)
		}
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

func (this *WindowManager) changeActiveTask(x int, y int) bool {
	changed := false

	sizeX, sizeY := this.screen.Size()
	bottom := sizeY - 1

	if y != bottom || sizeX - buttonSize < x {
		return changed
	}
	logicalX := int(math.Trunc(float64(x) / buttonSize))
	for i, t := range this.tasks {
		t.Active(false)
		if i == logicalX {
			t.Active(true)

			newWindows := append([]*Window{t}, this.tasks[0:i]...)
			newWindows = append(newWindows, this.tasks[i+1:]...)
			this.windows = newWindows
			changed = true
		}
	}
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

func (this *WindowManager) renderTaskBar() {
	sizeX, sizeY := this.screen.Size()
	bottom := sizeY - 1

	active := this.activeWindow()

	x := 0
	for i, t := range this.tasks {
		style := taskBarStyle
		if active != nil && active == t {
			style = activeTaskButtonStyle
		}
		num := []rune(t.GetTitle())
		for l := len(num); l < buttonSize; l++ {
			num = append([]rune{' '}, num[:]...)

		}
		x = i * 3
		this.screen.SetContent(x, bottom, num[0], num[1:], style)
	}
	if x != 0 {
		x += buttonSize
	}
	for ; x < sizeX - buttonSize; x++ {
		this.screen.SetContent(x, bottom, ' ', nil, taskBarStyle)
	}
	c := []rune{' ', '+', ' '}
	this.screen.SetContent(x, bottom, c[0], c[1:], activeTaskButtonStyle.Background(tcell.ColorRed))
}