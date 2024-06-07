package controller

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"

	"datasupervision/internal/service"
)

const (
	PostgresDBType = "PostgreSQL"
	MySQLDBType    = "MySQL"
)

type Controller struct {
	ConnectorDB ConnectorDB
}

type ConnectorDB interface {
	GetDatabases() ([]string, error)
	GetSchemas() ([]string, error)

	GetTables(schemaName string) ([]string, error)
	GetColumns(schemaName, tableName string) ([]service.TableColumn, error)

	GetTablesAndColumns(schemaName string) (map[string][]string, error)
	Join(schemaName, table1, column1, table2, column2 string) (*service.TableData, error)

	GetKeys(schemaName, tableName string) ([]service.TableKey, error)
	GetPrimaryKeys(schemaName, tableName string) ([]string, error)
	GetForeignKeys(schemaName, tableName string) ([]service.ForeignKey, error)

	InsertRow(schemaName, tableName string, row map[string]interface{}) error
	InsertRows(schemaName, tableName string, column []string, row [][]interface{}) error
	UpdateRow(schemaName, tableName string, row service.UpdateRow) error
	DeleteRow(schemaName, tableName string, row service.Row) error

	Select(query string, args ...any) (*service.TableData, error)
	SelectTableData(schemaName, tableName string) (*service.TableData, error)
	SelectWithFilter(schemaName, tableName, columnName, filterValue string) (*service.TableData, error)

	CreateTable(schemaName string, tableData service.CreatedTableData) error
	GetColumnTypes() []string
	DeleteTable(schemaName, tableName string) error

	BeginTransaction() error
	Rollback() error
	Commit() error
}

func (c *Controller) GetHealthcheck(w http.ResponseWriter, r *http.Request) {
	println("test")

	w.WriteHeader(http.StatusNotImplemented)
}

type Database struct {
	Databases []string
}

func (c *Controller) DatabasesHandler(w http.ResponseWriter, r *http.Request) {
	databases, err := c.ConnectorDB.GetDatabases()
	if err != nil {
		log.WithField("func", "DatabasesHandler").Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("html/databases.html")
	if err != nil {
		log.WithField("func", "DatabasesHandler").Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, Database{Databases: databases})
	if err != nil {
		log.WithField("func", "DatabasesHandler").Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type DashboardData struct {
	Schemas []string
}

func (c *Controller) SelectNewSchemas(w http.ResponseWriter, r *http.Request) {
	databasesName := chi.URLParam(r, "databaseName")
	log.Info(databasesName)

	http.Redirect(w, r, "/schemas", http.StatusSeeOther)
}

func (c *Controller) Schemas(w http.ResponseWriter, r *http.Request) {
	if c.ConnectorDB == nil {
		return
	}

	schemas, err := c.ConnectorDB.GetSchemas()
	if err != nil {
		log.WithField("func", "Schemas").Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("html/schemas.html"))

	data := DashboardData{
		Schemas: schemas,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.WithField("func", "Schemas").Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type TablesData struct {
	SchemaName string
	Tables     []string
}

func (c *Controller) TablesHandler(w http.ResponseWriter, r *http.Request) {
	if c.ConnectorDB == nil {
		return
	}

	schemaName := chi.URLParam(r, "schemaName")
	log.Info(schemaName)

	tables, err := c.ConnectorDB.GetTables(schemaName)
	if err != nil {
		log.WithField("func", "TablesHandler").Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("html/tables.html")
	if err != nil {
		log.WithField("func", "TablesHandler").Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := TablesData{SchemaName: schemaName, Tables: tables}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.WithField("func", "TablesHandler").Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type TableInfo struct {
	Columns []service.TableColumn `json:"columns"`
	Keys    []service.TableKey    `json:"keys"`
}

func (c *Controller) TableInfoHandler(w http.ResponseWriter, r *http.Request) {
	schemaName := chi.URLParam(r, "schemaName")
	log.Info(schemaName)

	tableName := strings.TrimPrefix(r.URL.Path, "/tableinfo/")

	columns, err := c.ConnectorDB.GetColumns(schemaName, tableName)
	if err != nil {
		log.Error()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	keys, err := c.ConnectorDB.GetKeys(schemaName, tableName)
	if err != nil {
		log.WithField("func", "TableInfoHandler").Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tableInfo := TableInfo{Columns: columns, Keys: keys}

	jsonData, err := json.Marshal(tableInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
