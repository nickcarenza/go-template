package template

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/google/uuid"
	"github.com/the-control-group/go-timeutils"
	"github.com/the-control-group/go-ttlcache"
)

var templateCache *ttlcache.TTLCache
var authxTokenCache *ttlcache.TTLCache
var sprigFuncs = sprig.FuncMap()

func init() {
	// Create template cache
	templateCache = ttlcache.NewTTLCache(15 * time.Minute)
	authxTokenCache = ttlcache.NewTTLCache(5 * time.Minute)
}

// TemplateFuncs ...
var TemplateFuncs = map[string]interface{}{
	"uuid": func() (string, error) {
		id, err := uuid.NewRandom()
		if err != nil {
			return "", err
		}
		return id.String(), nil
	},
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
		email = strings.ToLower(email)
		return email
	},
	"toLower": func(str string) string {
		return strings.ToLower(str)
	},
	"fingerprint": func(vars ...string) (fingerprint string) {
		fingerprint = strings.Join(vars, "_")
		re := regexp.MustCompile(`[^\p{L}0-9]`)
		fingerprint = re.ReplaceAllString(fingerprint, "_")
		fingerprint = strings.ToLower(fingerprint)
		return
	},
	"fingerprint_address": func(address, city, state, zip, plus4Code interface{}) string {
		var addressStr, cityStr, stateStr, zipStr, plus4CodeStr string
		addressStr, _ = address.(string)
		cityStr, _ = city.(string)
		stateStr, _ = state.(string)
		zipStr, _ = zip.(string)
		plus4CodeStr, _ = plus4Code.(string)
		fingerprint := strings.Join([]string{addressStr, cityStr, stateStr, zipStr, plus4CodeStr}, "_")
		re := regexp.MustCompile(`[^\p{L}0-9]`)
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
	"parseJSON": func(data interface{}) (interface{}, error) {
		var v interface{}
		var err error
		switch d := data.(type) {
		case []byte:
			err = json.Unmarshal(d, &v)
			return v, err
		case string:
			err = json.Unmarshal([]byte(d), &v)
			return v, err
		case bytes.Buffer:
			err = json.Unmarshal(d.Bytes(), &v)
			return v, err
		case io.Reader:
			var buf bytes.Buffer
			buf.ReadFrom(d)
			err = json.Unmarshal(buf.Bytes(), &v)
			return v, err
		}
		return nil, fmt.Errorf("TypeAssertionError")
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
		case time.Time:
			return v.Format(targetLayout), nil
		default:
			return "", fmt.Errorf("Invalid type for time.Unix in formatUnix")
		}
	},
	"formatUnixFull": func(targetLayout string, seconds interface{}, nanoseconds interface{}) (string, error) {
		var secs, nanos int64
		var err error
		secs, err = interfaceToInt64(seconds)
		if err != nil {
			return "", err
		}
		nanos, err = interfaceToInt64(nanoseconds)
		if err != nil {
			return "", err
		}
		return time.Unix(secs, nanos).Format(targetLayout), nil
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
	"addInt64": func(a, b int64) int64 {
		return a + b
	},
	"unquote": func(s string) string {
		if len(s) > 0 && s[0] == '"' {
			s = s[1:]
		}
		if len(s) > 0 && s[len(s)-1] == '"' {
			s = s[:len(s)-1]
		}
		return s
	},
	"getAuthXBearerToken": func(authxURL, authxToken, authorizationId string) (string, error) {
		var cacheKey = strings.Join([]string{authxURL, authxToken, authorizationId}, "::")
		cachedToken, _ := authxTokenCache.Get(cacheKey)
		if cachedTokenString, ok := cachedToken.(string); ok {
			return cachedTokenString, nil
		}
		var err error
		var graphqlQuery = fmt.Sprintf(`query {
			authorization(id: %q) {
				token(format:BEARER)
			}
		}`, authorizationId)
		var requestQuery = map[string]interface{}{
			"query": graphqlQuery,
		}
		var requestBody []byte
		requestBody, err = json.Marshal(requestQuery)
		var req *http.Request
		req, err = http.NewRequest("POST", authxURL, bytes.NewBuffer(requestBody))
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", authxToken)
		req.Header.Set("Content-Type", "application/json")
		var res *http.Response
		res, err = http.DefaultClient.Do(req)
		var body []byte
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		defer res.Body.Close()
		var tokenResponse struct {
			Errors []struct {
				Message string
			}
			Data struct {
				Authorization struct {
					Token string
				}
			}
		}
		err = json.Unmarshal(body, &tokenResponse)
		if err != nil {
			return "", err
		}
		if tokenResponse.Errors != nil && len(tokenResponse.Errors) > 0 {
			return "", fmt.Errorf("Authx error: %s", tokenResponse.Errors[0].Message)
		}
		var authxBearerToken = tokenResponse.Data.Authorization.Token
		tokenParts := strings.Split(strings.Split(authxBearerToken, " ")[1], ".")
		jwtBase64 := tokenParts[1]
		var jwtBytes []byte
		jwtBytes, err = base64.RawURLEncoding.DecodeString(jwtBase64)
		if err != nil {
			return "", err
		}
		var jwt struct {
			AID    string
			Scopes []string
			IAT    int64
			EXP    int64
			ISS    string
			SUB    string
			JTI    string
		}
		err = json.Unmarshal(jwtBytes, &jwt)
		if err != nil {
			return "", err
		}
		var expireAt = time.Duration(jwt.EXP-time.Now().Unix())*time.Second - time.Minute
		authxTokenCache.SetEx(cacheKey, authxBearerToken, expireAt)
		return authxBearerToken, nil
	},
	"cacheSet": func(key string, value interface{}, expire interface{}) (interface{}, error) {
		exp, err := timeutils.InterfaceToApproxBigDuration(expire)
		if err != nil {
			return value, err
		}
		return value, templateCache.SetEx(key, value, time.Duration(exp))
	},
	"cacheGet": func(key string) interface{} {
		v, _ := templateCache.Get(key)
		return v
	},
	"parseCIDR": func(cidr string) (*net.IPNet, error) {
		_, ipnet, err := net.ParseCIDR(cidr)
		return ipnet, err
	},
	"toApproxBigDuration": func(i interface{}) (timeutils.ApproxBigDuration, error) {
		return timeutils.InterfaceToApproxBigDuration(i)
	},
	"int":            sprigFuncs["int"],
	"int64":          sprigFuncs["int64"],
	"atoi":           sprigFuncs["atoi"],
	"b64dec":         sprigFuncs["b64dec"],
	"b64enc":         sprigFuncs["b64enc"],
	"ternary":        sprigFuncs["ternary"],
	"sha1sum":        sprigFuncs["sha1sum"],
	"sha256sum":      sprigFuncs["sha256sum"],
	"encryptAES":     sprigFuncs["encryptAES"],
	"decryptAES":     sprigFuncs["decryptAES"],
	"parseTime":      timeutils.ParseAny,
	"maybeParseTime": timeutils.ParseAnyMaybe,
	"formatAnyTime": func(targetLayout, input string) (string, error) {
		t, err := timeutils.ParseAny(input)
		if err != nil {
			return "", err
		}
		return t.Format(targetLayout), nil
	},
	"maybeFormatAnyTime": func(targetLayout, input string) *string {
		t := timeutils.ParseAnyMaybe(input)
		if t == nil {
			return nil
		}
		s := t.Format(targetLayout)
		return &s
	},
	"left": func(str string, n int) string {
		if len(str) <= n {
			return str
		}
		return str[:n]
	},
	"right": func(str string, n int) string {
		if len(str) <= n {
			return str
		}
		return str[len(str)-n:]
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
		} else if v, ok := i.(float64); ok {
			parsed[key] = v
		} else if v, ok := i.(int64); ok {
			parsed[key] = v
		} else if v, ok := i.(int); ok {
			parsed[key] = v
		} else if v, ok := i.(json.Number); ok {
			f, err := v.Float64()
			if err != nil {
				return nil, err
			}
			parsed[key] = f
		} else if v, ok := i.(bool); ok {
			parsed[key] = v
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
	var src string
	err = json.Unmarshal(data, &src)
	if err != nil {
		return err
	}
	t.Template, err = template.New("template").Funcs(TemplateFuncs).Parse(src)
	return
}

// ExecuteToString executes the template and returns the result as a string
func (t *Template) ExecuteToString(data interface{}) (string, error) {
	var tBuf bytes.Buffer
	var err = t.Execute(&tBuf, data)

	if err != nil {
		return "", err
	}

	return tBuf.String(), nil
}

// ExecuteToInt executes the template and returns the result as an int
func (t *Template) ExecuteToInt(data interface{}) (int, error) {
	var tBuf bytes.Buffer
	var err = t.Execute(&tBuf, data)

	if err != nil {
		return 0, err
	}

	return strconv.Atoi(tBuf.String())
}

// Parse is a shorthand for template.Parse using templatefuncs
func Parse(src string) (*Template, error) {
	t, err := template.New("template").Funcs(TemplateFuncs).Parse(src)
	if err != nil {
		return nil, err
	}
	return &Template{t}, nil
}

// Must is an feature copy of template.Must
func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}

func interfaceToInt64(i interface{}) (int64, error) {
	switch v := i.(type) {
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case int:
		return int64(v), nil
	case string:
		intVal, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		return int64(intVal), nil
	case json.Number:
		intVal, err := v.Int64()
		if err != nil {
			return 0, err
		}
		return intVal, nil
	default:
		return 0, fmt.Errorf("Unable to convert type to int64")
	}
}
