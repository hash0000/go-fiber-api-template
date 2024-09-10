package helpers

import "math"

func CalculateTotalPagesPaginationHelper(totalCount int64, limit int16) int {
	if limit == 0 {
		return 0
	}
	return int(math.Ceil(float64(totalCount) / float64(limit)))
}
