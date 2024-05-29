package controller

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type SQLInput struct {
	Query string `yaml:"query"`
}

func (c *Controller) SQLInput(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Body)

	var body SQLInput
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info(body.Query)

	tableData, err := c.ConnectorDB.Select(body.Query)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := TableDataResponse{
		Columns: tableData.Columns,
		Rows:    tableData.Rows,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// SELECT * FROM users;
