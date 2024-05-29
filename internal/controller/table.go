package controller

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

type TableDataResponse struct {
	TableName string
	Columns   []string
	Rows      [][]interface{}
}

type TableRow map[string]interface{}

func (c *Controller) TableDataHandler(w http.ResponseWriter, r *http.Request) {
	tableName := chi.URLParam(r, "tableName")
	log.Info(tableName)

	schemaName := chi.URLParam(r, "schemaName")
	log.Info(schemaName)

	tableData, err := c.ConnectorDB.SelectTableData(schemaName, tableName)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := TableDataResponse{
		TableName: tableName,
		Columns:   tableData.Columns,
		Rows:      tableData.Rows,
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

func convertToString(value interface{}) string {
	switch vv := value.(type) {
	case *interface{}:
		vvv := *vv
		t := reflect.TypeOf(vvv)
		if vvv == nil {
			return "null"
		}
		if t.Kind() == reflect.String {
			return fmt.Sprintf("%s", vvv)
		}
		return fmt.Sprintf("%#v", vvv)
	}
	return fmt.Sprintf("%#v", value)
}

func (c *Controller) DownloadCSVHandler(w http.ResponseWriter, r *http.Request) {
	tableName := chi.URLParam(r, "tableName")
	log.Info(tableName)

	schemaName := chi.URLParam(r, "schemaName")
	log.Info(schemaName)

	tableData, err := c.ConnectorDB.SelectTableData(schemaName, tableName)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := TableDataResponse{
		TableName: tableName,
		Columns:   tableData.Columns,
		Rows:      tableData.Rows,
	}

	// Устанавливаем заголовки для ответа
	w.Header().Set("Content-Disposition", "attachment; filename=data.csv")
	w.Header().Set("Content-Type", "text/csv")

	// Создаем новый CSV writer
	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Записываем заголовки колонок
	if err := writer.Write(data.Columns); err != nil {
		http.Error(w, "Unable to write CSV columns", http.StatusInternalServerError)
		return
	}

	// Записываем строки данных
	for _, row := range data.Rows {
		var stringRow []string
		for _, value := range row {
			stringRow = append(stringRow, convertToString(value))
		}
		if err := writer.Write(stringRow); err != nil {
			http.Error(w, "Unable to write CSV row", http.StatusInternalServerError)
			return
		}
	}
}

func (c *Controller) ApplyFilter(w http.ResponseWriter, r *http.Request) {
	schemaName := chi.URLParam(r, "schemaName")
	log.Info(schemaName)
	// Получаем параметры запроса
	r.ParseForm()
	tableName := r.FormValue("tableName")
	columnName := r.FormValue("columnName")
	filterValue := r.FormValue("filterValue")

	tableData, err := c.ConnectorDB.SelectWithFilter(schemaName, tableName, columnName, filterValue)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем данные в формате HTML
	response := TableDataResponse{
		TableName: tableName,
		Columns:   tableData.Columns,
		Rows:      tableData.Rows,
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
