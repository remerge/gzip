package gzip

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/klauspost/compress/gzip"
)

type gzipHandler struct {
	*Options
	gzPool sync.Pool
}

func newGzipHandler(level int, options ...Option) *gzipHandler {
	handler := &gzipHandler{
		Options: DefaultOptions,
		gzPool: sync.Pool{
			New: func() interface{} {
				gz, err := gzip.NewWriterLevel(ioutil.Discard, level)
				if err != nil {
					panic(err)
				}
				return gz
			},
		},
	}
	for _, setter := range options {
		setter(handler.Options)
	}
	return handler
}

func (g *gzipHandler) scheduleTriggerOnContextDone(done <-chan struct{}, cb func()) {
	if done == nil {
		return
	}

	go func(doneCh <-chan struct{}) {
		deadline := time.NewTicker(10 * time.Minute)
		defer deadline.Stop()
		for {
			select {
			case <-doneCh:
				break
			case <-deadline.C:
				break
			}
		}

		cb()
	}(done)
}

func (g *gzipHandler) Handle(c *gin.Context) {
	if fn := g.DecompressFn; fn != nil && c.Request.Header.Get("Content-Encoding") == "gzip" {
		fn(c)
	}

	if !g.shouldCompress(c.Request) {
		return
	}

	gz := g.gzPool.Get().(*gzip.Writer)
	gz.Reset(c.Writer)

	c.Header("Content-Encoding", "gzip")
	c.Header("Vary", "Accept-Encoding")
	c.Writer = &gzipWriter{ResponseWriter: c.Writer, writer: gz}

	g.scheduleTriggerOnContextDone(c.Request.Context().Done(), func() {
		gz.Reset(ioutil.Discard)
		g.gzPool.Put(gz)
	})

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
