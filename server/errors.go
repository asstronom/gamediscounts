package server

func e ( e error) map[string]string{
	return map[string]string{
		"error":e.Error(),
	}
}
