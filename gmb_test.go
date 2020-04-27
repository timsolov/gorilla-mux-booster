package gmb

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMacros(t *testing.T) {
	router := NewRouter()
	router.GET("/{param:@uuid@}/ok", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(mux.Vars(r)["param"]))
	})
	router.POST("/{param}/ok", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(mux.Vars(r)["param"]))
	})

	testCases := []struct {
		desc   string
		method string
		uri    string
		code   int
	}{
		{
			desc:   "uuid",
			method: http.MethodGet,
			uri:    "/8942d16c-9ba2-4569-8ecf-3f5a7773b90c/ok",
			code:   http.StatusOK,
		},
		{
			desc:   "not uuid",
			method: http.MethodGet,
			uri:    "/hi/ok",
			code:   http.StatusMethodNotAllowed,
		},
		{
			desc:   "param",
			method: http.MethodPost,
			uri:    "/hi/ok",
			code:   http.StatusOK,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// request
			r := httptest.NewRequest(tC.method, tC.uri, nil)
			// response
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)
			if !assert.Equal(t, tC.code, w.Result().StatusCode) {
				t.Log(w.Body.String())
			}

			router.ServeHTTP(w, r)
		})
	}
}
