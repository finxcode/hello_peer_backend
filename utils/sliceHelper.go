package utils

func Contains(ids []int, id int) bool {
	if len(ids) == 0 {
		return false
	}

	for _, value := range ids {
		if value == id {
			return true
		}
	}
	return false
}
