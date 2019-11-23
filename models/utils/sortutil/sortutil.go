package sortutil

//
type Int64Slice []int64

func (slice Int64Slice) Len() int {
	return len(slice)
}

func (slice Int64Slice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (slice Int64Slice) Less(i, j int) bool {
	return slice[i] < slice[j]
}
