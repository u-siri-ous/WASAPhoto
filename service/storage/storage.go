package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
)

func SavePhoto(file multipart.File, photoId uint64, userId uint64) error {
	_, err := os.Stat(filepath.Join(BasePath, "/", strconv.FormatUint(userId, 10), UploadedPhotoFolder))
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Join(BasePath, "/", strconv.FormatUint(userId, 10), UploadedPhotoFolder), 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	out, err := os.Create(filepath.Join(BasePath, "/", strconv.FormatUint(userId, 10), UploadedPhotoFolder, fmt.Sprintf("%d.jpg", photoId)))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}

	return nil
}

func DeletePhoto(currentUserId uint64, photoId uint64) error {
	err := os.Remove(filepath.Join(BasePath, "/", strconv.FormatUint(currentUserId, 10), UploadedPhotoFolder, fmt.Sprintf("%d.jpg", photoId)))
	return err
}
