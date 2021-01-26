package template

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"text/template"
)

func TestInterpolateMap(t *testing.T) {
	var tmpl = map[string]interface{}{
		"event_id": "{{ .event.id }}",
	}
	var data = map[string]interface{}{
		"event": map[string]interface{}{
			"id": "8D469E95-D2CA-4DF4-A67C-C141B51AFE99",
		},
	}
	res, err := InterpolateMap(data, tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	if res["event_id"] != "8D469E95-D2CA-4DF4-A67C-C141B51AFE99" {
		t.Error("event_id is wrong")
		return
	}
}

func TestTemplateFuncFormatTime(t *testing.T) {
	var tpl = `{{formatTime "2006-01-02" "Mon Jan 2 2006" "2020-11-23"}}`
	tmpl, err := template.New(t.Name()).Funcs(map[string]interface{}{
		"formatTime": TemplateFuncs["formatTime"],
	}).Parse(tpl)

	if err != nil {
		t.Error(err)
		return
	}

	var data = map[string]interface{}{}

	var tBuf bytes.Buffer
	err = tmpl.Execute(&tBuf, data)

	if err != nil {
		t.Error(err)
		return
	}

	if tBuf.String() != "Mon Nov 23 2020" {
		t.Error("Should format input time according to format")
	}
}

func TestTemplateFuncFormatUnix(t *testing.T) {
	var tpl = `{{formatUnix "Mon Jan 2 2006" 1606158230}}`
	tmpl, err := template.New(t.Name()).Funcs(map[string]interface{}{
		"formatUnix": TemplateFuncs["formatUnix"],
	}).Parse(tpl)

	if err != nil {
		t.Error(err)
		return
	}

	var data = map[string]interface{}{}

	var tBuf bytes.Buffer
	err = tmpl.Execute(&tBuf, data)

	if err != nil {
		t.Error(err)
		return
	}

	if tBuf.String() != "Mon Nov 23 2020" {
		t.Error("Should format input unix timestamp according to format")
	}
}

func TestTemplateFuncDict(t *testing.T) {
	var tpl = `{{ $d := dict "a" 1 "b" 2 "c" 3 }}{{ $d.b }}`
	tmpl, err := template.New(t.Name()).Funcs(map[string]interface{}{
		"dict": TemplateFuncs["dict"],
	}).Parse(tpl)

	if err != nil {
		t.Error(err)
		return
	}

	var data = map[string]interface{}{}

	var tBuf bytes.Buffer
	err = tmpl.Execute(&tBuf, data)

	if err != nil {
		t.Error(err)
		return
	}

	if tBuf.String() != "2" {
		t.Error("Should make a map from arguments")
	}
}

func TestTemplateFuncHttp(t *testing.T) {
	t.Skip()
}

func TestTemplateFuncSplit(t *testing.T) {
	// var tpl = `{{ split " " "first middle last" | printf "%s" }}`
	var tpl = `{{ $parts := split " " "first middle last" }}{{ printf "%s %s %s" (index $parts 2) (index $parts 1) (index $parts 0) }}`
	tmpl, err := template.New(t.Name()).Funcs(map[string]interface{}{
		"split": TemplateFuncs["split"],
	}).Parse(tpl)

	if err != nil {
		t.Error(err)
		return
	}

	var data = map[string]interface{}{}

	var tBuf bytes.Buffer
	err = tmpl.Execute(&tBuf, data)

	if err != nil {
		t.Error(err)
		return
	}

	if tBuf.String() != "last middle first" {
		t.Log(tBuf.String())
		t.Error("Should split input by delim")
	}
}

func TestTemplateFuncFirst(t *testing.T) {
	var tpl = `{{first .list}}`
	tmpl, err := template.New(t.Name()).Funcs(map[string]interface{}{
		"first": TemplateFuncs["first"],
	}).Parse(tpl)

	if err != nil {
		t.Error(err)
		return
	}

	var data = map[string]interface{}{
		"list": []string{"one", "two", "three"},
	}

	var tBuf bytes.Buffer
	err = tmpl.Execute(&tBuf, data)

	if err != nil {
		t.Error(err)
		return
	}

	if tBuf.String() != "one" {
		t.Error("Should return first element in list")
	}
}

func TestTemplateFuncFirstOfEmpty(t *testing.T) {
	var tpl = `{{first .list}}`
	tmpl, err := template.New(t.Name()).Funcs(map[string]interface{}{
		"first": TemplateFuncs["first"],
	}).Parse(tpl)

	if err != nil {
		t.Error(err)
		return
	}

	var data = map[string]interface{}{
		"list": []string{},
	}

	var tBuf bytes.Buffer
	err = tmpl.Execute(&tBuf, data)

	if err != nil {
		t.Error(err)
		return
	}

	if tBuf.String() != "<no value>" {
		t.Log(tBuf.String())
		t.Error("Should return first element in list")
	}
}

func TestTemplateFuncLast(t *testing.T) {
	var tpl = `{{last .list}}`
	tmpl, err := template.New(t.Name()).Funcs(map[string]interface{}{
		"last": TemplateFuncs["last"],
	}).Parse(tpl)

	if err != nil {
		t.Error(err)
		return
	}

	var data = map[string]interface{}{
		"list": []string{"one", "two", "three"},
	}

	var tBuf bytes.Buffer
	err = tmpl.Execute(&tBuf, data)

	if err != nil {
		t.Error(err)
		return
	}

	if tBuf.String() != "three" {
		t.Error("Should return last element in list")
	}
}

type UnmarshalTarget struct {
	Template *Template `json:"template"`
}

func TestUnmarshalJSONStr(t *testing.T) {
	var err error
	var jsondata = []byte(`{"template":"\"{{ .str }}\""}`)
	var tmpl UnmarshalTarget
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Template.Execute(&buf, map[string]interface{}{
		"str": "hello world!",
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != `"hello world!"` {
		t.Log(buf.String())
		t.Error("Unexpected template output")
	}
}

func TestUnmarshalJSONNumericLiteral(t *testing.T) {
	var err error
	var jsondata = []byte(`{"template":"{{ 5 }}"}`)
	var tmpl UnmarshalTarget
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Template.Execute(&buf, map[string]interface{}{})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != `5` {
		t.Error("Should result in 5")
		return
	}
}

func TestUnmarshalJSONMultiplyNum(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ multiply 5 .num }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"num": json.Number("5"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != `25` {
		t.Error("Should result in 25", buf.String())
	}
}

func TestUnmarshalJSONMultiplyNumericString(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ multiply 5 .num }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"num": "5",
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != `25` {
		t.Error("Should result in 25", buf.String())
	}
}

func TestUnmarshalJSONMultiplyFloat64(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ multiply 5 .num }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"num": float64(5),
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != `25` {
		t.Error("Should result in 25", buf.String())
	}
}

func TestEscapedToJSON(t *testing.T) {
	var err error
	var jsondata = []byte("\"{{ .Event.comments | toJSON }}\"")
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"Event": map[string]interface{}{
			"comments": "\"hello\"",
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != `"\"hello\""` {
		t.Error(`Should result in "\"hello\""`, buf.String())
	}
}

func TestToJSONNull(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ .Event.comments | toJSON }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"Event": map[string]interface{}{},
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != `null` {
		t.Errorf(`Unexpected result %q`, buf.String())
	}
}

func TestToJsonWithQuotes(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ .Event.comments | toJSON }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	var event map[string]interface{}
	err = json.Unmarshal([]byte(`{"comments":"quote \"this\""}`), &event)
	err = tmpl.Execute(&buf, map[string]interface{}{
		"Event": event,
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != `"quote \"this\""` {
		t.Error(`Unexpected result`, buf.String())
	}
}

func TestUnquote(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ .Event.comments | toJSON | unquote }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	var event map[string]interface{}
	err = json.Unmarshal([]byte(`{"comments":"quote \"this\""}`), &event)
	err = tmpl.Execute(&buf, map[string]interface{}{
		"Event": event,
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != `quote \"this\"` {
		t.Error(`Unexpected result`, buf.String())
	}
}

func TestJsonNumberMethod(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{.number.Float64 | printf \"%.0f\"}}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"number": json.Number("5"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != `5` {
		t.Error(`Unexpected result`, buf.String())
	}
}

func TestUUID(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ uuid }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{})
	if err != nil {
		t.Error(err)
		return
	}
	if len(buf.String()) != 36 {
		t.Error(`Unexpected result`, buf.String())
	}
}

func TestParseJSONBytes(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ (parseJSON .bytes).key }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"bytes": []byte(`{"key":"value"}`),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(buf.String())
	if buf.String() != "value" {
		t.Error(`Unexpected result`, buf.String())
	}
}

func TestParseJSONBuffer(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ (parseJSON .buf).key }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"buf": bytes.NewBufferString(`{"key":"value"}`),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(buf.String())
	if buf.String() != "value" {
		t.Error(`Unexpected result`, buf.String())
	}
}

func TestParseJSONString(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ (parseJSON .string).key }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"string": `{"key":"value"}`,
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(buf.String())
	if buf.String() != "value" {
		t.Error(`Unexpected result`, buf.String())
	}
}

func TestParseJSONIOReader(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ (parseJSON .ioreader).key }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"ioreader": bytes.NewReader([]byte(`{"key":"value"}`)),
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(buf.String())
	if buf.String() != "value" {
		t.Error(`Unexpected result`, buf.String())
	}
}

func TestCacheSetAndGet(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ with (cacheSet \"test\" 5 \"1m\") }}{{ end }}{{ cacheGet \"test\" }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != "5" {
		t.Error(`Unexpected result`, buf.String())
	}
}

func TestCacheCheckAndSet(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ $cachedVal := (cacheGet \"test2\") }}{{ if $cachedVal }}{{ print $cachedVal }}{{ else }}{{ $newVal := \"newVal\" }}{{ $_ := (cacheSet \"test2\" $newVal \"1m\") }}{{ print $newVal }}{{ end }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != "newVal" {
		t.Errorf(`Unexpected result %q`, buf.String())
	}
}

func TestGetAuthXBearerToken(t *testing.T) {
	t.SkipNow()
	var err error
	var jsondata = []byte(`"{{ getAuthXBearerToken \"url...\" \"token...\" \"userId...\" ` + "`" + `[\"**:**:**\"]` + "`" + `}}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{})
	if err != nil {
		t.Error(err)
		return
	}
	if !strings.HasPrefix(buf.String(), "Bearer") {
		t.Error("No bearer token returned")
	}
}
