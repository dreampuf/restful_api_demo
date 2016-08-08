package service

import (
	"net/http"
	"service/web.v1"
	"fmt"
	"time"
	"log"
	"db"
)

type WebService struct {
	handler http.Handler
	addr string
}

func NewWebService(host string, port uint, db db.DataSource) *WebService {
	return &WebService{
		handler: web_v1.NewWebV1Router(web_v1.NewEnv(db)),
		addr: fmt.Sprintf("%s:%d", host, port),
	}
}

func (w *WebService) Serve() {
	srv := &http.Server{
		Handler: w.handler,
		Addr: w.addr,

		WriteTimeout: 10 * time.Second,
		ReadTimeout: 10 * time.Second,
	}
	log.Println("WebService started...")
	log.Fatal(srv.ListenAndServe())
}

