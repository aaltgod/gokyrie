package team

type Bubble struct {
	name, ip     string
	graph        string
	width        int
	widthMargin  int
	height       int
	heightMargin int

	Active bool
}

func NewBubble() *Bubble {
	return &Bubble{}
}
