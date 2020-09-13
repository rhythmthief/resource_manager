package business

import (
	"fmt"
	m "library/internal/app/models"
	db "library/internal/pkg/dbutil"
	"net/http"
	"sync"
	"time"

	"github.com/Kamva/mgm"
	"github.com/dgrijalva/jwt-go"
	"github.com/prprprus/scheduler"
	"go.mongodb.org/mongo-driver/bson"
)

// ShowAllSessionsBusiness godoc
func ShowAllSessionsBusiness() (int, interface{}) {
	var code int
	var response interface{}
	sessionsFound := []m.Session{}

	//Searching for sessions
	_ = mgm.Coll(&m.Session{}).SimpleFind(&sessionsFound, bson.M{})

	if len(sessionsFound) == 0 {
		code, response = http.StatusNotFound, m.SessionNotFound
	} else {
		code, response = http.StatusOK, sessionsFound
	}

	return code, response
}

// SessionResCheckoutBusiness godoc
func SessionResCheckoutBusiness(resID string, sessID string, mux map[string]*sync.Mutex) (int, interface{}) {
	var err error
	var code int
	var response interface{}

	session := &m.Session{}

	mux["Sessions"].Lock()

	//Look up the session in the db
	if err = mgm.Coll(session).FindByID(sessID, session); err != nil {
		code, response = http.StatusNotFound, m.SessionNotFound
	} else {
		alreadyCheckedOut := false

		for _, res := range session.Resources {
			if res == resID {
				alreadyCheckedOut = true
				break
			}
		}

		if !alreadyCheckedOut {
			resource := &m.Resource{}

			mux["Resources"].Lock()

			//Look up the resource in the db
			if err = mgm.Coll(resource).FindByID(resID, resource); err != nil {
				code, response = http.StatusNotFound, m.ResourceNotFound
			} else {
				//Increment checkout counter
				resource.CheckedOut++

				//Updating resource db entry to checked out
				if err = mgm.Coll(resource).Update(resource); err != nil {
					code, response = http.StatusInternalServerError, m.InternalError
				} else {
					session.Resources = append(session.Resources, resID)

					//Updating session db entry to include new resource
					if err = mgm.Coll(session).Update(session); err != nil {
						code, response = http.StatusInternalServerError, m.InternalError
					} else {
						code, response = http.StatusOK, resource
					}
				}
			}
			mux["Resources"].Unlock()
		} else {
			code, response = http.StatusConflict, m.SessionResAlreadyCheckedOut
		}

	}
	mux["Sessions"].Unlock()

	return code, response
}

// SessionResCheckinBusiness godoc
func SessionResCheckinBusiness(sessID string, resID string, mux map[string]*sync.Mutex) (int, interface{}) {
	var err error
	var code int
	var response interface{}
	var i int

	session := &m.Session{}

	mux["Sessions"].Lock()

	//Look up the session in the db
	if err = mgm.Coll(session).FindByID(sessID, session); err != nil {
		code, response = http.StatusNotFound, m.SessionNotFound
	} else {
		found := false

		//Ensure that the resource is checked out by this session
		for i = 0; i < len(session.Resources); i++ {
			if session.Resources[i] == resID {
				found = true
				break
			}
		}

		if !found {
			code, response = http.StatusBadRequest, m.SessionResNotCheckedOut
		} else {
			resource := &m.Resource{}

			mux["Resources"].Lock()

			//Look up the resource in the db
			if err = mgm.Coll(resource).FindByID(resID, resource); err != nil {
				code, response = http.StatusNotFound, m.ResourceNotFound
			} else {
				resource.CheckedOut--

				//Release consumed subresources, if any
				for j := 0; j < len(session.Consumed); j++ {
					//See if a subresource of this resource has been checked out by the session
					if session.Consumed[j].ParentID == resID {
						for k := 0; k < len(resource.Fields); k++ {
							//Locate the subresource entry
							if session.Consumed[j].Key == resource.Fields[k].Key {
								//Since Value is an interface, have to parse it as integer before reassigning a value
								counter := resource.Fields[k].Value.(int32) + int32(session.Consumed[j].Amount)
								resource.Fields[k].Value = counter

								//Pop subresource record from session db entry
								session.Consumed[j] = session.Consumed[len(session.Consumed)-1]
								session.Consumed = session.Consumed[:len(session.Consumed)-1]

								break
							}
						}
					}
				}

				//Updating resource db entry to checked in
				if err = mgm.Coll(resource).Update(resource); err != nil {
					code, response = http.StatusInternalServerError, m.InternalError
				} else {

					//Updating resource list on the session now that the resource is officially checked in
					session.Resources[i] = session.Resources[len(session.Resources)-1]
					session.Resources = session.Resources[:len(session.Resources)-1]

					//Updating session db entry
					if err = mgm.Coll(session).Update(session); err != nil {
						code, response = http.StatusInternalServerError, m.InternalError
					} else {
						code, response = http.StatusOK, m.SessionResCheckedIn
					}
				}
			}
			mux["Resources"].Unlock()
		}
	}
	mux["Sessions"].Unlock()

	return code, response
}

// ConsumeSubResourceBusiness godoc
func ConsumeSubResourceBusiness(sessID string, resID string, subResKey string, mux map[string]*sync.Mutex) (int, interface{}) {
	var err error
	var i int
	var code int
	var response interface{}
	session := &m.Session{}

	mux["Sessions"].Lock()

	if err = mgm.Coll(session).FindByID(sessID, session); err != nil {
		code, response = http.StatusNotFound, m.SessionNotFound
	} else {
		found := false

		//Make sure the resource is checked out by the session
		for _, res := range session.Resources {
			if res == resID {
				found = true
				break
			}
		}

		if !found {
			code, response = http.StatusUnauthorized, m.SessionResNotCheckedOut
		} else {
			resource := &m.Resource{}

			mux["Resources"].Lock()

			//Find the resource db entry
			if err = mgm.Coll(resource).FindByID(resID, resource); err != nil {
				code, response = http.StatusNotFound, m.ResourceNotFound
			} else {
				found = false

				//Search Field objects for the given subresource
				for i = 0; i < len(resource.Fields); i++ {
					if resource.Fields[i].Type == "subresource" && resource.Fields[i].Key == subResKey {
						found = true
						break
					}
				}

				if !found {
					code, response = http.StatusNotFound, m.SessionSubResNotFound
				} else {
					if resource.Fields[i].Value.(int32) < 1 {
						code, response = http.StatusConflict, m.SessionSubResDepleted
					} else {
						//Found the subresource and it can be consumed

						//Since Value is an interface, have to parse it as integer before reassigning a value
						counter := resource.Fields[i].Value.(int32) - 1
						resource.Fields[i].Value = counter

						if mgm.Coll(resource).Update(resource); err != nil {
							code, response = http.StatusInternalServerError, m.InternalError
						} else {
							found = false

							for i = 0; i < len(session.Consumed); i++ {
								if session.Consumed[i].Key == subResKey {
									found = true
									session.Consumed[i].Amount++
									break
								}
							}

							if !found {
								//This is the first time the session consumes this subresource
								consumedRes := &m.SubResConsumed{
									ParentID: resID,
									Key:      subResKey,
									Amount:   1,
								}

								//Append a new consumed entry
								session.Consumed = append(session.Consumed, *consumedRes)
							}

							//Update db entry for the session
							if mgm.Coll(session).Update(session); err != nil {
								code, response = http.StatusInternalServerError, m.InternalError
							} else {
								code, response = http.StatusOK, m.SessionSubResConsumed
							}
						}
					}
				}
			}
			mux["Resources"].Unlock()
		}
	}
	mux["Sessions"].Unlock()

	return code, response
}

// ReleaseSubResourceBusiness godoc
func ReleaseSubResourceBusiness(sessID string, resID string, subResKey string, mux map[string]*sync.Mutex) (int, interface{}) {
	var err error
	var code int
	var response interface{}
	var i int
	session := &m.Session{}

	mux["Sessions"].Lock()

	if err = mgm.Coll(session).FindByID(sessID, session); err != nil {
		code, response = http.StatusNotFound, m.SessionNotFound
	} else {
		found := false

		//Make sure the resource is checked out by the session
		for i = 0; i < len(session.Consumed); i++ {
			if session.Consumed[i].Key == subResKey && session.Consumed[i].ParentID == resID {
				found = true
				break
			}
		}

		if !found {
			code, response = http.StatusUnauthorized, m.SessionSubResNotConsumed
		} else {
			resource := &m.Resource{}

			mux["Resources"].Lock()

			//Make sure the actual resource still exists
			if err = mgm.Coll(resource).FindByID(session.Consumed[i].ParentID, resource); err != nil {
				code, response = http.StatusNotFound, m.ResourceNotFound
			} else {
				session.Consumed[i].Amount--
				found = false

				//Pop subresource entry from consumed list if the amount consumed is 0
				if session.Consumed[i].Amount < 1 {
					session.Consumed[i] = session.Consumed[len(session.Consumed)-1]
					session.Consumed = session.Consumed[:len(session.Consumed)-1]
				}

				//Find subresource with this key and update
				for i = 0; i < len(resource.Fields); i++ {
					if resource.Fields[i].Key == subResKey {
						//Since Value is an interface, have to parse it as integer before reassigning a value
						counter := resource.Fields[i].Value.(int32) + 1
						resource.Fields[i].Value = counter
						found = true
						break
					}
				}

				if !found {
					code, response = http.StatusInternalServerError, m.InternalError
				} else {
					//Update session db entry
					if err = mgm.Coll(session).Update(session); err != nil {
						code, response = http.StatusInternalServerError, m.InternalError
					} else {
						//Update resource entry
						if err = mgm.Coll(resource).Update(resource); err != nil {
							code, response = http.StatusInternalServerError, m.InternalError
						} else {
							code, response = http.StatusOK, m.SessionSubResReleased
						}
					}
				}
			}
			mux["Resources"].Unlock()
		}
	}
	mux["Sessions"].Unlock()

	return code, response
}

// TerminateSessionBusiness godoc
// Caller is responsible for cancelling the termination job scheduled for the session using the returned jobID
// Args:	session db id
// Rets:	scheduled session termination jobID, error
func TerminateSessionBusiness(sessionID string, mux map[string]*sync.Mutex) (string, error) {
	var err error
	var jobID string
	session := &m.Session{}

	mux["Sessions"].Lock()

	//Retrieving session info
	if err = mgm.Coll(session).FindByID(sessionID, session); err == nil {
		//Release associated resources
		for _, resID := range session.Resources {
			resource := &m.Resource{}

			mux["Resources"].Lock()

			//Checking in resources
			if err = mgm.Coll(resource).FindByID(resID, resource); err == nil {
				resource.CheckedOut--

				//Release consumed subresources, if any
				for i := 0; i < len(session.Consumed); i++ {
					//See if a subresource of this resource has been checked out by the session
					if session.Consumed[i].ParentID == resID {
						for j := 0; j < len(resource.Fields); j++ {
							//Locate the subresource entry
							if session.Consumed[i].Key == resource.Fields[j].Key {
								//Since Value is an interface, have to parse it as integer before reassigning a value
								counter := resource.Fields[j].Value.(int32) + int32(session.Consumed[i].Amount)
								resource.Fields[i].Value = counter
							}
						}
					}
				}

				mgm.Coll(resource).Update(resource)
			} else {
				err = fmt.Errorf("error: resource not found")
			}

			mux["Resources"].Unlock()
		}

		jobID = session.JobID
		mgm.Coll(session).Delete(session)
	} else {
		err = fmt.Errorf("error: session not found")
	}

	mux["Sessions"].Unlock()

	return jobID, err
}

// CreateSessionBusiness godoc
func CreateSessionBusiness(requestData *m.SessionRequest, sessExt int, signingKey string, scheduler *scheduler.Scheduler, mux map[string]*sync.Mutex) (int, interface{}) {
	var err error
	var code int
	var response interface{}
	project := &m.Project{}

	mux["Projects"].Lock()

	//Looking for a project by the provided API key
	if err = mgm.Coll(project).First(bson.M{"apikey": requestData.APIKey}, project); err != nil {
		code, response = http.StatusNotFound, m.ProjectNotFound
	} else {
		//Assigning project relation to a new session
		newSession := &m.Session{
			Project: project.ID.Hex(),
		}

		mux["Sessions"].Lock()

		//Inserting a new session entry into the database
		mgm.Coll(newSession).Save(newSession)

		//Generating a JWT token for the session
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = newSession.ID.Hex() //Passing ObjectID of the session in the token
		claims["exp"] = time.Now().Add(time.Hour * time.Duration(sessExt)).Unix()

		//Schedule a termination task
		newSession.JobID = scheduler.Delay().Hour(sessExt).Do(TerminateSessionBusiness, newSession.ID.Hex(), mux)

		mgm.Coll(newSession).Update(newSession) //Updating session db entry to include job id

		if t, err := token.SignedString([]byte(signingKey)); err != nil {
			code, response = http.StatusInternalServerError, m.InternalError
		} else {
			//Returning the new token prefixed with Bearer in a response
			code, response = http.StatusOK, map[string]string{"token": "Bearer " + t}
		}

		mux["Sessions"].Unlock()
	}

	mux["Projects"].Unlock()

	return code, response
}

// RenewSessionBusiness creates new session timeout jobs for every session in the database. Meant to be executed at runtime to create timeouts for any previously created sessions
//Args:	scheduler, delay hours
//Rets:	error
func RenewSessionBusiness(token *jwt.Token, sessExt int, signingKey string, scheduler *scheduler.Scheduler, mux map[string]*sync.Mutex) (int, interface{}) {
	var err error
	var code int
	var response interface{}
	session := &m.Session{}

	claims := token.Claims.(jwt.MapClaims)
	sessID := claims["id"].(string)

	mux["Sessions"].Lock()

	if err = mgm.Coll(session).FindByID(sessID, session); err != nil {
		code, response = http.StatusNotFound, m.SessionNotFound
	} else {
		//Generating an updated expiration timestamp for a new bearer token
		claims["exp"] = time.Now().Add(time.Hour * time.Duration(sessExt)).Unix()

		//Cancel the old job, schedule a new one
		scheduler.CancelJob(session.JobID)
		session.JobID = scheduler.Delay().Hour(sessExt).Do(TerminateSessionBusiness, session.ID.Hex(), mux)

		//Update DB entry with new job id
		mgm.Coll(session).Update(session)

		//Create and send a new token
		if t, err := token.SignedString([]byte(signingKey)); err != nil {
			code, response = http.StatusInternalServerError, m.InternalError
		} else {
			//Returning the new token prefixed with Bearer in a response
			code, response = http.StatusOK, map[string]string{"token": "Bearer " + t}
		}
	}

	mux["Sessions"].Unlock()

	return code, response
}

// RecoverSessionsBusiness godoc
func RecoverSessionsBusiness(scheduler *scheduler.Scheduler, sessExt int, mux map[string]*sync.Mutex) error {
	var err error
	sessionsFound := []m.Session{}

	mux["Sessions"].Lock()

	//Searching for sessions
	_ = mgm.Coll(&m.Session{}).SimpleFind(&sessionsFound, bson.M{})

	//Creating expiration jobs for each session
	for _, s := range sessionsFound {
		s.JobID = scheduler.Delay().Hour(sessExt).Do(TerminateSessionBusiness, s.ID.Hex(), mux)
		mgm.Coll(&s).Update(&s) //At the moment each session's DB entry is updated manually as opposed to a single batch job because of mgm's limitations
	}

	mux["Sessions"].Unlock()

	return err
}

// CloseSessionBusiness godoc
func CloseSessionBusiness(sessID string, scheduler *scheduler.Scheduler, mux map[string]*sync.Mutex) (int, interface{}) {
	var err error
	var code int
	var response interface{}
	var jobID string

	//TerminateSession has its own lock on Sessions, no need to lock here

	//Verifying that the ObjectID contains 24 hexademical characters
	if !db.VerifyObjectIDString(sessID) {
		code, response = http.StatusBadRequest, m.InvalidID
	} else {
		if jobID, err = TerminateSessionBusiness(sessID, mux); err != nil {
			code, response = http.StatusNotFound, m.SessionNotFound
		} else {
			//Cancel scheduled termination job, if any
			scheduler.CancelJob(jobID)
			code, response = http.StatusOK, m.SessionTerminated
		}
	}

	return code, response
}
