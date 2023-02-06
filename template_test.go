package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
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

func TestInterpolateMapTypes(t *testing.T) {
	var tmpl = map[string]interface{}{
		"string": "test",
		"int":    1,
		"float":  0.5,
		"true":   true,
		"false":  false,
	}
	var data = map[string]interface{}{}
	res, err := InterpolateMap(data, tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	if res["string"] != "test" {
		t.Errorf("string is wrong %T:%[1]v", res["string"])
		return
	}
	if res["int"] != 1 {
		t.Errorf("int is wrong %T:%[1]v", res["int"])
		return
	}
	if res["float"] != 0.5 {
		t.Errorf("float is wrong %T:%[1]v", res["float"])
		return
	}
	if b, ok := res["true"].(bool); !ok || b != true {
		t.Errorf("true is wrong %T:%[1]v (%v,%v)", res["true"], b, ok)
		return
	}
	if b, ok := res["false"].(bool); !ok || b != false {
		t.Errorf("false is wrong %T:%[1]v (%v,%v)", res["false"], b, ok)
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

func TestUnmarshalJSONDivideNum(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ multiply 0.1 .num }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"num": json.Number("100"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != `10` {
		t.Error("Should result in 10", buf.String())
	}
}

func TestUnmarshalJSONDivideNumericString(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ multiply 0.1 .num }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"num": "100",
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != `10` {
		t.Error("Should result in 10", buf.String())
	}
}

func TestUnmarshalJSONDivideFloat64(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ multiply 0.1 .num }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"num": float64(100),
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != `10` {
		t.Error("Should result in 10", buf.String())
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
	if err != nil {
		t.Error(err)
		return
	}
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
	if err != nil {
		t.Error(err)
		return
	}
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

func TestCacheSetTwice(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ with (cacheSet \"test\" 5 \"1m\") }}{{ end }}{{ with (cacheSet \"test\" 5 \"1m\") }}{{ end }}{{ cacheGet \"test\" }}"`)
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

func TestDictLookup(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{- $cardBrandMap := dict \"mastercard\" \"mastercard\" \"MC\" \"mastercard\" \"MasterCard\" \"mastercard\" \"visa\" \"visa\" \"Visa\" \"visa\" \"VI\" \"visa\" \"discover\" \"discover\" \"Discover\" \"discover\" \"DI\" \"discover\" \"american_express\" \"american_express\" \"American Express\" \"american_express\" \"american express\" \"american_express\" \"AX\" \"american_express\" \"jcb\" \"jcb\" \"JCB\" \"jcb\" -}}{{- $cardBrand := (coalesce (index $cardBrandMap (coalesce .paymentOptionPpd.card_network .paymentOptionPpd.tokenResponse.type \"visa\")) \"visa\") -}}{{- $cardBrand -}}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"paymentOptionPpd": map[string]interface{}{
			"card_network": "visa",
			"tokenResponse": map[string]interface{}{
				"type": "VI",
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != "visa" {
		t.Errorf(`Unexpected result %q`, buf.String())
	}
}

func TestHttp(t *testing.T) {
	t.Skip()
	var err error
	var jsondata = []byte(`"{{- $headers := dict \"Accept-Version\" \"3\" -}}{{- $res := (http \"GET\" \"https://lookup.binlist.net/372723\" $headers ).Body | parseJSON -}}{{- index $res \"scheme\" -}}"`)
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
	if buf.String() != "amex" {
		t.Errorf(`Unexpected result %q`, buf.String())
	}
}

func TestToApproxBigDuration(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{- (3600000000000 | toApproxBigDuration).Pretty -}}"`)
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
	if buf.String() != "1h0m0s" {
		t.Errorf(`Unexpected result %q`, buf.String())
	}
}

func TestToApproxBigDurationJson(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{- (.n | toApproxBigDuration).Pretty -}}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"n": json.Number("3600000000000"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != "1h0m0s" {
		t.Errorf(`Unexpected result %q`, buf.String())
	}
}

func TestToApproxBigDurationMath(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{- addInt64 100 (\"1h\" | toApproxBigDuration | int64) -}}"`)
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
	if buf.String() != "3600000000100" {
		t.Errorf(`Unexpected result %q`, buf.String())
	}
}

func TestToApproxBigDurationMathToISO8601(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ addInt64 1617838340 ((\"1 day\" | toApproxBigDuration).Seconds | int64) | formatUnix \"2006-01-02T15:04:05Z\" }}"`)
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
	if buf.String() != "2021-04-08T16:32:20Z" {
		t.Errorf(`Unexpected result %q`, buf.String())
	}
}

func TestParseCIDRv6(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{- if le (len .ip) 15 -}}{{ .ip }}{{- else -}}{{- print .ip \"/64\" | parseCIDR }}{{- end -}}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"ip": "2601:201:4381:8a0:e830:3b3d:4b34:f2e3",
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != "2601:201:4381:8a0::/64" {
		t.Errorf(`Unexpected result %q`, buf.String())
	}
}

func TestParseCIDRv4(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{- if le (len .ip) 15 -}}{{ .ip }}{{- else -}}{{- print .ip \"/64\" | parseCIDR }}{{- end -}}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"ip": "96.230.197.226",
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != "96.230.197.226" {
		t.Errorf(`Unexpected result %q`, buf.String())
	}
}

func TestAdd(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ add 5 6}}"`)
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
	if buf.String() != "11" {
		t.Errorf(`Unexpected result %q`, buf.String())
	}
}

func TestAddJsonNumber(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ add (.x.Int64 | int) 6}}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"x": json.Number("5"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != "11" {
		t.Errorf(`Unexpected result %q`, buf.String())
	}
}

func TestSliceString(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ slice \"123456\" 0 1 }}"`)
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
	if buf.String() != "1" {
		t.Errorf(`Unexpected result %q`, buf.String())
	}
}

func TestExecuteToString(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ .key }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var str string
	str, err = tmpl.ExecuteToString(map[string]interface{}{
		"key": "value",
	})
	if err != nil {
		t.Error(err)
		return
	}
	if str != "value" {
		t.Fail()
	}
	str, err = tmpl.ExecuteToString(map[string]interface{}{
		"key": 5,
	})
	if err != nil {
		t.Error(err)
		return
	}
	if str != "5" {
		t.Fail()
	}
}

func TestExecuteToInt(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ .key }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var i int
	i, err = tmpl.ExecuteToInt(map[string]interface{}{
		"key": "4",
	})
	if err != nil {
		t.Error(err)
		return
	}
	if i != 4 {
		t.Fail()
	}
	i, err = tmpl.ExecuteToInt(map[string]interface{}{
		"key": 4,
	})
	if err != nil {
		t.Error(err)
		return
	}
	if i != 4 {
		t.Fail()
	}
}

func TestMaybeFormatAnyTimeExists(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ .sometimes_time | maybeFormatAnyTime \"2006-01-02\" }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"sometimes_time": "2006-01-02T15:04:05Z",
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != "2006-01-02" {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestMaybeFormatAnyTimeNoExists(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ .sometimes_time | maybeFormatAnyTime \"2006-01-02\" }}"`)
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
	if buf.String() != "<nil>" {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestFingerprintAddress(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ fingerprint_address \"1234 adams st.\" \"city\" \"state\" \"12345\" \"1234\" }}"`)
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
	if buf.String() != "1234_adams_st__city_state_12345_1234" {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestFingerprintAddressForeignCharacters(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ fingerprint_address \"台江区\" \"福州\" \"福建\" \"350000\" \"\" }}"`)
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
	if buf.String() != "台江区_福州_福建_350000_" {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestFingerprint(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ fingerprint \"1234 adams st.\" \"city\" \"state\" \"12345\" \"1234\" }}"`)
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
	if buf.String() != "1234_adams_st__city_state_12345_1234" {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestRight(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ right \"2019\" 2 }}"`)
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
	if buf.String() != "19" {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestLeft(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ left \"2019\" 2 }}"`)
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
	if buf.String() != "20" {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestBase64Enc(t *testing.T) {
	var err error
	var jsondata = []byte(`"Basic {{ print \"user\" \":\" | b64enc }}"`)
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
	if buf.String() != "Basic dXNlcjo=" {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestParseAnyTimeSub(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ ((parseTime \"2021-08-26 02:33:08.000\").Sub (parseTime \"2021-08-26 02:26:05.000\")).Seconds | int64 }}"`)
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
	if buf.String() != "423" {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestToAmount(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ (23.45435 | toAmount).ToString }}"`)
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
	if buf.String() != "23.45" {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestToAmountDollars(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ (23.99 | toAmount).Dollars }}"`)
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
	if buf.String() != "23" {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestBooleanTrue(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ .boolTrue | toJSON }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"boolTrue": true,
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != "true" {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestBooleanNull(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ .boolTrue | toJSON }}"`)
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
	if buf.String() != "null" {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestMapLoop(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ range $key, $value := .map }}\"prefix_{{ $key }}\": \"{{ $value }}\",{{ end }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"map": map[string]interface{}{
			"k": "v",
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != `"prefix_k": "v",` {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestMap(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ .map }}"`)
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"map": map[string]interface{}{
			"k": "v",
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != `map[k:v]` {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestHttpGet(t *testing.T) {
	t.Skip() // httpstat.us is down
	var err error
	var jsondata = []byte(`"{{ (http \"GET\" \"https://httpstat.us/200\" (dict)).StatusCode }}"`)
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
	if buf.String() != `200` {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestOnlyDigits(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ onlyDigits \"a123fg45\" }}"`)
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
	if buf.String() != `12345` {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestOnlyAlpha(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ onlyAlpha \"0a123bc45d\" }}"`)
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
	if buf.String() != `abcd` {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestNoSpace(t *testing.T) {
	var err error
	var jsondata = []byte(`"{{ nospace \" 1 2 3 4 5\" }}"`)
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
	if buf.String() != `12345` {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestInterpolateDisallowUnsafeRender(t *testing.T) {
	var err error

	// Create temp file to use as render template
	f, err := os.CreateTemp(``, `go.template.test.render.*.tmp`)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = f.WriteString(`{{ .data }}`)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(f.Name())

	_, err = Interpolate(map[string]interface{}{
		"data": "x",
	}, fmt.Sprintf(`{{ UNSAFE_render "%s" . }}`, f.Name()))

	if err == nil {
		t.Fail()
		return
	}
}

func TestInterpolateAllowUnsafeRender(t *testing.T) {
	var err error

	// Create temp file to use as render template
	f, err := os.CreateTemp(``, `go.template.test.render.*.tmp`)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = f.WriteString(`{{ .data }}`)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(f.Name())

	AllowUnsafeRender(true)
	defer AllowUnsafeRender(false)

	res, err := Interpolate(map[string]interface{}{
		"data": "x",
	}, fmt.Sprintf(`{{ UNSAFE_render "%s" . }}`, f.Name()))

	if err != nil {
		t.Error(err)
		return
	}

	if res != `x` {
		t.Log(res)
		t.Fail()
		return
	}
}

func TestJsonAllowUnsafeRender(t *testing.T) {
	var err error

	// Create temp file to use as render template
	f, err := os.CreateTemp(``, `go.template.test.render.*.tmp`)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = f.WriteString(`{{ .data }}`)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(f.Name())

	AllowUnsafeRender(true)
	defer AllowUnsafeRender(false)

	// Create main template to call partial by filename
	var jsondata = []byte(fmt.Sprintf(`"{{ UNSAFE_render \"%s\" . }}"`, f.Name()))
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"data": "x",
	})
	if err != nil {
		t.Error(err)
		return
	}
	if buf.String() != `x` {
		t.Log(buf.String())
		t.Fail()
	}
}

func TestJsonDisallowUnsafeRender(t *testing.T) {
	var err error

	// Create temp file to use as render template
	f, err := os.CreateTemp(``, `go.template.test.render.*.tmp`)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = f.WriteString(`{{ .data }}`)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(f.Name())

	// Create main template to call partial by filename
	var jsondata = []byte(fmt.Sprintf(`"{{ UNSAFE_render \"%s\" . }}"`, f.Name()))
	var tmpl *Template
	err = json.Unmarshal(jsondata, &tmpl)
	if err != nil {
		t.Error(err)
		return
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]interface{}{
		"data": "x",
	})
	if err == nil {
		t.Fail()
		return
	}
}
