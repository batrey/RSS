package handlers

import (
	"encoding/json"
	"net/http"
	"rss/app/models"
	"strings"
)

type BaseHandler struct {
	ArticleRep models.ArticlesRepo
}

func NewBaseHandler(articleRepo models.ArticlesRepo) *BaseHandler {
	return &BaseHandler{
		ArticleRep: articleRepo,
	}
}

func (db *BaseHandler) HandleTaskGetList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category, okcategory := r.URL.Query()["category"]
		cursor, ok := r.URL.Query()["cursor"]
		limit, okLimit := r.URL.Query()["limit"]

		if !okcategory || category[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !ok || len(cursor[0]) < 1 || !isNumeric(cursor[0]) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !okLimit || len(limit[0]) < 1 || !isNumeric(limit[0]) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tmp := category[0]
		categoryOk := strings.ToLower(tmp)
		categoryOk = CategoryCheck(categoryOk)

		resp, err := db.ArticleRep.PagnationArticles(categoryOk, cursor[0], limit[0])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

	}
}

func (db *BaseHandler) HandleTaskGetArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category, okcategory := r.URL.Query()["category"]
		cursor, ok := r.URL.Query()["cursor"]

		if !okcategory || category[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !ok || len(cursor[0]) < 1 || !isNumeric(cursor[0]) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tmp := category[0]
		categoryOk := strings.ToLower(tmp)
		categoryOk = CategoryCheck(categoryOk)
		resp, err := db.ArticleRep.GetOneArticle(categoryOk, cursor[0])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
}

func (db *BaseHandler) HandleTaskSendEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category, okcategory := r.URL.Query()["category"]
		cursor, okcursor := r.URL.Query()["cursor"]
		limit, okLimit := r.URL.Query()["limit"]
		addr, okaddr := r.URL.Query()["addr"]
		sender, oksender := r.URL.Query()["sender"]
		reciver, okreciver := r.URL.Query()["cursor"]

		if !okcategory || category[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !okreciver || reciver[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !oksender || sender[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !okaddr || addr[0] == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !okcursor || len(cursor[0]) < 1 || !isNumeric(cursor[0]) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !okLimit || len(limit[0]) < 1 || !isNumeric(limit[0]) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tmp := category[0]
		categoryOk := strings.ToLower(tmp)
		categoryOk = CategoryCheck(categoryOk)

		resp, err := db.ArticleRep.GetOneArticle(categoryOk, cursor[0])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tmpResp, err := json.Marshal(resp)
		err = Email(addr[0], reciver, sender[0], tmpResp)
		w.WriteHeader(http.StatusOK)
	}
}

// func (s *models.Server) respnd(w http.ResponseWriter, r *http.Request, resp interface{}, status int) {
// 	if data != nil {
// 		err := json.NewEncoder(w).Encode(resp)
// 		//TODO:
// 	}

// 	w.WriteHeader(status)
// 	w.WriteHeader(http.StatusOK)
// 	encoder := json.NewEncoder(w)
// 	encoder.SetEscapeHTML(false)
// 	if err := encoder.Encode(resp); err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(err.Error()))
// 	}
// }

// func (s *Server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
// 	return json.NewDecoder(r.Body).Decode(v)
// }
