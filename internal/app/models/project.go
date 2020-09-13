package models

import (
	a "library/internal/pkg/auth"

	"github.com/Kamva/mgm"
)

//ProjectSetting structure
type ProjectSetting struct {
	Key   string      `json:"key" example:"settingName" format:"string"`
	Value interface{} `json:"value"` //Any data type
}

//Project structure
type Project struct {
	mgm.DefaultModel `bson:",inline"` //Default mgm-defined fields

	Name      string           `json:"name" example:"project name" format:"string"`
	APIKey    string           `json:"apikey" example:"ba5e7c738d40bbeacdcad85191872171d917afa5d680e11590d281f8cb59ebe3" format:"string"`
	Resources []string         `json:"resources" example:"5f19a22e5b40abf84d198e53" format:"string"`
	Settings  []ProjectSetting `json:"settings"`
}

//ProjectRequest structure
type ProjectRequest struct {
	Name     string `json:"name" example:"project name" format:"string"`
	Settings []ProjectSetting
}

//UpdateAPIKey sets a new API key for the project
func (proj *Project) UpdateAPIKey() error {
	var err error

	proj.APIKey, err = a.GenerateKey(64)

	return err
}

//DeleteResource removes a resource id from the list
func (proj *Project) DeleteResource(resID string) error {
	var err error

	for i := 0; i < len(proj.Resources); i++ {
		if resID == proj.Resources[i] {
			//Order of elements doesn't matter, copy last element and pop
			proj.Resources[i] = proj.Resources[len(proj.Resources)-1]
			proj.Resources = proj.Resources[:len(proj.Resources)-1]
			break
		}
	}

	return err
}
