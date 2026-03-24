package utils

import (
	"MaximPLNV/json_repo/entities"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func NewJsonFileReader(fileName string) *JsonFileReader {
	fm := &JsonFileReader{}
	fm.fileName = fileName
	fm.fileAccessFlags = os.O_RDONLY
	return fm
}

type JsonFileReader struct {
	filterFn        func(*entities.BaseEntity) (bool, error)
	actionFn        func(*entities.BaseEntity)
	closeCn         chan struct{}
	fileName        string
	fileAccessFlags int
}

func (fm *JsonFileReader) ReadByLine() error {
	fm.closeCn = make(chan struct{}, 1)
	defer close(fm.closeCn)

	file, err := fm.openFile()
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	return fm.iterateFileLines(scanner)
}

func (fm *JsonFileReader) SetFilter(fn func(*entities.BaseEntity) (bool, error)) {
	fm.filterFn = fn
}

func (fm *JsonFileReader) SetAction(fn func(*entities.BaseEntity)) {
	fm.actionFn = fn
}

func (fm *JsonFileReader) StopReading() {
	fm.closeCn <- struct{}{}
}

func (fm *JsonFileReader) openFile() (*os.File, error) {
	file, err := os.OpenFile(fm.fileName, fm.fileAccessFlags, 0644)

	if err != nil {
		return nil, fmt.Errorf("There is no file with the name: %s", fm.fileName)
	}

	return file, nil
}

func (fm *JsonFileReader) iterateFileLines(sc *bufio.Scanner) error {
loop:
	for sc.Scan() {
		line := sc.Bytes()
		if err := fm.processLine(&line); err != nil {
			return err
		}

		select {
		case <-fm.closeCn:
			break loop
		default:
			continue loop
		}
	}

	return nil
}

func (fm *JsonFileReader) processLine(line *[]byte) error {
	line = fm.trimLine(line)
	if line == nil {
		return nil
	}

	e, err := fm.parseEntity(line)
	if err != nil {
		return err
	}

	return fm.executeLineRelatedLogic(e)
}

func (fm *JsonFileReader) trimLine(line *[]byte) *[]byte {
	start := bytes.IndexByte(*line, '{')
	end := bytes.LastIndexByte(*line, '}')

	if start == -1 || end == -1 {
		return nil
	}

	tLine := (*line)[start : end+1]
	return &tLine
}

func (fm *JsonFileReader) parseEntity(line *[]byte) (*entities.BaseEntity, error) {
	var e *entities.BaseEntity
	if err := json.Unmarshal(*line, &e); err != nil {
		return nil, fmt.Errorf("Line can't be parsed. Line: \"%s\"", *line)
	}

	return e, nil
}

func (fm *JsonFileReader) executeLineRelatedLogic(e *entities.BaseEntity) error {
	var err error
	isSel := true

	if fm.filterFn != nil {
		isSel, err = fm.filterFn(e)
	}

	if err == nil && isSel {
		fm.actionFn(e)
		return nil
	}

	return err
}
