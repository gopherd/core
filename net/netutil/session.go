package netutil

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gopherd/core/io/pagebuf"
	"github.com/gopherd/core/proto"
	"github.com/gopherd/core/text/resp"
	"github.com/gopherd/core/text/shell"
)

const (
	// max length of content: 1G
	MaxContentLength = 1 << 30

	hello = "hello"
)

var (
	ErrNotHandshaked  = errors.New("hello command required")
	ErrInvalidCommand = errors.New("invalid command")
)

// errno returns v's underlying uintptr, else 0.
//
// TODO: See comment in isClosedConnError.
func errno(v error) uintptr {
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.Uintptr {
		return uintptr(rv.Uint())
	}
	return 0
}

// isClosedConnError reports whether err is an error from use of a closed
// network connection.
func isClosedConnError(err error) bool {
	if err == nil {
		return false
	}

	// TODO: remove this string search and be more like the Windows
	// case below. That might involve modifying the standard library
	// to return better error types.
	str := err.Error()
	if strings.Contains(str, "use of closed network connection") {
		return true
	}

	// TODO(bradfitz): x/tools/cmd/bundle doesn't really support
	// build tags, so I can't make an http2_windows.go file with
	// Windows-specific stuff. Fix that and move this, once we
	// have a way to bundle this into std's net/http somehow.
	if runtime.GOOS == "windows" {
		if oe, ok := err.(*net.OpError); ok && oe.Op == "read" {
			if se, ok := oe.Err.(*os.SyscallError); ok && se.Syscall == "wsarecv" {
				const WSAECONNABORTED = 10053
				const WSAECONNRESET = 10054
				if n := errno(se.Err); n == WSAECONNRESET || n == WSAECONNABORTED {
					return true
				}
			}
		}
	}
	return false
}

// IsNetworkError returns whether the error is a network error or an EOF
func IsNetworkError(err error) bool {
	terr, ok := err.(net.Error)
	if err != io.EOF && err != io.ErrUnexpectedEOF && (!ok || !terr.Temporary()) && !isClosedConnError(err) {
		return false
	}
	return true
}

// timeoutReader wraps net.Conn as an io.Reader with timeout
type timeoutReader struct {
	conn    net.Conn
	timeout time.Duration
}

// Read implements io.Reader Read method
func (tr *timeoutReader) Read(p []byte) (n int, err error) {
	if tr.timeout > 0 {
		tr.conn.SetReadDeadline(time.Now().Add(tr.timeout))
	}
	return tr.conn.Read(p)
}

type reader struct {
	conn net.Conn
	bufr *bufio.Reader
	size int
}

func newReader(conn net.Conn, timeout time.Duration) *reader {
	return &reader{
		conn: conn,
		bufr: bufio.NewReader(&timeoutReader{
			conn:    conn,
			timeout: timeout,
		}),
		size: -1,
	}
}

// Len implements proto.Body Len method, -1 returned if no limit
func (b *reader) Len() int {
	return b.size
}

// Peek implements proto.Body Peek method
func (b *reader) Peek(n int) ([]byte, error) {
	if b.size >= 0 && b.size < n {
		return nil, io.EOF
	}
	return b.bufr.Peek(n)
}

// ReadByte implements io.ByteReader ReadByte method
func (b *reader) ReadByte() (c byte, err error) {
	if b.size == 0 {
		err = io.EOF
		return
	}
	c, err = b.bufr.ReadByte()
	if err == nil {
		if b.size > 0 {
			b.size--
		}
	}
	return
}

// Read implements io.Reader Read method
func (b *reader) Read(p []byte) (n int, err error) {
	if b.size == 0 {
		return 0, io.EOF
	}
	if b.size > 0 && len(p) > b.size {
		p = p[:b.size]
	}
	if len(p) > 0 {
		n, err = b.bufr.Read(p)
		if b.size > 0 {
			b.size -= n
		}
	}
	return
}

// Discard implements proto.Body Discard method
func (b *reader) Discard(n int) (discarded int, err error) {
	if b.size >= 0 && b.size < n {
		return 0, io.EOF
	}
	discarded, err = b.bufr.Discard(n)
	if b.size > 0 {
		b.size -= discarded
	}
	return
}

func (b *reader) discardAll() error {
	if b.size <= 0 {
		return nil
	}
	_, err := b.Discard(b.size)
	return err
}

// Option represents options of NewSession
type Option func(*option)

type option struct {
	timeout time.Duration
}

func defaultOption() option {
	return option{}
}

// WithTimeout specify read timeout of session
func WithTimeout(timeout time.Duration) Option {
	return func(opt *option) {
		opt.timeout = timeout
	}
}

// SessionEventHandler handles session events
type SessionEventHandler interface {
	OnOpen()                                // ready to read/write
	OnClose(err error)                      // session closed, err maybe nil
	OnHandshake(proto.ContentType) error    // session handshaked
	OnMessage(proto.Type, proto.Body) error // received a message
}

// Command represents textproto command
type Command interface {
	Name() string     // Name of command
	NArg() int        // Number of arguments
	Arg(i int) string // Arg returns ith argument
}

// CommandHandler handles textproto command
type CommandHandler interface {
	Commands() []string
	OnCommand(Command) error
}

// Session wraps network session
type Session struct {
	reader         *reader
	writer         *bufio.Writer
	handler        SessionEventHandler
	command        *resp.Command
	commandHandler CommandHandler

	// Handshake state
	handshaked  bool
	contentType proto.ContentType

	started  int32
	closed   int32
	wrunning int32

	errMu sync.RWMutex
	err   error

	mutex sync.Mutex
	cond  *sync.Cond
	pipe  *pagebuf.PageBuffer
	bufw  []byte
}

// NewSession creates a session
func NewSession(conn net.Conn, handler SessionEventHandler, options ...Option) *Session {
	var opt = defaultOption()
	for i := range options {
		options[i](&opt)
	}
	s := &Session{
		reader:  newReader(conn, opt.timeout),
		writer:  bufio.NewWriter(conn),
		handler: handler,
		pipe:    pagebuf.NewPageBuffer(),
	}
	if commandHandler, ok := handler.(CommandHandler); ok {
		s.commandHandler = commandHandler
	}
	s.cond = sync.NewCond(&s.mutex)
	s.bufw = make([]byte, s.pipe.PageSize())
	return s
}

// Conn returns the underlying connection
func (s *Session) Conn() net.Conn {
	return s.reader.conn
}

// ContentType returns type of content
func (s *Session) ContentType() proto.ContentType {
	return s.contentType
}

// Write implements io.Writer Write method, this IS NOT thread-safe.
func (s *Session) Write(p []byte) (n int, err error) {
	if s.IsClosed() {
		err = net.ErrClosed
		return
	}
	var (
		size         = len(p)
		maxWriteSize = s.pipe.PageSize() << 2
	)

	for n < size {
		end := n + maxWriteSize
		if end > size {
			end = size
		}
		s.mutex.Lock()
		var nn int
		nn, err = s.pipe.Write(p[n:end])
		buffered := s.pipe.Len()
		s.mutex.Unlock()
		n += nn
		if err != nil {
			return
		}
		if buffered == nn {
			s.cond.Signal()
		}
	}
	return
}

// Serve runs the read/write loops, it will block until the session closed
func (s *Session) Serve() bool {
	if !atomic.CompareAndSwapInt32(&s.started, 0, 1) {
		return false
	}
	var (
		readyWg sync.WaitGroup
		closeWg sync.WaitGroup
	)
	closeWg.Add(2)

	readyWg.Add(1)
	go s.writeLoop(&readyWg, &closeWg)
	readyWg.Wait()

	readyWg.Add(1)
	go s.readLoop(&readyWg, &closeWg)
	readyWg.Wait()

	s.handler.OnOpen()

	// Blcoking
	closeWg.Wait()

	s.errMu.RLock()
	err := s.err
	s.errMu.RUnlock()

	s.handler.OnClose(err)
	if err == nil {
		s.flush()
	}
	if proto.IsTextproto(s.contentType) {
		if err == nil || !IsNetworkError(err) {
			if err != nil {
				fmt.Fprintf(s.writer, "-connection closed because of %q\r\n", err.Error())
			} else {
				fmt.Fprintf(s.writer, "-connection closed\r\n")
			}
			s.writer.Flush()
		}
	}
	// close the underlying connection
	s.reader.conn.Close()

	return true
}

// IsClosed returns whether the session is closed
func (s *Session) IsClosed() bool {
	return atomic.LoadInt32(&s.closed) == 1
}

func (s *Session) setClosed(err error) {
	if err != nil {
		s.errMu.Lock()
		s.err = err
		s.errMu.Unlock()
	}
	atomic.StoreInt32(&s.closed, 1)
}

// Close closes the session
func (s *Session) Close(err error) {
	s.setClosed(err)
}

func (s *Session) readLoop(readyWg, closeWg *sync.WaitGroup) {
	readyWg.Done()
	for !s.IsClosed() {
		if err := s.underlyingRead(); err != nil {
			s.setClosed(err)
			break
		}
	}
	t := time.NewTicker(time.Millisecond)
	defer t.Stop()
	for range t.C {
		if atomic.LoadInt32(&s.wrunning) == 1 {
			s.cond.Signal()
		} else {
			break
		}
	}
	closeWg.Done()
}

func (s *Session) writeLoop(readyWg, closeWg *sync.WaitGroup) {
	atomic.StoreInt32(&s.wrunning, 1)
	readyWg.Done()
	for !s.IsClosed() {
		s.cond.L.Lock()
		for s.pipe.Len() == 0 && !s.IsClosed() {
			s.cond.Wait()
		}
		s.cond.L.Unlock()
		if s.IsClosed() {
			break
		}
		s.flush()
	}
	atomic.StoreInt32(&s.wrunning, 0)
	s.flush()
	closeWg.Done()
}

func (s *Session) flush() {
	written := 0
	for {
		s.cond.L.Lock()
		n, _ := s.pipe.Read(s.bufw)
		s.cond.L.Unlock()
		if n == 0 {
			break
		}
		written += n
		s.underlyingWrite(s.bufw[:n])
	}
	if written > 0 {
		s.writer.Flush()
	}
}

func (s *Session) underlyingWrite(p []byte) error {
	_, err := s.writer.Write(p)
	if err != nil {
		s.setClosed(err)
	}
	return err
}

func (s *Session) underlyingRead() error {
	// read type of message body
	s.reader.size = -1
	n, typ, err := proto.PeekType(s.reader)
	if err != nil {
		return err
	}
	if !s.handshaked {
		if err := s.handshake(typ); err != nil {
			msg := err.Error()
			var buf = make([]byte, 0, len(msg)+3)
			buf = append(buf, resp.ErrorType.Byte())
			buf = append(buf, msg...)
			buf = append(buf, '\r', '\n')
			s.Write(buf)
			return err
		}
		return nil
	}
	// If the session using textproto
	if proto.IsTextproto(s.contentType) && s.commandHandler != nil {
		if err := s.readCommand(); err != nil {
			return err
		}
		if s.command.Is(hello) {
			_, err := s.Write([]byte("-don't hello again\r\n"))
			return err
		}
		return s.commandHandler.OnCommand(s.command)
	}

	// discard read type
	if _, err := s.reader.Discard(n); err != nil {
		return err
	}
	// handle the message body
	size, err := proto.ReadSize(s.reader)
	if err != nil {
		return err
	}
	s.reader.size = size
	if err := s.handler.OnMessage(typ, s.reader); err != nil {
		return err
	}
	// discard unread bytes
	return s.reader.discardAll()
}

func (s *Session) readCommand() error {
	if s.command == nil {
		s.command = resp.NewCommand()
	} else {
		s.command.Request.Reset()
	}

	if err := s.command.Request.ReadFrom(s.reader.bufr); err != nil {
		return err
	}
	switch s.command.Request.Type {
	case resp.StringType:
		args, err := shell.Split(string(s.command.Request.Value()))
		if err != nil {
			return err
		}
		s.command.Request.SetArray(len(args))
		for i := range args {
			s.command.Request.Append(resp.StringType, []byte(args[i]))
		}
	case resp.ArrayType:
	default:
		return ErrInvalidCommand
	}
	if len(s.command.Request.Elements()) == 0 {
		return resp.ErrNumberOfArguments
	}
	return nil
}

// handshake message format:
//
// +hello <contentType>\r\n
//
// If everything is ok, "+hello <contentType>\r\n" responded.
// Otherwise, "-Error message\r\n" responded.
//
// example:
//
//	$ telnet 127.0.0.1 11001
//	Trying 127.0.0.1...
//	Connected to localhost.
//	Escape character is '^]'.
//	+hello 1
//	+hello 1
//
//	$ telnet 127.0.0.1 11001
//	Trying 127.0.0.1...
//	Connected to localhost.
//	Escape character is '^]'.
//	+hello
//	-not handshaked
func (s *Session) handshake(typ proto.Type) error {
	if err := s.readCommand(); err != nil {
		return err
	}
	name := strings.ToLower(s.command.Name())
	if name == "command" {
		// ignore "command" command before handshaking
		_, err := s.Write([]byte("+ignored\r\n"))
		return err
	}
	if name != hello {
		return ErrNotHandshaked
	}
	if s.command.NArg() != 1 {
		return resp.ErrNumberOfArguments
	}

	var contentType proto.ContentType
	t, err := strconv.Atoi(s.command.Arg(0))
	if err != nil {
		return err
	}
	contentType = proto.ContentType(t)
	if err := s.handler.OnHandshake(contentType); err != nil {
		return err
	}
	s.handshaked = true
	s.contentType = contentType
	var buf = make([]byte, 0, 16)
	buf = append(buf, resp.StringType.Byte())
	buf = append(buf, hello...)
	buf = append(buf, ' ')
	buf = strconv.AppendInt(buf, int64(contentType), 10)
	buf = append(buf, '\r', '\n')
	_, err = s.Write(buf)
	return err
}
