package main

import (
	"fmt"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"github.com/valyala/fasthttp"
)

// OptionsWithHandler is a struct with a mux and shared credentials
type OptionsWithHandler struct {
	Options *options
	Client  *twilio.RestClient
}

// NewMOptionsWithHandler returns a OptionsWithHandler for http requests
// with shared credentials
func NewMOptionsWithHandler(o *options) OptionsWithHandler {
	return OptionsWithHandler{
		o,
		twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: o.AccountSid,
			Password: o.AuthToken,
		}),
	}
}

// HandleFastHTTP is the router function
func (m OptionsWithHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/":
		m.ping(ctx)
	case "/sms":
		m.sms(ctx)
	case "/call":
		m.call(ctx)
	case "/callandsms":
		m.call(ctx)
		m.sms(ctx)
	default:
		ctx.Error("Not found", fasthttp.StatusNotFound)
	}
}

func (m OptionsWithHandler) ping(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "ping")
}

func (m OptionsWithHandler) sms(ctx *fasthttp.RequestCtx) {
	if !ctx.IsPost() {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	} else {
		if string(ctx.Request.Header.Peek("Content-Type")) != "application/json" {
			ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		} else {
			params := &openapi.CreateMessageParams{}
			// TODO change from hard-coded reciver to reciver retrieved from WebHook
			params.SetTo("+zzxxxxxxxx")
			params.SetFrom(m.Options.Sender)
			// TODO change from hard-coded message to message retrieved from JSON WebHook
			params.SetBody("Hello from Go!")

			resp, err := m.Client.Api.CreateMessage(params)
			if err != nil {
				fmt.Println(err.Error())
				err = nil
			} else {
				fmt.Println("Message Sid: " + *resp.Sid)
			}
		}
	}
}

func (m OptionsWithHandler) call(ctx *fasthttp.RequestCtx) {
	if !ctx.IsPost() {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	} else {
		if string(ctx.Request.Header.Peek("Content-Type")) != "application/json" {
			ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		} else {
			params := &openapi.CreateCallParams{}
			// TODO change from hard-coded reciver to reciver retrieved from WebHook
			params.SetTo("+zzxxxxxxxx")
			params.SetFrom(m.Options.Sender)
			// TODO change from hard-coded message to message retrieved from JSON WebHook
			params.SetTwiml("<response><say>Hello there!</say></response>")

			resp, err := m.Client.Api.CreateCall(params)
			if err != nil {
				fmt.Println(err.Error())
				err = nil
			} else {
				fmt.Println("Call Status: " + *resp.Status)
				fmt.Println("Call Sid: " + *resp.Sid)
				fmt.Println("Call Direction: " + *resp.Direction)
			}
		}
	}
}
