package stream

func ToListSelf[R any](s Stream) []R {
	return ToList[R](s, nil)
}

func ToList[R any](s Stream, fn func(s any) R) []R {
	var result []R
	for item := range s.source {
		if fn == nil {
			result = append(result, item.(R))
		} else {
			result = append(result, fn(item))
		}
	}
	return result
}
