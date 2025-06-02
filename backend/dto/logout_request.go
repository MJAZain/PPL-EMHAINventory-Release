package dto

type LogoutRequest struct {
	UserId uint `json:"user_id" binding:"required"`
}
