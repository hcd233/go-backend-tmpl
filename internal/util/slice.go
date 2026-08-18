package util

// MapErr 将切片元素逐一映射，任一元素映射失败时立即返回 nil 与该错误。
func MapErr[T any, R any](collection []T, iteratee func(item T, index int) (R, error)) ([]R, error) {
	result := make([]R, 0, len(collection))
	for i, item := range collection {
		r, err := iteratee(item, i)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	return result, nil
}
