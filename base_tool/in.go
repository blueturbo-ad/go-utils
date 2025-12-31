package basetool

// 基于基础类型的in操作
func In[T string | int | int64 | int32 | float64 | float32 | uint64 | uint32 | uint16, A []T](target T, arr A) bool {
	for _, a := range arr {
		if a == target {
			return true
		}
	}
	return false
}