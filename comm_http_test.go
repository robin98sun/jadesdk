// reference: https://medium.com/rungo/unit-testing-made-easy-in-go-25077669318

package jadesdk

import (
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"
	"github.com/ant0ine/go-json-rest/rest"
)

// RESTful Server

type TestRequest struct {
	ReqID              string `json:"reqId,omitempty"`
	ExpectedRetryTimes int    `json:"expectedRetryTimes,omitempty"`
	ExpectingResponse  string `json:"expectingRes,omitempty"`
}

func setupHttpServer(logger func(...interface{}), fatal func(...interface{})) {
	// Read command line flags

	port := 3333
	addr := "127.0.0.1"

	servedRequests := make(map[string]int)

	routes := []*rest.Route{
		rest.Put("/$jade$/testPut", func(w rest.ResponseWriter, r *rest.Request) {
			// retryCntInHeader := r.Header.Get("retry-count")
			// logger("the request indicator of retry:", retryCntInHeader)
			req := &TestRequest{}
			_, err := decodeRawRequestBody(r, req)
			if err != nil {
				w.WriteJson(newErrorResponse(err.Error()))
				return
			}
			// check the retry times
			if _, e := servedRequests[req.ReqID]; !e {
				servedRequests[req.ReqID] = 0
			} else {
				servedRequests[req.ReqID] += 1
			}
			if servedRequests[req.ReqID] < req.ExpectedRetryTimes {
				w.WriteJson(newErrorResponse("not reaching the expected retry times"))
				return
			}
			if req.ExpectingResponse == "retry count" {
				w.WriteJson(newSuccessResponse(servedRequests[req.ReqID]))
			} else if req.ExpectingResponse == "remote addr" {
				w.WriteJson(newSuccessResponse(r.RemoteAddr))
			} else {
				w.WriteJson(newSuccessResponse("OK"))
			}
		}),
	}

	// Build the RESTful server
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(routes...)
	if err != nil {
		fatal(err)
	}
	api.SetApp(router)
	logger("plankton is listening on ", addr+":"+strconv.Itoa(port))
	fatal(http.ListenAndServe(addr+":"+strconv.Itoa(port), api.MakeHandler()))
}

func TestHttpComm_Init(t *testing.T) {
	go setupHttpServer(log.Println, log.Fatal)
}

func TestHttpComm_Basic(t *testing.T) {
	sdk := NewJadeSDK()
	sdk.Verbose(false)
	serverNode := &Node{
		Addr:     "127.0.0.1",
		Port:     3333,
		Protocol: "http",
	}
	res, bytes, err := sdk.HTTPCommunicate(
		"test http communication",
		"http", "put", "/$jade$/testPut",
		serverNode,
		&TestRequest{
			ReqID: "rid:" + strconv.FormatInt(int64(time.Now().UnixNano()), 10),
		},
		0, -1,
	)
	if err != nil {
		t.Error("server responded an error:", err)
	} else {
		t.Log("done, server response:", res, ", bytes:", bytes)
	}
}

func TestHttpComm_RetryExpectedTimes(t *testing.T) {
	sdk := NewJadeSDK()
	sdk.Verbose(false)
	serverNode := &Node{
		Addr:     "127.0.0.1",
		Port:     3333,
		Protocol: "http",
	}
	res, bytes, err := sdk.HTTPCommunicate(
		"test http communication",
		"http", "put", "/$jade$/testPut",
		serverNode,
		&TestRequest{
			ReqID:              "rid:" + strconv.FormatInt(int64(time.Now().UnixNano()), 10),
			ExpectedRetryTimes: 3,
			ExpectingResponse:  "retry count",
		},
		0, -1,
	)
	if err != nil {
		t.Error("server responded an error:", err)
	} else {
		if res != float64(3) {
			t.Error("unexpected retry count:", res, ", which should exactly be 10")
		} else {
			t.Log("retry count is as expected:", res, ", bytes:", bytes)
		}
	}
}

func TestHttpComm_KeepAlive(t *testing.T) {
	sdk := NewJadeSDK()
	sdk.Verbose(false)
	serverNode := &Node{
		Addr:     "127.0.0.1",
		Port:     3333,
		Protocol: "http",
	}

	var clientAddr interface{}
	for i := 0; i < 10; i++ {
		res, _, err := sdk.HTTPCommunicate(
			"test http communication",
			"http", "put", "/$jade$/testPut",
			serverNode,
			&TestRequest{
				ReqID:              "rid:" + strconv.FormatInt(int64(time.Now().UnixNano()), 10),
				ExpectedRetryTimes: 0,
				ExpectingResponse:  "remote addr",
			},
			0, -1, 
		)
		if err != nil {
			t.Error("server responded an error:", err)
		} else {
			if clientAddr == nil {
				clientAddr = res
			} else {
				if clientAddr != res {
					t.Error("the http server is not keeping alive")
				}
			}
		}
	}
	t.Log("keep-alive is ok")
}

