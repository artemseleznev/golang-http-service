package server

import (
	"time"
)

type Server struct {
	reqTimeout time.Duration
}

func New(reqTimeout time.Duration) *Server {
	return &Server{reqTimeout: reqTimeout}
}
