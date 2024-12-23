// Package proto provides a proto codec
package proto // import "go.unistack.org/micro-codec-segmentio/v3/proto"

import (
	"github.com/segmentio/encoding/proto"
	pb "go.unistack.org/micro-proto/v3/codec"
	"go.unistack.org/micro/v3/codec"
	rutil "go.unistack.org/micro/v3/util/reflect"
	newproto "google.golang.org/protobuf/proto"
)

type protoCodec struct {
	opts codec.Options
}

var _ codec.Codec = &protoCodec{}

func (c *protoCodec) Marshal(v interface{}, opts ...codec.Option) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, nerr := rutil.StructFieldByTag(v, options.TagName, "flatten"); nerr == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		return m.Data, nil
	case *pb.Frame:
		return m.Data, nil
	}

	switch m := v.(type) {
	case proto.Message:
		return proto.Marshal(m)
	case newproto.Message:
		return proto.Marshal(m)
	case codec.RawMessage:
		return []byte(m), nil
	case *codec.RawMessage:
		return []byte(*m), nil
	default:
		return nil, codec.ErrInvalidMessage
	}
}

func (c *protoCodec) Unmarshal(d []byte, v interface{}, opts ...codec.Option) error {
	if v == nil || len(d) == 0 {
		return nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, nerr := rutil.StructFieldByTag(v, options.TagName, "flatten"); nerr == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		m.Data = d
		return nil
	case *pb.Frame:
		m.Data = d
		return nil
	case *codec.RawMessage:
		*m = append((*m)[0:0], d...)
		return nil
	case codec.RawMessage:
		copy(m, d)
		return nil
	}

	switch m := v.(type) {
	case proto.Message:
		return proto.Unmarshal(d, m)
	case newproto.Message:
		return proto.Unmarshal(d, m)
	default:
		return codec.ErrInvalidMessage
	}
}

func (c *protoCodec) String() string {
	return "proto"
}

func NewCodec(opts ...codec.Option) *protoCodec {
	return &protoCodec{opts: codec.NewOptions(opts...)}
}
