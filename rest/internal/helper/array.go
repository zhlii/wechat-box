package helper

func RemoveElement(arr []string, target string) []string {
	var result []string

	for _, value := range arr {
		if value != target {
			result = append(result, value)
		}
	}

	return result
}

func ContainsElement(arr []string, element string) bool {
	for _, v := range arr {
		if v == element {
			return true
		}
	}
	return false
}
