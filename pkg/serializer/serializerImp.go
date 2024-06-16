package serializer

import (
	"encoding/json"
	"io"
)

// JSONSerializer implements Serializer interface
type JSONSerializer struct{}

func (s *JSONSerializer) Serialize(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (s *JSONSerializer) Deserialize(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (s *JSONSerializer) SerializeToWriter(data interface{}, w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(data)
}

func (s *JSONSerializer) DeserializeFromReader(r io.Reader, v interface{}) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(v)
}
