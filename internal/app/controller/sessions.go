package controller

import (
	"library/internal/app/business"
	m "library/internal/app/models"
	db "library/internal/pkg/dbutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// CreateSession godoc
// @Summary Start a new session
// @Description Initiates a new session using a project-specific API key and assigns a JWT token to the session
// @Tags session
// @Accept json
// @Produce json
// @Param Key body models.SessionRequest true "Start a new session"
// @Success 200 {object} models.Msg
// @Failure 400 {object} models.Msg
// @Failure 404 {object} models.Msg
// @Failure 500 {object} models.Msg
// @Router /session [post]
func (controller *Controller) CreateSession(c echo.Context) error {
	var err error
	requestData := &m.SessionRequest{}

	if err = c.Bind(requestData); err != nil {
		err = c.JSON(http.StatusBadRequest, m.SessionValidateFailed)
	} else {
		err = c.JSON(business.CreateSessionBusiness(requestData, controller.DefSessExt, controller.SigningKey, controller.Scheduler, controller.Mux))
	}

	return err
}

// ShowAllSessions godoc
// @Summary Show all sessions
// @Description Returns a list of all ongoing sessions
// @Tags session
// @Accept json
// @Produce json
// @Success 200 {object} models.Session
// @Failure 404 {object} models.Msg
// @Router /session [get]
func (controller *Controller) ShowAllSessions(c echo.Context) error {
	return c.JSON(business.ShowAllSessionsBusiness())
}

// RenewSession godoc
// @Summary Renews a session
// @Description Renews a session, resetting the expiration time
// @Tags session
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Success 201 {object} models.Msg
// @Failure 400 {object} models.Msg
// @Failure 401 {object} models.Msg
// @Failure 404 {object} models.Msg
// @Router /session/authorized [put]
func (controller *Controller) RenewSession(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)

	return c.JSON(business.RenewSessionBusiness(token, controller.DefSessExt, controller.SigningKey, controller.Scheduler, controller.Mux))
}

// CloseSessionByToken godoc
// @Summary Terminate a session with a bearer token
// @Description Terminates a running session and releases associated resources. Meant to be used by the test suite to terminate the session. Authorized with a bearer token.
// @Tags session
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} models.Msg
// @Failure 401 {object} models.Msg
// @Failure 404 {object} models.Msg
// @Router /session/authorized [delete]
func (controller *Controller) CloseSessionByToken(c echo.Context) error {
	sessID := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)

	return c.JSON(business.CloseSessionBusiness(sessID, controller.Scheduler, controller.Mux))
}

// CloseSessionByID godoc
// @Summary Terminate a session by session id
// @Description Terminates a running session and releases associated resources. Meant to be used by a human operator with foreknowledge of the session ID within the database.
// @Tags session
// @Accept json
// @Produce json
// @Param id path string true "Session ObjectID"
// @Success 200 {object} models.Msg
// @Failure 400 {object} models.Msg
// @Failure 404 {object} models.Msg
// @Router /session/{id} [delete]
func (controller *Controller) CloseSessionByID(c echo.Context) error {
	sessID := c.Param("id")

	return c.JSON(business.CloseSessionBusiness(sessID, controller.Scheduler, controller.Mux))
}

// SessionResCheckout godoc
// @Summary Check out a resource
// @Description Checks out a resource and returns its information to the caller
// @Tags session
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path string true "Resource ObjectID"
// @Success 200 {object} models.Resource
// @Failure 400 {object} models.Msg
// @Failure 404 {object} models.Msg
// @Failure 409 {object} models.Msg
// @Failure 500 {object} models.Msg
// @Router /session/authorized/checkout/{id} [put]
func (controller *Controller) SessionResCheckout(c echo.Context) error {
	var err error
	sessID := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)
	resID := c.Param("id")

	//Verifying that the ObjectID contains 24 hexademical characters
	if !db.VerifyObjectIDString(resID) {
		err = c.JSON(http.StatusBadRequest, m.InvalidID)
	} else {
		err = c.JSON(business.SessionResCheckoutBusiness(resID, sessID, controller.Mux))
	}

	return err
}

// SessionResCheckin godoc
// @Summary Check in a resource
// @Description Checks in a resource previously checked out by the test suite
// @Tags session
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path string true "Resource ObjectID"
// @Success 200 {object} models.Msg
// @Failure 400 {object} models.Msg
// @Failure 404 {object} models.Msg
// @Failure 500 {object} models.Msg
// @Router /session/authorized/checkin/{id} [put]
func (controller *Controller) SessionResCheckin(c echo.Context) error {
	var err error
	sessID := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)
	resID := c.Param("id")

	//Verifying that the ObjectID contains 24 hexademical characters
	if !db.VerifyObjectIDString(resID) {
		err = c.JSON(http.StatusBadRequest, m.InvalidID)
	} else {
		err = c.JSON(business.SessionResCheckinBusiness(sessID, resID, controller.Mux))
	}

	return err
}

// ConsumeSubResource godoc
// @Summary Consume a subresource
// @Description Consumes a subresource belonging to a checked out resource
// @Tags session
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path string true "Resource ObjectID"
// @Param key path string true "Subresource key"
// @Success 200 {object} models.Msg
// @Failure 400 {object} models.Msg
// @Failure 401 {object} models.Msg
// @Failure 404 {object} models.Msg
// @Failure 409 {object} models.Msg
// @Failure 500 {object} models.Msg
// @Router /session/authorized/checkout/{id}/{key} [put]
func (controller *Controller) ConsumeSubResource(c echo.Context) error {
	var err error
	sessID := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)
	resID := c.Param("id")
	subResKey := c.Param("key")

	if !db.VerifyObjectIDString(resID) {
		err = c.JSON(http.StatusBadRequest, m.InvalidID)
	} else {
		err = c.JSON(business.ConsumeSubResourceBusiness(sessID, resID, subResKey, controller.Mux))
	}

	return err
}

// ReleaseSubResource godoc
// @Summary Release a subresource
// @Description Releases a subresource consumed by the current session
// @Tags session
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path string true "Resource ObjectID"
// @Param key path string true "Subresource key"
// @Success 200 {object} models.Msg
// @Failure 400 {object} models.Msg
// @Failure 401 {object} models.Msg
// @Failure 404 {object} models.Msg
// @Failure 500 {object} models.Msg
// @Router /session/authorized/checkin/{id}/{key} [put]
func (controller *Controller) ReleaseSubResource(c echo.Context) error {
	var err error
	sessID := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["id"].(string)
	resID := c.Param("id")
	subResKey := c.Param("key")

	if !db.VerifyObjectIDString(resID) {
		err = c.JSON(http.StatusBadRequest, m.InvalidID)
	} else {
		err = c.JSON(business.ReleaseSubResourceBusiness(sessID, resID, subResKey, controller.Mux))
	}

	return err
}
