package gzip

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/klauspost/compress/gzip"
)

type gzipHandler struct {
	*Options
	gzipLevel int
}

func newGzipHandler(level int, options ...Option) *gzipHandler {
	handler := &gzipHandler{
		Options:   DefaultOptions,
		gzipLevel: level,
	}
	for _, setter := range options {
		setter(handler.Options)
	}
	return handler
}

func (g *gzipHandler) Handle(c *gin.Context) {
	if fn := g.DecompressFn; fn != nil && c.Request.Header.Get("Content-Encoding") == "gzip" {
		fn(c)
	}

	if !g.shouldCompress(c.Request) {
		return
	}

	gz := g.mustNewGzipWriter()
	gz.Reset(c.Writer)

	c.Header("Content-Encoding", "gzip")
	c.Header("Vary", "Accept-Encoding")
	c.Writer = &gzipWriter{ResponseWriter: c.Writer, writer: gz}

	defer func() {
		gz.Close()
		c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
	}()
	c.Next()
}

func (g *gzipHandler) shouldCompress(req *http.Request) bool {
	if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") ||
		strings.Contains(req.Header.Get("Connection"), "Upgrade") ||
		strings.Contains(req.Header.Get("Accept"), "text/event-stream") {
		return false
	}

	if g.MatchSupportedRequestFn != nil {
		if ok, supported := g.MatchSupportedRequestFn(req); ok {
			return supported
		}
	}

	return !g.isPathExcluded(req.URL.Path)
}

func (g *gzipHandler) isPathExcluded(path string) bool {
	extension := filepath.Ext(path)
	return g.ExcludedExtensions.Contains(extension) ||
		g.ExcludedPaths.Contains(path) ||
		g.ExcludedPathesRegexs.Contains(path)
}

func (g *gzipHandler) mustNewGzipWriter() *gzip.Writer {
	gz, err := gzip.NewWriterLevel(ioutil.Discard, g.gzipLevel)
	if err != nil {
		panic(err)
	}
	return gz
}
