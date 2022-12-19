package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

// readJSON tries to read the body of a request and converts it into JSON
func (app *settings) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	// 5.9MiB
	maxBytes := 6206016
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	//--:REX you changed `err := dec.Decode(data)` -> `err := dec.Decode(&data)`
	err := dec.Decode(data)
	if err != nil {
		app.errorLog.Println("BRUHhhhhhhhhh")
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("error parsing json")
	}
	return nil
}

// writeJSON takes a response status code and aribitrary data and writes a json response to the client
func (app *settings) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	var out []byte
	if app.environment == "dev" {
		output, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			return err
		}
		out = output
	} else {
		output, err := json.Marshal(data)
		if err != nil {
			return err
		}
		out = output
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// errorJSON takes an error, and optionally a response status code, and generates and sends
// a json error response
func (app *settings) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	app.writeJSON(w, statusCode, payload)
	return nil
}

func (app *settings) GetRuntimeParams(path string) (RuntimeParms, error) {
	var config RuntimeParms
	contents, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	// this is where we unmarshal the contents into config
	err = json.Unmarshal(contents, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
