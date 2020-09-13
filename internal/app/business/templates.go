package business

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	m "library/internal/app/models"

	"github.com/Kamva/mgm"
	"github.com/Kamva/mgm/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateTemplateBusiness godoc
func CreateTemplateBusiness(requestData *m.TemplateRequest, mux map[string]*sync.Mutex) (int, interface{}) {
	var err error
	var response interface{}
	var code int

	//Transferring request into a database model
	newTemplate := &m.Template{
		Name:        requestData.Name,
		Description: requestData.Description,
		Fields:      requestData.Fields,
	}

	mux["Templates"].Lock()

out:
	for i := 0; i < len(newTemplate.Fields); i++ {
		for j := i + 1; j < len(newTemplate.Fields); j++ {
			//Verifying that key names are unique
			if newTemplate.Fields[i].Key == newTemplate.Fields[j].Key {
				code, response = http.StatusBadRequest, m.TemplateValidateFailed
				err = fmt.Errorf("")
				break out
			}

			//Verifying that each field has a type
			if newTemplate.Fields[i].Type == "" {
				code, response = http.StatusBadRequest, m.TemplateValidateFailed
				err = fmt.Errorf("")
				break out
			}
		}
	}

	if err == nil {
		//Verifying that the template name is unique
		if err = mgm.Coll(newTemplate).First(bson.M{"name": newTemplate.Name}, &m.Template{}); err == nil {
			//Not unique
			code, response = http.StatusConflict, m.TemplateExists
		} else {
			//Inserting into MongoDB
			if err = mgm.Coll(newTemplate).Save(newTemplate); err != nil {
				code, response = http.StatusInternalServerError, m.InternalError
			} else {
				//Success
				code, response = http.StatusCreated, newTemplate
			}
		}
	}

	mux["Templates"].Unlock()

	return code, response
}

// ShowAllTemplatesBusiness godoc
func ShowAllTemplatesBusiness() (int, interface{}) {
	var code int
	var response interface{}
	var templatesFound []m.Template

	//Searching for templates
	_ = mgm.Coll(&m.Template{}).SimpleFind(&templatesFound, bson.M{})

	//Verifying whether we found any templates
	if len(templatesFound) == 0 {
		code, response = http.StatusNotFound, m.TemplateNotFound
	} else {
		code, response = http.StatusOK, templatesFound
	}

	return code, response
}

// UpdateTemplateBusiness godoc
func UpdateTemplateBusiness(id string, requestData *m.TemplateRequest, mux map[string]*sync.Mutex) (int, interface{}) {
	var err error
	var response interface{}
	var code int
	template := &m.Template{}

	mux["Templates"].Lock()

	//Looking up a template under passed id
	if err = mgm.Coll(template).FindByID(id, template); err != nil {
		code, response = http.StatusNotFound, m.TemplateNotFound
	} else {
		//Checking whether the updated template name matches its old name
		if template.Name != requestData.Name {
			//Names don't match, make sure the name isn't taken by another template
			if err = mgm.Coll(template).First(bson.M{"name": requestData.Name}, &m.Template{}); err == nil {
				//Name taken by another template
				code, response = http.StatusConflict, m.TemplateExists
				err = fmt.Errorf("") //Standin non-empty error to fail next logic check
			} else {
				err = nil
			}
		}

	out:
		//Verifying that key names are unique
		for i := 0; i < len(requestData.Fields); i++ {
			for j := i + 1; j < len(requestData.Fields); j++ {
				if requestData.Fields[i].Key == requestData.Fields[j].Key {
					code, response = http.StatusBadRequest, m.TemplateValidateFailed
					err = fmt.Errorf("")
					break out
				}

				//Verifying that each field has a type
				if requestData.Fields[i].Type == "" {
					code, response = http.StatusBadRequest, m.TemplateValidateFailed
					err = fmt.Errorf("")
					break out
				}
			}
		}

		//No name conflicts detected, ready to update
		if err == nil {
			//Assigning respective data to a copy of the template
			template.Name = requestData.Name
			template.Description = requestData.Description
			template.Fields = requestData.Fields

			//Updating the template in the database
			if err = mgm.Coll(template).Update(template); err != nil {
				code, response = http.StatusInternalServerError, m.InternalError
			} else {
				code, response = http.StatusOK, template
			}
		}
	}

	mux["Templates"].Unlock()

	return code, response
}

// DeleteTemplateBusiness godoc
func DeleteTemplateBusiness(id string, mux map[string]*sync.Mutex) (int, interface{}) {
	var err error
	var response interface{}
	var code int

	//Generating Mongo ObjectID from passed id
	oid, _ := primitive.ObjectIDFromHex(id)

	mux["Templates"].Lock()

	//Attempting to find and delete a matching document
	res, err := mgm.Coll(&m.Template{}).DeleteOne(context.Background(), bson.M{field.ID: oid})
	if err != nil || res.DeletedCount == 0 {
		if err != nil {
			code, response = http.StatusInternalServerError, m.InternalError
		} else {
			code, response = http.StatusNotFound, m.TemplateNotFound
		}
	} else {
		code, response = http.StatusOK, m.TemplateDeleteSuccess
	}

	mux["Templates"].Unlock()

	return code, response
}
