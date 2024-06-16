package serializer

import "io"

type Serializer interface {
	Serialize(data interface{}) ([]byte, error)
	Deserialize(data []byte, v interface{}) error
	SerializeToWriter(data interface{}, w io.Writer) error
	DeserializeFromReader(r io.Reader, v interface{}) error
}
