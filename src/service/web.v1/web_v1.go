package web_v1

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/gorilla/context"
	"db"
	"log"
	"encoding/json"
)

type Env struct {
	db db.DataSource
	// session, etc.
}

func (e *Env) WithAPICtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Content-Type", "application/json")
		context.Set(r, "db", e.db)
		h.ServeHTTP(w, r)
		context.Clear(r)
	})
}

func NewEnv(db db.DataSource) *Env {
	return &Env{
		db: db,
	}
}

func NewWebV1Router(env *Env) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	router.Handle("/users", env.WithAPICtx(http.HandlerFunc(User)))
	router.Handle("/users/:user_id/relationships", env.WithAPICtx(http.HandlerFunc(UserRelationShip)))
	return router
}

func User(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "db").(db.DataSource)

	switch r.Method {
	case "GET":
		users, err := db.Users()
		if err != nil {
			log.Fatal(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		if err = json.NewEncoder(w).Encode(users); err != nil {
			log.Fatal(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
	case "POST":
		user := &db.User{}
		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		user, err = db.CreateUser(user.Name)
		if err != nil {
			log.Fatal(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(user); err != nil {
			log.Fatal(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

	default:
		http.Error(w, http.StatusText(405), 405)
		return
	}
}

func UserRelationShip(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "db").(db.DataSource)
}
