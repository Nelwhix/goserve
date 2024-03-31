package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

func FormatSSE(data string) (string, error) {
	m := map[string]string{
		"data": data,
	}

	buff := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buff)
	err := encoder.Encode(m)

	if err != nil {
		return "", err
	}

	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("data: %v\n\n", buff.String()))

	return sb.String(), nil
}
