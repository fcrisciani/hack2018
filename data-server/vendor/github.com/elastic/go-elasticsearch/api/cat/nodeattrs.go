// generated by github.com/elastic/go-elasticsearch/cmd/generator; DO NOT EDIT

package cat

import (
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/transport"
	"github.com/elastic/go-elasticsearch/util"
)

// NodeattrsOption is a non-required Nodeattrs option that gets applied to an HTTP request.
type NodeattrsOption func(r *transport.Request)

// WithNodeattrsFormat - a short version of the Accept header, e.g. json, yaml.
func WithNodeattrsFormat(format string) NodeattrsOption {
	return func(r *transport.Request) {
	}
}

// WithNodeattrsH - comma-separated list of column names to display.
func WithNodeattrsH(h []string) NodeattrsOption {
	return func(r *transport.Request) {
	}
}

// WithNodeattrsHelp - return help information.
func WithNodeattrsHelp(help bool) NodeattrsOption {
	return func(r *transport.Request) {
	}
}

// WithNodeattrsLocal - return local information, do not retrieve the state from master node (default: false).
func WithNodeattrsLocal(local bool) NodeattrsOption {
	return func(r *transport.Request) {
	}
}

// WithNodeattrsMasterTimeout - explicit operation timeout for connection to master node.
func WithNodeattrsMasterTimeout(masterTimeout time.Duration) NodeattrsOption {
	return func(r *transport.Request) {
	}
}

// WithNodeattrsS - comma-separated list of column names or column aliases to sort by.
func WithNodeattrsS(s []string) NodeattrsOption {
	return func(r *transport.Request) {
	}
}

// WithNodeattrsV - verbose mode. Display column headers.
func WithNodeattrsV(v bool) NodeattrsOption {
	return func(r *transport.Request) {
	}
}

// WithNodeattrsErrorTrace - include the stack trace of returned errors.
func WithNodeattrsErrorTrace(errorTrace bool) NodeattrsOption {
	return func(r *transport.Request) {
	}
}

// WithNodeattrsFilterPath - a comma-separated list of filters used to reduce the respone.
func WithNodeattrsFilterPath(filterPath []string) NodeattrsOption {
	return func(r *transport.Request) {
	}
}

// WithNodeattrsHuman - return human readable values for statistics.
func WithNodeattrsHuman(human bool) NodeattrsOption {
	return func(r *transport.Request) {
	}
}

// WithNodeattrsIgnore - ignores the specified HTTP status codes.
func WithNodeattrsIgnore(ignore []int) NodeattrsOption {
	return func(r *transport.Request) {
	}
}

// WithNodeattrsPretty - pretty format the returned JSON response.
func WithNodeattrsPretty(pretty bool) NodeattrsOption {
	return func(r *transport.Request) {
	}
}

// WithNodeattrsSourceParam - the URL-encoded request definition. Useful for libraries that do not accept a request body for non-POST requests.
func WithNodeattrsSourceParam(sourceParam string) NodeattrsOption {
	return func(r *transport.Request) {
	}
}

// Nodeattrs - see https://www.elastic.co/guide/en/elasticsearch/reference/5.x/cat-nodeattrs.html for more info.
//
// options: optional parameters.
func (c *Cat) Nodeattrs(options ...NodeattrsOption) (*NodeattrsResponse, error) {
	req := c.transport.NewRequest("GET")
	for _, option := range options {
		option(req)
	}
	resp, err := c.transport.Do(req)
	return &NodeattrsResponse{resp}, err
}

// NodeattrsResponse is the response for Nodeattrs.
type NodeattrsResponse struct {
	Response *http.Response
	// TODO: fill in structured response
}

// DecodeBody decodes the JSON body of the HTTP response.
func (r *NodeattrsResponse) DecodeBody() (util.MapStr, error) {
	return transport.DecodeResponseBody(r.Response)
}