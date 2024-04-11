package custom

type ResponseCommon struct {
	Success bool        `json:"success,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Toast   string      `json:"toast,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
