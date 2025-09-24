package fiberv2

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
)

// Unwrap は Huma コンテキストから Fiber v2 のコンテキストを取り出します。
func Unwrap(ctx huma.Context) *fiber.Ctx {
	for {
		if c, ok := ctx.(interface{ Unwrap() huma.Context }); ok {
			ctx = c.Unwrap()
			continue
		}
		break
	}
	if c, ok := ctx.(*fiberWrapper); ok {
		return c.Unwrap()
	}
	panic("not a fiberv2 context")
}

type fiberAdapter struct {
	tester requestTester
	router router
}

type fiberWrapper struct {
	op     *huma.Operation
	status int
	orig   *fiber.Ctx
	ctx    context.Context
}

var _ huma.Context = &fiberWrapper{}

func (c *fiberWrapper) Unwrap() *fiber.Ctx {
	return c.orig
}

func (c *fiberWrapper) Operation() *huma.Operation { return c.op }

func (c *fiberWrapper) Matched() string {
	if route := c.orig.Route(); route != nil {
		return route.Path
	}
	return ""
}

func (c *fiberWrapper) Context() context.Context { return c.ctx }
func (c *fiberWrapper) Method() string           { return c.orig.Method() }
func (c *fiberWrapper) Host() string             { return c.orig.Hostname() }
func (c *fiberWrapper) RemoteAddr() string       { return c.orig.Context().RemoteAddr().String() }

func (c *fiberWrapper) URL() url.URL {
	uri := c.orig.Context().URI()
	u, _ := url.Parse(uri.String())
	return *u
}

func (c *fiberWrapper) Param(name string) string  { return c.orig.Params(name) }
func (c *fiberWrapper) Query(name string) string  { return c.orig.Query(name) }
func (c *fiberWrapper) Header(name string) string { return c.orig.Get(name) }

func (c *fiberWrapper) EachHeader(cb func(name, value string)) {
	ctx := c.orig.Context()
	ctx.Request.Header.VisitAll(func(k, v []byte) {
		cb(string(k), string(v))
	})
}

func (c *fiberWrapper) BodyReader() io.Reader { return bytes.NewReader(c.orig.Body()) }

func (c *fiberWrapper) GetMultipartForm() (*multipart.Form, error) {
	return c.orig.MultipartForm()
}

func (c *fiberWrapper) SetReadDeadline(deadline time.Time) error {
	conn := c.orig.Context().Conn()
	if conn == nil {
		return nil
	}
	return conn.SetReadDeadline(deadline)
}

func (c *fiberWrapper) SetStatus(code int) {
	c.status = code
	c.orig.Status(code)
}

func (c *fiberWrapper) Status() int { return c.status }

func (c *fiberWrapper) AppendHeader(name, value string) { c.orig.Append(name, value) }
func (c *fiberWrapper) SetHeader(name, value string)    { c.orig.Set(name, value) }

func (c *fiberWrapper) BodyWriter() io.Writer { return c.orig.Context() }

func (c *fiberWrapper) TLS() *tls.ConnectionState { return c.orig.Context().TLSConnectionState() }

func (c *fiberWrapper) Version() huma.ProtoVersion {
	return huma.ProtoVersion{Proto: c.orig.Protocol()}
}

type router interface {
	Add(method string, path string, handlers ...fiber.Handler) fiber.Router
}

type requestTester interface {
	Test(*http.Request, ...int) (*http.Response, error)
}

type contextWrapperValue struct {
	Key   any
	Value any
}

type contextWrapper struct {
	values []*contextWrapperValue
	context.Context
}

var _ context.Context = &contextWrapper{}

func (c *contextWrapper) Value(key any) any {
	if v := c.Context.Value(key); v != nil {
		return v
	}
	for _, pair := range c.values {
		if pair.Key == key {
			return pair.Value
		}
	}
	return nil
}

func (a *fiberAdapter) Handle(op *huma.Operation, handler func(huma.Context)) {
	path := strings.ReplaceAll(op.Path, "{", ":")
	path = strings.ReplaceAll(path, "}", "")

	a.router.Add(op.Method, path, func(c *fiber.Ctx) error {
		var values []*contextWrapperValue
		c.Context().VisitUserValuesAll(func(key, value any) {
			values = append(values, &contextWrapperValue{Key: key, Value: value})
		})

		baseCtx := c.UserContext()
		if baseCtx == nil {
			baseCtx = context.Background()
		}

		handler(&fiberWrapper{
			op:   op,
			orig: c,
			ctx: &contextWrapper{
				values:  values,
				Context: baseCtx,
			},
		})
		return nil
	})
}

func (a *fiberAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := a.tester.Test(r)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		panic(err)
	}
	for k, v := range resp.Header {
		for _, item := range v {
			w.Header().Add(k, item)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

// New は Fiber v2 アプリケーションに Huma API を組み込みます。
func New(app *fiber.App, config huma.Config) huma.API {
	return huma.NewAPI(config, &fiberAdapter{tester: app, router: app})
}

// NewWithGroup は Fiber v2 の Router グループに Huma API を組み込みます。
func NewWithGroup(app *fiber.App, group fiber.Router, config huma.Config) huma.API {
	return huma.NewAPI(config, &fiberAdapter{tester: app, router: group})
}
