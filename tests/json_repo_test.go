package jsonrepo_test

import (
	jsonrepo "MaximPLNV/json_repo"
	"MaximPLNV/json_repo/entities"
	"fmt"
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
		filepath := GenerateValidJson(5, TEST_FILE_NAME, t)
		repo := jsonrepo.NewJsonRepo(filepath)
		ents, err := repo.GetByIds(param.ids)

		if err != nil {
			t.Fatal(err)
		}

		if len(*ents) != len(*param.ids) {
			t.Fatal("Error during entity receiving")
		}

		for _, ent := range *ents {
			if f := slices.Contains(*param.ids, ent.Id); !f {
				t.Fatal("Incorrect entity has been received")
			}
		}
	}
}

func TestGetBadIds(t *testing.T) {
	filepath := GenerateValidJson(5, TEST_FILE_NAME, t)
	repo := jsonrepo.NewJsonRepo(filepath)
	_, err := repo.GetByIds(nil)

	if err == nil {
		t.Fatal("Error should be displayed for incorrect ids parameter")
	}
}

func TestGetAll(t *testing.T) {
	filepath := GenerateValidJson(5, TEST_FILE_NAME, t)
	repo := jsonrepo.NewJsonRepo(filepath)
	ents, err := repo.GetAll()

	if err != nil {
		t.Fatal(err)
	}

	if len(*ents) != 5 {
		t.Fatal("Error during entity receiving")
	}
}

func TestGetByFn(t *testing.T) {
	filepath := GenerateValidJson(5, TEST_FILE_NAME, t)

	fn := func(e *entities.BaseEntity) (bool, error) {
		filterTime := time.Now().Add(-time.Hour)
		return filterTime.Before(e.UpdatedAt), nil
	}

	repo := jsonrepo.NewJsonRepo(filepath)
	ents, err := repo.GetByFilter(fn, -1)

	if err != nil {
		t.Fatal(err)
	}

	if len(*ents) != 5 {
		t.Fatal("Error during entity receiving")
	}
}

func TestGetByFnBadLimit(t *testing.T) {
	filepath := GenerateValidJson(5, TEST_FILE_NAME, t)
	repo := jsonrepo.NewJsonRepo(filepath)
	_, err := repo.GetByFilter(nil, 0)

	if err == nil {
		t.Fatal("Limit validation error should be received")
	}
}

func TestOpenUnexistingFile(t *testing.T) {
	fName := fmt.Sprintf("unexisting-%s.json", time.Now())
	repo := jsonrepo.NewJsonRepo(fName)
	ents, err := repo.GetAll()

	if err == nil || ents != nil {
		t.Fatal("Error should be displayed as file doesn't exist")
	}
}

func TestBadJsonParsing(t *testing.T) {
	fName := GenerateInvalidJson(5, TEST_FILE_NAME, t)
	repo := jsonrepo.NewJsonRepo(fName)
	ents, err := repo.GetAll()

	if err == nil || ents != nil {
		t.Fatal("Error should be displayed as file has unexpected stracture")
	}
}
