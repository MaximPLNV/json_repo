package utils

import (
	"MaximPLNV/json_repo/entities"
	"os"
)

func NewJsonFileWriter(f string) *JsonFileWriter {
	fm := &JsonFileWriter{}
	fm.fileName = f
	fm.fileAccessFlags = os.O_RDONLY
	return fm
}

type JsonFileWriter struct {
	filterFn        func(*entities.BaseEntity) (bool, error)
	actionFn        func(*entities.BaseEntity)
	closeCn         chan bool
	fileName        string
	fileAccessFlags int
}

func (fm *JsonFileWriter) WriteByLine() {
	// open file
	// read by line
	// skip if start or end
	// parse entity
	// execute filter logic
	// if true execute actionFn
	// if close chan receive signal then break loop and close
}
