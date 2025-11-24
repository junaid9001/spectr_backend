package utils

import "strconv"

//string to uint
func StringToUint(idString string) (uint, error) {

	id, err := strconv.ParseUint(idString, 10, 64)

	if err != nil {
		return 0, err
	}
	return uint(id), nil

}
