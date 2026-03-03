package jsonrepo_test

import (
	jsonrepo "MaximPLNV/json_repo"
	"MaximPLNV/json_repo/entities"
	"slices"
	"testing"
	"time"
)

func TestGetIds(t *testing.T) {
	var tests = []struct {
		ids *[]int
	}{
		{&[]int{3}},
		{&[]int{1, 3}},
	}

	for _, param := range tests {
		repo := jsonrepo.NewJsonRepo("test_file.json")
		ents, err := repo.GetByIds(param.ids)

		if err != nil {
			t.Error(err)
			return
		}

		if ents == nil || len(*ents) != len(*param.ids) {
			t.Error("Error during entity receiving")
			return
		}

		for _, ent := range *ents {
			if f := slices.Contains(*param.ids, ent.Id); !f {
				t.Error("Incorrect entity has been received")
			}
		}
	}
}

func TestGetBadIds(t *testing.T) {
	repo := jsonrepo.NewJsonRepo("test_file.json")
	_, err := repo.GetByIds(nil)

	if err == nil {
		t.Error("Error should be displayed for incorrect ids parameter")
	}
}

func TestGetAll(t *testing.T) {
	repo := jsonrepo.NewJsonRepo("test_file.json")
	ents, err := repo.GetAll()

	if err != nil {
		t.Error(err)
		return
	}

	if ents == nil || len(*ents) != 4 {
		t.Error("Error during entity receiving")
	}
}

func TestGetByFn(t *testing.T) {
	fn := func(e *entities.BaseEntity) (bool, error) {
		filterTime := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
		return filterTime.Before(e.UpdatedAt), nil
	}

	repo := jsonrepo.NewJsonRepo("test_file.json")
	ents, err := repo.GetByFilter(fn, -1)

	if err != nil {
		t.Error(err)
		return
	}

	if ents == nil || len(*ents) != 2 {
		t.Error("Error during entity receiving")
		return
	}

	for _, ent := range *ents {
		if f := slices.Contains([]int{1, 2}, ent.Id); !f {
			t.Error("Incorrect entity has been received")
		}
	}
}

func TestGetByFnBadLimit(t *testing.T) {
	repo := jsonrepo.NewJsonRepo("test_file.json")
	_, err := repo.GetByFilter(nil, 0)

	if err == nil {
		t.Error("Limit validation error should be received")
		return
	}
}
