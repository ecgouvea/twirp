// Code generated by protoc-gen-twirp v5.3.0, DO NOT EDIT.
// source: service.proto

/*
Package example is a generated twirp stub package.
This code was generated with github.com/twitchtv/twirp/protoc-gen-twirp v5.3.0.

It is generated from these files:
	service.proto
*/
package example

import bytes "bytes"
import strings "strings"
import context "context"
import fmt "fmt"
import ioutil "io/ioutil"
import http "net/http"

import jsonpb "github.com/golang/protobuf/jsonpb"
import proto "github.com/golang/protobuf/proto"
import twirp "github.com/twitchtv/twirp"
import ctxsetters "github.com/twitchtv/twirp/ctxsetters"

// Imports only used by utility functions:
import io "io"
import strconv "strconv"
import json "encoding/json"
import url "net/url"
import bufio "bufio"
import binary "encoding/binary"

// =====================
// Haberdasher Interface
// =====================

// A Haberdasher makes hats for clients.
type Haberdasher interface {
	// MakeHat produces a hat of mysterious, randomly-selected color!
	MakeHat(ctx context.Context, in *Size) (*Hat, error)

	MakeHats(ctx context.Context, in *MakeHatsReq) (HatStream, error)
}

// ===========================
// Haberdasher Protobuf Client
// ===========================

type haberdasherProtobufClient struct {
	client HTTPClient
	urls   [2]string
}

// NewHaberdasherProtobufClient creates a Protobuf client that implements the Haberdasher interface.
// It communicates using Protobuf and can be configured with a custom HTTPClient.
func NewHaberdasherProtobufClient(addr string, client HTTPClient) Haberdasher {
	prefix := urlBase(addr) + HaberdasherPathPrefix
	urls := [2]string{
		prefix + "MakeHat",
		prefix + "MakeHats",
	}
	if httpClient, ok := client.(*http.Client); ok {
		return &haberdasherProtobufClient{
			client: withoutRedirects(httpClient),
			urls:   urls,
		}
	}
	return &haberdasherProtobufClient{
		client: client,
		urls:   urls,
	}
}

func (c *haberdasherProtobufClient) MakeHat(ctx context.Context, in *Size) (*Hat, error) {
	ctx = ctxsetters.WithPackageName(ctx, "twitch.twirp.example")
	ctx = ctxsetters.WithServiceName(ctx, "Haberdasher")
	ctx = ctxsetters.WithMethodName(ctx, "MakeHat")
	out := new(Hat)
	err := doProtobufRequest(ctx, c.client, c.urls[0], in, out)
	return out, err
}

func (c *haberdasherProtobufClient) MakeHats(ctx context.Context, in *MakeHatsReq) (HatStream, error) {
	ctx = ctxsetters.WithPackageName(ctx, "twitch.twirp.example")
	ctx = ctxsetters.WithServiceName(ctx, "Haberdasher")
	ctx = ctxsetters.WithMethodName(ctx, "MakeHats")
	reqBodyBytes, err := proto.Marshal(in)
	if err != nil {
		return nil, clientError("failed to marshal proto request", err)
	}
	reqBody := bytes.NewBuffer(reqBodyBytes)
	if err = ctx.Err(); err != nil {
		return nil, clientError("aborted because context was done", err)
	}

	req, err := newRequest(ctx, c.urls[1], reqBody, "application/protobuf")
	if err != nil {
		return nil, clientError("could not build request", err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, clientError("failed to do request", err)
	}

	return &protoHatStreamReader{
		prs: protoStreamReader{
			r:       bufio.NewReader(resp.Body),
			c:       resp.Body,
			maxSize: 1 << 21, // 1GB
		},
	}, nil
}

// =======================
// Haberdasher JSON Client
// =======================

type haberdasherJSONClient struct {
	client HTTPClient
	urls   [2]string
}

// NewHaberdasherJSONClient creates a JSON client that implements the Haberdasher interface.
// It communicates using JSON and can be configured with a custom HTTPClient.
func NewHaberdasherJSONClient(addr string, client HTTPClient) Haberdasher {
	prefix := urlBase(addr) + HaberdasherPathPrefix
	urls := [2]string{
		prefix + "MakeHat",
		prefix + "MakeHats",
	}
	if httpClient, ok := client.(*http.Client); ok {
		return &haberdasherJSONClient{
			client: withoutRedirects(httpClient),
			urls:   urls,
		}
	}
	return &haberdasherJSONClient{
		client: client,
		urls:   urls,
	}
}

func (c *haberdasherJSONClient) MakeHat(ctx context.Context, in *Size) (*Hat, error) {
	ctx = ctxsetters.WithPackageName(ctx, "twitch.twirp.example")
	ctx = ctxsetters.WithServiceName(ctx, "Haberdasher")
	ctx = ctxsetters.WithMethodName(ctx, "MakeHat")
	out := new(Hat)
	err := doJSONRequest(ctx, c.client, c.urls[0], in, out)
	return out, err
}

func (c *haberdasherJSONClient) MakeHats(ctx context.Context, in *MakeHatsReq) (HatStream, error) {
	ctx = ctxsetters.WithPackageName(ctx, "twitch.twirp.example")
	ctx = ctxsetters.WithServiceName(ctx, "Haberdasher")
	ctx = ctxsetters.WithMethodName(ctx, "MakeHats")
	reqBodyBytes, err := proto.Marshal(in)
	if err != nil {
		return nil, clientError("failed to marshal proto request", err)
	}
	reqBody := bytes.NewBuffer(reqBodyBytes)
	if err = ctx.Err(); err != nil {
		return nil, clientError("aborted because context was done", err)
	}

	req, err := newRequest(ctx, c.urls[1], reqBody, "application/json")
	if err != nil {
		return nil, clientError("could not build request", err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, clientError("failed to do request", err)
	}

	jrs, err := newJSONStreamReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return &jsonHatStreamReader{
		jrs: jrs,
		c:   resp.Body,
	}, nil
}

// ==========================
// Haberdasher Server Handler
// ==========================

type haberdasherServer struct {
	Haberdasher
	hooks *twirp.ServerHooks
}

func NewHaberdasherServer(svc Haberdasher, hooks *twirp.ServerHooks) TwirpServer {
	return &haberdasherServer{
		Haberdasher: svc,
		hooks:       hooks,
	}
}

// writeError writes an HTTP response with a valid Twirp error format, and triggers hooks.
// If err is not a twirp.Error, it will get wrapped with twirp.InternalErrorWith(err)
func (s *haberdasherServer) writeError(ctx context.Context, resp http.ResponseWriter, err error) {
	writeError(ctx, resp, err, s.hooks)
}

// HaberdasherPathPrefix is used for all URL paths on a twirp Haberdasher server.
// Requests are always: POST HaberdasherPathPrefix/method
// It can be used in an HTTP mux to route twirp requests along with non-twirp requests on other routes.
const HaberdasherPathPrefix = "/twirp/twitch.twirp.example.Haberdasher/"

func (s *haberdasherServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	ctx = ctxsetters.WithPackageName(ctx, "twitch.twirp.example")
	ctx = ctxsetters.WithServiceName(ctx, "Haberdasher")
	ctx = ctxsetters.WithResponseWriter(ctx, resp)

	var err error
	ctx, err = callRequestReceived(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	if req.Method != "POST" {
		msg := fmt.Sprintf("unsupported method %q (only POST is allowed)", req.Method)
		err = badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, err)
		return
	}

	switch req.URL.Path {
	case "/twirp/twitch.twirp.example.Haberdasher/MakeHat":
		s.serveMakeHat(ctx, resp, req)
		return
	case "/twirp/twitch.twirp.example.Haberdasher/MakeHats":
		s.serveMakeHats(ctx, resp, req)
		return
	default:
		msg := fmt.Sprintf("no handler for path %q", req.URL.Path)
		err = badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, err)
		return
	}
}

func (s *haberdasherServer) serveMakeHat(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveMakeHatJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveMakeHatProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *haberdasherServer) serveMakeHatJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "MakeHat")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	reqContent := new(Size)
	unmarshaler := jsonpb.Unmarshaler{AllowUnknownFields: true}
	if err = unmarshaler.Unmarshal(req.Body, reqContent); err != nil {
		err = wrapErr(err, "failed to parse request json")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	// Call service method
	var respContent *Hat
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.MakeHat(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *Hat and nil error while calling MakeHat. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	var buf bytes.Buffer
	marshaler := &jsonpb.Marshaler{OrigName: true}
	if err = marshaler.Marshal(&buf, respContent); err != nil {
		err = wrapErr(err, "failed to marshal json response")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	ctx = ctxsetters.WithStatusCode(ctx, http.StatusOK)
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)

	respBytes := buf.Bytes()
	if n, err := resp.Write(respBytes); err != nil {
		msg := fmt.Sprintf("failed to write response, %d of %d bytes written: %s", n, len(respBytes), err.Error())
		twerr := twirp.NewError(twirp.Unknown, msg)
		callError(ctx, s.hooks, twerr)
	}
	callResponseSent(ctx, s.hooks)
}

func (s *haberdasherServer) serveMakeHatProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "MakeHat")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	resp.Header().Set("Content-Type", "application/protobuf")
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		err = wrapErr(err, "failed to read request body")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}
	reqContent := new(Size)
	if err = proto.Unmarshal(buf, reqContent); err != nil {
		err = wrapErr(err, "failed to parse request proto")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	// Call service method
	var respContent *Hat
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respContent, err = s.MakeHat(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil *Hat and nil error while calling MakeHat. nil responses are not supported"))
		return
	}

	ctx = callResponsePrepared(ctx, s.hooks)

	respBytes, err := proto.Marshal(respContent)
	if err != nil {
		err = wrapErr(err, "failed to marshal proto response")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	ctx = ctxsetters.WithStatusCode(ctx, http.StatusOK)
	resp.WriteHeader(http.StatusOK)
	if n, err := resp.Write(respBytes); err != nil {
		msg := fmt.Sprintf("failed to write response, %d of %d bytes written: %s", n, len(respBytes), err.Error())
		twerr := twirp.NewError(twirp.Unknown, msg)
		callError(ctx, s.hooks, twerr)
	}
	callResponseSent(ctx, s.hooks)
}

func (s *haberdasherServer) serveMakeHats(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serveMakeHatsJSON(ctx, resp, req)
	case "application/protobuf":
		s.serveMakeHatsProtobuf(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		s.writeError(ctx, resp, twerr)
	}
}

func (s *haberdasherServer) serveMakeHatsJSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
}

func (s *haberdasherServer) serveMakeHatsProtobuf(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	ctx = ctxsetters.WithMethodName(ctx, "MakeHats")
	ctx, err = callRequestRouted(ctx, s.hooks)
	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	resp.Header().Set("Content-Type", "application/protobuf")
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		err = wrapErr(err, "failed to read request body")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}
	reqContent := new(MakeHatsReq)
	if err = proto.Unmarshal(buf, reqContent); err != nil {
		err = wrapErr(err, "failed to parse request proto")
		s.writeError(ctx, resp, twirp.InternalErrorWith(err))
		return
	}

	// Call service method
	var respStream HatStream
	func() {
		defer func() {
			// In case of a panic, serve a 500 error and then panic.
			if r := recover(); r != nil {
				s.writeError(ctx, resp, twirp.InternalError("Internal service panic"))
				panic(r)
			}
		}()
		respStream, err = s.MakeHats(ctx, reqContent)
	}()

	if err != nil {
		s.writeError(ctx, resp, err)
		return
	}

	if respStream == nil {
		s.writeError(ctx, resp, twirp.InternalError("received a nil MakeHatsReq and nil error while calling MakeHats. nil responses are not supported"))
		return
	}

	respFlusher, canFlush := resp.(http.Flusher)

	ctx = callResponsePrepared(ctx, s.hooks)

	messages := proto.NewBuffer(nil)

	trailer := proto.NewBuffer(nil)
	_ = trailer.EncodeVarint((2 << 3) | 2) // field tag
	for {
		msg, err := respStream.Next(ctx)
		if err != nil {
			// TODO: figure out trailers' proto encoding beyond just a string
			if err == io.EOF {
				_ = trailer.EncodeStringBytes("OK")
			} else {
				_ = trailer.EncodeStringBytes(err.Error())
			}
			break
		}

		messages.Reset()
		_ = messages.EncodeVarint((1 << 3) | 2) // field tag
		err = messages.EncodeMessage(msg)
		if err != nil {
			err = wrapErr(err, "failed to marshal proto message")
			respStream.End(err)
			break
		}

		_, err = resp.Write(messages.Bytes())
		if err != nil {
			err = wrapErr(err, "failed to send proto message")
			respStream.End(err)
			break
		}

		if canFlush {
			respFlusher.Flush()
		}

		// TODO: Call a hook that we sent a message in a stream?
	}

	_, err = resp.Write(trailer.Bytes())
	if err != nil {
		// TODO: call error hook?
		err = wrapErr(err, "failed to write trailer")
		respStream.End(err)
	}
}

func (s *haberdasherServer) ServiceDescriptor() ([]byte, int) {
	return twirpFileDescriptor0, 0
}

func (s *haberdasherServer) ProtocGenTwirpVersion() string {
	return "v5.3.0"
}

// HatStream represents a stream of Hat messages.
type HatStream interface {
	Next(context.Context) (*Hat, error)
	End(error)
}

type protoHatStreamReader struct {
	prs protoStreamReader
}

func (r protoHatStreamReader) Next(context.Context) (*Hat, error) {
	out := new(Hat)
	err := r.prs.Read(out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r protoHatStreamReader) End(error) { _ = r.prs.c.Close() }

type jsonHatStreamReader struct {
	jrs *jsonStreamReader
	c   io.Closer
}

func (r jsonHatStreamReader) Next(context.Context) (*Hat, error) {
	out := new(Hat)
	err := r.jrs.Read(out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r jsonHatStreamReader) End(error) { _ = r.c.Close() }

// =====
// Utils
// =====

// HTTPClient is the interface used by generated clients to send HTTP requests.
// It is fulfilled by *(net/http).Client, which is sufficient for most users.
// Users can provide their own implementation for special retry policies.
//
// HTTPClient implementations should not follow redirects. Redirects are
// automatically disabled if *(net/http).Client is passed to client
// constructors. See the withoutRedirects function in this file for more
// details.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// TwirpServer is the interface generated server structs will support: they're
// HTTP handlers with additional methods for accessing metadata about the
// service. Those accessors are a low-level API for building reflection tools.
// Most people can think of TwirpServers as just http.Handlers.
type TwirpServer interface {
	http.Handler
	// ServiceDescriptor returns gzipped bytes describing the .proto file that
	// this service was generated from. Once unzipped, the bytes can be
	// unmarshalled as a
	// github.com/golang/protobuf/protoc-gen-go/descriptor.FileDescriptorProto.
	//
	// The returned integer is the index of this particular service within that
	// FileDescriptorProto's 'Service' slice of ServiceDescriptorProtos. This is a
	// low-level field, expected to be used for reflection.
	ServiceDescriptor() ([]byte, int)
	// ProtocGenTwirpVersion is the semantic version string of the version of
	// twirp used to generate this file.
	ProtocGenTwirpVersion() string
}

// WriteError writes an HTTP response with a valid Twirp error format.
// If err is not a twirp.Error, it will get wrapped with twirp.InternalErrorWith(err)
func WriteError(resp http.ResponseWriter, err error) {
	writeError(context.Background(), resp, err, nil)
}

// writeError writes Twirp errors in the response and triggers hooks.
func writeError(ctx context.Context, resp http.ResponseWriter, err error, hooks *twirp.ServerHooks) {
	// Non-twirp errors are wrapped as Internal (default)
	twerr, ok := err.(twirp.Error)
	if !ok {
		twerr = twirp.InternalErrorWith(err)
	}

	statusCode := twirp.ServerHTTPStatusFromErrorCode(twerr.Code())
	ctx = ctxsetters.WithStatusCode(ctx, statusCode)
	ctx = callError(ctx, hooks, twerr)

	resp.Header().Set("Content-Type", "application/json") // Error responses are always JSON (instead of protobuf)
	resp.WriteHeader(statusCode)                          // HTTP response status code

	respBody := marshalErrorToJSON(twerr)
	_, writeErr := resp.Write(respBody)
	if writeErr != nil {
		// We have three options here. We could log the error, call the Error
		// hook, or just silently ignore the error.
		//
		// Logging is unacceptable because we don't have a user-controlled
		// logger; writing out to stderr without permission is too rude.
		//
		// Calling the Error hook would confuse users: it would mean the Error
		// hook got called twice for one request, which is likely to lead to
		// duplicated log messages and metrics, no matter how well we document
		// the behavior.
		//
		// Silently ignoring the error is our least-bad option. It's highly
		// likely that the connection is broken and the original 'err' says
		// so anyway.
		_ = writeErr
	}

	callResponseSent(ctx, hooks)
}

// urlBase helps ensure that addr specifies a scheme. If it is unparsable
// as a URL, it returns addr unchanged.
func urlBase(addr string) string {
	// If the addr specifies a scheme, use it. If not, default to
	// http. If url.Parse fails on it, return it unchanged.
	url, err := url.Parse(addr)
	if err != nil {
		return addr
	}
	if url.Scheme == "" {
		url.Scheme = "http"
	}
	return url.String()
}

// getCustomHTTPReqHeaders retrieves a copy of any headers that are set in
// a context through the twirp.WithHTTPRequestHeaders function.
// If there are no headers set, or if they have the wrong type, nil is returned.
func getCustomHTTPReqHeaders(ctx context.Context) http.Header {
	header, ok := twirp.HTTPRequestHeaders(ctx)
	if !ok || header == nil {
		return nil
	}
	copied := make(http.Header)
	for k, vv := range header {
		if vv == nil {
			copied[k] = nil
			continue
		}
		copied[k] = make([]string, len(vv))
		copy(copied[k], vv)
	}
	return copied
}

// newRequest makes an http.Request from a client, adding common headers.
func newRequest(ctx context.Context, url string, reqBody io.Reader, contentType string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if customHeader := getCustomHTTPReqHeaders(ctx); customHeader != nil {
		req.Header = customHeader
	}
	req.Header.Set("Accept", contentType)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Twirp-Version", "v5.3.0")
	return req, nil
}

// JSON serialization for errors
type twerrJSON struct {
	Code string            `json:"code"`
	Msg  string            `json:"msg"`
	Meta map[string]string `json:"meta,omitempty"`
}

func (tj twerrJSON) toTwirpError() twirp.Error {
	errorCode := twirp.ErrorCode(tj.Code)
	if !twirp.IsValidErrorCode(errorCode) {
		msg := "invalid type returned from server error response: " + tj.Code
		return twirp.InternalError(msg)
	}

	twerr := twirp.NewError(errorCode, tj.Msg)
	for k, v := range tj.Meta {
		twerr = twerr.WithMeta(k, v)
	}
	return twerr
}

// marshalErrorToJSON returns JSON from a twirp.Error, that can be used as HTTP error response body.
// If serialization fails, it will use a descriptive Internal error instead.
func marshalErrorToJSON(twerr twirp.Error) []byte {
	// make sure that msg is not too large
	msg := twerr.Msg()
	if len(msg) > 1e6 {
		msg = msg[:1e6]
	}

	tj := twerrJSON{
		Code: string(twerr.Code()),
		Msg:  msg,
		Meta: twerr.MetaMap(),
	}

	buf, err := json.Marshal(&tj)
	if err != nil {
		buf = []byte("{\"type\": \"" + twirp.Internal + "\", \"msg\": \"There was an error but it could not be serialized into JSON\"}") // fallback
	}

	return buf
}

// errorFromResponse builds a twirp.Error from a non-200 HTTP response.
// If the response has a valid serialized Twirp error, then it's returned.
// If not, the response status code is used to generate a similar twirp
// error. See twirpErrorFromIntermediary for more info on intermediary errors.
func errorFromResponse(resp *http.Response) twirp.Error {
	statusCode := resp.StatusCode
	statusText := http.StatusText(statusCode)

	if isHTTPRedirect(statusCode) {
		// Unexpected redirect: it must be an error from an intermediary.
		// Twirp clients don't follow redirects automatically, Twirp only handles
		// POST requests, redirects should only happen on GET and HEAD requests.
		location := resp.Header.Get("Location")
		msg := fmt.Sprintf("unexpected HTTP status code %d %q received, Location=%q", statusCode, statusText, location)
		return twirpErrorFromIntermediary(statusCode, msg, location)
	}

	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return clientError("failed to read server error response body", err)
	}
	var tj twerrJSON
	if err := json.Unmarshal(respBodyBytes, &tj); err != nil {
		// Invalid JSON response; it must be an error from an intermediary.
		msg := fmt.Sprintf("Error from intermediary with HTTP status code %d %q", statusCode, statusText)
		return twirpErrorFromIntermediary(statusCode, msg, string(respBodyBytes))
	}

	return tj.toTwirpError()
}

// twirpErrorFromIntermediary maps HTTP errors from non-twirp sources to twirp errors.
// The mapping is similar to gRPC: https://github.com/grpc/grpc/blob/master/doc/http-grpc-status-mapping.md.
// Returned twirp Errors have some additional metadata for inspection.
func twirpErrorFromIntermediary(status int, msg string, bodyOrLocation string) twirp.Error {
	var code twirp.ErrorCode
	if isHTTPRedirect(status) { // 3xx
		code = twirp.Internal
	} else {
		switch status {
		case 400: // Bad Request
			code = twirp.Internal
		case 401: // Unauthorized
			code = twirp.Unauthenticated
		case 403: // Forbidden
			code = twirp.PermissionDenied
		case 404: // Not Found
			code = twirp.BadRoute
		case 429, 502, 503, 504: // Too Many Requests, Bad Gateway, Service Unavailable, Gateway Timeout
			code = twirp.Unavailable
		default: // All other codes
			code = twirp.Unknown
		}
	}

	twerr := twirp.NewError(code, msg)
	twerr = twerr.WithMeta("http_error_from_intermediary", "true") // to easily know if this error was from intermediary
	twerr = twerr.WithMeta("status_code", strconv.Itoa(status))
	if isHTTPRedirect(status) {
		twerr = twerr.WithMeta("location", bodyOrLocation)
	} else {
		twerr = twerr.WithMeta("body", bodyOrLocation)
	}
	return twerr
}
func isHTTPRedirect(status int) bool {
	return status >= 300 && status <= 399
}

// wrappedError implements the github.com/pkg/errors.Causer interface, allowing errors to be
// examined for their root cause.
type wrappedError struct {
	msg   string
	cause error
}

func wrapErr(err error, msg string) error { return &wrappedError{msg: msg, cause: err} }
func (e *wrappedError) Cause() error      { return e.cause }
func (e *wrappedError) Error() string     { return e.msg + ": " + e.cause.Error() }

// clientError adds consistency to errors generated in the client
func clientError(desc string, err error) twirp.Error {
	return twirp.InternalErrorWith(wrapErr(err, desc))
}

// badRouteError is used when the twirp server cannot route a request
func badRouteError(msg string, method, url string) twirp.Error {
	err := twirp.NewError(twirp.BadRoute, msg)
	err = err.WithMeta("twirp_invalid_route", method+" "+url)
	return err
}

// The standard library will, by default, redirect requests (including POSTs) if it gets a 302 or
// 303 response, and also 301s in go1.8. It redirects by making a second request, changing the
// method to GET and removing the body. This produces very confusing error messages, so instead we
// set a redirect policy that always errors. This stops Go from executing the redirect.
//
// We have to be a little careful in case the user-provided http.Client has its own CheckRedirect
// policy - if so, we'll run through that policy first.
//
// Because this requires modifying the http.Client, we make a new copy of the client and return it.
func withoutRedirects(in *http.Client) *http.Client {
	copy := *in
	copy.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if in.CheckRedirect != nil {
			// Run the input's redirect if it exists, in case it has side effects, but ignore any error it
			// returns, since we want to use ErrUseLastResponse.
			err := in.CheckRedirect(req, via)
			_ = err // Silly, but this makes sure generated code passes errcheck -blank, which some people use.
		}
		return http.ErrUseLastResponse
	}
	return &copy
}

// doProtobufRequest is common code to make a request to the remote twirp service.
func doProtobufRequest(ctx context.Context, client HTTPClient, url string, in, out proto.Message) (err error) {
	reqBodyBytes, err := proto.Marshal(in)
	if err != nil {
		return clientError("failed to marshal proto request", err)
	}
	reqBody := bytes.NewBuffer(reqBodyBytes)
	if err = ctx.Err(); err != nil {
		return clientError("aborted because context was done", err)
	}

	req, err := newRequest(ctx, url, reqBody, "application/protobuf")
	if err != nil {
		return clientError("could not build request", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return clientError("failed to do request", err)
	}

	defer func() {
		cerr := resp.Body.Close()
		if err == nil && cerr != nil {
			err = clientError("failed to close response body", cerr)
		}
	}()

	if err = ctx.Err(); err != nil {
		return clientError("aborted because context was done", err)
	}

	if resp.StatusCode != 200 {
		return errorFromResponse(resp)
	}

	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return clientError("failed to read response body", err)
	}
	if err = ctx.Err(); err != nil {
		return clientError("aborted because context was done", err)
	}

	if err = proto.Unmarshal(respBodyBytes, out); err != nil {
		return clientError("failed to unmarshal proto response", err)
	}
	return nil
}

// doJSONRequest is common code to make a request to the remote twirp service.
func doJSONRequest(ctx context.Context, client HTTPClient, url string, in, out proto.Message) (err error) {
	reqBody := bytes.NewBuffer(nil)
	marshaler := &jsonpb.Marshaler{OrigName: true}
	if err = marshaler.Marshal(reqBody, in); err != nil {
		return clientError("failed to marshal json request", err)
	}
	if err = ctx.Err(); err != nil {
		return clientError("aborted because context was done", err)
	}

	req, err := newRequest(ctx, url, reqBody, "application/json")
	if err != nil {
		return clientError("could not build request", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return clientError("failed to do request", err)
	}

	defer func() {
		cerr := resp.Body.Close()
		if err == nil && cerr != nil {
			err = clientError("failed to close response body", cerr)
		}
	}()

	if err = ctx.Err(); err != nil {
		return clientError("aborted because context was done", err)
	}

	if resp.StatusCode != 200 {
		return errorFromResponse(resp)
	}

	unmarshaler := jsonpb.Unmarshaler{AllowUnknownFields: true}
	if err = unmarshaler.Unmarshal(resp.Body, out); err != nil {
		return clientError("failed to unmarshal json response", err)
	}
	if err = ctx.Err(); err != nil {
		return clientError("aborted because context was done", err)
	}
	return nil
}

// Call twirp.ServerHooks.RequestReceived if the hook is available
func callRequestReceived(ctx context.Context, h *twirp.ServerHooks) (context.Context, error) {
	if h == nil || h.RequestReceived == nil {
		return ctx, nil
	}
	return h.RequestReceived(ctx)
}

// Call twirp.ServerHooks.RequestRouted if the hook is available
func callRequestRouted(ctx context.Context, h *twirp.ServerHooks) (context.Context, error) {
	if h == nil || h.RequestRouted == nil {
		return ctx, nil
	}
	return h.RequestRouted(ctx)
}

// Call twirp.ServerHooks.ResponsePrepared if the hook is available
func callResponsePrepared(ctx context.Context, h *twirp.ServerHooks) context.Context {
	if h == nil || h.ResponsePrepared == nil {
		return ctx
	}
	return h.ResponsePrepared(ctx)
}

// Call twirp.ServerHooks.ResponseSent if the hook is available
func callResponseSent(ctx context.Context, h *twirp.ServerHooks) {
	if h == nil || h.ResponseSent == nil {
		return
	}
	h.ResponseSent(ctx)
}

// Call twirp.ServerHooks.Error if the hook is available
func callError(ctx context.Context, h *twirp.ServerHooks, err twirp.Error) context.Context {
	if h == nil || h.Error == nil {
		return ctx
	}
	return h.Error(ctx, err)
}

type protoStreamReader struct {
	r *bufio.Reader
	c io.Closer

	maxSize int
}

func (r protoStreamReader) Read(msg proto.Message) error {
	// Get next field tag.
	tag, err := binary.ReadUvarint(r.r)
	if err != nil {
		return err
	}

	const (
		msgTag     = (1 << 3) | 2
		trailerTag = (2 << 3) | 2
	)

	if tag == trailerTag {
		_ = r.c.Close()
		return io.EOF
	}

	if tag != msgTag {
		return fmt.Errorf("invalid field tag: %v", tag)
	}

	// This is a real message. How long is it?
	l, err := binary.ReadUvarint(r.r)
	if err != nil {
		return err
	}
	if int(l) < 0 || int(l) > r.maxSize {
		return io.ErrShortBuffer
	}
	buf := make([]byte, int(l))

	// Go ahead and read a message.
	_, err = io.ReadFull(r.r, buf)
	if err != nil {
		return err
	}

	err = proto.Unmarshal(buf, msg)
	if err != nil {
		return err
	}
	return nil
}

type jsonStreamReader struct {
	dec               *json.Decoder
	unmarshaler       *jsonpb.Unmarshaler
	messageStreamDone bool
}

func newJSONStreamReader(r io.Reader) (*jsonStreamReader, error) {
	// stream should start with {"messages":[
	dec := json.NewDecoder(r)
	t, err := dec.Token()
	if err != nil {
		return nil, err
	}
	delim, ok := t.(json.Delim)
	if !ok || delim != '{' {
		return nil, fmt.Errorf("missing leading { in JSON stream, found %q", t)
	}

	t, err = dec.Token()
	if err != nil {
		return nil, err
	}
	key, ok := t.(string)
	if !ok || key != "messages" {
		return nil, fmt.Errorf("missing \"messages\" key in JSON stream, found %q", t)
	}

	t, err = dec.Token()
	if err != nil {
		return nil, err
	}
	delim, ok = t.(json.Delim)
	if !ok || delim != '[' {
		return nil, fmt.Errorf("missing [ to open messages array in JSON stream, found %q", t)
	}

	return &jsonStreamReader{
		dec:         dec,
		unmarshaler: &jsonpb.Unmarshaler{AllowUnknownFields: true},
	}, nil
}

func (r *jsonStreamReader) Read(msg proto.Message) error {
	if !r.messageStreamDone && r.dec.More() {
		return r.unmarshaler.UnmarshalNext(r.dec, msg)
	}

	// else, we hit the end of the message stream. finish up the array, and then read the trailer.
	r.messageStreamDone = true
	t, err := r.dec.Token()
	if err != nil {
		return err
	}
	delim, ok := t.(json.Delim)
	if !ok || delim != ']' {
		return fmt.Errorf("missing end of message array in JSON stream, found %q", t)
	}

	t, err = r.dec.Token()
	if err != nil {
		return err
	}
	key, ok := t.(string)
	if !ok || key != "trailer" {
		return fmt.Errorf("missing trailer after messages in JSON stream, found %q", t)
	}

	var tj twerrJSON
	err = r.dec.Decode(&tj)
	if err != nil {
		return err
	}

	if tj.Code == "stream_complete" {
		return io.EOF
	}

	return tj.toTwirpError()
}

var twirpFileDescriptor0 = []byte{
	// 234 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x29, 0x29, 0xcf, 0x2c, 0x49,
	0xce, 0xd0, 0x2b, 0x29, 0xcf, 0x2c, 0x2a, 0xd0, 0x4b, 0xad, 0x48, 0xcc, 0x2d, 0xc8, 0x49, 0x55,
	0x72, 0xe6, 0x62, 0xf6, 0x48, 0x2c, 0x11, 0x12, 0xe2, 0x62, 0x29, 0xce, 0xac, 0x4a, 0x95, 0x60,
	0x54, 0x60, 0xd4, 0x60, 0x0d, 0x02, 0xb3, 0x85, 0x44, 0xb8, 0x58, 0x93, 0xf3, 0x73, 0xf2, 0x8b,
	0x24, 0x98, 0x14, 0x18, 0x35, 0x38, 0x83, 0x20, 0x1c, 0x90, 0xca, 0xbc, 0xc4, 0xdc, 0x54, 0x09,
	0x66, 0xb0, 0x20, 0x98, 0xad, 0x24, 0xc7, 0xc5, 0x12, 0x0c, 0xd2, 0x21, 0xc6, 0xc5, 0x96, 0x99,
	0x97, 0x9c, 0x91, 0x5a, 0x0c, 0x35, 0x07, 0xca, 0x53, 0x72, 0xe4, 0xe2, 0xf6, 0x4d, 0xcc, 0x4e,
	0xf5, 0x48, 0x2c, 0x29, 0x0e, 0x4a, 0x2d, 0xc4, 0xa5, 0x4c, 0x48, 0x8a, 0x8b, 0xa3, 0xb0, 0x34,
	0x31, 0xaf, 0x24, 0xb3, 0xa4, 0x12, 0x6c, 0x27, 0x6b, 0x10, 0x9c, 0x6f, 0x34, 0x9b, 0x91, 0x8b,
	0xdb, 0x23, 0x31, 0x29, 0xb5, 0x28, 0x25, 0xb1, 0x38, 0x23, 0xb5, 0x48, 0xc8, 0x81, 0x8b, 0x1d,
	0x6a, 0xa4, 0x90, 0x94, 0x1e, 0x36, 0x9f, 0xe9, 0x81, 0x5c, 0x24, 0x25, 0x89, 0x5d, 0x0e, 0xa4,
	0xcd, 0x8b, 0x8b, 0x03, 0xe6, 0x28, 0x21, 0x45, 0xec, 0xca, 0x90, 0x1c, 0x8d, 0xc7, 0x24, 0x03,
	0x46, 0x27, 0xce, 0x28, 0x76, 0xa8, 0x40, 0x12, 0x1b, 0x38, 0xb4, 0x8d, 0x01, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x78, 0xf6, 0xd2, 0xc4, 0x7e, 0x01, 0x00, 0x00,
}
