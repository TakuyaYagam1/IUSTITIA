package request

func ArchiveListParams(limit, offset *int) (int, int) {
	l := 50
	if limit != nil {
		l = *limit
	}
	o := 0
	if offset != nil {
		o = *offset
	}
	return l, o
}
