package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func DecodeJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

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
func GetFormValue(formData map[string][]string, key string) (string, error) {
	values, ok := formData[key]
	if !ok || len(values) == 0 {
		return "", fmt.Errorf("missing or empty field: %s", key)
	}
	return values[0], nil
}

type AutoTime struct {
	time.Time
}

func (t *AutoTime) Scan(value any) error {
	if value == nil {
		t.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v
		return nil
	default:
		return fmt.Errorf("cannot scan %T into AutoTime", value)
	}
}

/*
OPTIONAL: allow writing back to DB if needed
*/
func (t AutoTime) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}

/*
AUTOMATIC JSON FORMATTING
*/
func (t AutoTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte(`""`), nil
	}

	return []byte(fmt.Sprintf(
		`"%s"`,
		t.Format("Jan 2 2006 15:04"),
	)), nil
}

func Strtoint64(amount string)(int64,error){
	value,err := strconv.ParseInt(amount,10,64)
	if err != nil {
		log.Println("Failed to convert to float")
		return 0,err
	}

	return value,nil
}


func Strtofloat(amount string)(float64,error){
	value,err := strconv.ParseFloat(amount,64)
	if err != nil {
		log.Println("Failed to convert to float")
		return 0,err
	}

	return value,nil
}

func Floattostr(amount float64)(string){
	value := strconv.FormatFloat(amount,'f',2,64)
	return value
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
