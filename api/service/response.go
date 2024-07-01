package service

import "net/http"

func ServerResponse(w http.ResponseWriter, dataInterface interface{}) {
	bytes, err := SerializeJSON(dataInterface)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
