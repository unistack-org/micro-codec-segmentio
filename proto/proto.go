// Package proto provides a proto codec
package proto

import (
	"io"
	"io/ioutil"

	// nolint: staticcheck
	oldproto "github.com/golang/protobuf/proto"
	"github.com/segmentio/encoding/proto"
	"github.com/unistack-org/micro/v3/codec"
	newproto "google.golang.org/protobuf/proto"
)

type protoCodec struct{}

func (c *protoCodec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case nil:
		return nil, nil
	case *codec.Frame:
		return m.Data, nil
	case newproto.Message, oldproto.Message, proto.Message:
		return proto.Marshal(m)
	}
	return nil, codec.ErrInvalidMessage
}

func (c *protoCodec) Unmarshal(d []byte, v interface{}) error {
	var err error
	if d == nil {
		return nil
	}
	switch m := v.(type) {
	case nil:
		return nil
	case *codec.Frame:
		m.Data = d
	case newproto.Message, oldproto.Message, proto.Message:
		err = proto.Unmarshal(d, m)
	default:
		err = codec.ErrInvalidMessage
	}
	return err
}

func (c *protoCodec) ReadHeader(conn io.Reader, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *protoCodec) ReadBody(conn io.Reader, b interface{}) error {
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
	case oldproto.Message, newproto.Message, proto.Message:
		buf, err := ioutil.ReadAll(conn)
		if err != nil {
			return err
		}
		return proto.Unmarshal(buf, m)
	}
	return codec.ErrInvalidMessage
}

func (c *protoCodec) Write(conn io.Writer, m *codec.Message, b interface{}) error {
	switch m := b.(type) {
	case nil:
		return nil
	case *codec.Frame:
		_, err := conn.Write(m.Data)
		return err
	case oldproto.Message, newproto.Message, proto.Message:
		buf, err := proto.Marshal(m)
		if err != nil {
			return err
		}
		_, err = conn.Write(buf)
		return err
	}
	return codec.ErrInvalidMessage
}

func (c *protoCodec) String() string {
	return "proto"
}

func NewCodec() codec.Codec {
	return &protoCodec{}
}
