package jsonrepo_test

import (
	"MaximPLNV/json_repo/entities"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var TEST_FILE_NAME string = "test_file.json"

func GenerateValidJson(recsCount int, fName string, t *testing.T) string {
	content := genereteJsonEntities(recsCount, t)
	return writeFile(*content, fName, t)
}

func GenerateInvalidJson(recsCount int, fName string, t *testing.T) string {
	content := generateBadEntities(recsCount)
	return writeFile(*content, fName, t)
}

func genereteJsonEntities(recsCount int, t *testing.T) *[]byte {
	content := []byte("[\n")
	ents := genereteEntities(recsCount)

	for i, e := range *ents {
		bEnt, err := json.Marshal(e)

		if err != nil {
			t.Fatal(err)
		}

		if i != recsCount-1 {
			bEnt = append(bEnt, []byte(",\n")...)
		} else {
			bEnt = append(bEnt, []byte("\n")...)
		}

		content = append(content, bEnt...)
	}

	content = append(content, ']')
	return &content
}

func genereteEntities(recsCount int) *[]entities.BaseEntity {
	ents := make([]entities.BaseEntity, 0)
	for i := range recsCount {
		ent := entities.BaseEntity{}
		ent.Id = i
		ent.CreatedAt = time.Now()
		ent.UpdatedAt = time.Now()

		ents = append(ents, ent)
	}

	return &ents
}

func writeFile(content []byte, fName string, t *testing.T) string {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, fName)

	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatal(err)
	}

	return filePath
}

func generateBadEntities(recsCount int) *[]byte {
	content := []byte("[\n")

	for i := range recsCount {
		line := fmt.Sprintf("{\"id\": %d, \"created_at\": \"null\", \"updated_at\": \"null\"}", i)
		bLine := []byte(line)

		if i != recsCount-1 {
			bLine = append(bLine, []byte(",\n")...)
		} else {
			bLine = append(bLine, []byte("\n")...)
		}

		content = append(content, bLine...)
	}

	content = append(content, ']')
	return &content
}
