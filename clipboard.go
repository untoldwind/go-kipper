package kipper

// Atom to access (only relevant for X11)
type Atom uint8

const (
	// AtomClipboard default clipboard access
	AtomClipboard Atom = iota
	// AtomPrimary primary selection on X11 (will be ignored on other systems)
	AtomPrimary
	// AtomSecondary secondary selection on X11 (will be ignored on other systems)
	AtomSecondary
)

// Clipboard interface to the systems clipboard
type Clipboard interface {
	Get(atom Atom) (string, error)
}

// NewClipboard create a new Clipboard interface for a givven program name
// (that helps other programs understand where stuff is coming from)
func NewClipboard(programName string) (Clipboard, error) {
	return newClipboard(programName)
}
