package controller

import (
	"html/template"
	"io"
	"library/internal/app/business"
	"sync"

	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/prprprus/scheduler"
)

//Controller structure
type Controller struct {
	Scheduler  *scheduler.Scheduler
	SigningKey string                 //JWT token signing key
	DefSessExt int                    //0-23; default session extension time in hours
	Mux        map[string]*sync.Mutex //Mutex map
	Renderer   *TemplateMap           //Map of templates used for the public views
}

//NewController returns a controller reference
func NewController(config map[string]interface{}) *Controller {
	c := &Controller{}
	c.Scheduler = config["Scheduler"].(*scheduler.Scheduler)
	c.SigningKey = config["SigningKey"].(string)
	c.DefSessExt = config["DefSessExt"].(int)

	//Initialize mutex map
	c.Mux = map[string]*sync.Mutex{
		"Templates": {},
		"Projects":  {},
		"Resources": {},
		"Sessions":  {},
	}

	//Initialize renderer
	tempMap := map[string]*template.Template{
		"index":       template.Must(template.ParseFiles("internal/app/views/base.html", "internal/app/views/index.html")),
		"collections": template.Must(template.ParseFiles("internal/app/views/base.html", "internal/app/views/collections.html")),
	}

	c.Renderer = &TemplateMap{
		templates: tempMap,
	}

	//Recover any sessions which were running before termination
	business.RecoverSessionsBusiness(c.Scheduler, c.DefSessExt, c.Mux)

	return c
}

//TemplateMap structure, keeps track of view templates used for the WebUI
type TemplateMap struct {
	templates map[string]*template.Template
}

//Render export
func (t *TemplateMap) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	var err error

	if template, ok := t.templates[name]; !ok {
		err = fmt.Errorf("error: %s template not found", name)
	} else {
		err = template.ExecuteTemplate(w, "base.html", data)
	}

	return err
}
