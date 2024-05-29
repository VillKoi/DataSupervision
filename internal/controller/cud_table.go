package controller

import (
	"datasupervision/internal/service"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

type CreatedData struct {
	Schema     string
	Tables     template.JS `json:"tables"`
	TablesJSON template.JS
}

func (c *Controller) RenderCreateTable(w http.ResponseWriter, r *http.Request) {
	schemaName := chi.URLParam(r, "schemaName")
	log.Info(schemaName)

	tables, err := c.ConnectorDB.GetTablesAndColumns(schemaName)
	if err != nil {
		log.Error("Error fetching tables and columns:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("html/create_table.html")
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonString, err := json.Marshal(tables)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tt := []string{}
	for table := range tables {
		tt = append(tt, table)
	}

	ttJsonString, err := json.Marshal(tt)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := CreatedData{
		Schema:     schemaName,
		Tables:     template.JS(ttJsonString),
		TablesJSON: template.JS(jsonString),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) CreateTable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	schemaName := chi.URLParam(r, "schemaName")
	log.Info(schemaName)
	if schemaName == "" {
		http.Error(w, "Schema not specified", http.StatusBadRequest)
		return
	}

	var tableData service.CreatedTableData
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&tableData)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.ConnectorDB.CreateTable(schemaName, tableData)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"Table created successfully"}`))
}

func (c *Controller) GetColumnTypesHandler(w http.ResponseWriter, r *http.Request) {
	columnTypes := c.ConnectorDB.GetColumnTypes()

	jsonData, err := json.Marshal(columnTypes)
	if err != nil {
		log.WithField("func", "GetColumnTypesHandler").Error("Error marshalling column types:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (c *Controller) DeleteTableHandler(w http.ResponseWriter, r *http.Request) {
	schemaName := chi.URLParam(r, "schemaName")
	log.Info(schemaName)

	tableName := r.URL.Query().Get("name")
	if tableName == "" {
		http.Error(w, "Missing table name", http.StatusBadRequest)
		return
	}

	err := c.ConnectorDB.DeleteTable(schemaName, tableName)
	if err != nil {
		log.Error(err)
		http.Error(w, "Error deleting table: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
