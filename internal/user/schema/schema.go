package schema

type UpdateRequest struct {
	UserID uint   `json:"-"`
	Name   string `json:"name"`
}
