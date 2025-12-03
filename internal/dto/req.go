package dto

// post
type RegisterReq struct {
	Username string `json:"username" binding:"required,alphanum,min=4,max=20"`
	Password string `json:"password" binding:"required,min=8,max=32"`
	RePass   string `json:"repass"   binding:"required,eqfield=Password"`
	Email    string `json:"email"    binding:"required,email"`
	Age      int    `json:"age"      binding:"omitempty,gte=1,lte=130"`
	Role     int    `json:"role"     binding:"required,oneof=1 2 3"`
}

// get
type PageReq struct {
	Page     int `form:"page"     binding:"required,min=1"`
	PageSize int `form:"size"     binding:"omitempty,oneof=10 20 50 100"`
}

type DeleteReq struct {
	IDs []int64 `json:"ids" binding:"required,dive,gt=0"` // dive
}

type GetUserReq struct {
	Age int `form:"age"      binding:"omitempty,gte=1,lte=130"`
}
