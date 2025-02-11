package model

type Example struct {
	Name   string `query:"name" json:"name" validate:"required"`
	Age    int    `query:"age" json:"age" validate:"gte=0,lte=130"`
	Email  string `query:"email" json:"email" validate:"required,email"`
	Active bool   `query:"active" json:"active"`
}
