package flashmessage

import (
	html "github.com/multiverse-os/starshipyard/framework/html"
)

// Notifications may be a better name
//

type FlashMessage struct {
	Type    FlashType
	Message string
}
type FlashMessages []FlashMessage

type FlashType int

const (
	ErrorType FlashType = iota
	AlertType
	WarningType
	SuccessType
	AnnouncementType
)

// TODO: Needs to migrate to new HTML Element delecaration
func (self FlashMessages) String() (output string) {
	if len(self) > 0 {
		for _, message := range self {
			output = html.Div.Class("container", "flash-messages").Text(message.String()).String()
		}
	}
	return output
}

// TODO: These types may be best definable so they can be established based on
// CSS framework.
func (self FlashType) String() string {
	switch self {
	case ErrorType:
		return "error"
	case WarningType:
		return "warning"
	case SuccessType:
		return "success"
	case AnnouncementType:
		return "announcement"
	default: // AlertType
		return "alert"
	}
}

func (self FlashMessage) String() string {
	return html.Div.Class("flash-message", self.Type.String()).Text(self.Message).String()
}
