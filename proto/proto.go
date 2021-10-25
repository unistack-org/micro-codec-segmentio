// Package proto provides a proto codec
package proto

import (
	"io"

	"github.com/segmentio/encoding/proto"
	"go.unistack.org/micro/v3/codec"
	rutil "go.unistack.org/micro/v3/util/reflect"
	newproto "google.golang.org/protobuf/proto"
)

type protoCodec struct {
	opts codec.Options
}

var _ codec.Codec = &protoCodec{}

const (
	flattenTag = "flatten"
)

func (c *protoCodec) Marshal(v interface{}, opts ...codec.Option) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if nv, nerr := rutil.StructFieldByTag(v, options.TagName, flattenTag); nerr == nil {
		v = nv
	}

	if m, ok := v.(*codec.Frame); ok {
		return m.Data, nil
	}

	switch m := v.(type) {
	case proto.Message:
		return proto.Marshal(m)
	case newproto.Message:
		return proto.Marshal(m)
	}

	return nil, codec.ErrInvalidMessage
}

func (c *protoCodec) Unmarshal(d []byte, v interface{}, opts ...codec.Option) error {
	if v == nil || len(d) == 0 {
		return nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if nv, nerr := rutil.StructFieldByTag(v, options.TagName, flattenTag); nerr == nil {
		v = nv
	}

	if m, ok := v.(*codec.Frame); ok {
		m.Data = d
		return nil
	}

	switch m := v.(type) {
	case proto.Message:
		return proto.Unmarshal(d, m)
	case newproto.Message:
		return proto.Unmarshal(d, m)

	}

	return codec.ErrInvalidMessage
}

func (c *protoCodec) ReadHeader(conn io.Reader, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *protoCodec) ReadBody(conn io.Reader, v interface{}) error {
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

func (c *protoCodec) Write(conn io.Writer, m *codec.Message, v interface{}) error {
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

func (c *protoCodec) String() string {
	return "proto"
}

func NewCodec(opts ...codec.Option) *protoCodec {
	return &protoCodec{opts: codec.NewOptions(opts...)}
}
