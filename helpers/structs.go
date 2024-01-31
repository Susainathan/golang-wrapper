package helpers

import "github.com/google/uuid"

type FormData struct {
	Data string `form:"data"`
}

type ExceptionStruct struct {
	Id           uuid.UUID `json:"Id"`
	ExceptionKey string    `json:"exceptionKey"`
	Message      string    `json:"message"`
	Debug        string    `json:"debug"`
}

type OutputData struct {
	Event       string                 `json:"event"`
	EventType   string                 `json:"event_type"`
	AppId       string                 `json:"app_id"`
	UserId      string                 `json:"user_id"`
	MessageId   string                 `json:"message_id"`
	PageTitle   string                 `json:"page_title"`
	PageURL     string                 `json:"page_url"`
	BrowserLang string                 `json:"browser_language"`
	ScreenSize  string                 `json:"screen_size"`
	Attributes  map[string]interface{} `json:"attributes"`
	UserTraits  map[string]interface{} `json:"traits"`
}

type ResponseStruct struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
