package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	db "rss/app/db"
	"rss/app/handlers"
	"rss/app/server"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func middleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)

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
	//Connect to redis
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	//Connect to postgres
	d := db.DataBase{}
	database, err := d.ConnectDb()
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}
	//close db connection
	defer database.Conn.Close()

	envTime := os.Getenv("PULL_TIME")
	pullTime, err := time.ParseDuration(envTime)
	fmt.Println("PostgreSQL and Redis connected successfully...")
	ticker := time.NewTicker(pullTime * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Getting sky Articles")
				server.TaskSky(*database)
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
				server.TaskBbc(*database)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	listHandle := http.HandlerFunc(handlers.GetList(*database))
	oneHandle := http.HandlerFunc(handlers.GetArticle(*database))
	mailHandler := http.HandlerFunc(handlers.ShareEmail(*database))
	mux.Handle("/list", middleWare(listHandle))
	mux.Handle("/one", middleWare(oneHandle))
	mux.Handle("/email", middleWare(mailHandler))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Server started on port 8080")
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
