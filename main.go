package main

import (
	"log"
	"os"

	"github.com/valyala/fasthttp"
)

type options struct {
	AccountSid string
	AuthToken  string
	Receiver   string
	Sender     string
}

func main() {
	opts := options{
		AccountSid: os.Getenv("AccountSid"),
		AuthToken:  os.Getenv("AuthToken"),
		Receiver:   os.Getenv("Receiver"),
		Sender:     os.Getenv("Sender"),
	}

	if opts.AccountSid == "" || opts.AuthToken == "" || opts.Sender == "" {
		log.Fatal("'SID', 'TOKEN' and 'SENDER' environment variables need to be set")
	}

	o := NewMOptionsWithHandler(&opts)
	err := fasthttp.ListenAndServe(":9090", o.HandleFastHTTP)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
