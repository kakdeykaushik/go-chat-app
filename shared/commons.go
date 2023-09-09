package shared

import "log"

func RemoveIndex[T any](s []T, i int) []T {
	return append(s[:i], s[i+1:]...)
}

func HandleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s - %v\n", msg, err)
	}
}
