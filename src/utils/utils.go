package utils

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"strconv"
	"time"
)

func GenerateFileName() string {
	fileTime := time.Now()
	fileName := fmt.Sprintf(
		"%d/%d/%d/%s.jsonl",
		fileTime.Year(),
		fileTime.Month(),
		fileTime.Day(),
		uuid.New(),
	)
	return fileName
}


func GetEnvOrString(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func GetEnvOrInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		val, err := strconv.Atoi(value)
		PanicOnError(err, fmt.Sprintf("unable to convert value for env var %s to int", key))
		return val
	}
	return defaultVal
}

func GetEnvOrPanic(key string) string {
	value, exists := os.LookupEnv(key)
	if exists == false {
		err := errors.New(fmt.Sprintf("Key does not exist: %s", key))
		PanicOnError(err, "Unable to get environment variable")
	}
	return value
}

func PanicOnError(err error, message string) {
	if err != nil {
		log.Panicf("%s: %s", message, err)
	}
}

func ContinueOnError(err error, message string) {
	if err != nil {
		log.Printf("%s: %s", message, err)
	}
}