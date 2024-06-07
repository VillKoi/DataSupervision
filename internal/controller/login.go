package controller

import (
	"datasupervision/internal/mysql"
	"datasupervision/internal/postgres"
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type LoginForm struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("html/authorization.html")
		if err != nil {
			log.WithField("func", "Login").Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		loginForm := LoginForm{
			Host:     "localhost",
			Port:     "5432",
			Database: "blog",
			Username: "postgres",
			Password: "postgres",
		}

		err = tmpl.Execute(w, loginForm)
		if err != nil {
			log.WithField("func", "Login").Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (c *Controller) Login2(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.WithField("func", "Login2").Error(err)
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		host := r.FormValue("host")
		port := r.FormValue("port")
		dbType := r.FormValue("dbType")
		database := r.FormValue("database")
		username := r.FormValue("username")
		password := r.FormValue("password")

		if dbType == PostgresDBType {
			connectorDB, err := postgres.NewAuth(&postgres.AuthConfig{
				Host:     host,
				Port:     port,
				Database: database,
				Username: username,
				Password: password,
				SSLMode:  "disable",
			})
			if err != nil {
				log.WithField("func", "Login2").Error(err)
				http.Error(w, "Authentication failed", http.StatusUnauthorized)
				return
			}

			c.ConnectorDB = connectorDB

			http.Redirect(w, r, "/schemas", http.StatusSeeOther)
			return
		}

		if dbType == MySQLDBType {
			connectorDB, err := mysql.NewConnect(&mysql.AuthConfig{
				Host:     host,
				Port:     port,
				Database: database,
				Username: username,
				Password: password,
				SSLMode:  "disable",
			})
			if err != nil {
				log.WithField("func", "Login2").Error(err)
				http.Error(w, "Authentication failed", http.StatusUnauthorized)
				return
			}

			c.ConnectorDB = connectorDB

			http.Redirect(w, r, "/schemas", http.StatusSeeOther)
			return
		}

		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}
}
