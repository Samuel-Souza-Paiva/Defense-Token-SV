package defense

import (
	"bytes"
	"encoding/json"
)

func ToBuffer(data any) *bytes.Buffer {
  b, _ := json.Marshal(data)

  buffer := bytes.NewBuffer(b)

  return buffer
}
