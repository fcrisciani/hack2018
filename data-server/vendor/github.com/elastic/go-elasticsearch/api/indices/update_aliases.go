// generated by github.com/elastic/go-elasticsearch/cmd/generator; DO NOT EDIT

package indices

import (
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/transport"
	"github.com/elastic/go-elasticsearch/util"
)

// UpdateAliasesOption is a non-required UpdateAliases option that gets applied to an HTTP request.
type UpdateAliasesOption func(r *transport.Request)

// WithUpdateAliasesMasterTimeout - specify timeout for connection to master.
func WithUpdateAliasesMasterTimeout(masterTimeout time.Duration) UpdateAliasesOption {
	return func(r *transport.Request) {
	}
}

// WithUpdateAliasesTimeout - request timeout.
func WithUpdateAliasesTimeout(timeout time.Duration) UpdateAliasesOption {
	return func(r *transport.Request) {
	}
}

// WithUpdateAliasesErrorTrace - include the stack trace of returned errors.
func WithUpdateAliasesErrorTrace(errorTrace bool) UpdateAliasesOption {
	return func(r *transport.Request) {
	}
}

// WithUpdateAliasesFilterPath - a comma-separated list of filters used to reduce the respone.
func WithUpdateAliasesFilterPath(filterPath []string) UpdateAliasesOption {
	return func(r *transport.Request) {
	}
}

// WithUpdateAliasesHuman - return human readable values for statistics.
func WithUpdateAliasesHuman(human bool) UpdateAliasesOption {
	return func(r *transport.Request) {
	}
}

// WithUpdateAliasesIgnore - ignores the specified HTTP status codes.
func WithUpdateAliasesIgnore(ignore []int) UpdateAliasesOption {
	return func(r *transport.Request) {
	}
}

// WithUpdateAliasesPretty - pretty format the returned JSON response.
func WithUpdateAliasesPretty(pretty bool) UpdateAliasesOption {
	return func(r *transport.Request) {
	}
}

// WithUpdateAliasesSourceParam - the URL-encoded request definition. Useful for libraries that do not accept a request body for non-POST requests.
func WithUpdateAliasesSourceParam(sourceParam string) UpdateAliasesOption {
	return func(r *transport.Request) {
	}
}

// UpdateAliases - APIs in Elasticsearch accept an index name when working against a specific index, and several indices when applicable. See https://www.elastic.co/guide/en/elasticsearch/reference/5.x/indices-aliases.html for more info.
//
// body: the definition of "actions" to perform.
//
// options: optional parameters.
func (i *Indices) UpdateAliases(body map[string]interface{}, options ...UpdateAliasesOption) (*UpdateAliasesResponse, error) {
	req := i.transport.NewRequest("POST")
	for _, option := range options {
		option(req)
	}
	resp, err := i.transport.Do(req)
	return &UpdateAliasesResponse{resp}, err
}

// UpdateAliasesResponse is the response for UpdateAliases.
type UpdateAliasesResponse struct {
	Response *http.Response
	// TODO: fill in structured response
}

// DecodeBody decodes the JSON body of the HTTP response.
func (r *UpdateAliasesResponse) DecodeBody() (util.MapStr, error) {
	return transport.DecodeResponseBody(r.Response)
}