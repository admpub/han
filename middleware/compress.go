package middleware

import (
	"bufio"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/admpub/han"
	"github.com/admpub/han/engine"
)

type (
	// GzipConfig defines the config for Gzip middleware.
	GzipConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper han.Skipper `json:"-"`

		// Gzip compression level.
		// Optional. Default value -1.
		Level int `json:"level"`
	}

	gzipWriter struct {
		io.Writer
		engine.Response
	}
)

var (
	// DefaultGzipConfig is the default Gzip middleware config.
	DefaultGzipConfig = &GzipConfig{
		Skipper: han.DefaultSkipper,
		Level:   -1,
	}
)

func (w *gzipWriter) WriteHeader(code int) {
	if code == http.StatusNoContent {
		w.Header().Del(han.HeaderContentEncoding)
	}
	w.WriteHeader(code)
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	if w.Header().Get(han.HeaderContentType) == `` {
		w.Header().Set(han.HeaderContentType, http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

func (w *gzipWriter) Flush() error {
	return w.Writer.(*gzip.Writer).Flush()
}

func (w *gzipWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.Response.(http.Hijacker).Hijack()
}

func (w *gzipWriter) CloseNotify() <-chan bool {
	return w.Response.(http.CloseNotifier).CloseNotify()
}

// Gzip returns a middleware which compresses HTTP response using gzip compression
// scheme.
func Gzip(config ...*GzipConfig) han.MiddlewareFunc {
	if len(config) < 1 || config[0] == nil {
		return GzipWithConfig(DefaultGzipConfig)
	}
	return GzipWithConfig(config[0])
}

// GzipWithConfig return Gzip middleware with config.
// See: `Gzip()`.
func GzipWithConfig(config *GzipConfig) han.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultGzipConfig.Skipper
	}
	if config.Level == 0 {
		config.Level = DefaultGzipConfig.Level
	}
	scheme := `gzip`

	return func(h han.Handler) han.Handler {
		return han.HandlerFunc(func(c han.Context) error {
			if config.Skipper(c) {
				return h.Handle(c)
			}
			resp := c.Response()
			resp.Header().Add(han.HeaderVary, han.HeaderAcceptEncoding)
			if strings.Contains(c.Request().Header().Get(han.HeaderAcceptEncoding), scheme) {
				resp.Header().Add(han.HeaderContentEncoding, scheme)
				rw := resp.Writer()
				w, err := gzip.NewWriterLevel(rw, config.Level)
				if err != nil {
					return err
				}
				defer func() {
					if resp.Size() == 0 {
						if resp.Header().Get(han.HeaderContentEncoding) == scheme {
							resp.Header().Del(han.HeaderContentEncoding)
						}
						// We have to reset response to it's pristine state when
						// nothing is written to body or error is returned.
						// See issue #424, #407.
						resp.SetWriter(rw)
						w.Reset(ioutil.Discard)
					}
					w.Close()
				}()
				resp.SetWriter(&gzipWriter{Writer: w, Response: resp})
			}
			return h.Handle(c)
		})
	}
}
