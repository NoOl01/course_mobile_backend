package common

type NewUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Email struct {
	Email string `json:"email"`
}

type Otp struct {
	Email string `json:"email"`
	Otp   int    `json:"otp"`
}

type Password struct {
	Password string `json:"password"`
}
