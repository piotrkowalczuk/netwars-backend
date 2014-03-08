package database

type RepositoryInterface interface {
	FindOne() interface{}
	Find() []interface{}
}
