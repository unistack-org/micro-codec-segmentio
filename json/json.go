// Package json provides a json codec
package json // import "go.unistack.org/micro-codec-segmentio/v3/json"

import (
	"io"

	"github.com/segmentio/encoding/json"
	pb "go.unistack.org/micro-proto/v3/codec"
	"go.unistack.org/micro/v3/codec"
	rutil "go.unistack.org/micro/v3/util/reflect"
)

var (
	DefaultMarshalOptions = JsonMarshalOptions{
		EscapeHTML:  true,
		SortMapKeys: true,
	}

	DefaultUnmarshalOptions = JsonUnmarshalOptions{
		ZeroCopy: true,
	}
)

var _ codec.Codec = &jsonCodec{}

const (
	flattenTag = "flatten"
)

type JsonMarshalOptions struct {
	EscapeHTML      bool
	SortMapKeys     bool
	TrustRawMessage bool
}

type JsonUnmarshalOptions struct {
	DisallowUnknownFields                bool
	DontCopyNumber                       bool
	DontCopyRawMessage                   bool
	DontCopyString                       bool
	DontMatchCaseInsensitiveStructFields bool
	UseNumber                            bool
	ZeroCopy                             bool
}

type jsonCodec struct {
	opts        codec.Options
	encodeFlags json.AppendFlags
	decodeFlags json.ParseFlags
}

func getMarshalFlags(o JsonMarshalOptions) json.AppendFlags {
	var encodeFlags json.AppendFlags

	if o.EscapeHTML {
		encodeFlags |= json.EscapeHTML
	}

	if o.SortMapKeys {
		encodeFlags |= json.SortMapKeys
	}

	if o.TrustRawMessage {
		encodeFlags |= json.TrustRawMessage
	}

	return encodeFlags
}

func getUnmarshalFlags(o JsonUnmarshalOptions) json.ParseFlags {
	var decodeFlags json.ParseFlags

	if o.DisallowUnknownFields {
		decodeFlags |= json.DisallowUnknownFields
	}

	if o.DontCopyNumber {
		decodeFlags |= json.DontCopyNumber
	}

	if o.DontCopyRawMessage {
		decodeFlags |= json.DontCopyRawMessage
	}

	if o.DontCopyString {
		decodeFlags |= json.DontCopyString
	}

	if o.DontMatchCaseInsensitiveStructFields {
		decodeFlags |= json.DontMatchCaseInsensitiveStructFields
	}

	if o.UseNumber {
		decodeFlags |= json.UseNumber
	}

	if o.ZeroCopy {
		decodeFlags |= json.ZeroCopy
	}

	return decodeFlags
}

func (c *jsonCodec) Marshal(v interface{}, opts ...codec.Option) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if nv, err := rutil.StructFieldByTag(v, options.TagName, flattenTag); err == nil {
		v = nv
	}

	switch m := v.(type) {
	case *codec.Frame:
		return m.Data, nil
	case *pb.Frame:
		return m.Data, nil
	}

	marshalOptions := DefaultMarshalOptions
	if options.Context != nil {
		if f, ok := options.Context.Value(marshalOptionsKey{}).(JsonMarshalOptions); ok {
			marshalOptions = f
		}
	}

	var err error
	var buf []byte
	buf, err = json.Append(buf, v, getMarshalFlags(marshalOptions))
	return buf, err
}

func (c *jsonCodec) Unmarshal(b []byte, v interface{}, opts ...codec.Option) error {
	if len(b) == 0 || v == nil {
		return nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if nv, err := rutil.StructFieldByTag(v, options.TagName, flattenTag); err == nil {
		v = nv
	}

	switch m := v.(type) {
	case *codec.Frame:
		m.Data = b
		return nil
	case *pb.Frame:
		m.Data = b
		return nil
	}

	unmarshalOptions := DefaultUnmarshalOptions
	if options.Context != nil {
		if f, ok := options.Context.Value(unmarshalOptionsKey{}).(JsonUnmarshalOptions); ok {
			unmarshalOptions = f
		}
	}

	_, err := json.Parse(b, v, getUnmarshalFlags(unmarshalOptions))
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

func NewCodec(opts ...codec.Option) *jsonCodec {
	return &jsonCodec{opts: codec.NewOptions(opts...)}
}
