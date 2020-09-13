package controller

import (
	"net/http"

	"library/internal/app/business"
	m "library/internal/app/models"
	db "library/internal/pkg/dbutil"

	"github.com/labstack/echo/v4"
)

// CreateProject godoc
// @Summary Create a new project
// @Description Create a new project
// @Tags project
// @Accept json
// @Produce json
// @Param project body models.ProjectRequest true "Add new project"
// @Success 201 {object} models.Project
// @Failure 400 {object} models.ProjectRequest
// @Failure 409 {object} models.Msg
// @Failure 500 {object} models.Msg
// @Router /project [post]
func (controller *Controller) CreateProject(c echo.Context) error {
	var err error
	requestData := &m.ProjectRequest{}

	//Validating the passed JSON structure
	if err = c.Bind(requestData); err != nil {
		err = c.JSON(http.StatusBadRequest, m.ProjectValidateFailed)
	} else {
		err = c.JSON(business.CreateProjectBusiness(requestData, controller.Mux))
	}

	return err
}

// ShowAllProjects godoc
// @Summary Show all projects
// @Description Returns all projects stored in the database
// @Tags project
// @Accept json
// @Produce json
// @Success 200 {object} models.Project
// @Failure 404 {object} models.Msg
// @Router /project [get]
func (controller *Controller) ShowAllProjects(c echo.Context) error {
	return c.JSON(business.ShowAllProjectsBusiness())
}

// UpdateAPIKey godoc
// @Summary Update project API key
// @Description Allows the user to update an API key associated with a project
// @Tags project
// @Accept json
// @Produce json
// @Param id path string true "Project ObjectID"
// @Success 200 {object} models.Project
// @Failure 400 {object} models.Msg
// @Failure 404 {object} models.Msg
// @Failure 500 {object} models.Msg
// @Router /project/{id}/newkey [put]
func (controller *Controller) UpdateAPIKey(c echo.Context) error {
	var err error
	id := c.Param("id")

	//Verifying that the ID contains 24 hexademical characters
	if !db.VerifyObjectIDString(id) {
		err = c.JSON(http.StatusBadRequest, m.InvalidID)
	} else {
		err = c.JSON(business.UpdateAPIKeyBusiness(id, controller.Mux))
	}

	return err
}

// UpdateProject godoc
// @Summary Update project contents
// @Description Allows the user to update any information stored in a project
// @Tags project
// @Accept json
// @Produce json
// @Param id path string true "Project ObjectID"
// @Param project body models.ProjectRequest true "Update project"
// @Success 200 {object} models.Project
// @Failure 400 {object} models.ProjectRequest
// @Failure 404 {object} models.Msg
// @Failure 409 {object} models.Msg
// @Failure 500 {object} models.Msg
// @Router /project/{id} [put]
func (controller *Controller) UpdateProject(c echo.Context) error {
	var err error
	id := c.Param("id")
	requestData := &m.ProjectRequest{}

	//Verifying that the ObjectID contains 24 hexademical characters
	if !db.VerifyObjectIDString(id) {
		err = c.JSON(http.StatusBadRequest, m.InvalidID)
	} else {
		//Validating the new project structure
		if err = c.Bind(requestData); err != nil {
			err = c.JSON(http.StatusBadRequest, m.ProjectValidateFailed)
		} else {
			err = c.JSON(business.UpdateProjectBusiness(id, requestData, controller.Mux))
		}
	}

	return err
}

// DeleteProject godoc
// @Summary Delete project by ID
// @Description Allows the user to delete a project using its ID. If some resource associated with the deleted resource has no other project associations, the resource is also deleted.
// @Tags project
// @Accept json
// @Produce json
// @Param id path string true "Project ObjectID"
// @Success 200 {object} models.Msg
// @Failure 400 {object} models.Msg
// @Failure 404 {object} models.Msg
// @Failure 500 {object} models.Msg
// @Router /project/{id} [delete]
func (controller *Controller) DeleteProject(c echo.Context) error {
	var err error
	id := c.Param("id")

	//Verifying that the ID contains 24 hexademical characters
	if !db.VerifyObjectIDString(id) {
		err = c.JSON(http.StatusBadRequest, m.InvalidID)
	} else {
		err = c.JSON(business.DeleteProjectBusiness(id, controller.Mux))
	}

	return err
}
