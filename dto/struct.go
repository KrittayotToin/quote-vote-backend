package dto

// UserStruct for registration (includes full name)
type UserStruct struct {
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginStruct for login (only email and password)
type LoginStruct struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type QuoteStruct struct {
	Text      string `json:"text"`
	Author    string `json:"author"`
	Votes     int    `json:"votes"`
	CreatedBy uint   `json:"created_by"`
}
