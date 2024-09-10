package helpers

import (
	"go-fiber-api-template/app/common/constants"
	"log/slog"
	"os"
)

func IsExistsByPathHelper(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		slog.Info("No such file in dir", slog.String("location", GetFileLine(constants.DeepCallerConstant.Common)), slog.Any("info", err))
		return false, nil
	}
	return false, err
}
