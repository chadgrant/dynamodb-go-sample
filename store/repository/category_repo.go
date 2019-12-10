package repository

type CategoryRepository interface {
	GetAll() ([]string, error)
	Upsert(category string) error
	Delete(category string) error
}
