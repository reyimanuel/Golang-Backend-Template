package user

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest is received as multipart/form-data.
type RegisterRequest struct {
	Username string `form:"name"          binding:"required"`
	Email    string `form:"email"         binding:"required,email"`
	Password string `form:"password"      binding:"required,min=6"`
}

type Response struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type PendingStudentResponse struct {
	UserID     uint   `json:"user_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	NIM        string `json:"nim"`
	Kredensial string `json:"kredensial"`
}

// RegisterWithKRSRequest is received as multipart/form-data.
// The KRS image is handled separately via ctx.FormFile("krs").
type RegisterWithKRSRequest struct {
	Email    string `form:"email"    binding:"required,email"`
	Password string `form:"password" binding:"required,min=6"`
}

// KRSPreviewResponse is returned after a successful KRS-based registration
// so the student can confirm which data was extracted from their document.
type KRSPreviewResponse struct {
	UserID       uint   `json:"user_id"`
	Name         string `json:"name"`
	NIM          string `json:"nim"`
	ProgramStudi string `json:"program_studi"`
	Angkatan     int    `json:"angkatan"`
}
