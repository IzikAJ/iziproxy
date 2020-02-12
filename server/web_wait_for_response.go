package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/izikaj/iziproxy/shared"
)

type waitForResponseParams struct {
	core    *Server
	req     *shared.Request
	signal  *CodeSignal
	w       *http.ResponseWriter
	timeout time.Duration
}

func (server *commonWebHelpers) waitForResponse(params waitForResponseParams) (err error) {
	select {
	case <-*params.signal:
		server.processResponse(params)

	case <-time.After(params.timeout):
		params.core.Stats.timeout()
		server.writeFailResponse(params.w, http.StatusGatewayTimeout, "TIMEOUT ERROR")
	}
	return
}

func (server *commonWebHelpers) processResponse(params waitForResponseParams) (err error) {
	params.core.Lock()
	d, ok := params.core.pool[params.req.ID]
	params.core.Unlock()

	if ok {
		resp := d.Response

		if resp.Status == 0 {
			params.core.Stats.fail()
			writeFailResponse(params.w, http.StatusBadGateway, "EMPTY RESPONSE FROM CLIENT")
			return
		}
		fmt.Printf("> [%d] %s\n", resp.Status, (*d).Request.Path)

		w := *params.w
		for _, header := range resp.Headers {
			for _, value := range header.Value {
				w.Header().Set(header.Name, value)
			}
		}

		w.WriteHeader(resp.Status)
		w.Write(resp.Body)
		params.core.Lock()
		delete(params.core.pool, params.req.ID)
		params.core.Unlock()
		params.core.Stats.complete()
	} else {
		params.core.Stats.fail()
		writeFailResponse(params.w, http.StatusBadGateway, "NO RESPONSE FROM CLIENT")
	}
	return
}
