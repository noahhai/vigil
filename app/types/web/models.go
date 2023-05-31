package web

type UserAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Email         string  `json:"email"`
	Username      string  `json:"username"`
	Password      string  `json:"password"`
	PhoneNumber   *string `json:"phone_number"`
	Token         string  `json:"token"`
	LoginFailures string  `json:"login_failures"`
	ResetCode     string  `json:"reset_code"`

	NotificationPhone  bool `json:"notification_phone"`
	NotificationEmail  bool `json:"notification_email"`
	NotificationMobile bool `json:"notification_mobile"`
}

type Token struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Email    string `json:"email"`
}


type AgentMessage struct {
	Name string `json:"name"`
	Duration string `json:"duration"`
	Status string `json:"status"`
	Username string `json:"username"`
	Token string `json:"token"`
}