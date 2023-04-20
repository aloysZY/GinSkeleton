package data_type

type Page struct {
	Page  float64 `form:"page" json:"page"`
	Limit float64 `form:"limit" json:"limit"`
}
