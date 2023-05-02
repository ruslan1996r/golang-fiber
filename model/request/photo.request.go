package request

type PhotoCrateRequest struct {
	CategoryId uint `json:"category_id" form:"category_id" validate:"required"`
}
