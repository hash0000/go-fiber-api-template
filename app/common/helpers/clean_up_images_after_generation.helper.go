package helpers

import (
	"fmt"
	"go-fiber-api-template/app/common/constants"
	"log/slog"
	"os"
	"path/filepath"
)

func CleanUpImagesAfterGenerationHelper(chapterRange int, sessionID string) {
	imageArray := make([]string, chapterRange)

	for i := range chapterRange {
		imageArray[i] = filepath.Join(constants.StoreDir, sessionID, fmt.Sprintf("image_%d.jpg", i))
	}

	for _, fileName := range imageArray {
		err := os.Remove(fileName)
		if err != nil {
			slog.Error("Failed to delete file", slog.String("file_name", fileName), slog.String("location", GetFileLine(constants.DeepCallerConstant.Modules)), slog.Any("info", err))
		} else {
			slog.Debug(fmt.Sprintf("Successfully deleted %s", fileName))
		}
	}
}
