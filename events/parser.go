package events

import (
	"encoding/json"
	"log"
)

func parseEvent[T any](rawEvent []byte) (event *T, err error) {
	event = new(T)

	err = json.Unmarshal(rawEvent, event)

	if err != nil {
		log.Println("Error parsing event:", err)
		return nil, err
	}

	return event, err
}
