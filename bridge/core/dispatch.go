package core

var (
	dispatchQueue chan func()
)

func init() {
	dispatchQueue = make(chan func(), 1)
}

func Dispatch(fn func()) {
	dispatchQueue <- fn
}
