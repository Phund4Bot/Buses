package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ParseFromFile[T any](name string) (T, error) {
	var obj T

	file, err := os.Open(name)
	if err != nil {
		return obj, fmt.Errorf("cannot open file: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return obj, fmt.Errorf("cannot read file: %v", err)
	}

	err = json.Unmarshal(data, &obj)
	if err != nil {
		return obj, fmt.Errorf("cannot unmarshall data: %v", err)
	}

	return obj, nil
}

func CreateFileConnection(name string) (*os.File, error) {
	var (
		file *os.File
		err  error
	)

	file, err = os.Open(name)
	if err == nil {
		err = os.Remove(name)
		if err != nil {
			return nil, fmt.Errorf("cannot clear old logs: %v", err)
		}
	}

	file, err = os.Create(name)
	if err != nil {
		return nil, fmt.Errorf("cannot open or create file: %v", err)
	}
	return file, nil
}
