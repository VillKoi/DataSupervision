package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	ServerPort     string        `yaml:"server_port"`
	BasePath       string        `yaml:"base_path"`
	RequestTimeout time.Duration `yaml:"request_timeout"`
	Validation     bool          `yaml:"validation"`
}

type Server struct {
	basePath string

	requestTimeout time.Duration
	router         http.Handler

	Controller *Controller
	Server     *http.Server
	// authManager   *auth.Manager
}

func NewServer(config *Config) (*Server, error) {
	if config == nil {
		log.Fatal("config file is empty for new server")

		return nil, errors.New("config file is empty for new server")
	}

	controller := &Controller{}

	s := &Server{
		basePath:       config.BasePath,
		requestTimeout: config.RequestTimeout,
		Controller:     controller,
		// authManager:    authManager,
	}

	mux, errRouter := s.NewRouter(true)
	if errRouter != nil {
		log.WithError(errRouter).Fatal("could not create router")
	}

	s.Server = &http.Server{
		Addr:              config.ServerPort,
		ReadTimeout:       config.RequestTimeout,
		ReadHeaderTimeout: config.RequestTimeout,
		Handler:           mux,
	}

	return s, nil
}

func (s *Server) NewRouter(validation bool) (http.Handler, error) {
	mux := chi.NewRouter()

	mux.Use(middleware.NoCache)

	mux.HandleFunc("/start", s.Controller.Login)
	mux.HandleFunc("/login", s.Controller.Login2)

	mux.HandleFunc("/databases", s.Controller.DatabasesHandler)
	mux.HandleFunc("/databases/{databaseName}/schemas", s.Controller.SelectNewSchemas)
	mux.HandleFunc("/schemas", s.Controller.Schemas)
	mux.HandleFunc("/{schemaName}/tables", s.Controller.TablesHandler)

	mux.HandleFunc("/{schemaName}/tableinfo/{tableName}", s.Controller.TableInfoHandler)
	mux.HandleFunc("/{schemaName}/tabledata/{tableName}", s.Controller.TableDataHandler)

	mux.HandleFunc("/{schemaName}/download/json/{tableName}", s.Controller.TableDataHandler)
	mux.HandleFunc("/{schemaName}/download/csv/{tableName}", s.Controller.DownloadCSVHandler)

	mux.HandleFunc("/{schemaName}/applyfilter", s.Controller.ApplyFilter)
	mux.HandleFunc("/{schemaName}/sql-input", s.Controller.SQLInput)

	mux.HandleFunc("/{schemaName}/insert-row", s.Controller.InsertRowHandler)
	mux.HandleFunc("/{schemaName}/insert-rows/json", s.Controller.InsertRowsHandlerJson)
	mux.HandleFunc("/{schemaName}/insert-rows/csv", s.Controller.InsertRowsHandlerCVS)

	mux.HandleFunc("/{schemaName}/update-row", s.Controller.UpdateRowHandler)
	mux.HandleFunc("/{schemaName}/delete-row", s.Controller.DeleteRowHandler)

	mux.HandleFunc("/welcome", s.Controller.Welcome)

	mux.HandleFunc("/{schemaName}/pre-join", s.Controller.PreJoinHandler)
	mux.HandleFunc("/{schemaName}/join", s.Controller.JoinHandler)

	mux.HandleFunc("/{schemaName}/schema", s.Controller.CreateSchema)
	mux.HandleFunc("/{schemaName}/graph", s.Controller.Graph)

	mux.HandleFunc("/{schemaName}/render-create-table", s.Controller.RenderCreateTable)
	mux.HandleFunc("/{schemaName}/create-table", s.Controller.CreateTable)
	mux.HandleFunc("/column-types", s.Controller.GetColumnTypesHandler)

	mux.HandleFunc("/{schemaName}/delete-table", s.Controller.DeleteTableHandler)

	mux.HandleFunc("/{schemaName}/begin-transaction", s.Controller.BeginTransaction)
	mux.HandleFunc("/{schemaName}/rollback", s.Controller.Rollback)
	mux.HandleFunc("/{schemaName}/commit", s.Controller.Commit)

	fsCss := http.FileServer(http.Dir("./css"))
	mux.Handle("/css/*", http.StripPrefix("/css", fsCss))

	fsJS := http.FileServer(http.Dir("./js"))
	mux.Handle("/js/*", http.StripPrefix("/js", fsJS))

	s.router = mux

	return mux, nil
}
