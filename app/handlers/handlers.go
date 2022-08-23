package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rss/app/models"
	"rss/app/server"
	"strings"
	"time"
)

type BaseHandler struct {
	ArticleRep models.ArticlesRepo
}

func NewBaseHandler(articleRepo models.ArticlesRepo) *BaseHandler {
	return &BaseHandler{
		ArticleRep: articleRepo,
	}
}

// Returns multi Articles
func (db *BaseHandler) HandleTaskGetList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category, okCategory := r.URL.Query()["category"]
		cursor, ok := r.URL.Query()["cursor"]
		limit, okLimit := r.URL.Query()["limit"]

		if !okCategory || category[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Something went wrong with the Category", http.StatusBadRequest)
			log.Println("Error with Category in  HandleTaskGetList err ")
			return
		}

		if !ok || len(cursor[0]) < 1 || !isNumeric(cursor[0]) {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "No Cursor provided", http.StatusBadRequest)
			log.Println("Error with Cursor in  HandleTaskGetList err ")
			return
		}

		if !okLimit || len(limit[0]) < 1 || !isNumeric(limit[0]) {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Something went wrong with the  limit", http.StatusBadRequest)
			log.Println("Error with LIMIT in  HandleTaskGetList err ")
			return
		}

		tmp := category[0]
		categoryOk := strings.ToLower(tmp)
		categoryOk = CategoryCheck(categoryOk)

		resp, err := db.ArticleRep.PaginationArticles(categoryOk, cursor[0], limit[0])
		if err != nil {
			log.Printf("Error Reading Articles in HandleTaskGetList func PaginationArticles err %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(resp); err != nil {
			log.Printf("Error Encoding Response in func HandleTaskGetList err %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

	}
}

// Returns One Article
func (db *BaseHandler) HandleTaskGetArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category, okCategory := r.URL.Query()["category"]
		cursor, ok := r.URL.Query()["cursor"]

		if !okCategory || category[0] == "" {
			log.Println("Error  with in category func HandleTaskGetArticle err")
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Something went wrong with the category", http.StatusBadRequest)
			return
		}
		if !ok || len(cursor[0]) < 1 || !isNumeric(cursor[0]) {
			log.Println("Error  with in cursor func HandleTaskGetArticle err")
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Something went wrong with the Cursor", http.StatusBadRequest)
			return
		}

		tmp := category[0]
		categoryOk := strings.ToLower(tmp)
		categoryOk = CategoryCheck(categoryOk)
		resp, err := db.ArticleRep.GetOneArticle(categoryOk, cursor[0])
		if err != nil {
			log.Printf("Error  with in  func HandleTaskGetArticle  with GetOneArticle err %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(resp); err != nil {
			log.Printf("Error  with in Encode func HandleTaskGetArticle err %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
}

// Sends email
func (db *BaseHandler) HandleTaskSendEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category, okCategory := r.URL.Query()["category"]
		cursor, okCursor := r.URL.Query()["cursor"]
		sender, okSender := r.URL.Query()["sender"]
		receiver, okReceiver := r.URL.Query()["To"]

		if !okCategory || category[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Error  with in category func HandleTaskSendEmail")
			http.Error(w, "Something went wrong with the Category", http.StatusBadRequest)
			return
		}
		if !okReceiver || receiver[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Error  with in receiver func HandleTaskSendEmail")
			http.Error(w, "Something went wrong with the Receiver", http.StatusBadRequest)
			return
		}
		if !okSender || sender[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Error  with in sender func HandleTaskSendEmail")
			http.Error(w, "Something went wrong with the Sender ", http.StatusBadRequest)
			return
		}

		if !okCursor || len(cursor[0]) < 1 || !isNumeric(cursor[0]) {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Error  with in cursor func HandleTaskSendEmail")
			http.Error(w, "Something went wrong with the Cursor", http.StatusBadRequest)
			return
		}

		tmp := category[0]
		categoryOk := strings.ToLower(tmp)
		categoryOk = CategoryCheck(categoryOk)

		resp, err := db.ArticleRep.GetOneArticle(categoryOk, cursor[0])
		if err != nil {
			log.Printf("Error  with  GetOneArticle func HandleTaskSendEmail err: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tmpResp, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Error  with  Marshal func HandleTaskSendEmail err: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = Email(receiver, sender[0], tmpResp)
		if err != nil {
			log.Printf("Error  with  Sending Email func HandleTaskSendEmail err: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (db *BaseHandler) HandlerPullArticles(pullTime time.Duration) {
	t := server.TaskNewBaseHandle(db.ArticleRep)
	ticker := time.NewTicker(pullTime * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Getting sky Articles")
				t.TaskSky()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Getting BBC Articles ")
				t.TaskBbc()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
