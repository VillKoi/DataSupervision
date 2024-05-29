package controller

import (
	"bytes"
	"datasupervision/internal/service"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

func (c *Controller) InsertRowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	schemaName := chi.URLParam(r, "schemaName")

	tableName := r.URL.Query().Get("tableName")

	var row service.Row
	if err := json.NewDecoder(r.Body).Decode(&row); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := c.ConnectorDB.InsertRow(schemaName, tableName, row.Row)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Row inserted successfully"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (c *Controller) UpdateRowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	schemaName := chi.URLParam(r, "schemaName")

	tableName := r.URL.Query().Get("tableName")

	var row service.UpdateRow
	if err := json.NewDecoder(r.Body).Decode(&row); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := c.ConnectorDB.UpdateRow(schemaName, tableName, row)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Row edited successfully"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (c *Controller) DeleteRowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	schemaName := chi.URLParam(r, "schemaName")

	tableName := r.URL.Query().Get("tableName")

	var row service.Row
	if err := json.NewDecoder(r.Body).Decode(&row); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Delete row into table %s\n", tableName)
	fmt.Println("Data:", row.Row)

	err := c.ConnectorDB.DeleteRow(schemaName, tableName, row)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Row delete successfully"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (c *Controller) InsertRowsHandlerJson(w http.ResponseWriter, r *http.Request) {
	schemaName := chi.URLParam(r, "schemaName")
	log.Info(schemaName)

	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("file")
	if err != nil {
		log.WithField("func", "InsertRowsHandler").Error(err)
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		log.WithField("func", "InsertRowsHandler").Error(err)
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	var insertRows service.InsertRows
	err = json.Unmarshal(buf.Bytes(), &insertRows)
	if err != nil {
		log.WithField("func", "InsertRowsHandler").Error(err)
		http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed JSON: %+v\n", insertRows)

	err = c.ConnectorDB.InsertRows(schemaName, insertRows.TableName, insertRows.Columns, insertRows.Rows)
	if err != nil {
		log.WithField("func", "InsertRowsHandler").Error(err)
		http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"success"}`))
}

func convertToInterfaceArray(record []string) []interface{} {
	var row []interface{}
	for _, field := range record {
		row = append(row, field)
	}
	return row
}

func (c *Controller) InsertRowsHandlerCVS(w http.ResponseWriter, r *http.Request) {
	schemaName := chi.URLParam(r, "schemaName")
	tableName := r.URL.Query().Get("tableName")

	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true

	// Чтение первой строки как названия колонок
	columns, err := reader.Read()
	if err != nil {
		http.Error(w, "Error reading CSV file", http.StatusInternalServerError)
		return
	}

	var rows [][]interface{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.WithField("func", "InsertRowsHandler").Errorf("reader.Read(): %s", err.Error())
			http.Error(w, "Error reading CSV file", http.StatusInternalServerError)
			return
		}
		rows = append(rows, convertToInterfaceArray(record))
	}

	fmt.Printf("Parsed JSON: %+v\n", rows)

	err = c.ConnectorDB.InsertRows(schemaName, tableName, columns, rows)
	if err != nil {
		log.WithField("func", "InsertRowsHandler").Error(err)
		http.Error(w, "Unable to parse JSON", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"success"}`))
}
