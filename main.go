package websocket

import (
	"errors"
	"io"
	"sync"

	"github.com/fibergo/fastws"
	"github.com/gofiber/fiber"
)

// Config ...
type Config struct {
	// Protocols are the supported protocols.
	Protocols []string

	// Origin is used to limit the clients coming from the defined origin
	Origin string

	// Compress defines whether using compression or not.
	// TODO
	Compress bool
}

func Upgrade(handler func(*Conn), config ...Config) func(*fiber.Ctx) {
	// Init config
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}
	if len(cfg.Protocols) == 0 {
		cfg.Protocols = []string{""}
	}
	if cfg.Origin == "" {
		cfg.Origin = ""
	}
	upgrader := &fastws.Upgrader{
		Handler: func(fconn *fastws.Conn) {
			conn := acquireConn(fconn)
			handler(conn)
			releaseConn(conn)
		},
	}
	return func(ctx *fiber.Ctx) {
		upgrader.Upgrade(ctx.Fasthttp)
	}
}

func ReleaseFrame(fr *Frame) {
	// release fastws
	fastws.ReleaseFrame(fr.Frame)
	// release fiber frame
	releaseFrame(fr)
}

type Frame struct {
	*fastws.Frame
}

func (fr *Frame) CopyTo(fr2 *Frame) {
	fr.Frame.CopyTo(fr2.Frame)
}

type Conn struct {
	*fastws.Conn
}

func (c *Conn) NextFrame() (fr *Frame, err error) {
	ffr, err := c.Conn.NextFrame()
	fr = acquireFrame(ffr)
	return fr, err
}
func (c *Conn) ReadFrame(fr *Frame) (nn int, err error) {
	return c.Conn.ReadFrame(fr.Frame)
}
func (c *Conn) ReadFull(b []byte, fr *Frame) ([]byte, error) {
	return c.Conn.ReadFull(b, fr.Frame)
}
func (c *Conn) ReplyClose(fr *Frame) (err error) {
	return c.Conn.ReplyClose(fr.Frame)
}
func (c *Conn) WriteFrame(fr *Frame) (int, error) {
	return c.Conn.WriteFrame(fr.Frame)
}

// Frame pool
var poolFrame = sync.Pool{
	New: func() interface{} {
		return new(Frame)
	},
}

// Acquire Frame from pool
func acquireFrame(ffr *fastws.Frame) *Frame {
	fr := poolFrame.Get().(*Frame)
	fr.Frame = ffr
	return fr
}

// Return Frame to pool
func releaseFrame(fr *Frame) {
	fr.Frame = nil
	poolFrame.Put(fr)
}

// Conn pool
var poolConn = sync.Pool{
	New: func() interface{} {
		return new(Conn)
	},
}

// Acquire Conn from pool
func acquireConn(fconn *fastws.Conn) *Conn {
	conn := poolConn.Get().(*Conn)
	conn.Conn = fconn
	return conn
}

// Return Conn to pool
func releaseConn(conn *Conn) {
	conn.Conn = nil
	poolConn.Put(conn)
}

type StatusCode uint16

const (
	// StatusNone is used to let the peer know nothing happened.
	StatusNone StatusCode = 1000
	// StatusGoAway peer's error.
	StatusGoAway = 1001
	// StatusProtocolError problem with the peer's way to communicate.
	StatusProtocolError = 1002
	// StatusNotAcceptable when a request is not acceptable
	StatusNotAcceptable = 1003
	// StatusReserved when a reserved field have been used
	StatusReserved = 1004
	// StatusNotConsistent IDK
	StatusNotConsistent = 1007
	// StatusViolation a violation of the protocol happened
	StatusViolation = 1008
	// StatusTooBig payload bigger than expected
	StatusTooBig = 1009
	// StatuseExtensionsNeeded IDK
	StatuseExtensionsNeeded = 1010
	// StatusUnexpected IDK
	StatusUnexpected = 1011
)

type Code uint8

const (
	// CodeContinuation defines the continuation code
	CodeContinuation Code = 0x0
	// CodeText defines the text code
	CodeText Code = 0x1
	// CodeBinary defines the binary code
	CodeBinary Code = 0x2
	// CodeClose defines the close code
	CodeClose Code = 0x8
	// CodePing defines the ping code
	CodePing Code = 0x9
	// CodePong defines the pong code
	CodePong Code = 0xA
)

type Mode uint8

const (
	// ModeText defines to use a text mode
	ModeText Mode = iota
	// ModeBinary defines to use a binary mode
	ModeBinary
)

var (
	// EOF represents an io.EOF error.
	EOF = io.EOF
)

var (
	// ErrCannotUpgrade shows up when an error ocurred when upgrading a connection.
	ErrCannotUpgrade = errors.New("cannot upgrade connection")
)
