package jsonrepo

import (
	"MaximPLNV/json_repo/entities"
	"MaximPLNV/json_repo/interfaces"
	"MaximPLNV/json_repo/utils"
	"fmt"
	"slices"
)

func NewJsonRepo(fName string) *JsonRepo {
	r := &JsonRepo{}
	r.reader = utils.NewJsonFileReader(fName)
	r.writer = utils.NewJsonFileWriter(fName)
	return r
}

type JsonRepo struct {
	reader interfaces.Reader
	writer interfaces.Writer
}

func (jr *JsonRepo) GetByIds(ids *[]int) (*[]entities.BaseEntity, error) {
	if ids == nil || len(*ids) == 0 {
		return nil, fmt.Errorf("Incorrect input parameter <ids>")
	}

	f := func(e *entities.BaseEntity) (bool, error) {
		return slices.Contains(*ids, e.Id), nil
	}

	return jr.GetByFilter(f, len(*ids))
}

func (jr *JsonRepo) GetByFilter(filterFn func(*entities.BaseEntity) (bool, error), l int) (*[]entities.BaseEntity, error) {
	if l <= 0 && l != -1 {
		return nil, fmt.Errorf("Limit value should be more than 0 or equal to -1 [l: %d]", l)
	}

	var ents []entities.BaseEntity

	actFn := func(e *entities.BaseEntity) {
		ents = append(ents, *e)

		if l != -1 && len(ents) >= l {
			jr.reader.StopReading()
		}
	}

	jr.reader.SetFilter(filterFn)
	jr.reader.SetAction(actFn)

	if err := jr.reader.ReadByLine(); err != nil {
		return nil, err
	}

	return &ents, nil
}

func (jr *JsonRepo) GetAll() (*[]entities.BaseEntity, error) {
	return jr.GetByFilter(nil, -1)
}

func (jr *JsonRepo) Add(es *[]entities.BaseEntity) (*[]entities.BaseEntity, error) {
	// Reader to read file and return line by line
	// Writer to write lines to temp file
	// If there is free slot between records then add new task
	// If no write already existing line
	// If file is ended and there are still tasks to inser then add at the end
	// Close writed file and rename it
	return nil, nil
}

func (jr *JsonRepo) Update(es *[]entities.BaseEntity) error {
	// Reader to read file and return line by line
	// Writer to write lines to temp file
	// If there is task with correct id then update task
	// If no write already existing line
	// If file is ended and there are still tasks to update then show error
	// Close writed file and rename it
	return nil
}

func (jr *JsonRepo) Delete(ids *[]int) error {
	// Reader to read file and return line by line
	// Writer to write lines to temp file
	// If there is task with correct id then skip it
	// If no write already existing line
	// If file is ended and there are still tasks to delete then show error
	// Close writed file and rename it
	return nil
}
