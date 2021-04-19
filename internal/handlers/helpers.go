package handlers

func removeDuplicateValues(objects []int64) []int64 {
	keys := make(map[int64]bool)
	list := []int64{}

	// If the key(values of the slice) is not equal to the already present value in new slice (list)
	// then we append it, else we jump on another element.
	for _, entry := range objects {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
