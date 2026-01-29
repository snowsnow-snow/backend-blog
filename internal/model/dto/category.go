package dto

type CategoryReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
