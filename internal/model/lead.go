package model

type CreateLeadRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=100"`
	Email   string `json:"email" validate:"required,email"`
	Phone   string `json:"phone" validate:"required,min=6,max=20"`
	Company string `json:"company,omitempty" validate:"omitempty,max=100"`
	AppType string `json:"app_type" validate:"required,oneof=mobile web saas corporate other"`
	Message string `json:"message,omitempty" validate:"omitempty,max=1000"`
}
