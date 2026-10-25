package api

func pagination(limit *int, offset *int) (int, int) {
	pageLimit := 50
	if limit != nil {
		pageLimit = *limit
	}

	var pageOffset int
	if offset != nil {
		pageOffset = *offset
	}

	return pageLimit, pageOffset
}
