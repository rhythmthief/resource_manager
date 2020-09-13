package router

import (
	"fmt"
	"library/internal/app/business"
	"library/internal/app/controller"
	m "library/internal/app/models"
	a "library/internal/pkg/auth"

	"github.com/Kamva/mgm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prprprus/scheduler"
	swag "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/bson"

	//Swagger documentation
	_ "library/docs/swagger"
)

//New sets up middleware and registers handlers for given routes/paths
func New(conf *business.Config) *echo.Echo {
	var signingKey string
	keySetting := &m.GlobalSetting{}
	contConfig := map[string]interface{}{}

	//Attempting to locate an old JWT signing key
	if err := mgm.Coll(keySetting).First(bson.M{"key": "signingKey"}, keySetting); err == nil {
		//Reading back previously stored JWT signing key
		signingKey = fmt.Sprintf("%v", keySetting.Value)
	} else {
		//Generating a new JWT signing key
		signingKey, _ = a.GenerateKey(256)

		//Inserting the new signing key into the database
		keySetting.Key = "signingKey"
		keySetting.Value = signingKey
		mgm.Coll(keySetting).Save(keySetting)
	}

	//Setting controller configuration
	contConfig["Scheduler"], _ = scheduler.NewScheduler(5000)
	contConfig["SigningKey"] = signingKey
	contConfig["DefSessExt"] = conf.SessExt

	e := echo.New()
	c := controller.NewController(contConfig)
	e.Renderer = c.Renderer

	//Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] method=${method}, uri=${uri}, status=${status}\n"}))
	e.Use(middleware.Recover())

	//Swagger documentation
	e.GET("/swagger/*any", swag.WrapHandler)
	e.Static("/docs/images", "docs/images")

	//API v1
	v1 := e.Group("/v1")
	{
		template := v1.Group("/template")
		{
			template.POST("", c.CreateTemplate)
			template.GET("", c.ShowAllTemplates)
			template.PUT("/:id", c.UpdateTemplate)
			template.DELETE("/:id", c.DeleteTemplate)
		}

		project := v1.Group("/project")
		{
			project.POST("", c.CreateProject)
			project.GET("", c.ShowAllProjects)
			project.PUT("/:id/newkey", c.UpdateAPIKey)
			project.PUT("/:id", c.UpdateProject)
			project.DELETE("/:id", c.DeleteProject)
		}

		session := v1.Group("/session")
		{
			session.GET("", c.ShowAllSessions)
			session.POST("", c.CreateSession)
			session.DELETE("/:id", c.CloseSessionByID)

			sessionRestricted := session.Group("/authorized")
			{
				//Require authentication with the signed JWT token
				sessionRestricted.Use(middleware.JWT([]byte(signingKey)))

				sessionRestricted.PUT("", c.RenewSession)
				sessionRestricted.DELETE("", c.CloseSessionByToken)

				//Passed ids designate db ids for a resource the session wants to interact with
				sessionRestricted.PUT("/checkout/:id", c.SessionResCheckout)
				sessionRestricted.PUT("/checkin/:id", c.SessionResCheckin)

				sessionRestricted.PUT("/checkout/:id/:key", c.ConsumeSubResource)
				sessionRestricted.PUT("/checkin/:id/:key", c.ReleaseSubResource)
			}
		}

		resource := v1.Group("/resource")
		{
			resource.POST("", c.CreateResource)
			resource.GET("", c.ShowAllResources)
			resource.GET("/:id", c.ShowResourcesByPrj)
			resource.PUT("/:id", c.UpdateResource)
			resource.DELETE("/:id", c.DeleteResource)
		}
	}

	//WebUI
	e.GET("/", c.UIIndex)
	e.GET("/collections/:collname", c.UIShowCollection)
	e.GET("/collections/:collname/:id", c.UIShowCollection)

	return e
}
