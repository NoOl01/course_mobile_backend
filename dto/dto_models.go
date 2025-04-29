package dto

type RegisterDto struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SendOtpDto struct {
	Email string `json:"email"`
}

type ResetPasswordDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPasswordWithCodeDto struct {
	Email string `json:"email"`
	Code  int    `json:"code"`
}

type ProfileDto struct {
	FirstName      *string `json:"first_name,omitempty"`
	LastName       *string `json:"last_name,omitempty"`
	Address        *string `json:"address,omitempty"`
	BankRequisites *string `json:"bank_requisites,omitempty"`
	Email          *string `json:"email,omitempty"`
}

type CartDto struct {
	ProductId int64 `json:"product_id"`
}
