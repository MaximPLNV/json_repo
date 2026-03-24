package utils

import (
	"MaximPLNV/json_repo/entities"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	actionFn        func(*entities.BaseEntity) (*[]byte, error)
	closeCn         chan struct{}
	fileName        string
	tempFilePath    string
	tempFilePattern string
	fileAccessFlags int
}

func (fm *JsonFileWriter) SetFilter(fn func(*entities.BaseEntity) (bool, error)) {
	fm.filterFn = fn
}

func (fm *JsonFileWriter) SetAction(fn func(*entities.BaseEntity) (*[]byte, error)) {
	fm.actionFn = fn
}

func (fm *JsonFileWriter) StopReading() {
	fm.closeCn <- struct{}{}
}

func (fm *JsonFileWriter) WriteByLine() error {
	file, fErr := fm.openFile()
	tmp, tErr := fm.createTempFile()

	if fErr != nil {
		return fErr
	} else if tErr != nil {
		return tErr
	}

	defer tmp.Close()
	defer os.Remove(tmp.Name())

	fm.closeCn = make(chan struct{}, 1)
	defer close(fm.closeCn)

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(tmp)

	return fm.processLineByLine(scanner, writer)
}

func (fm *JsonFileWriter) openFile() (*os.File, error) {
	file, openErr := os.OpenFile(fm.fileName, fm.fileAccessFlags, 0644)

	if openErr != nil {
		return nil, fmt.Errorf("File \"%s\" doesn't exist", fm.fileName)
	}

	return file, nil
}

func (fm *JsonFileWriter) createTempFile() (*os.File, error) {
	tmp, tmpErr := os.CreateTemp("", fm.tempFilePattern)

	if tmpErr != nil {
		return nil, errors.New("File modification issues")
	}

	fm.tempFilePath = tmp.Name()
	return tmp, nil
}

func (fm *JsonFileWriter) processLineByLine(sc *bufio.Scanner, wr *bufio.Writer) error {
loop:
	for sc.Scan() {
		line := sc.Bytes()
		if err := fm.processLine(&line, wr); err != nil {
			return err
		}

		select {
		case <-fm.closeCn:
			break loop
		default:
			continue loop
		}
	}
	return fm.saveChanges(wr)
}

func (fm *JsonFileWriter) processLine(line *[]byte, wr *bufio.Writer) error {
	start := bytes.IndexByte(*line, '{')
	end := bytes.LastIndexByte(*line, '}')

	if start == -1 || end == -1 || start > end {
		return errors.New("File reading issue")
	}

	begining := (*line)[:start]
	data := (*line)[start : end+1]
	ending := (*line)[end+1:]

	processed, err := fm.processCrearData(&data)
	if err != nil {
		return err
	}

	result := append(begining, data...)
	result = append(*processed, ending...)
	wr.Write(result)
	return nil
}

func (fm *JsonFileWriter) processCrearData(line *[]byte) (*[]byte, error) {
	var e entities.BaseEntity
	if err := json.Unmarshal(*line, &e); err != nil {
		return nil, fmt.Errorf("Line can't be parsed. Line: \"%s\"", *line)
	}

	selected, sErr := fm.filterFn(&e)
	if sErr != nil {
		return nil, sErr
	} else if selected {
		return fm.actionFn(&e)
	}

	return line, nil
}

func (fm *JsonFileWriter) saveChanges(wr *bufio.Writer) error {
	wr.Flush()

	if err := os.Rename(fm.tempFilePath, fm.fileName); err != nil {
		return errors.New("File modification issues")
	}
	return nil
}
