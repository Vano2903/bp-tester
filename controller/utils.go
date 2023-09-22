package controller

import (
	"math/rand"
	"os"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func writeToFile(path string, content []byte) error {
	dockerfileWriter, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer dockerfileWriter.Close()

	if _, err = dockerfileWriter.Write(content); err != nil {
		return err
	}
	return nil
}
