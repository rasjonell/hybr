package orchestration

import "fmt"

func ToEventData(eventType EventType, v any, extras ...map[string]string) *EventChannelData {
	extra := make(map[string]string)
	if len(extras) != 0 {
		extra = extras[0]
	}

	var data string
	switch v.(type) {
	case string:
		data = fmt.Sprintf("%s", v)
	case int:
		data = fmt.Sprintf("%d", v)
	default:
		fmt.Printf("Unkown type: %t, %+v\n", v, v)
		data = fmt.Sprintf("%v", v)
	}

	return &EventChannelData{
		Data:      data,
		Extras:    extra,
		EventType: eventType,
	}
}
