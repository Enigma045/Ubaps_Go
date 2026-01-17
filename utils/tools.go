package utils

import (
	"fmt"
	"net/http"
)

func Formdata(r *http.Request) (map[string][]string, error) {
	if r == nil {
		return nil, fmt.Errorf("request is nil")
	}

	// Limit request body size (optional, pass w if you want MaxBytesReader)
	// r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB

	formData := make(map[string][]string)

	// Check if multipart/form-data
	if r.Header.Get("Content-Type") != "" &&
		r.Header.Get("Content-Type")[:19] == "multipart/form-data" {
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
			return nil, fmt.Errorf("failed to parse multipart form: %w", err)
		}
		for key, values := range r.MultipartForm.Value {
			if len(values) == 0 {
				continue
			}
			formData[key] = values
		}
	} else {
		// normal form
		if err := r.ParseForm(); err != nil {
			return nil, fmt.Errorf("failed to parse form: %w", err)
		}
		for key, values := range r.Form {
			if len(values) == 0 {
				continue
			}
			formData[key] = values
		}
	}

	return formData, nil
}

// helper to safely get first value from form map
func getFormValue(formData map[string][]string, key string) (string, error) {
	values, ok := formData[key]
	if !ok || len(values) == 0 {
		return "", fmt.Errorf("missing or empty field: %s", key)
	}
	return values[0], nil
}

// func BuildInsertFromMap(
// 	table string,
// 	data map[string][]string,
// 	allowed map[string]bool,
// ) (string, []any, error) {

// 	if table == "" {
// 		return "", nil, fmt.Errorf("table name required")
// 	}

// 	cols := []string{}
// 	vals := []any{}
// 	holders := []string{}

// 	i := 1
// 	for col, v := range data {
// 		if !allowed[col] || len(v) == 0 {
// 			continue
// 		}

// 		cols = append(cols, col)
// 		vals = append(vals, v[0])
// 		holders = append(holders, fmt.Sprintf("$%d", i))
// 		i++
// 	}

// 	if len(cols) == 0 {
// 		return "", nil, fmt.Errorf("no valid fields to insert")
// 	}

// 	query := fmt.Sprintf(
// 		"INSERT INTO %s (%s) VALUES (%s)",
// 		table,
// 		strings.Join(cols, ", "),
// 		strings.Join(holders, ", "),
// 	)

// 	return query, vals, nil
// }
