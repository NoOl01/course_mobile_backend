package images

import (
	"fmt"
	"os"
	"path/filepath"
)

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

func FindAvatar(id int64) string {
	filename := fmt.Sprintf("avatar_%d.png", id)
	avatarPath := filepath.Join("./upload/avatars", filename)

	if _, err := os.Stat(avatarPath); err == nil {
		// Возвращаем только имя файла
		return filename
	}
	return ""
}
