package main

import (
	"context"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/rgraphql/magellan/schema"
	nserver "github.com/rgraphql/magellan/server"
	"github.com/sirupsen/logrus"
)

// Server implements a magellan-based websocket server.
type Server struct {
	ctx context.Context
	le  *logrus.Entry
	// magellanServer is the magellan server
	magellanServer *nserver.Server
}

// NewServer builds a new server.
func NewServer(ctx context.Context, le *logrus.Entry, schema *schema.Schema) *Server {
	magellanServer := nserver.NewServer(schema)
	return &Server{le: le, ctx: ctx, magellanServer: magellanServer}
}

// ServeHTTP serves a websocket session over http.
func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	le := s.le
	ctx := s.ctx
	conn, _, _, err := ws.UpgradeHTTP(req, rw)
	if err != nil {
		// If a basic tcp connection is made, this error is returned.
		// Prevent infrastructure DOS by TCP connections by ignoring this.
		isInvalidPacket := err == ws.ErrProtocolMaskRequired
		if !isInvalidPacket && !isNormalCloseError(err) {
			le.WithError(err).Warn("error upgrading websocket conn")
		}
		return
	}

	sessCtx, sessCtxCancel := context.WithCancel(ctx)
	defer sessCtxCancel()

	le.Debug("session started")
	sess := NewSession(
		sessCtx,
		le,
		s,
		conn,
	)
	if err := sess.Execute(); err != nil {
		if isNormalCloseError(err) {
			le.Debug("session exited")
		} else {
			le.WithError(err).Warn("session exited with error")
		}
	}
}

// _ is a type assertion
var _ http.Handler = ((*Server)(nil))
