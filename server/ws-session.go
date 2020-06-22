package main

import (
	"context"
	"io/ioutil"
	"net"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	demo "github.com/rgraphql/magellan-soyuz-demo/pb"
	genresolvers "github.com/rgraphql/magellan-soyuz-demo/server/resolve"
	"github.com/rgraphql/magellan-soyuz-demo/server/simple"
	"github.com/rgraphql/magellan/resolver"
	nserver "github.com/rgraphql/magellan/server"
	"github.com/rgraphql/rgraphql"
	"github.com/sirupsen/logrus"
)

// maxMessageSize is the maximum message size
const maxMessageSize = 2e6

// Session is a websocket session.
type Session struct {
	// s is the server
	s *Server
	// ctx is the session context
	ctx context.Context
	// ctxCancel cancels the session context
	ctxCancel context.CancelFunc
	// le is the log entry
	le *logrus.Entry
	// conn is the websocket conn
	conn net.Conn
	// magellanSession is the magellan session
	magellanSession *nserver.Session
	// sendCh is the outgoing channel
	sendCh <-chan *rgraphql.RGQLServerMessage

	// mtx guards the below fields
	mtx sync.Mutex
	// closed indicates the session is closed
	closed bool
}

// NewSession constructs a new session.
func NewSession(
	ctx context.Context,
	le *logrus.Entry,
	s *Server,
	c net.Conn,
) *Session {
	// Read "hello" packet w/ deadline.
	sess := &Session{
		s:    s,
		le:   le,
		conn: c,
	}
	sess.ctx, sess.ctxCancel = context.WithCancel(ctx)
	sendCh := make(chan *rgraphql.RGQLServerMessage, 10)
	rootRes := &simple.RootResolver{}
	sess.magellanSession = s.magellanServer.BuildSession(
		ctx,
		sendCh,
		func(r *resolver.Context) {
			genresolvers.ResolveRootQuery(r, rootRes)
		},
	)
	sess.sendCh = sendCh
	return sess
}

// Execute is the main session management routine.
func (s *Session) Execute() error {
	defer func() {
		s.Close()
		_ = s.conn.Close()
	}()

	ctx := s.ctx
	le := s.le

	go s.executeKeepAlive(ctx)
	go s.executeSendCh(ctx)

	for {
		h, r, err := wsutil.NextReader(s.conn, ws.StateServerSide)
		if err != nil {
			return err
		}
		if h.OpCode.IsControl() {
			return wsutil.ControlFrameHandler(s.conn, ws.StateServerSide)(h, r)
		}

		data, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		/*
			if mt != websocket.MessageBinary {
				return 0, nil, errors.Errorf("invalid websocket message type %d", mt)
			}
		*/
		if len(data) > maxMessageSize {
			return errors.Errorf("message size %d > maximum %d bytes", len(data), maxMessageSize)
		}

		if err := s.handleIncMessage(data); err != nil {
			le.WithError(err).Warn("error handling incoming message")
		}

		select {
		case <-ctx.Done():
			return context.Canceled
		default:
			continue
		}
	}
}

// executeSendCh executes the send channel.
func (s *Session) executeSendCh(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case sg := <-s.sendCh:
			if err := s.writeRGQLMessage(sg); err != nil {
				s.le.WithError(err).Warn("dropped outgoing message")
			}
		}
	}
}

// executeKeepAlive manages the keep-alive timer.
func (s *Session) executeKeepAlive(ctx context.Context) {
	pingTicker := time.NewTicker(time.Second * 5)
	defer pingTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-pingTicker.C:
			s.writePing()
		}
	}
}

// handleIncMessage handles an incoming message.
func (s *Session) handleIncMessage(data []byte) error {
	incMsg := &demo.RPCMessage{}
	if err := proto.Unmarshal(data, incMsg); err != nil {
		return err
	}
	return s.handleRPCMessage(incMsg)
}

// handleRPCMessage handles an incoming rpc message.
func (s *Session) handleRPCMessage(msg *demo.RPCMessage) error {
	s.le.WithField("rpc-type", msg.GetRpcId().String()).Debug("rx rpc message")
	switch msg.GetRpcId() {
	case demo.RPC_RPC_RGQLClientMessage:
		s.le.Debug("handling rgql client message")
		s.magellanSession.HandleMessage(msg.GetRgqlClientMessage())
	default:
	}
	return nil
}

// writePing writes a ping message to the session
func (s *Session) writePing() {
	s.writeRGQLMessage(&rgraphql.RGQLServerMessage{})
}

// writeMessage writes a message to the stream.
func (s *Session) writeMessage(data []byte) error {
	wt := wsutil.NewWriter(s.conn, ws.StateServerSide, ws.OpBinary)
	_, err := wt.Write(data)
	if err != nil {
		return err
	}

	return wt.Flush()
}

// writeRGQLMessage writes a rgraphql message
func (s *Session) writeRGQLMessage(msg *rgraphql.RGQLServerMessage) error {
	data, err := proto.Marshal(&demo.RPCMessage{
		RpcId:             demo.RPC_RPC_RGQLServerMessage,
		RgqlServerMessage: msg,
	})
	if err != nil {
		return err
	}
	return s.writeMessage(data)
}

// Close closes the session.
func (s *Session) Close() {
	s.ctxCancel()
	s.conn.Close() // called multiple times
}
