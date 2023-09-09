package domain

type Storage interface {
	Get(K string, k any) (any, error)
	List() ([]any, error)
	Save(V any) error
	Delete(K string) error
}
