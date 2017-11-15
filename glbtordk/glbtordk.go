package glbtordk

import (
  "bytes"
  "encoding/json"
  "fmt"
  "github.com/codegangsta/cli"
  "github.com/vulcand/vulcand/plugin"
  "io/ioutil"
  "net/http"
  "strings"
)

const Type = "glbtordk"

func GetSpec() *plugin.MiddlewareSpec {
  return &plugin.MiddlewareSpec{
    Type:      Type,   // A short name for the middleware
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


  gitlabheader := r.Header.Get(g.cfg.Header)

  if gitlabheader != "" {
    oldBody, _ := ioutil.ReadAll(r.Body)

    newData := make(map[string]interface{})
    newData["options"] = make(map[string]string)
    newData["options"].(map[string]string)["gitlabpayload"] = string(oldBody[:])

    params := r.URL.Query()
    argstring := params.Get("argString")
    if argstring != "" {
      args := strings.Split(argstring,"-")
      for _, element := range args {
        if element != "" {
          myoption := strings.Fields(element)
          newData["options"].(map[string]string)[myoption[0]] = myoption[1]
        }
      }
    }

    newBody, _ := json.Marshal(newData)
    r.Body = ioutil.NopCloser(bytes.NewReader(newBody))
    r.ContentLength = int64(len(newBody))
  }

  // Pass the request to the next middleware in chain
  g.next.ServeHTTP(w, r)
}

// This function is optional but handy, used to check input parameters when creating new middlewares
func New(header string) (*GlbtordkMiddleware, error) {
    if header == "" {
        return nil, fmt.Errorf("header can not be empty")
    }
    return &GlbtordkMiddleware{Header: header}, nil
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

// FromOther Will be called by Vulcand when engine or API will read the middleware from the serialized format.
// It's important that the signature of the function will be exactly the same, otherwise Vulcand will
// fail to register this middleware.
// The first and the only parameter should be the struct itself, no pointers and other variables.
// Function should return middleware interface and error in case if the parameters are wrong.
func FromOther(g GlbtordkMiddleware) (plugin.Middleware, error) {
  return New(g.Header)
}

// FromCli constructs the middleware from the command line
func FromCli(c *cli.Context) (plugin.Middleware, error) {
  return New(c.String("header"))
}

// CliFlags will be used by Vulcand construct help and CLI command for the vctl command
func CliFlags() []cli.Flag {
  return []cli.Flag{
    cli.StringFlag{
      Name:  "header",
      Usage: "Gitlab header to match against http request to ensure it's coming from gitlab",
    },
  }
}
