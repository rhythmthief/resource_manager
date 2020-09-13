package models

//Msg type used for messages
type Msg map[string]interface{}

//InternalError error
var InternalError = Msg{"message": "internal error"}

//TemplateNotFound error
var TemplateNotFound = Msg{"message": "template not found"}

//TemplateExists error
var TemplateExists = Msg{"message": "template with this name already exists"}

//TemplateValidateFailed error
var TemplateValidateFailed = Msg{
	"description":  "cannot be blank",
	"fields":       "cannot be blank, field keys have to be unique, field type cannot be empty",
	"templatename": "cannot be blank",
}

//InvalidID error
var InvalidID = Msg{"message": "ID is a hexademical string of length 24; check template, project, session or resource ID"}

//TemplateDeleteSuccess message
var TemplateDeleteSuccess = Msg{"message": "template deleted successfully"}

//ProjectValidateFailed error
var ProjectValidateFailed = Msg{
	"name":     "cannot be blank",
	"settings": "cannot be blank",
}

//ProjectExists error
var ProjectExists = Msg{"message": "project with this name already exists"}

//ProjectNotFound error
var ProjectNotFound = Msg{"message": "project not found"}

//ProjectDeleteSuccess message
var ProjectDeleteSuccess = Msg{"message": "project deleted successfully"}

//SessionValidateFailed error
var SessionValidateFailed = Msg{
	"apikey": "cannot be blank",
}

//SigningKeyError error
var SigningKeyError = Msg{
	"message": "JWT signing key corrupted or missing",
}

//SessionNotFound error
var SessionNotFound = Msg{"message": "session not found"}

//SessionTerminated message
var SessionTerminated = Msg{"message": "session terminated"}

//SessionResNotCheckedOut error
var SessionResNotCheckedOut = Msg{"message": "resource not checked out by the session"}

//SessionResAlreadyCheckedOut error
var SessionResAlreadyCheckedOut = Msg{"message": "resource already checked out"}

//SessionResCheckedIn message
var SessionResCheckedIn = Msg{"message": "checked in"}

//SessionSubResNotFound error
var SessionSubResNotFound = Msg{"message": "requested subresource could not be found under provided resource"}

//SessionSubResNotConsumed error
var SessionSubResNotConsumed = Msg{"message": "session hasn't consumed this subresource"}

//SessionSubResDepleted error
var SessionSubResDepleted = Msg{"message": "subresource already depleted"}

//SessionSubResConsumed message
var SessionSubResConsumed = Msg{"message": "subresource successfully consumed"}

//SessionSubResReleased message
var SessionSubResReleased = Msg{"message": "subresource successfully released"}

//ResourceValidateFailed error
var ResourceValidateFailed = Msg{
	"name":        "cannot be blank",
	"description": "cannot be blank",
	"templateid":  "cannot be blank",
	"projectid":   "cannot be blank",
	"fields":      "cannot be blank",
}

//ResourceTemplateMismatch error
var ResourceTemplateMismatch = Msg{"message": "mismatched base template, cannot validate resource structure"}

//ResourceTemplateMismatchRequired error
var ResourceTemplateMismatchRequired = Msg{"message": "mismatched base template, cannot validate resource structure; make sure that required fields are not empty and subresource-type Fields have integers in Value"}

//ResourceExists error
var ResourceExists = Msg{"message": "resource with this name already exists"}

//ResourceNotFound error
var ResourceNotFound = Msg{"message": "resource not found"}

//ResourceDeleteSuccess message
var ResourceDeleteSuccess = Msg{"message": "resource deleted successfully"}

//ResourceProjectIDDuplicate error
var ResourceProjectIDDuplicate = Msg{"message": "duplicate project ids not allowed"}

//ResourceUpdateCheckedOut error
var ResourceUpdateCheckedOut = Msg{"message": "resource is checked out; cannot update a checked out resource"}

//PageNotFound error
var PageNotFound = Msg{"message": "page not found"}
