package models

//Field structure
type Field struct {
	Type     string      `json:"type" example:"subresource" format:"string"`
	Required bool        `json:"required" example:"true" format:"boolean"`
	Key      string      `json:"key" example:"someKey" format:"string"`
	Value    interface{} `json:"value"` //Any data type
}
