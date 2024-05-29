package controller

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

type JoinRequest struct {
	Table1  string `json:"table1"`
	Column1 string `json:"column1"`
	Table2  string `json:"table2"`
	Column2 string `json:"column2"`
}

type Results struct {
	Columns []string
	Rows    [][]interface{}
}

type JoinData struct {
	Schema     string
	Tables     map[string][]string `json:"tables"`
	TablesJSON template.JS
	Results    Results
}

func (c *Controller) PreJoinHandler(w http.ResponseWriter, r *http.Request) {
	schemaName := chi.URLParam(r, "schemaName")
	log.Info(schemaName)

	tables, err := c.ConnectorDB.GetTablesAndColumns(schemaName)
	if err != nil {
		log.Error("Error fetching tables and columns:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("html/join.html")
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

	data := JoinData{
		Schema:     schemaName,
		Tables:     tables,
		TablesJSON: template.JS(jsonString),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Error("Error executing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) JoinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	schemaName := chi.URLParam(r, "schemaName")
	log.Info(schemaName)

	table1 := r.FormValue("table1")
	column1 := r.FormValue("column1")
	table2 := r.FormValue("table2")
	column2 := r.FormValue("column2")

	results, err := c.ConnectorDB.Join(table1, column1, table2, column2)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tables, err := c.ConnectorDB.GetTablesAndColumns(schemaName)
	if err != nil {
		log.Error("Error fetching tables and columns:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonString, err := json.Marshal(tables)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := JoinData{
		Schema:     schemaName,
		Tables:     tables,
		TablesJSON: template.JS(jsonString),
		Results: Results{
			Columns: results.Columns,
			Rows:    results.Rows,
		},
	}

	tmpl, err := template.ParseFiles("html/join.html")
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Error("Error executing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
