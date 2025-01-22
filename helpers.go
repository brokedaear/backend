package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// jsonWrap wraps a json message response before it gets sent out. This makes
// it easier for the requester to read the JSON data they get back. You can
// imagine that this "envelops" the json data that will go into it.
// TODO: Add example
type jsonWrap map[string]any

// readJSON reads JSON from a request, and places the result in a dst of type
// any. This MUST be a reference and not a pass by value. Furthermore, this
// type will be of map[string]any. The error handling logic here is a more
// effective way to "triage" errors ( view documentation for a link to the
// definition of triage) and issues/that/arise/from decoding JSON data.
// This is a variation of the readJSON function in the book
// "Let's Go Further" by Alex Edwards,
// page 90. What's key is that the method call must use a reference for
// dst, since that is typically what is consumed in a JSON decode.
func (app *app) readJSON(
	w http.ResponseWriter,
	r *http.Request,
	dst any,
) error {
	// Maximum amount of bytes our JSON reader will accept, which is 1MB.
	maxBytes := 1_048_576

	// Limit the size of the incoming request.
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)

	// Disallows unknown fields in JSON body. Once detected, this straight-up
	// returns an error.
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)

	// Here, errors are sorted out, in other words triaged. This makes the
	// error returns more readable for the API consumer, which is important
	// when we want to debug.
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf(
				"body contains badly-formed JSON (at character %d)",
				syntaxError.Offset,
			)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf(
					"body contains incorrect JSON type for field %q",
					unmarshalTypeError.Field,
				)
			}

			return fmt.Errorf(
				"body contains incorrect JSON type (at character %d)",
				unmarshalTypeError.Offset,
			)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// In the event a field cannot be mapped to a destination key,
		// an error occurs in the form of "json: unknown field".
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf(
				"body must not be larger than %d bytes",
				maxBytesError.Limit,
			)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// writeJSON ingests a data, map of map[string]any, and writes it to a
// http.ResponseWriter stream. It returns an error in case
// there was one writing to the stream.
func (app *app) writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
	headers http.Header,
) error {

	js, err := jsonBuilder(data)
	if err != nil {
		return err
	}

	// Add any headers
	for k, v := range headers {
		w.Header()[k] = v
	}

	// Add headers, then write to the output stream.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)

	if err != nil {
		return err
	}

	return nil
}

// jsonBuilder builds a JSON that can then be written to
// a http.ResponseWriter stream. The parameter "data", is a
// map[string]any
func jsonBuilder(data any) ([]byte, error) {
	js, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Convenience \n for terminal view
	js = append(js, '\n')

	return js, nil
}
