package utils

// GetOffsetLimit get offset and limit
func GetOffsetLimit(page int64, size int64) (int64, int64) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	var limit int64 = 1
	if size > 1 {
		limit = size
	}

	offset := (page - 1) * limit

	return offset, limit
}
