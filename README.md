# Resource Manager (Library)

## Summary

The library manages access to accounts and devices used by various test suites within the QAA group.

## Project Tree

``` bash
.
├── build
│   └── deploy
├── conf
├── docs
│   ├── images
│   └── swagger
├── internal
│   ├── app
│   │   ├── business
│   │   ├── controller
│   │   ├── models
│   │   ├── router
│   │   ├── static
│   │   └── views
│   │       └── partials
│   └── pkg
│       ├── auth
│       └── dbutil
└── scripts

```
* **Makefile**: provides convenience targets such as _clean_ and _image_ used in deploying a docker image to a production server.
* **main.go**: main entry point for the library application.
* **build**:  directory contains files used within the makefile to produce the library, docker image.
* **conf**: directory contains project configuration files.
* **docs**: contains images and files supporting project and/or swagger, API documentation.
* **internal**: contains all code internal to (only used by) this library project.
* **app**: application code.
* **pkg**: common code used by the application.
* **scripts**:  tools that support working with the project.

## Building and Deploying

The project is deployed and run using docker.  To build a docker image for the project, simply change into the root of the library, project directory and execute the following command:

``` bash
make image
```

This cleans the project, builds the required swagger documentation and then makes the docker image for the project, with a tag that can be pushed to the project, gitlab docker container registry.

### Deploying Locally

You can deploy locally, for development (testing) without pushing to the gitlab registory.  Simply do the following from within the root of your project:

``` bash
make image
make compose-up
```

Note:  'make compose-up' will run the application as a deamon and it will stay up until you bring it down, or stop the application.  You can stop the application also using make as follows:

``` bash
make compose-dn
```
