package response

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

var EmptyResponse struct{}

type ErrorPayload struct {
	Traceback string `json:"traceback,omitempty"`
	Message   string `json:"message,omitempty"`
}

func Response(w http.ResponseWriter, r *http.Request, httpStatusHeader int, data interface{}, errors interface{}, meta interface{}, code int) {
	apiResponse := struct {
		Data   interface{} `json:"data,omitempty"`
		Errors interface{} `json:"errors,omitempty"`
		Meta   interface{} `json:"meta,omitempty"`
	}{
		data,
		errors,
		meta,
	}
	if code >= http.StatusOK && code <= 299 {
		logrus.WithFields(logrus.Fields{
			"request_uri": r.URL.Path,
			"response":    apiResponse,
			"status_code": httpStatusHeader,
		}).Infoln("API Response Success")
	} else {
		logrus.WithFields(logrus.Fields{
			"request_uri": r.URL.Path,
			"response":    apiResponse,
			"status_code": httpStatusHeader,
		}).Errorln("API Response Error")
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusHeader)

	_ = json.NewEncoder(w).Encode(apiResponse)
}
