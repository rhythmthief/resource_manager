package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"library/internal/app/business"
	m "library/internal/app/models"
	db "library/internal/pkg/dbutil"
)

// CreateTemplate godoc
// @Summary Create a new template
// @Description Create a new template with a passed json. When creating a template with a subresource (streams, et cetera), designate a Field element's type as "subresource" and use an integer Value
// @Tags template
// @Accept json
// @Produce json
// @Param template body models.TemplateRequest true "Add template"
// @Success 201 {object} models.Template
// @Failure 400 {object} models.TemplateRequest
// @Failure 409 {object} models.Msg
// @Failure 500 {object} models.Msg
// @Router /template [post]
func (controller *Controller) CreateTemplate(c echo.Context) error {
	var err error
	requestData := &m.TemplateRequest{}

	//Validating the passed JSON structure
	if err = c.Bind(requestData); err != nil {
		err = c.JSON(http.StatusBadRequest, m.TemplateValidateFailed)
	} else {
		err = c.JSON(business.CreateTemplateBusiness(requestData, controller.Mux))
	}

	return err
}

// ShowAllTemplates godoc
// @Summary Show all templates
// @Description Returns all templates stored in the database
// @Tags template
// @Accept json
// @Produce json
// @Success 200 {object} models.Template
// @Failure 404 {object} models.Msg
// @Router /template [get]
func (controller *Controller) ShowAllTemplates(c echo.Context) error {
	return c.JSON(business.ShowAllTemplatesBusiness())
}

// UpdateTemplate godoc
// @Summary Update template contents
// @Description Allows the user to update any information stored in a template, including the customizable fields
// @Tags template
// @Accept json
// @Produce json
// @Param id path string true "Template ObjectID"
// @Param template body models.TemplateRequest true "Update template"
// @Success 200 {object} models.Template
// @Failure 400 {object} models.Msg
// @Failure 404 {object} models.Msg
// @Failure 409 {object} models.Msg
// @Failure 500 {object} models.Msg
// @Router /template/{id} [put]
func (controller *Controller) UpdateTemplate(c echo.Context) error {
	var err error
	id := c.Param("id")
	requestData := &m.TemplateRequest{}

	//Verifying that the ObjectID contains 24 hexademical characters
	if !db.VerifyObjectIDString(id) {
		err = c.JSON(http.StatusBadRequest, m.InvalidID)
	} else {
		//Validating the new template structure
		if err = c.Bind(requestData); err != nil {
			err = c.JSON(http.StatusBadRequest, m.TemplateValidateFailed)
		} else {
			err = c.JSON(business.UpdateTemplateBusiness(id, requestData, controller.Mux))
		}
	}

	return err
}

// DeleteTemplate godoc
// @Summary Delete template by ID
// @Description Allows the user to delete a template using its ID
// @Tags template
// @Accept json
// @Produce json
// @Param id path string true "Template ObjectID"
// @Success 200 {object} models.Msg
// @Failure 400 {object} models.Msg
// @Failure 404 {object} models.Msg
// @Failure 500 {object} models.Msg
// @Router /template/{id} [delete]
func (controller *Controller) DeleteTemplate(c echo.Context) error {
	var err error
	id := c.Param("id")

	//Verifying that the ID contains 24 hexademical characters
	if !db.VerifyObjectIDString(id) {
		err = c.JSON(http.StatusBadRequest, m.InvalidID)
	} else {
		err = c.JSON(business.DeleteTemplateBusiness(id, controller.Mux))
	}

	return err
}
