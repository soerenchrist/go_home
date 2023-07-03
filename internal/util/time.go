package util

import "time"

func GetTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

func ValidateTimestamp(timestamp string) error {
	_, err := time.Parse(time.RFC3339, timestamp)
	return err
}
