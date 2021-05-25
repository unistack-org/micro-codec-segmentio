// Package json provides a json codec
package json

import (
	"io"

	"github.com/segmentio/encoding/json"
	"github.com/unistack-org/micro/v3/codec"
	rutil "github.com/unistack-org/micro/v3/util/reflect"
)

var (
	JsonMarshaler = &Marshaler{}

	JsonUnmarshaler = &Unmarshaler{
		ZeroCopy: true,
	}
)

const (
	flattenTag = "flatten"
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

func (c *jsonCodec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case nil:
		return nil, nil
	case *codec.Frame:
		return m.Data, nil
	}

	if nv, nerr := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); nerr == nil {
		v = nv
	}

	var err error
	var buf []byte
	buf, err = json.Append(buf, v, c.encodeFlags)
	return buf, err
}

func (c *jsonCodec) Unmarshal(b []byte, v interface{}) error {
	if len(b) == 0 {
		return nil
	}
	switch m := v.(type) {
	case nil:
		return nil
	case *codec.Frame:
		m.Data = b
		return nil
	}

	if nv, nerr := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); nerr == nil {
		v = nv
	}
	_, err := json.Parse(b, v, c.decodeFlags)
	return err
}

func (c *jsonCodec) ReadHeader(conn io.Reader, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *jsonCodec) ReadBody(conn io.Reader, v interface{}) error {
	if v == nil {
		return nil
	}

	buf, err := io.ReadAll(conn)
	if err != nil {
		return err
	} else if len(buf) == 0 {
		return nil
	}
	return c.Unmarshal(buf, v)
}

func (c *jsonCodec) Write(conn io.Writer, m *codec.Message, v interface{}) error {
	if v == nil {
		return nil
	}

	buf, err := c.Marshal(v)
	if err != nil {
		return err
	} else if len(buf) == 0 {
		return codec.ErrInvalidMessage
	}

	_, err = conn.Write(buf)
	return err
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
