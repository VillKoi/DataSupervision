package controller

import (
	"datasupervision/internal/service"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

type DatabaseSchema struct {
	SchemaName string
	Tables     map[string]Table
}

type Table struct {
	Name        string
	Columns     []service.TableColumn
	PrimaryKeys []string
	ForeignKeys []service.ForeignKey
}

func (c *Controller) GetSchema(schemaName string) (*DatabaseSchema, error) {
	tables, err := c.ConnectorDB.GetTables(schemaName)
	if err != nil {
		return nil, err
	}

	schema := &DatabaseSchema{
		SchemaName: schemaName,
		Tables:     make(map[string]Table),
	}

	for _, table := range tables {
		columns, err := c.ConnectorDB.GetColumns(schemaName, table)
		if err != nil {
			return nil, err
		}

		primaryKeys, err := c.ConnectorDB.GetPrimaryKeys(schemaName, table)
		if err != nil {
			return nil, err
		}

		dbforeignKeys, err := c.ConnectorDB.GetForeignKeys(schemaName, table)
		if err != nil {
			return nil, err
		}

		foreignKeys := make([]service.ForeignKey, 0, len(dbforeignKeys))
		for _, key := range dbforeignKeys {
			foreignKeys = append(foreignKeys, service.ForeignKey{
				Column:           key.Column,
				ReferencedTable:  key.ReferencedTable,
				ReferencedColumn: key.ReferencedColumn,
			})
		}

		schema.Tables[table] = Table{
			Name:        table,
			Columns:     columns,
			PrimaryKeys: primaryKeys,
			ForeignKeys: foreignKeys,
		}
	}

	log.Info(schema)
	return schema, nil
}

func (c *Controller) CreateSchema(w http.ResponseWriter, r *http.Request) {
	schemaName := chi.URLParam(r, "schemaName")
	log.Info(schemaName)

	schema, err := c.GetSchema(schemaName)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("html/schema.html")
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, schema)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) Graph(w http.ResponseWriter, r *http.Request) {
	schemaName := chi.URLParam(r, "schemaName")
	log.Info(schemaName)

	schema, err := c.GetSchema(schemaName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("html/graph.html")
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, struct{ SchemaJSON string }{SchemaJSON: string(schemaJSON)})
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
