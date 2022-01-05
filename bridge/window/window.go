package window

import "sync"

type Module struct {
	handles []string

	mu sync.Mutex
}

func (m *Module) All() (ret []string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ret = make([]string, len(m.handles))
	copy(ret, m.handles)
	return ret
}

func (m *Module) Create(options Options) (string, error) {
	// call C.window_create
	// add handle to m.handles
	// return handle or error
	return "", nil
}

type Options struct {
	AlwaysOnTop bool
	Frameless   bool
	Fullscreen  bool
	Size        Size
	MinSize     Size
	MaxSize     Size
	Maximized   bool
	Position    Position
	Resizable   bool
	Title       string
	Transparent bool
	Visible     bool
	Center      bool
	Icon        string // bytestream callback
	URL         string
	HTML        string
	Script      string
}

type Position struct {
	X float64
	Y float64
}

type Size struct {
	Width  float64
	Height float64
}
