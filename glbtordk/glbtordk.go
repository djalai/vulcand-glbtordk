package glbtordk

import (
        "bytes"
	"encoding/json"
        "fmt"
        "io"
        "net/http"
        "net/url"
        "regexp"
        "strconv"
        "strings"
        "github.com/urfave/cli"
        "github.com/vulcand/oxy/utils"
        "github.com/vulcand/vulcand/plugin"
)

const Type = "glbtordk"

func GetSpec() *plugin.MiddlewareSpec {
       return &plugin.MiddlewareSpec{
               Type:      Type,       // A short name for the middleware
               FromOther: FromOther,  // Tells vulcand how to create middleware from another one
               FromCli:   FromCli,    // Tells vulcand how to create middleware from CLI
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
       data, err := r.Header.Get(g.cfg.Header)
       // Reject the request by writing forbidden response
       if data != nil {
	       params, err := r.URL.Query()

	       if len(params) != 0 {
		oldbody, err := ioutil.ReadAll(r.Body);
		var olddata map[string]interface{}
		err = json.Unmarshal(body, &olddata)
		if err != nil {
			
		}
	       }

	       argstring := r.URL.Query().Get("argString")
	       if len(hostname) != 0 {
		   io.WriteString(w, hostname)  // or
		   w.Write([]byte(hostname))
	       }
               w.WriteHeader(http.StatusForbidden)
               io.WriteString(w, "Forbidden")
               return
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
