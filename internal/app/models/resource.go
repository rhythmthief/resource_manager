package models

import "github.com/Kamva/mgm"

//ResourceUpdateRequest structure
type ResourceUpdateRequest struct {
	Name        string   `json:"name" example:"resource name" format:"string"`
	Description string   `json:"description" example:"resource description" format:"string"`
	Projects    []string `json:"projects" example:"5f19a22e5b40abf84d198e53" format:"string"` //Name of the project this resource is associated with
	Fields      []Field  `json:"fields"`
	Active      bool     `json:"active" example:"true" format:"boolean"`
}

//ResourceRequest structure
type ResourceRequest struct {
	Name        string   `json:"name" example:"resource name" format:"string"`
	Description string   `json:"description" example:"resource description" format:"string"`
	TemplateID  string   `json:"templateid" example:"5f19a22e5b40abf84d198e53" format:"string"`
	Projects    []string `json:"projects" example:"5f19a22e5b40abf84d198e53" format:"string"` //Name of the project this resource is associated with
	Fields      []Field  `json:"fields"`
}

//Resource structure
type Resource struct {
	mgm.DefaultModel `bson:",inline"` //Default mgm-defined fields

	Name        string   `json:"name" example:"resource name" format:"string"`
	Description string   `json:"description" example:"resource description" format:"string"`
	TemplateID  string   `json:"templateid" example:"5f19a22e5b40abf84d198e53" format:"string"`
	Projects    []string `json:"projects" example:"5f19a22e5b40abf84d198e53" format:"string"` //Name of the project this resource is associated with
	Fields      []Field  `json:"fields"`
	CheckedOut  int      `json:"checkedout" example:"0" format:"boolean"` //Number of sessions which have this resource checked out
	Active      bool     `json:"active" example:"true" format:"boolean"`
}

//DeleteProject removes a project id from the list
func (res *Resource) DeleteProject(projID string) error {
	var err error

	for i := 0; i < len(res.Projects); i++ {
		if projID == res.Projects[i] {
			//Order of elements doesn't matter, copy last element and pop
			res.Projects[i] = res.Projects[len(res.Projects)-1]
			res.Projects = res.Projects[:len(res.Projects)-1]
			break
		}
	}

	return err
}
