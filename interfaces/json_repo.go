package interfaces

import "MaximPLNV/json_repo/entities"

type JsonRepo interface {
	GetByIds(*[]int) (*[]entities.BaseEntity, error)
	GetByFilter(func(*entities.BaseEntity) (bool, error)) (*[]entities.BaseEntity, error)
	GetAll() (*[]entities.BaseEntity, error)
	Add(*[]entities.BaseEntity) (*[]entities.BaseEntity, error)
	Update(*[]entities.BaseEntity) error
	Delete(*[]int) error
}
