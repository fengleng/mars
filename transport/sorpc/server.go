package sorpc

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"
)

// ErrServerClosed is returned by the Server's Serve, ListenAndServe after a call to Shutdown or Close.
var ErrServerClosed = errors.New("http: Server closed")

const (
	// ReaderBuffsize is used for bufio reader.
	ReaderBuffsize = 1024
	// WriterBuffsize is used for bufio writer.
	WriterBuffsize = 1024
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "rpcx context value " + k.name }

var (
	// RemoteConnContextKey is a context key. It can be used in
	// services with context.WithValue to access the connection arrived on.
	// The associated value will be of type net.Conn.
	RemoteConnContextKey = &contextKey{"remote-conn"}
	// StartRequestContextKey records the start time
	StartRequestContextKey = &contextKey{"start-parse-request"}
	// StartSendRequestContextKey records the start time
	StartSendRequestContextKey = &contextKey{"start-send-request"}
	// TagContextKey is used to record extra info in handling services. Its value is a map[string]interface{}
	TagContextKey = &contextKey{"service-tag"}
	// HttpConnContextKey is used to store http connection.
	HttpConnContextKey = &contextKey{"http-conn"}
)

// Server is rpcx server that use TCP or UDP.
type Server struct {
	ln                 net.Listener
	readTimeout        time.Duration
	writeTimeout       time.Duration
	gatewayHTTPServer  *http.Server
	DisableHTTPGateway bool // should disable http invoke or not.
	DisableJSONRPC     bool // should disable json rpc or not.

	serviceMapMu sync.RWMutex
	serviceMap   map[string]*service

	mu         sync.RWMutex
	activeConn map[net.Conn]struct{}
	doneChan   chan struct{}
	seq        uint64

	inShutdown int32
	onShutdown []func(s *Server)
	onRestart  []func(s *Server)

	// TLSConfig for creating tls tcp connection.
	tlsConfig *tls.Config
	// BlockCrypt for kcp.BlockCrypt
	options map[string]interface{}

	// CORS options
	//corsOptions *CORSOptions

	//Plugins PluginContainer

	// AuthFunc can be used to auth.
	//AuthFunc func(ctx context.Context, req *protocol.Message, token string) error

	handlerMsgNum int32

	HandleServiceError func(error)
}

// NewServer returns a server.
func NewServer() *Server {
	s := &Server{
		//Plugins:    &pluginContainer{},
		options:    make(map[string]interface{}),
		activeConn: make(map[net.Conn]struct{}),
		doneChan:   make(chan struct{}),
		serviceMap: make(map[string]*service),
	}

	//for _, op := range options {
	//	op(s)
	//}

	if s.options["TCPKeepAlivePeriod"] == nil {
		s.options["TCPKeepAlivePeriod"] = 3 * time.Minute
	}
	return s
}

// Address returns listened address.
func (s *Server) Address() net.Addr {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.ln == nil {
		return nil
	}
	return s.ln.Addr()
}

// ActiveClientConn returns active connections.
func (s *Server) ActiveClientConn() []net.Conn {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]net.Conn, 0, len(s.activeConn))
	for clientConn := range s.activeConn {
		result = append(result, clientConn)
	}
	return result
}

// SendMessage a request to the specified client.
// The client is designated by the conn.
// conn can be gotten from context in services:
//
//   ctx.Value(RemoteConnContextKey)
//
// servicePath, serviceMethod, metadata can be set to zero values.
func (s *Server) SendMessage(conn net.Conn, servicePath, serviceMethod string, metadata map[string]string, data []byte) error {
	//ctx := share.WithValue(context.Background(), StartSendRequestContextKey, time.Now().UnixNano())
	//s.Plugins.DoPreWriteRequest(ctx)
	//
	//req := protocol.GetPooledMsg()
	//req.SetMessageType(protocol.Request)
	//
	//seq := atomic.AddUint64(&s.seq, 1)
	//req.SetSeq(seq)
	//req.SetOneway(true)
	//req.SetSerializeType(protocol.SerializeNone)
	//req.ServicePath = servicePath
	//req.ServiceMethod = serviceMethod
	//req.Metadata = metadata
	//req.Payload = data
	//
	//b := req.EncodeSlicePointer()
	//_, err := conn.Write(*b)
	//protocol.PutData(b)
	//
	//s.Plugins.DoPostWriteRequest(ctx, req, err)
	//protocol.FreeMsg(req)
	return nil
}

func (s *Server) getDoneChan() <-chan struct{} {
	return s.doneChan
}

// startShutdownListener start a new goroutine to notify SIGTERM
// and SIGHUP signals and handle them gracefully
func (s *Server) startShutdownListener() {
	//go func(s *Server) {
	//	log.Info("server pid:", os.Getpid())
	//
	//	// channel to receive notifications of SIGTERM and SIGHUP
	//	ch := make(chan os.Signal, 1)
	//	signal.Notify(ch, syscall.SIGTERM, syscall.SIGHUP)
	//
	//	// custom functions to handle signal SIGTERM and SIGHUP
	//	var customFuncs []func(s *Server)
	//
	//	switch <-ch {
	//	case syscall.SIGTERM:
	//		customFuncs = append(s.onShutdown, func(s *Server) {
	//			s.Shutdown(context.Background())
	//		})
	//	case syscall.SIGHUP:
	//		customFuncs = append(s.onRestart, func(s *Server) {
	//			s.Restart(context.Background())
	//		})
	//	}
	//
	//	for _, fn := range customFuncs {
	//		fn(s)
	//	}
	//}(s)
}

// Serve starts and listens RPC requests.
// It is blocked until receiving connections from clients.
func (s *Server) Serve(network, address string) (err error) {
	s.startShutdownListener()
	//var ln net.Listener
	//ln, err = s.makeListener(network, address)
	//if err != nil {
	//	return
	//}
	//
	//if network == "http" {
	//	s.serveByHTTP(ln, "")
	//	return nil
	//}
	//
	//if network == "ws" || network == "wss" {
	//	s.serveByWS(ln, "")
	//	return nil
	//}
	//
	//// try to start gateway
	//ln = s.startGateway(network, ln)

	return s.serveListener(nil)
}

// ServeListener listens RPC requests.
// It is blocked until receiving connections from clients.
func (s *Server) ServeListener(network string, ln net.Listener) (err error) {
	//s.startShutdownListener()
	//if network == "http" {
	//	s.serveByHTTP(ln, "")
	//	return nil
	//}
	//
	//// try to start gateway
	//ln = s.startGateway(network, ln)
	//
	//return s.serveListener(ln)
	return nil
}

// serveListener accepts incoming connections on the Listener ln,
// creating a new service goroutine for each.
// The service goroutines read requests and then call services to reply to them.
func (s *Server) serveListener(ln net.Listener) error {

	return nil
}
