package web

type ProductCreateRequest struct {
	Name string `validate:"required,min=1,max=100" json:"name"`
}
