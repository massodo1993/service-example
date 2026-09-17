package health

import "net/http"

func Handler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		return
	}
}
