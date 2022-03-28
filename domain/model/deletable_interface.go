package model

type Deletable interface {
	Delete() error
}
