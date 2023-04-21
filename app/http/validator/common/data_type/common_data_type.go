package data_type

// 这个其实已经么用了，如果不传入，在后面paginate的时候进行了处理
type Page struct {
	Page  uint `form:"page" json:"page"`
	Limit uint `form:"limit" json:"limit"`
}
