package events

import (
	"log"
	"remy/mouse"
)

func HandleEvent(rawMessage []byte) error {
	event, err := parseEvent[mouse.MouseEvent](rawMessage)

	if err != nil {
		log.Println("Error parsing message:", err)
		return err
	}

	switch event.Type {
	case mouse.Movement:
		movementEvent, err := parseEvent[mouse.MouseMovementEvent](rawMessage)
		if err != nil {
			return err
		}
		handleMouseMovement(movementEvent)
	case mouse.Click:
		clickEvent, err := parseEvent[mouse.MouseClickEvent](rawMessage)
		if err != nil {
			return err
		}
		handleMouseClick(clickEvent)
	default:
		log.Println("Invalid event type:", event.Type)
	}

	return nil
}

func handleMouseMovement(event *mouse.MouseMovementEvent) {
	mouse.MoveByDelta(event.X, event.Y)
}

func handleMouseClick(event *mouse.MouseClickEvent) {
	mouse.ClickButton(event.Button)
}
