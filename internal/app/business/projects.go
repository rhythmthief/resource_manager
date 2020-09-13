package business

import (
	"fmt"
	m "library/internal/app/models"
	"net/http"
	"sync"

	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateProjectBusiness godoc
func CreateProjectBusiness(requestData *m.ProjectRequest, mux map[string]*sync.Mutex) (int, interface{}) {
	var err error
	var code int
	var response interface{}

	//Transferring request into a database model
	newProject := &m.Project{
		Name:     requestData.Name,
		Settings: requestData.Settings,
	}

	mux["Projects"].Lock()

	//Generating an API key for the new project
	newProject.UpdateAPIKey()

	//Verifying that the project name is unique
	if err = mgm.Coll(newProject).First(bson.M{"name": newProject.Name}, &m.Project{}); err == nil {
		//Not unique
		code, response = http.StatusConflict, m.ProjectExists
	} else {
		//Inerting into MongoDB
		if err = mgm.Coll(newProject).Save(newProject); err != nil {
			code, response = http.StatusInternalServerError, m.InternalError
		} else {
			//Success
			code, response = http.StatusCreated, newProject
		}
	}

	mux["Projects"].Unlock()

	return code, response
}

// ShowAllProjectsBusiness godoc
func ShowAllProjectsBusiness() (int, interface{}) {
	var code int
	var response interface{}
	projectsFound := []m.Project{}

	//Searching for projects
	_ = mgm.Coll(&m.Project{}).SimpleFind(&projectsFound, bson.M{})

	//Verifying whether we found any projects
	if len(projectsFound) == 0 {
		code, response = http.StatusNotFound, m.ProjectNotFound
	} else {
		code, response = http.StatusOK, projectsFound
	}

	return code, response
}

// UpdateAPIKeyBusiness godoc
func UpdateAPIKeyBusiness(id string, mux map[string]*sync.Mutex) (int, interface{}) {
	var err error
	var code int
	var response interface{}
	project := &m.Project{}

	mux["Projects"].Lock()

	//Looking up a project under passed id
	if err = mgm.Coll(project).FindByID(id, project); err != nil {
		code, response = http.StatusNotFound, m.ProjectNotFound
	} else {
		project.UpdateAPIKey()

		//Updating the project in the database
		if err = mgm.Coll(project).Update(project); err != nil {
			code, response = http.StatusInternalServerError, m.InternalError
		} else {
			code, response = http.StatusOK, project
		}
	}

	mux["Projects"].Unlock()

	return code, response
}

// UpdateProjectBusiness godoc
func UpdateProjectBusiness(id string, requestData *m.ProjectRequest, mux map[string]*sync.Mutex) (int, interface{}) {
	var err error
	var code int
	var response interface{}
	project := &m.Project{}

	mux["Projects"].Lock()

	//Looking up a project under passed id
	if err = mgm.Coll(project).FindByID(id, project); err != nil {
		code, response = http.StatusNotFound, m.ProjectNotFound
	} else {
		//Checking whether the updated project name matches its old name
		if project.Name != requestData.Name {
			//Names don't match, make sure the name isn't taken by another project
			if err = mgm.Coll(project).First(bson.M{"name": requestData.Name}, &m.Project{}); err == nil {
				//Name taken by another project
				code, response = http.StatusConflict, m.ProjectExists
				err = fmt.Errorf("") //Standin non-empty error to fail next logic check
			} else {
				err = nil
			}
		}

		//No name conflicts detected, ready to update
		if err == nil {
			//Assigning respective data to a copy of the project
			project.Name = requestData.Name
			project.Settings = requestData.Settings

			//Updating the project in the database
			if err = mgm.Coll(project).Update(project); err != nil {
				code, response = http.StatusInternalServerError, m.InternalError
			} else {
				code, response = http.StatusOK, project
			}
		}
	}

	mux["Projects"].Unlock()

	return code, response
}

// DeleteProjectBusiness godoc
func DeleteProjectBusiness(id string, mux map[string]*sync.Mutex) (int, interface{}) {
	var err error
	var code int
	var response interface{}
	project := &m.Project{}

	mux["Projects"].Lock()

	//Locate the project
	if err = mgm.Coll(project).FindByID(id, project); err != nil {
		code, response = http.StatusNotFound, m.ProjectNotFound
	} else {
		resource := &m.Resource{}

		mux["Resources"].Lock()

		for _, resID := range project.Resources {
			//Note: not breaking if a resounce is not found -- this edge case would indicate an internal server error, but catching and handling it would be useless as opposed to going through with project deletion
			if err = mgm.Coll(resource).FindByID(resID, resource); err == nil {
				//Deleting project id from a list within an associated resource
				resource.DeleteProject(id)

				//Delete resource if there are no more project associations for it
				if len(resource.Projects) == 0 {
					if err = mgm.Coll(resource).Delete(resource); err != nil {
						code, response = http.StatusInternalServerError, m.InternalError
						break
					}
				} else {
					//Update resource if it's not deleted
					if err = mgm.Coll(resource).Update(resource); err != nil {
						code, response = http.StatusInternalServerError, m.InternalError
						break
					}
				}
			}
		}

		mux["Resources"].Unlock()

		if err == nil {
			if err = mgm.Coll(project).Delete(project); err != nil {
				code, response = http.StatusInternalServerError, m.InternalError
			} else {
				code, response = http.StatusOK, m.ProjectDeleteSuccess
			}
		}
	}

	mux["Projects"].Unlock()

	return code, response
}
