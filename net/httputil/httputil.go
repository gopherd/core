package httputil

import (
	"bytes"
	"crypto"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/gopherd/core/typeconv"
)

var (
	ErrUnableToMarshalForm = errors.New("unable to marshal form")
	ErrSign                = errors.New("failed to verify sign")
)

const (
	Unused      = 0
	charsetUTF8 = "charset=utf-8"

	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJSONCharsetUTF8       = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; " + charsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; " + charsetUTF8
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEApplicationFormCharsetUTF8       = MIMEApplicationForm + "; " + charsetUTF8
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEApplicationMsgpack               = "application/msgpack"
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; " + charsetUTF8
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"
)

const (
	HeaderAccept                        = "Accept"
	HeaderAcceptEncoding                = "Accept-Encoding"
	HeaderAllow                         = "Allow"
	HeaderAuthorization                 = "Authorization"
	HeaderContentDisposition            = "Content-Disposition"
	HeaderContentEncoding               = "Content-Encoding"
	HeaderContentLength                 = "Content-Length"
	HeaderContentType                   = "Content-Type"
	HeaderCookie                        = "Cookie"
	HeaderSetCookie                     = "Set-Cookie"
	HeaderIfModifiedSince               = "If-Modified-Since"
	HeaderLastModified                  = "Last-Modified"
	HeaderLocation                      = "Location"
	HeaderUpgrade                       = "Upgrade"
	HeaderVary                          = "Vary"
	HeaderWWWAuthenticate               = "WWW-Authenticate"
	HeaderXForwardedProto               = "X-Forwarded-Proto"
	HeaderXHTTPMethodOverride           = "X-HTTP-Method-Override"
	HeaderXForwardedFor                 = "X-Forwarded-For"
	HeaderXRealIP                       = "X-Real-IP"
	HeaderServer                        = "Server"
	HeaderOrigin                        = "Origin"
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	HeaderStrictTransportSecurity = "Strict-Transport-Security"
	HeaderXContentTypeOptions     = "X-Content-Type-Options"
	HeaderXXSSProtection          = "X-XSS-Protection"
	HeaderXFrameOptions           = "X-Frame-Options"
	HeaderContentSecurityPolicy   = "Content-Security-Policy"
	HeaderXCSRFToken              = "X-CSRF-Token"
)

type Middleware interface {
	Apply(http.Handler) http.Handler
}

// MiddlewareFunc implements Middleware interface
type MiddlewareFunc func(http.Handler) http.Handler

func (m MiddlewareFunc) Apply(h http.Handler) http.Handler {
	return m(h)
}

type Result struct {
	Response   *http.Response
	StatusCode int
	Data       []byte
	Error      error
}

func (result Result) Ok() bool { return result.StatusCode == http.StatusOK }
func (result Result) Status() string {
	if result.Response == nil {
		if result.Error != nil {
			return result.Error.Error()
		}
		return ""
	}
	return result.Response.Status
}

func Get(url string) Result {
	resp, err := http.Get(url)
	return readResultFromResponse(resp, err)
}

func PostForm(url string, values url.Values) Result {
	resp, err := http.PostForm(url, values)
	return readResultFromResponse(resp, err)
}

func readResultFromResponse(resp *http.Response, err error) Result {
	result := Result{
		Response: resp,
		Error:    err,
	}
	if err != nil {
		return result
	}
	result.StatusCode = resp.StatusCode
	defer resp.Body.Close()
	result.Data, result.Error = ioutil.ReadAll(resp.Body)
	return result
}

type responseOptions struct {
	status      int
	acceptType  string
	contentType string
}

func newResponseOptions() *responseOptions {
	return &responseOptions{
		status:      http.StatusOK,
		contentType: MIMETextPlain,
	}
}

type ResponseOptions func(opts *responseOptions)

func mergeOptions(opts *responseOptions, options ...ResponseOptions) {
	for _, o := range options {
		o(opts)
	}
}

func WithStatus(status int) ResponseOptions {
	return func(opts *responseOptions) {
		opts.status = status
	}
}

func WithAcceptType(acceptType string) ResponseOptions {
	return func(opts *responseOptions) {
		opts.acceptType = acceptType
	}
}

func WithContentType(contentType string) ResponseOptions {
	return func(opts *responseOptions) {
		opts.contentType = contentType
	}
}

func Response(w http.ResponseWriter, body any, options ...ResponseOptions) error {
	var opts = newResponseOptions()
	mergeOptions(opts, options...)
	if body != nil {
		var marshaler MarshalFunc
		if strings.Contains(opts.contentType, MIMEApplicationJSON) {
			marshaler = json.Marshal
		} else if strings.Contains(opts.contentType, MIMEApplicationXML) {
			marshaler = xml.Marshal
		} else if strings.Contains(opts.contentType, MIMEApplicationForm) {
			marshaler = marshalForm
		} else {
			marshaler = typeconv.ToBytes
		}
		b, err := marshaler(body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
			return err
		}
		w.Header().Set("Content-Type", opts.contentType)
		w.WriteHeader(opts.status)
		_, err = w.Write(b)
		return err
	}
	w.WriteHeader(opts.status)
	return nil
}

func JSONResponse(w http.ResponseWriter, value any, options ...ResponseOptions) error {
	return Response(w, value, append(options, WithContentType(MIMEApplicationJSONCharsetUTF8))...)
}

func XMLResponse(w http.ResponseWriter, value any, options ...ResponseOptions) error {
	return Response(w, value, append(options, WithContentType(MIMEApplicationXMLCharsetUTF8))...)
}

func FormResponse(w http.ResponseWriter, value any, options ...ResponseOptions) error {
	return Response(w, value, append(options, WithContentType(MIMEApplicationFormCharsetUTF8))...)
}

func TextResponse(w http.ResponseWriter, value string, options ...ResponseOptions) error {
	return Response(w, value, append(options, WithContentType(MIMETextPlain))...)
}

type MarshalFunc func(any) ([]byte, error)

type FormMarshaler interface {
	MarshalForm() ([]byte, error)
}

func marshalForm(v any) ([]byte, error) {
	if marshaler, ok := v.(FormMarshaler); ok {
		return marshaler.MarshalForm()
	}
	return nil, ErrUnableToMarshalForm
}

type Signer func(key, keyField, signField string, args url.Values) (string, error)

type byKey [][2]string

func (by byKey) Len() int           { return len(by) }
func (by byKey) Less(i, j int) bool { return by[i][0] < by[j][0] }
func (by byKey) Swap(i, j int)      { by[i], by[j] = by[j], by[i] }

func HashSigner(hasher hash.Hash) Signer {
	return func(key, keyField, signField string, args url.Values) (string, error) {
		params := make([][2]string, 0, len(args))
		for k, vs := range args {
			if len(vs) > 0 && k != signField {
				params = append(params, [2]string{k, vs[0]})
			}
		}
		sort.Sort(byKey(params))
		var buf bytes.Buffer
		for _, pair := range params {
			buf.WriteString(pair[0])
			buf.WriteByte('=')
			buf.WriteString(pair[1])
			buf.WriteByte('&')
		}
		buf.WriteString(keyField)
		buf.WriteByte('=')
		buf.WriteString(key)
		hasher.Write(buf.Bytes())
		sum := hasher.Sum(nil)
		sign := fmt.Sprintf("%x", sum)
		return sign, nil
	}
}

func MD5Signer() Signer    { return HashSigner(crypto.MD5.New()) }
func SHA256Signer() Signer { return HashSigner(crypto.SHA256.New()) }

func VerifySign(signer Signer, key, keyField, signField string, args url.Values) error {
	expected, err := signer(key, keyField, signField, args)
	if err != nil {
		return err
	}
	if args == nil || args.Get(signField) != expected {
		return ErrSign
	}
	return nil
}

func PostFormJSON(httpc *http.Client, url_ string, data url.Values, res any) error {
	response, err := httpc.PostForm(url_, data)
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return err
	}
	if response.StatusCode >= 400 {
		return errors.New(response.Status)
	}
	if res == nil {
		return nil
	}
	return json.NewDecoder(response.Body).Decode(res)
}

type Pager struct {
	Total int64 `json:"total"`
	Curr  int   `json:"curr"`
	Limit int   `json:"limit"`
	List  any   `json:"list"`
}
