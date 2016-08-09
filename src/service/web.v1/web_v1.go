package web_v1

import (
	"db"
	"encoding/json"
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const (
	ERR_NO_RECORD = "no rows in result set"
)

type APIError struct {
	StatusCode int
	Message    string
	Exception  error
}

func (a *APIError) Error() string {
	return fmt.Sprintf("[%d] %s: %s", a.StatusCode, http.StatusText(a.StatusCode), a.Exception.Error())
}

type Env struct {
	ds db.DataSource
	// session, etc.
}

func (e *Env) WithAPICtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r == nil {
				return
			}
			switch t := r.(type) {
			case APIError:
				log.Println(t.Error())
				http.Error(w, http.StatusText(t.StatusCode), t.StatusCode)
			default:
				panic(r)
			}
		}()

		r.Header.Set("Content-Type", "application/json")
		context.Set(r, "ds", e.ds)
		h.ServeHTTP(w, r)
		context.Clear(r)
	})
}

func NewEnv(ds db.DataSource) *Env {
	return &Env{
		ds: ds,
	}
}

func NewWebV1Router(env *Env) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	router.Handle("/users", env.WithAPICtx(http.HandlerFunc(User)))
	router.Handle("/users/{user_id:[0-9]+}/relationships", env.WithAPICtx(http.HandlerFunc(UserRelationShip)))
	router.Handle("/users/{user_id:[0-9]+}/relationships/{other_user_id:[0-9]+}", env.WithAPICtx(http.HandlerFunc(PutUserRelationShip)))
	return router
}

func User(w http.ResponseWriter, r *http.Request) {
	ds := context.Get(r, "ds").(db.DataSource)

	switch r.Method {
	case "GET":
		users, err := ds.Users()
		if err != nil {
			panic(APIError{500, "DB error", err})
		}
		if err = json.NewEncoder(w).Encode(users); err != nil {
			panic(err)
		}
	case "POST":
		user := &db.User{}
		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			panic(APIError{400, "Invalidate post data", err})
		}
		user, err = ds.CreateUser(user.Name)
		if err != nil {
			panic(APIError{500, "DB error", err})
		}
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(user); err != nil {
			panic(APIError{500, "Serialize fail ", err})
		}

	default:
		http.Error(w, http.StatusText(405), 405)
		return
	}
}

func UserRelationShip(w http.ResponseWriter, r *http.Request) {
	ds := context.Get(r, "ds").(db.DataSource)
	vars := mux.Vars(r)
	userId, err := strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	switch r.Method {
	case "GET":
		relationships, err := ds.UserRelationShips(userId)
		if err != nil {
			panic(APIError{500, "DB error", err})
		}
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(relationships); err != nil {
			panic(APIError{500, "DB error", err})
		}
	default:
		http.Error(w, http.StatusText(405), 405)
	}
}

func PutUserRelationShip(w http.ResponseWriter, r *http.Request) {
	ds := context.Get(r, "ds").(db.DataSource)
	vars := mux.Vars(r)
	uid, err := strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil {
		panic(APIError{400, "Invalidate user_id", err})
	}
	oid, err := strconv.ParseInt(vars["other_user_id"], 10, 64)
	if err != nil {
		panic(APIError{400, "Invalidate other_user_id", err})
	}
	switch r.Method {
	case "PUT":
		stateObj := &struct{ State string }{}
		err := json.NewDecoder(r.Body).Decode(stateObj)
		if err != nil || stateObj.State == db.RELATIONSHIP_MATCHED {
			panic(APIError{400, "Invalidate post data", err})
		}

		myRS, err := ds.CreateOrUpdateRelationShip(uid, oid, stateObj.State, db.RELATIONSHIP_TYPE_RS)
		if err != nil {
			panic(APIError{500, "DB error", err})
		}
		hisRSIsEmpty := false
		hisRS, err := ds.UserRelationShip(oid, uid, db.RELATIONSHIP_TYPE_RS)
		if err != nil {
			if strings.Contains(err.Error(), ERR_NO_RECORD) {
				hisRSIsEmpty = true
			} else {
				panic(APIError{500, "DB error", err})
			}
		}

		w.WriteHeader(http.StatusOK)

		if myRS.State == db.RELATIONSHIP_LIKE && !hisRSIsEmpty && hisRS.State == db.RELATIONSHIP_LIKE {
			myNewRS, _ := ds.CreateOrUpdateRelationShip(uid, oid, db.RELATIONSHIP_MATCHED, db.RELATIONSHIP_TYPE_RS)
			ds.CreateOrUpdateRelationShip(oid, uid, db.RELATIONSHIP_MATCHED, db.RELATIONSHIP_TYPE_RS)
			if err = json.NewEncoder(w).Encode(*myNewRS); err != nil {
				panic(APIError{500, "Serialize fail ", err})
			}
			return
		} else if myRS.State == db.RELATIONSHIP_DISLIKE && !hisRSIsEmpty && hisRS.State == db.RELATIONSHIP_MATCHED {
			ds.CreateOrUpdateRelationShip(oid, uid, db.RELATIONSHIP_LIKE, db.RELATIONSHIP_TYPE_RS)
		}

		if err = json.NewEncoder(w).Encode(*myRS); err != nil {
			panic(APIError{500, "Serialize fail ", err})
		}
		return

	default:
		http.Error(w, http.StatusText(405), 405)
		return
	}
}
