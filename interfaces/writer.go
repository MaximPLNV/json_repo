package interfaces

import "MaximPLNV/json_repo/entities"

type Writer interface {
	SetFilter(func(*entities.BaseEntity) (bool, error))
	SetAction(func(*entities.BaseEntity) (*[]byte, error))
	StopReading()
	WriteByLine() error
}
