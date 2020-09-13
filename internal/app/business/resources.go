package business

import (
	"fmt"
	"net/http"
	"sync"

	m "library/internal/app/models"
	db "library/internal/pkg/dbutil"

	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateResourceBusiness godoc
func CreateResourceBusiness(requestData *m.ResourceRequest, mux map[string]*sync.Mutex) (int, interface{}) {
	var err error
	var code int
	var response interface{}
	var projects []m.Project

	mux["Resources"].Lock()

	//Verifying that the resource name is unique
	if err = mgm.Coll(&m.Resource{}).First(bson.M{"name": requestData.Name}, &m.Resource{}); err == nil {
		//Not unique
		code, response = http.StatusConflict, m.ResourceExists
	} else {
		//Verifying that the ObjectIDs for project and template are valid, and that there is in fact a project in the array
		if len(requestData.Projects) == 0 || !db.VerifyObjectIDString(requestData.TemplateID) || !db.VerifyObjectIDString(requestData.Projects) {
			code, response = http.StatusBadRequest, m.InvalidID
		} else {
			template := &m.Template{}

			mux["Templates"].Lock()

			//Locating the template associated with the resource
			if err = mgm.Coll(template).FindByID(requestData.TemplateID, template); err != nil {
				code, response = http.StatusNotFound, m.TemplateNotFound
			} else {
				mux["Projects"].Lock()

				//Making sure all project ids are unique
				for i := 0; i < len(requestData.Projects); i++ {
					for j := i + 1; j < len(requestData.Projects); j++ {
						if requestData.Projects[i] == requestData.Projects[j] {
							code, response = http.StatusBadRequest, m.ResourceProjectIDDuplicate
							err = fmt.Errorf("")
							break
						}
					}
				}

				if err == nil {
					//Ensure that all projects associated with the resource are real
					for i, projID := range requestData.Projects {
						//Caching mongo models for later use
						projects = append(projects, m.Project{})
						if err = mgm.Coll(&m.Project{}).FindByID(projID, &projects[i]); err != nil {
							code, response = http.StatusNotFound, m.ProjectNotFound
							break
						}
					}
				}

				//No issues with projects
				if err == nil {
					//2nd round of validation -- against the base template
					if len(template.Fields) != len(requestData.Fields) {
						//Catching a general template mismatch
						code, response = http.StatusBadRequest, m.ResourceTemplateMismatch
					} else {
					out:
						for i := 0; i < len(requestData.Fields); i++ {
							//Ensuring that fields in template and resource have the same key names, required flags, type and order
							if (template.Fields[i].Key != requestData.Fields[i].Key) || (template.Fields[i].Required != requestData.Fields[i].Required) || (template.Fields[i].Type != requestData.Fields[i].Type) {
								err = fmt.Errorf("") //Filler error for logics
								break
							} else {

								//Ensuring that the subresource value is an integer
								if requestData.Fields[i].Type == "subresource" {
									switch requestData.Fields[i].Value.(type) {
									//Ensuring that the subresource value is an integer
									case float64:
										//Integer cast
										requestData.Fields[i].Value = int(requestData.Fields[i].Value.(float64))
										break

									default:
										err = fmt.Errorf("")
										break out
									}
								}

								//Ensuring that required fields are not empty
								if template.Fields[i].Required {
									switch t := requestData.Fields[i].Value.(type) {
									case string:
										if t == "" {
											err = fmt.Errorf("")
											break out
										}
									//Case when value is []
									case []interface{}:
										if len(t) == 0 {
											err = fmt.Errorf("")
											break out
										}
									//Case when value is {}
									case map[string]interface{}:
										if len(t) == 0 {
											err = fmt.Errorf("")
											break out
										}
									default:
										if t == nil {
											err = fmt.Errorf("")
											break out
										}
									}
								}
							}
						}

						//Catching a template mismatch where required fields are empty
						if err != nil {
							code, response = http.StatusBadRequest, m.ResourceTemplateMismatchRequired
						} else {
							//Transferring request into a database model
							newResource := &m.Resource{
								Name:        requestData.Name,
								Description: requestData.Description,
								TemplateID:  requestData.TemplateID,
								Projects:    requestData.Projects,
								CheckedOut:  0,
								Active:      true,
								Fields:      requestData.Fields,
							}

							//Insert the new resource into the db
							if err = mgm.Coll(newResource).Save(newResource); err != nil {
								code, response = http.StatusInternalServerError, m.InternalError
							} else {
								//Updating resource lists within projects associated with this resource
								if err == nil {
									//Going over cached project models and inserting a new resource id
									for _, proj := range projects {
										proj.Resources = append(proj.Resources, newResource.ID.Hex())
										if err = mgm.Coll(&proj).Update(&proj); err != nil {
											code, response = http.StatusNotFound, m.InternalError
											break
										}
									}

									if err == nil {
										//Success
										code, response = http.StatusCreated, newResource
									}
								}
							}
						}
					}
				}
				mux["Projects"].Unlock()
			}
			mux["Templates"].Unlock()
		}
	}
	mux["Resources"].Unlock()

	return code, response
}

// ShowAllResourcesBusiness godoc
func ShowAllResourcesBusiness() (int, interface{}) {
	var code int
	var response interface{}

	var resourcesFound []m.Resource

	//Searching for resources
	_ = mgm.Coll(&m.Resource{}).SimpleFind(&resourcesFound, bson.M{})

	//Verifying whether we found any resources
	if len(resourcesFound) == 0 {
		code, response = http.StatusNotFound, m.ResourceNotFound
	} else {
		code, response = http.StatusOK, resourcesFound
	}

	return code, response
}

// ShowResourcesByPrjBusiness godoc
func ShowResourcesByPrjBusiness(projID string, mux map[string]*sync.Mutex) (int, interface{}) {
	var err error
	var code int
	var response interface{}
	var resourcesFound []m.Resource
	project := &m.Project{}

	mux["Projects"].Lock()

	//Searching for the project
	if err = mgm.Coll(project).FindByID(projID, project); err != nil {
		code, response = http.StatusNotFound, m.ProjectNotFound
	} else {

		if len(project.Resources) == 0 {
			code, response = http.StatusNotFound, m.ResourceNotFound
		} else {
			mux["Resources"].Lock()

			//mgm cannot do bulk search by id natively, have to send individual queries
			for i, resID := range project.Resources {
				resourcesFound = append(resourcesFound, m.Resource{})
				mgm.Coll(&m.Resource{}).FindByID(resID, &resourcesFound[i])
			}

			if len(resourcesFound) == 0 {
				code, response = http.StatusInternalServerError, m.InternalError
			} else {
				code, response = http.StatusOK, resourcesFound
			}

			mux["Resources"].Unlock()
		}
	}
	mux["Projects"].Unlock()

	return code, response
}

// DeleteResourceBusiness godoc
func DeleteResourceBusiness(resID string, mux map[string]*sync.Mutex) (int, interface{}) {
	var err error
	var code int
	var response interface{}

	resource := &m.Resource{}

	mux["Resources"].Lock()

	//Attempting to find a macthing resource
	if err = mgm.Coll(resource).FindByID(resID, resource); err != nil {
		code, response = http.StatusNotFound, m.ResourceNotFound
	} else {
		project := &m.Project{}

		mux["Projects"].Lock()

		//Going over every associated project and deleting references to this resource
		for _, projID := range resource.Projects {
			if err = mgm.Coll(project).FindByID(projID, project); err != nil {
				break
			} else {
				project.DeleteResource(projID)

				if err = mgm.Coll(project).Update(project); err != nil {
					break
				}
			}
		}

		mux["Projects"].Unlock()

		if err != nil {
			code, response = http.StatusInternalServerError, m.InternalError
		} else {
			if err = mgm.Coll(resource).Delete(resource); err != nil {
				code, response = http.StatusInternalServerError, m.InternalError
			} else {
				code, response = http.StatusOK, m.ResourceDeleteSuccess
			}
		}
	}
	mux["Resources"].Unlock()

	return code, response
}

// UpdateResourceBusiness godoc
func UpdateResourceBusiness(resID string, requestData *m.ResourceUpdateRequest, mux map[string]*sync.Mutex) (int, interface{}) {
	var err error
	var code int
	var response interface{}
	var newProjects []m.Project
	var oldProjects []m.Project
	resource := &m.Resource{}

	mux["Resources"].Lock()

	//Looking up db entry for the existing resource
	if err = mgm.Coll(resource).FindByID(resID, resource); err != nil {
		code, response = http.StatusNotFound, m.ResourceNotFound
	} else {
		//Making sure the resource isn't checked out
		if resource.CheckedOut > 0 {
			code, response = http.StatusConflict, m.ResourceUpdateCheckedOut
		} else {
			//Checking whether old name mathes new name
			if resource.Name != requestData.Name {
				//Making sure there's no conflict with another resource's name
				if err = mgm.Coll(resource).First(bson.M{"name": requestData.Name}, &m.Resource{}); err == nil {
					//There is a conflict
					code, response = http.StatusConflict, m.ResourceExists
					err = fmt.Errorf("") //Standin non-empty error to fail next logic check
				} else {
					err = nil
				}
			}

			mux["Projects"].Lock()

			//Making sure all project ids are valid and unique
			if !db.VerifyObjectIDString(requestData.Projects) {
				code, response = http.StatusBadRequest, m.InvalidID
				err = fmt.Errorf("")
			} else {
				for i := 0; i < len(requestData.Projects); i++ {
					for j := i + 1; j < len(requestData.Projects); j++ {
						if requestData.Projects[i] == requestData.Projects[j] {
							code, response = http.StatusBadRequest, m.ResourceProjectIDDuplicate
							err = fmt.Errorf("")
							break
						}
					}
				}
			}

			if err == nil {
				//Ensure that all projects now associated with the resource are real
				for i, projID := range requestData.Projects {
					//Caching mongo models for later use
					newProjects = append(newProjects, m.Project{})
					if err = mgm.Coll(&m.Project{}).FindByID(projID, &newProjects[i]); err != nil {
						code, response = http.StatusNotFound, m.ProjectNotFound
						break
					}
				}
			}

			//No issues with the project id, validate structure
			if err == nil {
				//Validate against the old resource structure, which has already been validated on creation
				if len(resource.Fields) != len(requestData.Fields) {
					//Catching a general template mismatch
					code, response = http.StatusBadRequest, m.ResourceTemplateMismatch
				} else {
				out:
					//Ensuring that fields in old and new resource have the same key names and order
					for i := 0; i < len(resource.Fields); i++ {
						if (resource.Fields[i].Key != requestData.Fields[i].Key) || (resource.Fields[i].Required != requestData.Fields[i].Required) || (resource.Fields[i].Type != requestData.Fields[i].Type) {
							err = fmt.Errorf("") //Filler error for logics
							break
						} else {

							//Ensuring that the subresource value is an integer
							if requestData.Fields[i].Type == "subresource" {
								switch requestData.Fields[i].Value.(type) {
								case float64:
									//Integer cast
									requestData.Fields[i].Value = int(requestData.Fields[i].Value.(float64))
									break

								default:
									code, response = http.StatusBadRequest, m.TemplateValidateFailed
									err = fmt.Errorf("")
									break out
								}
							}

							//Ensuring that required fields are not empty
							if requestData.Fields[i].Required {
								switch t := requestData.Fields[i].Value.(type) {
								case string:
									if t == "" {
										err = fmt.Errorf("")
										break out
									}
								//Case when value is []
								case []interface{}:
									if len(t) == 0 {
										err = fmt.Errorf("")
										break out
									}
								//Case when value is {}
								case map[string]interface{}:
									if len(t) == 0 {
										err = fmt.Errorf("")
										break out
									}
								default:
									if t == nil {
										err = fmt.Errorf("")
										break out
									}
								}
							}
						}
					}

					//Catching a template mismatch where required fields are empty
					if err != nil {
						code, response = http.StatusBadRequest, m.ResourceTemplateMismatchRequired
					} else {
						//Finding all projects associated with this resource in the past
						for i, projID := range resource.Projects {
							//Caching mongo models for later use
							oldProjects = append(oldProjects, m.Project{})
							if err = mgm.Coll(&m.Project{}).FindByID(projID, &oldProjects[i]); err != nil {
								code, response = http.StatusNotFound, m.ProjectNotFound
								break
							}
						}

						//Transferring request into a database model
						resource.Name = requestData.Name
						resource.Description = requestData.Description
						resource.Projects = requestData.Projects
						resource.Fields = requestData.Fields
						resource.Active = requestData.Active

						//Update resource in the db
						if err = mgm.Coll(resource).Update(resource); err != nil {
							code, response = http.StatusInternalServerError, m.InternalError
						} else {
							//Updating resource lists within projects
							if err = db.UpdateProjectResourceLists(resource.ID.Hex(), oldProjects, newProjects); err != nil {
								code, response = http.StatusInternalServerError, m.InternalError
							} else {
								code, response = http.StatusOK, resource
							}
						}
					}
				}
			}
			mux["Projects"].Unlock()
		}
	}
	mux["Resources"].Unlock()

	return code, response
}
