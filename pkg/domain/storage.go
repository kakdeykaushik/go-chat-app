package domain

type Storage interface {
	Get(K string) (any, error)
	List() ([]any, error)
	Save(K, V any) error
	Delete(K string) error
}
