package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
	"github.com/julienschmidt/httprouter"
)

type Response struct {
	Ip          string
	Port        string
	ForwardedIp string
}

func writeIp(w http.ResponseWriter, res *Response) {
	json, _ := json.Marshal(res)
	log.Println("Received request: " + string(json))
	fmt.Fprintf(w, string(json))
}

func writeError(w http.ResponseWriter, err error) {
	fmt.Fprintf(w, "{ \"error\": \"%s\" }", err)
}

func GetIp(req *http.Request) (res *Response, err error) {
	ip, port, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		err = errors.New("Unable to parse ip address: " + ip)
		return
	}

	// This will only be defined when site is accessed via non-anonymous proxy
	// and takes precedence over RemoteAddr
	// Header.Get is case-insensitive
	forward := req.Header.Get("X-Forwarded-For")

	res = &Response{
		Ip:          ip,
		Port:        port,
		ForwardedIp: forward,
	}
	return
}

// https://blog.golang.org/context/userip/userip.go
func getIpEndpoint(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	res, err := GetIp(req)
	if err != nil {
		writeError(w, err)
		return
	}
	writeIp(w, res)
}

func main() {
	app := cli.NewApp()
	app.Name = "WhatIsMyIp"
	app.Usage = "./WhatIsMyIP [options]"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "ip", Value: "0.0.0.0", EnvVar: "WHATISMYIP_IP"},
		cli.IntFlag{Name: "port", Value: 80, EnvVar: "WHATISMYIP_PORT"},
	}
	app.Action = func(c *cli.Context) {
		ip := c.String("ip")
		port := c.String("port")

		r := httprouter.New()

		r.GET("/", getIpEndpoint)

		// Add a handler on /test
		r.GET("/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			// Simply write some test data for now
			fmt.Fprint(w, "OK\n")
		})

		l, err := net.Listen("tcp", ip+":"+port)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Started webserver on: " + "http://" + ip + ":" + port)

		// Start the blocking server loop.
		log.Fatal(http.Serve(l, r))
	}
	app.Run(os.Args)
}
