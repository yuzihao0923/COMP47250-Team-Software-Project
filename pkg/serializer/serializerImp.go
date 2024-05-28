package serializer

import (
	"encoding/json" //Only json for now
)

// JSONSerializer implement Serializer interface
type JSONSerializer struct{}

func (s *JSONSerializer) Serialize(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (s *JSONSerializer) Deserialize(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
