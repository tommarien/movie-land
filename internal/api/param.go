package api

import (
	"net/http"
	"strconv"
)

func getIntParam(r *http.Request, name string) (int, error) {
	param := r.PathValue(name)

	value, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}

	return value, nil
}
