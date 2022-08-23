package models

type ArticlesRepo interface {
	AddArticles(category string, article interface{}) (err error)
	PaginationArticles(category string, cursor string, limit string) (map[string]interface{}, error)
	GetOneArticle(category string, id string) (article interface{}, err error)
}
