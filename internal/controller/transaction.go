package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

func (c *Controller) BeginTransaction(w http.ResponseWriter, r *http.Request) {
	schemaName := chi.URLParam(r, "schemaName")
	log.Info(schemaName)

	err := c.ConnectorDB.BeginTransaction()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"Begin Transaction successfully"}`))

}

func (c *Controller) Rollback(w http.ResponseWriter, r *http.Request) {
	schemaName := chi.URLParam(r, "schemaName")
	log.Info(schemaName)

	err := c.ConnectorDB.Rollback()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"Rollback successfully"}`))
}

func (c *Controller) Commit(w http.ResponseWriter, r *http.Request) {
	schemaName := chi.URLParam(r, "schemaName")
	log.Info(schemaName)

	err := c.ConnectorDB.Commit()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"Commit successfully"}`))
}
