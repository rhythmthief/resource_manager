package business

import (
	"encoding/json"
	m "library/internal/app/models"
	"net/http"

	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/bson"
)

// UIShowCollectionBusiness godoc
func UIShowCollectionBusiness(coll string, itemID string) (int, interface{}) {
	var code int
	var response interface{}

	var dataMap []map[string]interface{} //Map used to pass data into a template
	var marsh []byte

	switch coll {
	case "projects":
		//Searching for projects
		projectsFound := []m.Project{}
		_ = mgm.Coll(&m.Project{}).SimpleFind(&projectsFound, bson.M{})
		marsh, _ = json.Marshal(projectsFound)

		//verify that target item exists
		if itemID != "" {
			if err := mgm.Coll(&m.Project{}).FindByID(itemID, &m.Project{}); err != nil {
				code, response = http.StatusNotFound, m.ProjectNotFound
			}
		}

	case "resources":
		//Searching for resources
		resourcesFound := []m.Resource{}
		_ = mgm.Coll(&m.Resource{}).SimpleFind(&resourcesFound, bson.M{})
		marsh, _ = json.Marshal(resourcesFound)

		//verify that target item exists
		if itemID != "" {
			if err := mgm.Coll(&m.Resource{}).FindByID(itemID, &m.Resource{}); err != nil {
				code, response = http.StatusNotFound, m.ResourceNotFound
			}
		}

	case "templates":
		templatesFound := []m.Template{}
		_ = mgm.Coll(&m.Template{}).SimpleFind(&templatesFound, bson.M{})
		marsh, _ = json.Marshal(templatesFound)

		//verify that target item exists
		if itemID != "" {
			if err := mgm.Coll(&m.Template{}).FindByID(itemID, &m.Template{}); err != nil {
				code, response = http.StatusNotFound, m.TemplateNotFound
			}
		}

	case "sessions":
		sessionsFound := []m.Session{}
		_ = mgm.Coll(&m.Session{}).SimpleFind(&sessionsFound, bson.M{})
		marsh, _ = json.Marshal(sessionsFound)

		//verify that target item exists
		if itemID != "" {
			if err := mgm.Coll(&m.Session{}).FindByID(itemID, &m.Session{}); err != nil {
				code, response = http.StatusNotFound, m.SessionNotFound
			}
		}

	default:
		code, response = http.StatusNotFound, m.PageNotFound
		break
	}

	if code == 0 {
		code = 200
		//Have to marshal into json and unmarshal into a map
		json.Unmarshal(marsh, &dataMap)
		response = dataMap
	}

	return code, response
}
