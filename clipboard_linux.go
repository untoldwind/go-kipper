package kipper

// #cgo pkg-config: x11
//
// #include "clipboard_linux.go.c"
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

type linuxClipboard struct {
	programName string
	display     *C.Display
	window      C.Window
}

func newClipboard(programName string) (Clipboard, error) {
	display := C.XOpenDisplay(nil)
	if display == nil {
		return nil, errors.New("Failed to open default X11 display")
	}
	root := C.XDefaultRootWindow(display)
	screen := C.XDefaultScreenOfDisplay(display)
	black := screen.black_pixel
	window := C.XCreateSimpleWindow(display, root, 0, 0, 1, 1, 0, black, black)

	C.XSelectInput(display, window, C.PropertyChangeMask)

	clipboard := &linuxClipboard{
		programName: programName,
		display:     display,
		window:      window,
	}

	runtime.SetFinalizer(clipboard, (*linuxClipboard).close)

	return clipboard, nil
}

func (c *linuxClipboard) close() {
	C.XDestroyWindow(c.display, c.window)
}

func (c *linuxClipboard) convertAtom(atom Atom) (C.Atom, error) {
	switch atom {
	case AtomClipboard:
		return C._clipboard_atom(c.display), nil
	case AtomPrimary:
		return C._primary_atom(), nil
	case AtomSecondary:
		return C._seconary_atom(), nil
	}
	return 0, errors.New("Invalid atom")
}

func (c *linuxClipboard) Get(atom Atom) (string, error) {
	xAtom, err := c.convertAtom(atom)
	if err != nil {
		return "", err
	}
	result := C._get_selection_text(c.display, c.window, xAtom)
	if result == nil {
		return "", nil
	}
	defer C.free(unsafe.Pointer(result))

	return C.GoString(result), nil
}
