package models

import (
	"github.com/Kamva/mgm"
)

//Template structure
type Template struct {
	mgm.DefaultModel `bson:",inline"` //Default mgm-defined fields

	Name        string  `json:"name" example:"template name" format:"string"`
	Description string  `json:"description" example:"template description" format:"string"`
	Fields      []Field `json:"fields"`
}

//TemplateRequest structure
type TemplateRequest struct {
	Name        string  `json:"name" example:"template name" format:"string"`
	Description string  `json:"description" example:"template description" format:"string"`
	Fields      []Field `json:"fields"`
}
