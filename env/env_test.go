package env_test

import (
	"http-service-example/env"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestInitEnv_Success(t *testing.T) {
	testData := []struct {
		testName string
		envs     map[string]string
		res      *env.Env
	}{
		{
			testName: "nothing is set",
			envs:     nil,
			res: &env.Env{
				HttpPort:   env.DefaultHttpPort,
				ReqTimeout: env.DefaultReqTimeout,
				MaxConns:   env.DefaultMaxConns,
			},
		},
		{
			testName: "all is set",
			envs: map[string]string{
				"HTTP_PORT":                       "8080",
				"REQUEST_TIMEOUT_IN_MILLISECONDS": "222",
				"MAX_CONNECTIONS":                 "5",
			},
			res: &env.Env{
				HttpPort:   "8080",
				ReqTimeout: 222 * time.Millisecond,
				MaxConns:   5,
			},
		},
	}
	for _, td := range testData {
		td := td
		t.Run(td.testName, func(t *testing.T) {
			for k, v := range td.envs {
				if err := os.Setenv(k, v); err != nil {
					t.Fatal(err)
				}
			}
			res, err := env.InitEnv()
			if err != nil {
				t.Fatalf("expected nil but got error (%v)", err)
			}
			if !reflect.DeepEqual(td.res, res) {
				t.Fatalf("expected %+v but got %+v", td.res, res)
			}
		})
	}
}

func TestInitEnv_Error(t *testing.T) {
	testData := []struct {
		testName string
		envs     map[string]string
		errMsg   string
	}{
		{
			testName: "REQUEST_TIMEOUT_IN_MILLISECONDS is invalid",
			envs: map[string]string{
				"HTTP_PORT":                       "8080",
				"REQUEST_TIMEOUT_IN_MILLISECONDS": "invalid_value",
				"MAX_CONNECTIONS":                 "5",
			},
			errMsg: "REQUEST_TIMEOUT_IN_MILLISECONDS invalid",
		},
		{
			testName: "MAX_CONNECTIONS is invalid",
			envs: map[string]string{
				"HTTP_PORT":                       "8080",
				"REQUEST_TIMEOUT_IN_MILLISECONDS": "222",
				"MAX_CONNECTIONS":                 "invalid_value",
			},
			errMsg: "MAX_CONNECTIONS invalid",
		},
	}
	for _, td := range testData {
		td := td
		t.Run(td.testName, func(t *testing.T) {
			for k, v := range td.envs {
				if err := os.Setenv(k, v); err != nil {
					t.Fatal(err)
				}
			}
			res, err := env.InitEnv()
			if res != nil {
				t.Fatalf("expected nil but got %+v", res)
			}
			if err == nil {
				t.Fatal("expected error but got nil")
			}
		})
	}
}
