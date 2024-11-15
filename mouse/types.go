package mouse

type MouseEventType string
type MouseButton string

const (
	Movement MouseEventType = "movement"
	Click    MouseEventType = "click"
	Error    MouseEventType = "error"
)

const (
	RightClick MouseButton = "right"
	LeftClick  MouseButton = "left"
)

type MouseEvent struct {
	Type MouseEventType `json:"type"`
}

type MouseMovementEvent struct {
	MouseEvent
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type MouseClickEvent struct {
	MouseEvent
	Button MouseButton
}
