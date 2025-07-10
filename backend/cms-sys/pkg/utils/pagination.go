package utils

func CalculateOffset(page, limit int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * limit
}
