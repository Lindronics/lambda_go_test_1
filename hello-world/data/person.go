package data

// Name is a simple name
type Name struct {
	FirstName       string   `json:"first_name" validate:"required"`
	AdditionalNames []string `json:"additional_names"`
	LastName        string   `json:"last_name" validate:"required"`
}
