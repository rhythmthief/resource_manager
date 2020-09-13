package models

import (
	"github.com/Kamva/mgm"
)

//SubResConsumed structure
type SubResConsumed struct {
	ParentID string `json:"parentid" example:"5f19a22e5b40abf84d198e53" format:"string"` //MongoDB ID of the parent resource this subresource belongs to
	Key      string `json:"key" example:"keyName" format:"string"`                       //Subresource key
	Amount   int    `json:"amount" example:"2" format:"integer"`                         //How many subresources of this type have been consumed by the session
}

//SessionRequest structure
type SessionRequest struct {
	APIKey string `json:"apikey" example:"R_l7fU2h7ROa8W62xmpTo-FUSVadckpxzga_QWXvY2tsAapPff46d9JR9Fvn7wosx6Y0wfw9dsvuMgb3GSZKNg==" format:"string"`
}

//Session structure
type Session struct {
	mgm.DefaultModel `bson:",inline"` //Default mgm-defined fields

	JobID     string           `json:"jobid" example:"6b5df1009a275789275b76e6540d3d74" format:"string"` //Session's ID within the scheduler
	Project   string           `json:"project" example:"5f19a22e5b40abf84d198e53" format:"string"`       //Project this session is associated with
	Resources []string         `json:"resources" example:"5f19a22e5b40abf84d198e53" format:"string"`     //List of resources checked out by the session
	Consumed  []SubResConsumed `json:"consumed"`
}
