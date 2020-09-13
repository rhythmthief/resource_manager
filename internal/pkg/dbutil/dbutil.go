package dbutil

import (
	m "library/internal/app/models"
	"regexp"

	"github.com/Kamva/mgm"
)

/*VerifyObjectIDString verifies a string intended to be used as a MongoDB ObjectID
Args:	objectid as string
Rets:	error*/
func VerifyObjectIDString(id interface{}) bool {
	var isValid = true
	reg, _ := regexp.Compile("[0-9abcdef]") //Regex expression to match hexadecimal chars

	switch t := id.(type) {
	case string:
		//Verifying that the ID contains 24 hexademical characters
		if len(t) != 24 || len(reg.FindAllString(t, 25)) != 24 {
			isValid = false
		}
	//Case when an array of ids is passed
	case []string:
		for _, str := range t {
			if len(str) != 24 || len(reg.FindAllString(str, 25)) != 24 {
				isValid = false
				break
			}
		}
	}

	return isValid
}

//UpdateProjectResourceLists godoc
func UpdateProjectResourceLists(resID string, oldProjects []m.Project, newProjects []m.Project) error {

	/* It does not seem that setting arrays to only include unique elements is possible in mongo with mgm / mongo-go-driver, so we have to iterate manually to ensure uniqueness of resource entries per-project. Potentially slow, but should not matter given the nature of the operation (update is not expected to take place often) and real-world project sizes */

	var err error
	var keep bool
	var projOldRemove []m.Project //Tracks whether old projects are being removed from resource association
	//var projNewFoundFlag []bool   //Tracks whether new projects list contains any old projects
	projNewFoundFlag := make([]bool, len(newProjects))

	//Checking whether any old projects are being reused
	for i := 0; i < len(oldProjects); i++ {
		keep = false
		for j := 0; j < len(newProjects); j++ {
			if oldProjects[i].ID.Hex() == newProjects[j].ID.Hex() {
				keep = true
				projNewFoundFlag[j] = true
				break
			}
		}

		//Couldn't find an old project in the new list, making it a candidate for removal
		if !keep {
			projOldRemove = append(projOldRemove, oldProjects[i])
		}
	}

	//Removing resource entry from old projects which aren't being reused
	for _, proj := range projOldRemove {
		proj.DeleteResource(resID)

		if err = mgm.Coll(&proj).Update(&proj); err != nil {
			break
		}
	}

	//Adding resource entry to new projects which weren't among the old projects
	if err == nil {
		for i := 0; i < len(projNewFoundFlag); i++ {
			if !projNewFoundFlag[i] {
				//This project wasn't associated with the resource in the past, have to update
				newProjects[i].Resources = append(newProjects[i].Resources, resID)

				if err = mgm.Coll(&newProjects[i]).Update(&newProjects[i]); err != nil {
					break
				}
			}
		}
	}

	return err
}
