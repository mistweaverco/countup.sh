package colors

type Colors struct {
	Black         string
	Red           string
	Green         string
	Yellow        string
	Blue          string
	Magenta       string
	Cyan          string
	White         string
	BrightBlack   string
	BrightRed     string
	BrightGreen   string
	BrightYellow  string
	BrightBlue    string
	BrightMagenta string
	BrightCyan    string
	BrightWhite   string
}

func DefaultColors() (c Colors) {
	c.Black = "#42444d"
	c.Red = "#fc2e51"
	c.Green = "#25a45c"
	c.Yellow = "#ff9369"
	c.Blue = "#3375fe"
	c.Magenta = "#9f7efe"
	c.Cyan = "#4483aa"
	c.White = "#cdd3e0"
	c.BrightBlack = "#8f9aae"
	c.BrightRed = "#ff637f"
	c.BrightGreen = "#3fc56a"
	c.BrightYellow = "#f9c858"
	c.BrightBlue = "#10b0fe"
	c.BrightMagenta = "#ff78f8"
	c.BrightCyan = "#5fb9bc"
	c.BrightWhite = "#ffffff"
	return c
}