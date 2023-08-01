package server

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func (server *Server) InitTestRestAPI() {
	server.ECHO_SERVER.GET("/tests/domain/reachibility", server.TestDomainReachibility)
	server.ECHO_SERVER.GET("/tests/domain/ping", server.TestDomainPing)
}

// GET /tests/domain/reachibility?domain=example.com
func (server *Server) TestDomainReachibility(c echo.Context) error {
	domain := c.QueryParam("domain")
	if domain == "" {
		return c.JSON(400, map[string]interface{}{
			"message": "domain query parameter is required",
		})
	}
	// call domain:3333 /ping endpoint
	// set timeout to 5 seconds
	http_client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := http_client.Get("http://" + domain + ":"+strconv.Itoa(server.PORT)+"/tests/domain/ping")
	if err != nil {
		return c.JSON(500, map[string]interface{}{
			"message": "domain is not reachable",
		})
	}
	defer resp.Body.Close()
	// if response is 200, return true
	if resp.StatusCode == 200 {
		resp_body_byte, err := io.ReadAll(resp.Body)
		if err == nil && string(resp_body_byte) == "pong" {
			return c.JSON(200, map[string]interface{}{
				"message": "domain is reachable",
			})
		}
	}
	// else return false
	return c.JSON(500, map[string]interface{}{
		"message": "domain is not reachable",
	})
}

// GET /tests/domain/ping
func (server *Server) TestDomainPing(c echo.Context) error {
	return c.String(200, "pong")
}

