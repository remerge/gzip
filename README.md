# GZIP middleware for Gin

Gin middleware to enable `GZIP` support.

## Usage

Download and install it:

```sh
go get github.com/gin-contrib/gzip
```

Import it in your code:

```go
import "github.com/gin-contrib/gzip"
```

Canonical example:

```go
package main

import (
  "fmt"
  "net/http"
  "time"

  "github.com/gin-contrib/gzip"
  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  r.Use(gzip.Gzip(gzip.DefaultCompression))
  r.GET("/ping", func(c *gin.Context) {
    c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
  })

  // Listen and Server in 0.0.0.0:8080
  if err := r.Run(":8080"); err != nil {
    log.Fatal(err)
  }
}
```

Customized Excluded Extensions

```go
package main

import (
  "fmt"
  "net/http"
  "time"

  "github.com/gin-contrib/gzip"
  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions([]string{".pdf", ".mp4"})))
  r.GET("/ping", func(c *gin.Context) {
    c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
  })

  // Listen and Server in 0.0.0.0:8080
  if err := r.Run(":8080"); err != nil {
    log.Fatal(err)
  }
}
```

Customized Excluded Paths

```go
package main

import (
  "fmt"
  "net/http"
  "time"

  "github.com/gin-contrib/gzip"
  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{"/api/"})))
  r.GET("/ping", func(c *gin.Context) {
    c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
  })

  // Listen and Server in 0.0.0.0:8080
  if err := r.Run(":8080"); err != nil {
    log.Fatal(err)
  }
}
```

Customized Excluded Paths

```go
package main

import (
  "fmt"
  "net/http"
  "time"

  "github.com/gin-contrib/gzip"
  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPathsRegexs([]string{".*"})))
  r.GET("/ping", func(c *gin.Context) {
    c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
  })

  // Listen and Server in 0.0.0.0:8080
  if err := r.Run(":8080"); err != nil {
    log.Fatal(err)
  }
}
```
