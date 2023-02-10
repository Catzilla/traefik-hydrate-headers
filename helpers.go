package traefik_hydrate_headers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func compactJson(jsonB []byte) ([]byte, error) {

	var buff *bytes.Buffer = new(bytes.Buffer)
	errCompact := json.Compact(buff, jsonB)
	if errCompact != nil {
		newErr := fmt.Errorf("failure encountered compacting json := %v", errCompact)
		return []byte{}, newErr
	}

	b, err := io.ReadAll(buff)
	if err != nil {
		readErr := fmt.Errorf("read buffer error encountered := %v", err)
		return []byte{}, readErr
	}

	return b, nil
}

func contains[T comparable](elems []T, needle T) bool {
	for _, el := range elems {
		if el == needle {
			return true
		}
	}

	return false
}
