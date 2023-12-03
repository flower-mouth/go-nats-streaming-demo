package router

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"wbLab0/internal/models"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/homePage.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
}

func IdPage(w http.ResponseWriter, r *http.Request) {
	needId := r.URL.Query().Get("id")
	if _, ok := models.Cache[needId]; ok {
		b, _ := json.Marshal(models.Cache[needId])
		_, err := w.Write(b)
		if err != nil {
			log.Println(err)
		}
	} else {
		_, err := w.Write([]byte("Record not found"))
		if err != nil {
			log.Println(err)
		}
	}
}

func DataListPage(w http.ResponseWriter, r *http.Request) {
	outputArray := make([]models.Order, 0)
	for _, elem := range models.Cache {
		outputArray = append(outputArray, elem)
	}

	b, _ := json.Marshal(outputArray)
	_, err := w.Write(b)
	if err != nil {
		log.Println(err)
	}
}
