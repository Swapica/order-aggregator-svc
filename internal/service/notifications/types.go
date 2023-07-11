package notifications

type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Data struct {
	ACta   string `json:"acta"`
	AImg   string `json:"aimg"`
	AMsg   string `json:"amsg"`
	ASub   string `json:"asub"`
	Type   string `json:"type"`
	Secret string `json:"secret,omitempty"`
	ATime  string `json:"atime,omitempty"`
}

type NotificationPayload struct {
	Notification Notification `json:"notification"`
	Data         Data         `json:"data,omitempty"`
	Recipients   string       `json:"recipients"`
}

type ApiPayload struct {
	VerificationProof string `json:"verificationProof"`
	Identity          string `json:"identity"`
	Sender            string `json:"sender"`
	Source            string `json:"source"`
	Recipient         string `json:"recipient"`
}
