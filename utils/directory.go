package utils

import (
	"bufio"
	"errors"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreatePath(path string) (bool, error) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return false, err
	}
	return true, nil
}

func CreateFile(path, filename string) (bool, error) {
	f, err := os.Create(path + filename)
	if err != nil {
		return false, err
	}
	defer f.Close()
	return true, nil
}

func WriteFile(path, filename, content string) error {
	s := []byte(content)
	f, err := os.OpenFile(path+filename, os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(s)
	if err != nil {
		return err
	}
	return nil
}

func ReadFile(path, filename string) (string, error) {
	file, err := os.Open(path + filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "" {
			return scanner.Text(), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", errors.New("content not found")
}
