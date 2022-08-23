package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type ArticleRepo struct {
	Conn *sql.DB
}

func NewArticleRepo(db *sql.DB) *ArticleRepo {
	return &ArticleRepo{
		Conn: db,
	}
}

func (db *ArticleRepo) Close() error {
	return db.Conn.Close()
}

// Connects to the DB
func ConnectDb() *sql.DB {
	var err error
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB"))
	Conn, err := sql.Open("postgres", url)
	if err != nil {
		log.Printf("could not connect to postgres database: %v", err)
		return Conn
	}

	err = Conn.Ping()
	if err != nil {
		log.Printf("Error Pinging  DB err %s", err)
		return Conn
	}

	fmt.Println(url)
	return Conn
}

// Adds Article to the DB.
func (db *ArticleRepo) AddArticles(category string, article interface{}) (err error) {
	tmp, err := json.Marshal(article)
	if err != nil {
		return err
	}
	_, err = db.Conn.Exec("INSERT INTO articles (category,article) VALUES ($1,$2)", category, tmp)
	if err != nil {
		log.Printf("Error Inserting Articles err %s", err)
		return err
	}

	return nil

}

// Gets Multiple Articles with Cursor (id in the DB) and limit positive int
func (db *ArticleRepo) PaginationArticles(category string, cursor string, limit string) (map[string]interface{}, error) {
	articles := make(map[string]interface{})
	rows, err := db.Conn.Query("SELECT id,article FROM articles WHERE category = $1 and id > $2 ORDER BY id ASC LIMIT $3", category, cursor, limit)
	if err != nil {
		log.Printf("Error Reading Querying in Pagination Articles err %s", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id string
		var article []byte

		err := rows.Scan(&id, &article)
		if err != nil {
			log.Printf("Error Reading response  in Pagination Articles err %s", err)
			return articles, err
		}

		answer := issue(article)
		articles[id] = answer
	}
	return articles, nil
}

func (db *ArticleRepo) GetOneArticle(category string, id string) (article interface{}, err error) {
	var tmp []byte
	err = db.Conn.QueryRow("SELECT article FROM articles WHERE category = $1 and id = $2", category, id).Scan(&tmp)
	if err != nil {
		log.Printf("Error Getting one Article err %s", err)
		return nil, err
	}

	return issue(tmp), nil
}

func issue(value []byte) interface{} {
	var obj interface{}
	err := json.Unmarshal(value, &obj)
	if err != nil {
		log.Printf("Error Unmarshalling  one Article err %s", err)
	}
	return obj
}
