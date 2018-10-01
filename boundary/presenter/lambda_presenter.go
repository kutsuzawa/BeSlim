package presenter

// ConsoleViewer defines view layer methods(DIP)
type ConsoleViewer interface {
	View(models []ConsoleModel)
}
