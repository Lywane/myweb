package http

import (
	"time"
	"strings"
	"fmt"
	"net/http"
)

func RecoveryHandler(c *Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Log("ERROR", "[panic]", err)
			c.DieWithHttpStatus(http.StatusInternalServerError)
		}
	}()
	c.Next()
}

func LogHandler(c *Context) {
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	ips := c.Request.Header.Get("X-Forwarded-For")
	ip := ""
	if ips != "" {
		ip = strings.Split(ips, ",")[0]
	}
	if ip == "" {
		ip = c.Request.Header.Get("X-Real-Ip")
	}
	if ip == "" {
		ip = c.Request.RemoteAddr
	}

	c.Next()

	traceId := c.Request.Header.Get("Kelp-Traceid")
	uuid := c.Request.Header.Get("uuid")
	end := time.Now()
	latency := end.Sub(start)
	method := c.Request.Method
	resp := string(c.responseData)
	if len(resp) > 500 {
		resp = fmt.Sprintf("response is too large (with %d bytes, head is %s)", len(resp), resp[0:100]+"...")
	}
	req := string(c.Body())
	if raw != "" {
		path = path + "?" + raw
	}

	log.Log(
		"REQ",
		ip, // remote ip
		end.Format("2006/01/02 15:04:05"),
		latency.Nanoseconds()/int64(time.Millisecond),
		str(method),
		str(path),
		str(traceId), // trace id
		str(uuid),    // uuid
		`"""`+str(req)+`"""`,
		`"""`+str(resp)+`"""`,
	)
}

func str(v string) string {
	if v == "" {
		return "-"
	}
	return v
}
