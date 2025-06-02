package dto

type RegisterRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone" binding:"required"`        // Phone wajib diisi
	FullName string `json:"full_name" binding:"required"`    // FullName wajib diisi
	Role     string `json:"role" binding:"required"`         // Role wajib diisi
	NIP      string `json:"nip" binding:"required,alphanum"` // NIP wajib diisi dan alphanumeric
	Active   bool   `json:"active" binding:"required"`       // Active wajib diisi
}
