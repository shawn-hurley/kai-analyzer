package codec

import (
	"io"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/go-logr/logr"
)

// Connection wraps a pair of unidirectional streams as an io.ReadWriteCloser.
type Connection struct {
	Input  io.ReadCloser
	Output io.WriteCloser
}

// Read implements io.ReadWriteCloser.
func (c *Connection) Read(p []byte) (n int, err error) {
	return c.Input.Read(p)
}

// Write implements io.ReadWriteCloser.
func (c *Connection) Write(p []byte) (n int, err error) {
	return c.Output.Write(p)
}

// Close closes c's underlying ReadCloser and WriteCloser.
func (c *Connection) Close() error {
	rerr := c.Input.Close()
	werr := c.Output.Close()
	if rerr != nil {
		return rerr
	}
	return werr
}

type Codec interface {
	rpc.ClientCodec
	rpc.ServerCodec
}

type codec struct {
	rpc.ClientCodec
	rpc.ServerCodec
	logger logr.Logger
}

func (c *codec) ReadRequestHeader(r *rpc.Request) error {
	c.logger.V(3).Info("read request header", "request", r)
	err := c.ServerCodec.ReadRequestHeader(r)
	c.logger.V(3).Info("finished request header", "err", err)
	return err
}

func (c *codec) ReadRequestBody(r any) error {
	c.logger.V(3).Info("read request body", "request", r)
	err := c.ServerCodec.ReadRequestBody(r)
	c.logger.V(3).Info("finished request body", "err", err)
	return err
}

func (c *codec) WriteResponse(r *rpc.Response, v any) error {
	c.logger.V(3).Info("writing response", "response", r, "object", v)
	err := c.ServerCodec.WriteResponse(r, v)
	c.logger.V(3).Info("finished write response", "err", err)
	return err
}

func (c *codec) Close() error {
	err := c.ClientCodec.Close()
	if err != nil {
		return err
	}
	err = c.ServerCodec.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewCodec(connection Connection, logger logr.Logger) Codec {
	return &codec{
		ClientCodec: jsonrpc.NewClientCodec(&connection), ServerCodec: jsonrpc.NewServerCodec(&connection),
		logger: logger.WithName("json codec"),
	}
}
