package main

import (
	"encoding/json"
	"maps"
	"net/http"
	"strconv"
	"errors"

	"github.com/gorilla/mux"
)

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

type envelope map[string]any

func (app *application) writeJSON(
	w http.ResponseWriter,
	status int,
	data envelope,
	headers http.Header,
) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// curl -X POST "https://api.groq.com/openai/v1/chat/completions" -H "Content-Type: application/json" -H "Authorization: Bearer $API_KEY_GROQ_IDEAGEN" -d '{"model": "openai/gpt-oss-20b", "messages": [{"role": "user", "content": "what year now?"}]}'
