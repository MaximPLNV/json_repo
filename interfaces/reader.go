package interfaces

import "MaximPLNV/json_repo/entities"

type Reader interface {
	ReadByLine() error
	SetFilter(func(*entities.BaseEntity) (bool, error))
	SetAction(fn func(*entities.BaseEntity))
	StopReading()
}
