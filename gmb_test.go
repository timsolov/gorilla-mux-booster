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
	RegisterRegex("code", "[0-9A-Z]+")
	router.GET("/{param:@uuid@}/ok", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("uuid:" + mux.Vars(r)["param"]))
	})
	router.GET("/{param:@code@}/ok", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("code:" + mux.Vars(r)["param"]))
	})
	router.POST("/{param}/ok", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("other:" + mux.Vars(r)["param"]))
	})

	testCases := []struct {
		desc   string
		method string
		uri    string
		code   int
		resp   string
	}{
		{
			desc:   "uuid",
			method: http.MethodGet,
			uri:    "/8942d16c-9ba2-4569-8ecf-3f5a7773b90c/ok",
			code:   http.StatusOK,
			resp:   "uuid:8942d16c-9ba2-4569-8ecf-3f5a7773b90c",
		},
		{
			desc:   "code",
			method: http.MethodGet,
			uri:    "/JJQQWUE123123/ok",
			code:   http.StatusOK,
			resp:   "code:JJQQWUE123123",
		},
		{
			desc:   "not uuid and not code",
			method: http.MethodGet,
			uri:    "/hi/ok",
			code:   http.StatusMethodNotAllowed,
		},
		{
			desc:   "param",
			method: http.MethodPost,
			uri:    "/hi/ok",
			code:   http.StatusOK,
			resp:   "other:hi",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			// request
			r := httptest.NewRequest(tc.method, tc.uri, nil)
			// response
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)
			if !assert.Equal(t, tc.code, w.Result().StatusCode) {
				t.Log(w.Body.String())
			}

			if len(tc.resp) > 0 {
				assert.Equal(t, tc.resp, w.Body.String())
			}

			router.ServeHTTP(w, r)
		})
	}
}
