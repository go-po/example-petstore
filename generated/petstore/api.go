// Package petstore provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package petstore

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Error defines model for Error.
type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// NewPet defines model for NewPet.
type NewPet struct {
	Name string  `json:"name"`
	Tag  *string `json:"tag,omitempty"`
}

// Pet defines model for Pet.
type Pet struct {
	// Embedded struct due to allOf(#/components/schemas/NewPet)
	NewPet
	// Embedded fields due to inline allOf schema
	Id int64 `json:"id"`
}

// FindPetsParams defines parameters for FindPets.
type FindPetsParams struct {

	// tags to filter by
	Tags *[]string `json:"tags,omitempty"`

	// maximum number of results to return
	Limit *int32 `json:"limit,omitempty"`
}

// AddPetJSONBody defines parameters for AddPet.
type AddPetJSONBody NewPet

// AddPetRequestBody defines body for AddPet for application/json ContentType.
type AddPetJSONRequestBody AddPetJSONBody

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(req *http.Request, ctx context.Context) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A callback for modifying requests which are generated before sending over
	// the network.
	RequestEditor RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = http.DefaultClient
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditor = fn
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// FindPets request
	FindPets(ctx context.Context, params *FindPetsParams) (*http.Response, error)

	// AddPet request  with any body
	AddPetWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error)

	AddPet(ctx context.Context, body AddPetJSONRequestBody) (*http.Response, error)

	// DeletePet request
	DeletePet(ctx context.Context, id int64) (*http.Response, error)

	// FindPetById request
	FindPetById(ctx context.Context, id int64) (*http.Response, error)
}

func (c *Client) FindPets(ctx context.Context, params *FindPetsParams) (*http.Response, error) {
	req, err := NewFindPetsRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(req, ctx)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) AddPetWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewAddPetRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(req, ctx)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) AddPet(ctx context.Context, body AddPetJSONRequestBody) (*http.Response, error) {
	req, err := NewAddPetRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(req, ctx)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) DeletePet(ctx context.Context, id int64) (*http.Response, error) {
	req, err := NewDeletePetRequest(c.Server, id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(req, ctx)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) FindPetById(ctx context.Context, id int64) (*http.Response, error) {
	req, err := NewFindPetByIdRequest(c.Server, id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(req, ctx)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

// NewFindPetsRequest generates requests for FindPets
func NewFindPetsRequest(server string, params *FindPetsParams) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/pets")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	queryValues := queryUrl.Query()

	if params.Tags != nil {

		if queryFrag, err := runtime.StyleParam("form", true, "tags", *params.Tags); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

	}

	if params.Limit != nil {

		if queryFrag, err := runtime.StyleParam("form", true, "limit", *params.Limit); err != nil {
			return nil, err
		} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
			return nil, err
		} else {
			for k, v := range parsed {
				for _, v2 := range v {
					queryValues.Add(k, v2)
				}
			}
		}

	}

	queryUrl.RawQuery = queryValues.Encode()

	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewAddPetRequest calls the generic AddPet builder with application/json body
func NewAddPetRequest(server string, body AddPetJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewAddPetRequestWithBody(server, "application/json", bodyReader)
}

// NewAddPetRequestWithBody generates requests for AddPet with any type of body
func NewAddPetRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/pets")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryUrl.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	return req, nil
}

// NewDeletePetRequest generates requests for DeletePet
func NewDeletePetRequest(server string, id int64) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParam("simple", false, "id", id)
	if err != nil {
		return nil, err
	}

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/pets/%s", pathParam0)
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewFindPetByIdRequest generates requests for FindPetById
func NewFindPetByIdRequest(server string, id int64) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParam("simple", false, "id", id)
	if err != nil {
		return nil, err
	}

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/pets/%s", pathParam0)
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		if !strings.HasSuffix(baseURL, "/") {
			baseURL += "/"
		}
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

type findPetsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Pet
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r findPetsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r findPetsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type addPetResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Pet
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r addPetResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r addPetResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type deletePetResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r deletePetResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r deletePetResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type findPetByIdResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Pet
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r findPetByIdResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r findPetByIdResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// FindPetsWithResponse request returning *FindPetsResponse
func (c *ClientWithResponses) FindPetsWithResponse(ctx context.Context, params *FindPetsParams) (*findPetsResponse, error) {
	rsp, err := c.FindPets(ctx, params)
	if err != nil {
		return nil, err
	}
	return ParseFindPetsResponse(rsp)
}

// AddPetWithBodyWithResponse request with arbitrary body returning *AddPetResponse
func (c *ClientWithResponses) AddPetWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*addPetResponse, error) {
	rsp, err := c.AddPetWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseAddPetResponse(rsp)
}

func (c *ClientWithResponses) AddPetWithResponse(ctx context.Context, body AddPetJSONRequestBody) (*addPetResponse, error) {
	rsp, err := c.AddPet(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseAddPetResponse(rsp)
}

// DeletePetWithResponse request returning *DeletePetResponse
func (c *ClientWithResponses) DeletePetWithResponse(ctx context.Context, id int64) (*deletePetResponse, error) {
	rsp, err := c.DeletePet(ctx, id)
	if err != nil {
		return nil, err
	}
	return ParseDeletePetResponse(rsp)
}

// FindPetByIdWithResponse request returning *FindPetByIdResponse
func (c *ClientWithResponses) FindPetByIdWithResponse(ctx context.Context, id int64) (*findPetByIdResponse, error) {
	rsp, err := c.FindPetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return ParseFindPetByIdResponse(rsp)
}

// ParseFindPetsResponse parses an HTTP response from a FindPetsWithResponse call
func ParseFindPetsResponse(rsp *http.Response) (*findPetsResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &findPetsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Pet
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json"):
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseAddPetResponse parses an HTTP response from a AddPetWithResponse call
func ParseAddPetResponse(rsp *http.Response) (*addPetResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &addPetResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Pet
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json"):
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseDeletePetResponse parses an HTTP response from a DeletePetWithResponse call
func ParseDeletePetResponse(rsp *http.Response) (*deletePetResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &deletePetResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json"):
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseFindPetByIdResponse parses an HTTP response from a FindPetByIdWithResponse call
func ParseFindPetByIdResponse(rsp *http.Response) (*findPetByIdResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &findPetByIdResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Pet
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json"):
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /pets)
	FindPets(ctx echo.Context, params FindPetsParams) error

	// (POST /pets)
	AddPet(ctx echo.Context) error

	// (DELETE /pets/{id})
	DeletePet(ctx echo.Context, id int64) error

	// (GET /pets/{id})
	FindPetById(ctx echo.Context, id int64) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// FindPets converts echo context to params.
func (w *ServerInterfaceWrapper) FindPets(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params FindPetsParams
	// ------------- Optional query parameter "tags" -------------

	err = runtime.BindQueryParameter("form", true, false, "tags", ctx.QueryParams(), &params.Tags)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter tags: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.FindPets(ctx, params)
	return err
}

// AddPet converts echo context to params.
func (w *ServerInterfaceWrapper) AddPet(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.AddPet(ctx)
	return err
}

// DeletePet converts echo context to params.
func (w *ServerInterfaceWrapper) DeletePet(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameter("simple", false, "id", ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeletePet(ctx, id)
	return err
}

// FindPetById converts echo context to params.
func (w *ServerInterfaceWrapper) FindPetById(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameter("simple", false, "id", ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.FindPetById(ctx, id)
	return err
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}, si ServerInterface) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET("/pets", wrapper.FindPets)
	router.POST("/pets", wrapper.AddPet)
	router.DELETE("/pets/:id", wrapper.DeletePet)
	router.GET("/pets/:id", wrapper.FindPetById)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}
