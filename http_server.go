package jadesdk

import (
	// RESTful Server
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"

	// others
	// "encoding/json"
	"flag"
	"fmt"
	"strconv"
	// "strings"
)

func (j *JadeSDK) CreateHTTPServer() {
	// Read command line flags
	portArg := flag.Int("port", 8080, "port to listen")
	addrArg := flag.String("addr", "0.0.0.0", "ip address to listen")
	flag.Parse()

	port := *portArg
	addr := *addrArg
	// Read Env configurations
	j.ReadConfFromEnv()
	if j.Conf.SelfNode.Addr != "" {
		addr = j.Conf.SelfNode.Addr
	}
	if j.Conf.SelfNode.Port != 0 {
		port = j.Conf.SelfNode.Port
	}

	routes := []*rest.Route{}
	// JADE API
	routes = append(routes, j.createJadeRoutes()...)
	// Application API
	routes = append(routes, j.createAppRoutes()...)

	// Build the RESTful server
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(routes...)
	if err != nil {
		j.log.Fatal(err)
	}
	api.SetApp(router)
	fmt.Println("plankton is listening on port", port)
	j.log.Fatal(http.ListenAndServe(addr+":"+strconv.Itoa(port), api.MakeHandler()))
}
