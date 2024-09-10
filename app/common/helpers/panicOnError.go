package helpers

import (
	"database/sql"
	"errors"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/lib/pq"
	"go-fiber-api-template/app/common/constants"
	"log/slog"
)

func CheckSqlError(err error) (constants.SqlQueryStatusType, error) {
	if err != nil {
		if err == qrm.ErrNoRows || err == sql.ErrNoRows {
			return constants.NotFound, nil
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return constants.NotUnique, nil
			}
		}

		slog.Error("Unknown error while SQL querying", slog.String("location", GetFileLine(constants.DeepCallerConstant.Common)), slog.Any("info", err))
		return constants.Unknown, err
	}

	return constants.Success, nil
}
