package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func BuildUrl(url string, path ...string) string {
	return fmt.Sprintf("%s%s", url, strings.Join(path, ""))
}

func HttpPostJson_Json[REQUEST any, RESPONSE any](
	url string,
	request *REQUEST,
) (*RESPONSE, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer SilentClose(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response RESPONSE
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
