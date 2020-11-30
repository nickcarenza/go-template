package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// TemplateFuncs ...
var TemplateFuncs = map[string]interface{}{
	"toJSON": func(v interface{}) string {
		a, _ := json.Marshal(v)
		return string(a)
	},
	"now": func(layout string) string {
		return time.Now().Format(layout)
	},
	"timestamp": func() int64 {
		return time.Now().Unix()
	},
	"env": func(key string) string {
		return os.Getenv(key)
	},
	"trim": func(v, cutset string) string {
		return strings.Trim(v, cutset)
	},
	"multiply": func(x interface{}, y interface{}) float64 {
		var xFloat float64
		var yFloat float64
		switch v := x.(type) {
		case int:
			xFloat = float64(v)
		case float64:
			xFloat = v
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return 0
			}
			xFloat = f
		case json.Number:
			f, err := v.Float64()
			if err != nil {
				return 0
			}
			xFloat = f
		default:
			xFloat = 0
		}
		switch v := y.(type) {
		case int:
			yFloat = float64(v)
		case float64:
			yFloat = v
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return 0
			}
			yFloat = f
		case json.Number:
			f, err := v.Float64()
			if err != nil {
				return 0
			}
			yFloat = f
		default:
			yFloat = 0
		}
		return xFloat * yFloat
	},
	"ge": func(x interface{}, y interface{}) (bool, error) {
		var xFloat float64
		var yFloat float64
		switch v := x.(type) {
		case int:
			xFloat = float64(v)
		case float64:
			xFloat = v
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return false, err
			}
			xFloat = f
		case json.Number:
			f, err := v.Float64()
			if err != nil {
				return false, err
			}
			xFloat = f
		default:
			return false, fmt.Errorf(`unsupported type %T found in x-value of comparison "ge"`, x)
		}
		switch v := y.(type) {
		case int:
			yFloat = float64(v)
		case float64:
			yFloat = v
		case string:
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return false, err
			}
			yFloat = f
		case json.Number:
			f, err := v.Float64()
			if err != nil {
				return false, err
			}
			yFloat = f
		default:
			return false, fmt.Errorf(`unsupported type %T found in y-value of comparison "ge"`, y)
		}
		return xFloat >= yFloat, nil
	},
	"normalize_email": func(email string) string {
		// get everything before the first instance of '+'
		email = strings.Split(email, "+")[0]
		// get everything before the first instance of '@'
		email = strings.Split(email, "@")[0]
		email = strings.ReplaceAll(email, ".", "")
		re := regexp.MustCompile("[0-9]")
		email = re.ReplaceAllString(email, "")
		email = strings.TrimSpace(email)
		return email
	},
	"fingerprint_address": func(address, city, state, zip, plus4Code interface{}) string {
		var addressStr, cityStr, stateStr, zipStr, plus4CodeStr string
		addressStr, _ = address.(string)
		cityStr, _ = city.(string)
		stateStr, _ = state.(string)
		zipStr, _ = zip.(string)
		plus4CodeStr, _ = plus4Code.(string)
		fingerprint := strings.Join([]string{addressStr, cityStr, stateStr, zipStr, plus4CodeStr}, "_")
		re := regexp.MustCompile("[^a-zA-Z0-9_]")
		fingerprint = re.ReplaceAllString(fingerprint, "_")
		fingerprint = strings.ToLower(fingerprint)
		return fingerprint
	},
	"dict": func(keysAndValues ...interface{}) map[interface{}]interface{} {
		var dict = map[interface{}]interface{}{}
		for i, s := range keysAndValues {
			if i%2 != 0 {
				dict[keysAndValues[i-1]] = s
			}
		}
		return dict
	},
	"http": func(method, url string, headers map[interface{}]interface{}) (*http.Response, error) {
		var req *http.Request
		var err error
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
		if headers != nil {
			for k, v := range headers {
				req.Header.Set(k.(string), v.(string))
			}
		}
		return http.DefaultClient.Do(req)
	},
	"parseJSON": func(data []byte) (interface{}, error) {
		var v interface{}
		var err error
		err = json.Unmarshal(data, &v)
		return v, err
	},
	"formatTime": func(srcLayout, targetLayout, input string) (string, error) {
		t, err := time.Parse(srcLayout, input)
		if err != nil {
			return "", err
		}
		return t.Format(targetLayout), nil
	},
	"formatUnix": func(targetLayout string, input interface{}) (string, error) {
		switch v := input.(type) {
		case int64:
			return time.Unix(v, 0).Format(targetLayout), nil
		case float64:
			return time.Unix(int64(v), 0).Format(targetLayout), nil
		case int:
			return time.Unix(int64(v), 0).Format(targetLayout), nil
		case string:
			intVal, err := strconv.Atoi(v)
			if err != nil {
				return "", err
			}
			return time.Unix(int64(intVal), 0).Format(targetLayout), nil
		case json.Number:
			intVal, err := v.Int64()
			if err != nil {
				return "", err
			}
			return time.Unix(intVal, 0).Format(targetLayout), nil
		default:
			return "", fmt.Errorf("Invalid type for time.Unix in formatUnix")
		}
	},
	"split": func(sep, input string) []string {
		return strings.Split(input, sep)
	},
	"first": func(input interface{}) interface{} {
		list := interfaceSlice(input)
		if list == nil {
			return nil
		}
		if len(list) == 0 {
			return nil
		}
		return list[0]
	},
	"last": func(input interface{}) interface{} {
		list := interfaceSlice(input)
		if list == nil {
			return nil
		}
		if len(list) == 0 {
			return nil
		}
		return list[len(list)-1]
	},
	"coalesce": func(values ...interface{}) interface{} {
		for _, v := range values {
			if v != nil {
				return v
			}
		}
		return nil
	},
	"sortMap": func(list []interface{}, sortKey string, dir string) ([]interface{}, error) {
		if dir != "asc" && dir != "desc" {
			return nil, fmt.Errorf("Invalid sort direction for sortMap")
		}

		var err error
		sort.SliceStable(list, func(i, j int) bool {
			var a, b string
			switch _a := list[i].(map[string]interface{})[sortKey].(type) {
			case float64:
				a = fmt.Sprintf("%.0f", _a)
			case int64:
				a = fmt.Sprintf("%d", _a)
			case string:
				a = _a
			case nil:
				a = ""
			default:
				err = fmt.Errorf("Invalid value type in sortMap: %T", _a)
				return false
			}
			switch _b := list[j].(map[string]interface{})[sortKey].(type) {
			case float64:
				b = fmt.Sprintf("%.0f", _b)
			case int64:
				b = fmt.Sprintf("%d", _b)
			case string:
				b = _b
			case nil:
				b = ""
			default:
				err = fmt.Errorf("Invalid value type in sortMap: %T", _b)
				return false
			}
			if dir == "asc" {
				return a < b
			}
			return a > b
		})
		if err != nil {
			return nil, err
		}
		return list, nil
	},
	"add": func(a, b int) int {
		return a + b
	},
}

func interfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

// Interpolate simplifies interpolating a template string with data
func Interpolate(data interface{}, text string) (string, error) {
	tmpl, err := template.New("test").Funcs(TemplateFuncs).Parse(text)

	if err != nil {
		return text, err
	}

	var tBuf bytes.Buffer
	err = tmpl.Execute(&tBuf, data)

	if err != nil {
		return text, err
	}

	return tBuf.String(), nil
}

// InterpolateMap interpolates a recursive map
func InterpolateMap(data interface{}, templateMap map[string]interface{}) (map[string]interface{}, error) {
	var parsed = map[string]interface{}{}
	for key, i := range templateMap {
		if v, ok := i.(string); ok {
			str, err := Interpolate(data, v)
			if err != nil {
				return nil, err
			}
			parsed[key] = str
		} else if v, ok := i.(map[string]interface{}); ok {
			deepParsed, err := InterpolateMap(data, v)
			if err != nil {
				return nil, err
			}
			parsed[key] = deepParsed
			// } else if v, ok := i.([]map[string]interface{}); ok {
			// 	deepParsed, err := InterpolateMap(data, v)
			// 	if err != nil {
			// 		return nil, err
			// 	}
			// 	parsed[key] = deepParsed
		} else {
			parsed[key] = v
		}
	}
	return parsed, nil
}

// Template is a wrapper that implements unmarshalJSON
type Template struct {
	*template.Template
}

// UnmarshalJSON implementation for Template
func (t *Template) UnmarshalJSON(data []byte) (err error) {
	t.Template, err = template.New("template").Funcs(TemplateFuncs).Parse(string(data))
	return
}
