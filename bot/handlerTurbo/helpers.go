package handlerTurbo

import "fmt"

func makeLocalKey(key string, userID int) string {
	return fmt.Sprintf(key, userID)
}
