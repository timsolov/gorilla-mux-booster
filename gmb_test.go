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
	RegisterAliases(map[string]string{
		"{user_id}": "{user_id:@uuid@}",
	})
	router.GET("/{param:@uuid@}/ok", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("uuid:" + mux.Vars(r)["param"]))
	})
	router.GET("/{param:@code@}/ok", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("code:" + mux.Vars(r)["param"]))
	})
	router.POST("/{param}/ok", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("other:" + mux.Vars(r)["param"]))
	})
	router.GET("/users/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("user_id:" + mux.Vars(r)["user_id"]))
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
		{
			desc:   "/users/{user_id}",
			method: http.MethodGet,
			uri:    "/users/ae6f8837-c6e4-45cc-9e66-5e63bd749866",
			code:   http.StatusOK,
			resp:   "user_id:ae6f8837-c6e4-45cc-9e66-5e63bd749866",
		},
		{
			desc:   "/users/bad-user",
			method: http.MethodGet,
			uri:    "/users/bad-user",
			code:   http.StatusNotFound,
			resp:   "404 page not found\n",
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

func TestAlias(t *testing.T) {
	const LinksMuxPattern = `[0-9a-zA-Z]{8}-[0-9a-zA-Z]{4}-[0-9a-zA-Z]{4}-[0-9a-zA-Z]{4}-[0-9a-zA-Z]{12}`
	router := NewRouter()
	RegisterRegex("link", LinksMuxPattern)
	RegisterAliases(map[string]string{
		"{host_id}": "{host_id:@uuid@}",
		"{file_id}": "{file_id:@uuid@}",
		"{link_id}": "{link_id:@link@}",
	})
	router.GET("/{link_id}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("link_id:" + mux.Vars(r)["link_id"]))
	})
	router.GET("/link/{file_id}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("file_id:" + mux.Vars(r)["file_id"]))
	})

	testCases := []struct {
		desc   string
		method string
		uri    string
		code   int
		resp   string
	}{
		{
			desc:   "file_id",
			method: http.MethodGet,
			uri:    "/link/d8e0f69a-2680-4b09-9b05-d31a1306a1db",
			code:   http.StatusOK,
			resp:   "file_id:d8e0f69a-2680-4b09-9b05-d31a1306a1db",
		},
		{
			desc:   "link_id",
			method: http.MethodGet,
			uri:    "/2zpERrAb-apKS-8vuI-34NU-LXwOGJEsdcO6",
			code:   http.StatusOK,
			resp:   "link_id:2zpERrAb-apKS-8vuI-34NU-LXwOGJEsdcO6",
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

func TestGroup(t *testing.T) {
	router := NewRouter()
	router.GET("/index", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("index ok"))
	})
	group := router.Group("/group")
	group.GET("/index", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("group index ok"))
	})
	groupuuid := router.Group("/{@uuid@}")
	groupuuid.GET("/index", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("uuid index ok"))
	})

	cases := []struct {
		desc   string
		method string
		uri    string
		code   int
		resp   string
	}{
		{
			desc:   "404",
			method: http.MethodGet,
			uri:    "/group/index.html",
			code:   http.StatusNotFound,
			resp:   "",
		},
		{
			desc:   "index",
			method: http.MethodGet,
			uri:    "/index",
			code:   http.StatusOK,
			resp:   "index ok",
		},
		{
			desc:   "group index",
			method: http.MethodGet,
			uri:    "/group/index",
			code:   http.StatusOK,
			resp:   "group index ok",
		},
		{
			desc:   "uuid group index",
			method: http.MethodGet,
			uri:    "/ba0f8187-b82e-44a2-91dc-648c7f143bb5/index",
			code:   http.StatusOK,
			resp:   "uuid index ok",
		},
	}

	for _, tc := range cases {
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
