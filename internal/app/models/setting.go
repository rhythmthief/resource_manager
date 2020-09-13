package models

import "github.com/Kamva/mgm"

//GlobalSetting structure
type GlobalSetting struct {
	mgm.DefaultModel `bson:",inline"` //Default mgm-defined fields

	Key   string      `json:"key" example:"settingName" format:"string"`
	Value interface{} `json:"value"` //Any data type
}
