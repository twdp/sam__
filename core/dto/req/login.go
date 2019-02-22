package req

type EmailLoginDto struct {
	
	AppKey string `json:"app_key"`
	
	Secret string `json:"secret"`

	Email string `json:"email"`

	Password string `json:"password"`

	// 哪个端
	Terminal int8
}
