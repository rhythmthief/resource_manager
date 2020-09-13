package controller

import (
	"library/internal/app/business"
	m "library/internal/app/models"
	db "library/internal/pkg/dbutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

// CreateResource godoc
// @Summary Create a new resource
// @Description Create a new resource with a passed json
// @Tags resource
// @Accept json
// @Produce json
// @Param resource body models.ResourceRequest true "Add resource"
// @Success 201 {object} models.Resource
// @Failure 400 {object} models.Msg
// @Failure 409 {object} models.Msg
// @Failure 500 {object} models.Msg
// @Router /resource [post]
func (controller *Controller) CreateResource(c echo.Context) error {
	var err error
	requestData := &m.ResourceRequest{}

	//Validating the passed JSON structure
	if err = c.Bind(requestData); err != nil {
		err = c.JSON(http.StatusBadRequest, m.ResourceValidateFailed)
	} else {
		err = c.JSON(business.CreateResourceBusiness(requestData, controller.Mux))
	}

	return err
}

// ShowAllResources godoc
// @Summary Show all resources
// @Description Returns all resources stored in the database
// @Tags resource
// @Accept json
// @Produce json
// @Success 200 {object} models.Resource
// @Failure 404 {object} models.Msg
// @Router /resource [get]
func (controller *Controller) ShowAllResources(c echo.Context) error {
	return c.JSON(business.ShowAllResourcesBusiness())
}

// ShowResourcesByPrj godoc
// @Summary Show project resources
// @Description Returns all resources associated with a particular project stored in the database
// @Tags resource
// @Accept json
// @Produce json
// @Param id path string true "Project ObjectID"
// @Success 200 {object} models.Resource
// @Failure 400 {object} models.Msg
// @Failure 404 {object} models.Msg
// @Router /resource/{id} [get]
func (controller *Controller) ShowResourcesByPrj(c echo.Context) error {
	var err error
	projID := c.Param("id")

	//Verifying that the ObjectID contains 24 hexademical characters
	if !db.VerifyObjectIDString(projID) {
		err = c.JSON(http.StatusBadRequest, m.InvalidID)
	} else {
		err = c.JSON(business.ShowResourcesByPrjBusiness(projID, controller.Mux))
	}

	return err
}

// DeleteResource godoc
// @Summary Delete resource by ID
// @Description Allows the user to delete a resource using its ID
// @Tags resource
// @Accept json
// @Produce json
// @Param id path string true "Resource ObjectID"
// @Success 200 {object} models.Msg
// @Failure 400 {object} models.Msg
// @Failure 404 {object} models.Msg
// @Failure 500 {object} models.Msg
// @Router /resource/{id} [delete]
func (controller *Controller) DeleteResource(c echo.Context) error {
	var err error
	resID := c.Param("id")

	if !db.VerifyObjectIDString(resID) {
		err = c.JSON(http.StatusBadRequest, m.InvalidID)
	} else {
		err = c.JSON(business.DeleteResourceBusiness(resID, controller.Mux))
	}

	return err
}

// UpdateResource godoc
// @Summary Update an existing resource
// @Description Update a resource with a passed json. Can be used to assign an existing resource to a different project.
// @Tags resource
// @Accept json
// @Produce json
// @Param id path string true "Resource ObjectID"
// @Param resource body models.ResourceUpdateRequest true "Update resource"
// @Success 200 {object} models.Resource
// @Failure 400 {object} models.Msg
// @Failure 404 {object} models.Msg
// @Failure 409 {object} models.Msg
// @Failure 500 {object} models.Msg
// @Router /resource/{id} [put]
func (controller *Controller) UpdateResource(c echo.Context) error {
	var err error
	resID := c.Param("id")
	requestData := &m.ResourceUpdateRequest{}

	//Verifying that the ObjectID contains 24 hexadecimal characters
	if !db.VerifyObjectIDString(resID) {
		err = c.JSON(http.StatusBadRequest, m.InvalidID)
	} else {
		//Validating the passed JSON structure
		if err = c.Bind(requestData); err != nil {
			err = c.JSON(http.StatusBadRequest, m.ResourceValidateFailed)
		} else {
			err = c.JSON(business.UpdateResourceBusiness(resID, requestData, controller.Mux))
		}
	}

	return err
}
