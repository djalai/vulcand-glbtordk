package glbtordk

import (
  "strings"
  "encoding/json"
  "fmt"
  "github.com/codegangsta/cli"
  "github.com/vulcand/vulcand/plugin"
  "io/ioutil"
  "net/http"
)

const Type = "glbtordk"

func GetSpec() *plugin.MiddlewareSpec {
  return &plugin.MiddlewareSpec{
    Type:  Type,   // A short name for the middleware
    FromOther: FromOther,  // Tells vulcand how to create middleware from another one
    FromCli:   FromCli,  // Tells vulcand how to create middleware from CLI
    CliFlags:  CliFlags(), // Vulcand will add this flags CLI command
  }
}

// glbtordkMiddleware struct holds configuration parameters and is used to
// serialize/deserialize the configuration from storage engines.
type GlbtordkMiddleware struct {
  Header string
}

// glbtordk middleware handler
type GlbtordkHandler struct {
  cfg  GlbtordkMiddleware
  next http.Handler
}

// This function will be called each time the request hits the location with this middleware activated
func (g *GlbtordkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  data := r.Header.Get(g.cfg.Header)

  if data != "" {
    params := r.URL.Query()

    if len(params) != 0 {
      oldBody, err := ioutil.ReadAll(r.Body)
      var data map[string]interface{}
      json.Unmarshal(oldBody, &data)
      newData := make(map[string]interface{},len(params)+1)
      for k, v := range params {
        newData[k] = v
      }
      newData["payload"] = data
      newBody, _ := json.Marshal(newData)
      ioutil.NopCloser(strings.NewReader(newBody))
      r.ContentLength = int64(len(newBody))
    }
  }
  // Pass the request to the next middleware in chain
  g.next.ServeHTTP(w, r)
}

// This function is important, it's called by vulcand to create a new handler from the middleware config and put it into the
// middleware chain. Note that we need to remember 'next' handler to call
func (c *GlbtordkMiddleware) NewHandler(next http.Handler) (http.Handler, error) {
  return &GlbtordkHandler{next: next, cfg: *c}, nil
}

// String() will be called by loggers inside Vulcand and command line tool.
func (c *GlbtordkMiddleware) String() string {
  return fmt.Sprintf("gitlab webhook header=%v", c.Header)
}
