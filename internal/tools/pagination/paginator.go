package pagination

func Pagination(page, pageSize int) (limit, offset int) {
	if pageSize < 0 || pageSize > 100 {
		pageSize = 10
	}
	if page <= 0 {
		return pageSize, 0
	} else {
		return pageSize, (page - 1) * pageSize
	}
}
