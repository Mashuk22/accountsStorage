package store

type Store interface {
	Account() AccountRepository
}
