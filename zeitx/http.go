package zeitx

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Useful http related consts
const (
	ContentType         = "Content-Type"
	XContentTypeOptions = "X-Content-Type-Options"
	XPrettyPrint        = "X-Pretty-Print"
	WWWAuthenticate     = "WWW-Authenticate"

	ApplicationJSON            = "application/json"
	ApplicationJSONCharsetUTF8 = "application/json; charset=utf-8"
	ApplicationOctetStream     = "application/octet-stream"
)

// OkJSON marshalls the data in json format and writes them to w with status 200
func OkJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	JSON(w, r, data, http.StatusOK)
}

// JSON marshalls the data in json format and writes them to w with http status
// if a "pretty" queryparam is present it pretty prints the output setting
// the response X-Pretty-Print header
func JSON(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	vars := r.URL.Query()
	pretty := vars["pretty"]
	var (
		js  []byte
		err error
	)
	start := time.Now()
	if pretty != nil {
		js, err = json.MarshalIndent(data, "", "    ")
		w.Header().Set(XPrettyPrint, "1")
	} else {
		js, err = json.Marshal(data)
	}
	if err != nil {
		APIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("Marshalled response in %fs", time.Since(start).Seconds())
	w.Header().Set(ContentType, ApplicationJSONCharsetUTF8)
	start = time.Now()
	w.WriteHeader(status)
	w.Write(js)
	log.Printf("Wrote response in %fs", time.Since(start).Seconds())
}

// Logger is a middleware to dump useful logging info about a request
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		x, _ := httputil.DumpRequest(r, true)
		f := func() {
			log.Printf("%s\t%s\t%fs\t%q", r.Method, r.RequestURI, time.Since(start).Seconds(), x)
		}
		defer f()
		next.ServeHTTP(w, r)
	})
}

// JSONEnforce is a middleware that enforces an "application/json" or
// "application/json; charset=utf-8" Content-Type
func JSONEnforce(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PATCH" || r.Method == "POST" || r.Method == "PUT" {
			if !hasContentType(r, ApplicationJSONCharsetUTF8) && !hasContentType(r, ApplicationJSON) {
				APIError(w, http.StatusUnsupportedMediaType, fmt.Sprintf("only %s is supported", ApplicationJSON))
				return
			}
			if r.Body == nil {
				APIError(w, http.StatusBadRequest, "null body not allowed")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func hasContentType(r *http.Request, mimetype string) bool {
	contentType := r.Header.Get(ContentType)
	if contentType == "" {
		return mimetype == ApplicationOctetStream
	}
	for _, v := range strings.Split(contentType, ",") {
		t, _, err := mime.ParseMediaType(v)
		if err != nil {
			break
		}
		if t == mimetype {
			return true
		}
	}
	return false
}

// BasicAuth is a middleware that manages basic auth authentication
func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(WWWAuthenticate, `Basic realm="Restricted"`)
		uname, hash, ok := r.BasicAuth()
		if !ok {
			Unauthorized(w, "Unauthorized")
		}
		if err := verify(uname, hash); err != nil {
			Unauthorized(w, "Unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	}
}

func verify(uname string, upass string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(upass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword(hash, []byte(upass)); err != nil {
		return err
	}
	log.Print(uname)
	return nil
}
