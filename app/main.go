package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	database "rss/app/db"
	"rss/app/handlers"

	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

type Server struct {
	Db     *sql.DB
	Router *http.ServeMux
}

func middleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			//Add some auth here
		}

		next.ServeHTTP(w, r)

	})
}

func main() {
	//Load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db := database.ConnectDb()
	articlesRepo := database.NewArticleRepo(db)
	h := handlers.NewBaseHandler(articlesRepo)

	//close db connection
	defer articlesRepo.Close()
	fmt.Println("PostgreSQL connected successfully...")

	//Create a go routine that  pulls data every x amount
	envTime := os.Getenv("PULL_TIME")
	pullTime, err := time.ParseDuration(envTime)
	h.HandlerPullArticles(pullTime)

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	listHandle := http.HandlerFunc(h.HandleTaskGetList())
	oneHandle := http.HandlerFunc(h.HandleTaskGetArticle())
	mailHandler := http.HandlerFunc(h.HandleTaskSendEmail())
	mux.Handle("/list", middleWare(listHandle))
	mux.Handle("/one", middleWare(oneHandle))
	mux.Handle("/email", middleWare(mailHandler))

	//Gracefully shutdown server when interrupt is given
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	<-done
	log.Println("Server Stopped on port 8080")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")

}
