package controller

import (
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Welcome struct {
	Content string
}

func (c *Controller) Welcome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("html/welcome.html")
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, Welcome{Content: "Добро пожаловать"})
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
