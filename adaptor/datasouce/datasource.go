package datasouce

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/ac-kurniawan/proxy/core"
	"github.com/ac-kurniawan/proxy/library"
)

type Datasource struct {
	Host        string
	Port        string
	TokenPrefix string
	Trace       library.AppTrace
}

// Call implements core.IDataSource.
func (d *Datasource) Call(ctx context.Context, method string, path string, token *string, payload map[string]interface{}) (map[string]interface{}, error) {
	ctx, span := d.Trace.StartTrace(ctx, "DATASOURCE - Call")
	defer d.Trace.EndTrace(span)
	var client = &http.Client{}
	val := url.Values{}
	for key, value := range payload {
		val.Add(key, fmt.Sprintf("%v", value))
	}
	data := bytes.NewBufferString(val.Encode())

	req, _ := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s:%s%s", d.Host, d.Port, path), data)
	if method == http.MethodPost {
		contentType := multipart.NewWriter(data)
		contentType.Close()
		req.Header.Set("Content-Type", contentType.FormDataContentType())
	}
	if token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", d.TokenPrefix, *token))
	}
	res, err := client.Do(req)
	if err != nil {
		d.Trace.TraceError(span, err)
		return nil, err
	}
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		d.Trace.TraceError(span, err)
		return nil, err
	}
	var output map[string]interface{}
	err = json.Unmarshal(bytes, &output)
	if err != nil {
		d.Trace.TraceError(span, err)
		return nil, err
	}
	return output, nil
}

func NewDatasource(module Datasource) core.IDataSource {
	return &module
}
