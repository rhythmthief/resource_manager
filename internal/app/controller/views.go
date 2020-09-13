package controller

import (
	"encoding/json"
	"fmt"
	"library/internal/app/business"
	m "library/internal/app/models"
	"net/http"

	"github.com/Kamva/mgm"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

// UIIndex godoc
func (controller *Controller) UIIndex(c echo.Context) error {
	var err error
	var dataMap []map[string]interface{} //Map used to pass data into a template
	var marsh []byte

	/*
		Note: this is a stand-in implementation, for now it just defaults to the projects collection. A custom index page seems unnecessary at the moment.
	*/

	//Searching for projects
	projectsFound := []m.Project{}
	_ = mgm.Coll(&m.Project{}).SimpleFind(&projectsFound, bson.M{})
	marsh, _ = json.Marshal(projectsFound)

	//Have to marshal into json and unmarshal into a map
	json.Unmarshal(marsh, &dataMap)

	err = c.Render(http.StatusOK, "collections", map[string]interface{}{
		"coll":   "projects",
		"itemID": "",
		"docs":   dataMap,
	})

	return err
}

// UIShowCollection godoc
func (controller *Controller) UIShowCollection(c echo.Context) error {
	var err error
	coll := c.Param("collname")
	itemID := c.Param("id")
	fmt.Println(itemID)

	code, response := business.UIShowCollectionBusiness(coll, itemID)

	if code != 200 {
		err = c.JSON(code, response)
	} else {
		err = c.Render(http.StatusOK, "collections", map[string]interface{}{
			"coll":   coll,
			"itemID": itemID,
			"docs":   response,
		})
	}

	return err
}
