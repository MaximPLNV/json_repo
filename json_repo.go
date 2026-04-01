package jsonrepo

import (
	"MaximPLNV/json_repo/entities"
	"MaximPLNV/json_repo/interfaces"
	"MaximPLNV/json_repo/utils"
	"encoding/json"
	"fmt"
	"slices"
	"time"
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
	idCounter := 0
	insertedCount := 0
	esLen := len(*es)

	f := func(e *entities.BaseEntity) (bool, error) {
		if e.Id != idCounter && esLen > insertedCount {
			return true, nil
		}
		idCounter++
		return false, nil
	}

	act := func(e *entities.BaseEntity) (*[]byte, error) {
		result := make([]byte, 0)
		stopId := e.Id
		for i := idCounter; i < stopId && insertedCount < esLen; i++ {
			recToInsert := (*es)[insertedCount]
			recToInsert.Id = idCounter
			recToInsert.CreatedAt = time.Now()
			recToInsert.UpdatedAt = time.Now()
			bRecToInsert, err := json.Marshal(recToInsert)
			if err != nil {
				return nil, err
			}
			bRecToInsert = append(bRecToInsert, []byte(",\n")...)
			result = append(result, bRecToInsert...)
			idCounter++
			insertedCount++
		}

		bEntity, err := json.Marshal(e)
		if err != nil {
			return nil, err
		}
		result = append(result, bEntity...)
		idCounter++
		return &result, nil
	}

	postActFn := func() (*[]byte, error) {
		result := make([]byte, 0)

		if esLen == insertedCount {
			return &result, nil
		}

		for i := idCounter; insertedCount < esLen; i++ {
			result = append(result, []byte(",\n")...)
			recToInsert := (*es)[insertedCount]
			recToInsert.Id = idCounter
			recToInsert.CreatedAt = time.Now()
			recToInsert.UpdatedAt = time.Now()
			bRecToInsert, err := json.Marshal(recToInsert)
			if err != nil {
				return nil, err
			}
			result = append(result, bRecToInsert...)
			idCounter++
			insertedCount++
		}
		return &result, nil
	}

	jr.writer.SetFilter(f)
	jr.writer.SetAction(act)
	jr.writer.SetPostAction(postActFn)

	if err := jr.writer.WriteByLine(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (jr *JsonRepo) Update(es *[]entities.BaseEntity) error {
	ids := jr.getEntitiesIds(es)

	f := func(e *entities.BaseEntity) (bool, error) {
		if slices.Contains(*ids, e.Id) {
			return true, nil
		}

		return false, nil
	}

	act := func(e *entities.BaseEntity) (*[]byte, error) {
		//TODO
		return nil, nil
	}

	jr.writer.SetFilter(f)
	jr.writer.SetAction(act)

	if err := jr.writer.WriteByLine(); err != nil {
		return err
	}

	return nil
}

func (jr *JsonRepo) Delete(ids *[]int) error {

	f := func(e *entities.BaseEntity) (bool, error) {
		if slices.Contains(*ids, e.Id) {
			return true, nil
		}

		return false, nil
	}

	act := func(e *entities.BaseEntity) (*[]byte, error) {
		result := make()
		return nil, nil
	}

	jr.writer.SetFilter(f)
	jr.writer.SetAction(act)

	if err := jr.writer.WriteByLine(); err != nil {
		return err
	}

	return nil
}

func (jr *JsonRepo) getEntitiesIds(es *[]entities.BaseEntity) *[]int {
	ids := make([]int, len(*es))

	for i, e := range *es {
		ids[i] = e.Id
	}

	return &ids
}
