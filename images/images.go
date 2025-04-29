package images

import "os"

func CheckFolder(filename string) {
	file, err := os.Stat(filename)
	if os.IsNotExist(err) || !file.IsDir() {
		CreateFolder(filename)
	}
}

func CreateFolder(fileName string) {
	err := os.Mkdir(fileName, 0755)
	if err != nil {
		return
	}
}
