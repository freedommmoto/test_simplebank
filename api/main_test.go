package api

import (
	db "github.com/freedommmoto/test_simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, store db.Store) *Server {
	server := NewServer(store)
	return server
}
