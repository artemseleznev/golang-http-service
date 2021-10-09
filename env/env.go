package env

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	DefaultHttpPort   = "8080"
	DefaultReqTimeout = 1000 * time.Millisecond
	DefaultMaxConns   = 100
)

type Env struct {
	HttpPort   string
	ReqTimeout time.Duration
	MaxConns   int
}

func InitEnv() (*Env, error) {
	env := new(Env)
	env.HttpPort = os.Getenv("HTTP_PORT")
	if env.HttpPort == "" {
		env.HttpPort = DefaultHttpPort
	}

	requestTimeoutEnv := os.Getenv("REQUEST_TIMEOUT_IN_MILLISECONDS")
	if requestTimeoutEnv == "" {
		env.ReqTimeout = DefaultReqTimeout
	} else {
		reqTimeoutInt, err := strconv.Atoi(requestTimeoutEnv)
		if err != nil {
			return nil, fmt.Errorf("could not parse REQUEST_TIMEOUT_IN_MILLISECONDS env: %w", err)
		}
		env.ReqTimeout = time.Duration(reqTimeoutInt) * time.Millisecond
	}

	maxConns := os.Getenv("MAX_CONNECTIONS")
	if maxConns == "" {
		env.MaxConns = DefaultMaxConns
	} else {
		maxConnsInt, err := strconv.Atoi(maxConns)
		if err != nil {
			return nil, fmt.Errorf("could not parse MAX_CONNECTIONS env: %w", err)
		}
		env.MaxConns = maxConnsInt
	}
	log.Printf("got envs: %+v", env)
	return env, nil
}
