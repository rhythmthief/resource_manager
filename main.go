package main

import (
	"fmt"
	"library/internal/app/business"
	"library/internal/app/router"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Kamva/mgm"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
	Author:
		Dennis Windham

	Summary:
		Library main package
*/

var mongoURI string = "mongodb://mongodb:27017"

// @title Library API
// @version 1.0.0
// @description Resource management library.
// @BasePath /v1
func main() {
	var conf business.Config

	//Get application configuration
	if _, err := toml.DecodeFile("conf/library.toml", &conf); err != nil {
		//Can't read config, log to console
		fmt.Printf("%s couldn't open configuration file: %s", time.Now().Format(time.RFC3339), err)
	} else {
		//Configure logging
		if logFile, err := os.OpenFile(conf.Logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600); err != nil {
			fmt.Printf("%s couldn't create or open log file: %s", time.Now().Format(time.RFC3339), err)
		} else {
			log.SetOutput(logFile) //Set default logger output to an open file

			if conf.SessExt < 1 || 23 < conf.SessExt {
				log.Println("error: invalid sessext in config file; value 0<sessext<24 (integer) allowed")
			} else {
				//Configure Mongo access
				if conf.Runmode == "dev" {
					mongoURI = "mongodb://localhost:27017"
				}

				if err := mgm.SetDefaultConfig(nil, conf.DBName,
					options.Client().ApplyURI(mongoURI)); err != nil {
					log.Println("error: invalid mongo configuration")
				} else {
					//Create a new router instance
					r := router.New(&conf)

					//Set logging destination
					r.Logger.SetOutput(log.Writer())
					log.Println("starting the server...")

					//Start the server
					r.Logger.Fatal(r.Start(":" + strconv.Itoa(conf.Port)))
				}
			}
		}
	}
}
