// Package json provides a json codec
package json

import (
	"io"
	"io/ioutil"

	"github.com/segmentio/encoding/json"
	"github.com/unistack-org/micro/v3/codec"
)

var (
	JsonMarshaler = &Marshaler{}

	JsonUnmarshaler = &Unmarshaler{
		ZeroCopy: true,
	}
)

type Marshaler struct {
	EscapeHTML      bool
	SortMapKeys     bool
	TrustRawMessage bool
}

type Unmarshaler struct {
	DisallowUnknownFields                bool
	DontCopyNumber                       bool
	DontCopyRawMessage                   bool
	DontCopyString                       bool
	DontMatchCaseInsensitiveStructFields bool
	UseNumber                            bool
	ZeroCopy                             bool
}

type jsonCodec struct {
	encodeFlags json.AppendFlags
	decodeFlags json.ParseFlags
}

func (c *jsonCodec) Marshal(b interface{}) ([]byte, error) {
	switch m := b.(type) {
	case nil:
		return nil, nil
	case *codec.Frame:
		return m.Data, nil
	}

	var err error
	var buf []byte
	buf, err = json.Append(buf, b, c.encodeFlags)
	return buf, err
}

func (c *jsonCodec) Unmarshal(b []byte, v interface{}) error {
	if b == nil {
		return nil
	}
	switch m := v.(type) {
	case nil:
		return nil
	case *codec.Frame:
		m.Data = b
		return nil
	}

	_, err := json.Parse(b, v, c.decodeFlags)
	return err
}

func (c *jsonCodec) ReadHeader(conn io.ReadWriter, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *jsonCodec) ReadBody(conn io.ReadWriter, b interface{}) error {
	switch m := b.(type) {
	case nil:
		return nil
	case *codec.Frame:
		buf, err := ioutil.ReadAll(conn)
		if err != nil {
			return err
		}
		m.Data = buf
		return nil
	}

	return json.NewDecoder(conn).Decode(b)
}

func (c *jsonCodec) Write(conn io.ReadWriter, m *codec.Message, b interface{}) error {
	switch m := b.(type) {
	case nil:
		return nil
	case *codec.Frame:
		_, err := conn.Write(m.Data)
		return err
	}

	return json.NewEncoder(conn).Encode(b)
}

func (c *jsonCodec) String() string {
	return "json"
}

func NewCodec() codec.Codec {
	var encodeFlags json.AppendFlags
	var decodeFlags json.ParseFlags

	if JsonMarshaler.EscapeHTML {
		encodeFlags |= json.EscapeHTML
	}

	if JsonMarshaler.SortMapKeys {
		encodeFlags |= json.SortMapKeys
	}

	if JsonMarshaler.TrustRawMessage {
		encodeFlags |= json.TrustRawMessage
	}

	if JsonUnmarshaler.DisallowUnknownFields {
		decodeFlags |= json.DisallowUnknownFields
	}

	if JsonUnmarshaler.DontCopyNumber {
		decodeFlags |= json.DontCopyNumber
	}

	if JsonUnmarshaler.DontCopyRawMessage {
		decodeFlags |= json.DontCopyRawMessage
	}

	if JsonUnmarshaler.DontCopyString {
		decodeFlags |= json.DontCopyString
	}

	if JsonUnmarshaler.DontMatchCaseInsensitiveStructFields {
		decodeFlags |= json.DontMatchCaseInsensitiveStructFields
	}

	if JsonUnmarshaler.UseNumber {
		decodeFlags |= json.UseNumber
	}

	if JsonUnmarshaler.ZeroCopy {
		decodeFlags |= json.ZeroCopy
	}

	return &jsonCodec{encodeFlags: encodeFlags, decodeFlags: decodeFlags}
}
