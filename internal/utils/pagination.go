package utils

func PaginationCalculation(page int, limit int, totalData int) (int, interface{}) {

	metaData := make(map[string]int)

	offset := (page - 1) * limit

	totalPages := int(totalData) / limit
	if int(totalData)%limit > 0 {
		totalPages++
	}

	metaData["total_data"] = totalData
	metaData["total_pages"] = totalPages

	return offset, metaData
}