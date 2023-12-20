package filesystem

import "os"

func AbsPathToLocalStorage(partToConcat string) string {
	localStoragePath := os.Getenv("LOCAL_STORAGE_PATH")
	return localStoragePath + "/" + partToConcat
}
