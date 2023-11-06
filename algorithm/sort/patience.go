package sort

// https://zh.wikipedia.org/zh-tw/%E8%80%90%E5%BF%83%E6%8E%92%E5%BA%8F
func Patience(arr []int) []int {
	var cache [][]int

	for _, n := range arr {
		if len(cache) == 0 {
			cache = append(cache, []int{n})
			continue
		}

		//for li, l := range cache {
		//	for pi, p := range l {
		//		if pi == 0 && n > p {
		//			break
		//		}
		//
		//	}
		//}

	}
	return arr
}

func reverse(s []int) []int {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
