# gorilla-mux-booster
The simple wrapper to improve features of awesome [github.com/gorilla/mux](https://github.com/gorilla/mux). The package's gmb.NewRouter fully compatible with mux.NewRouter and you can start use it simply by replacing path in your import section to `github.com/timsolov/gorilla-mux-booster`.

### Features
* add support for idiomatic right way to add CORS and other middlewares over the routes.
    ```go
    // was:
    r := mux.NewRouter()
    http.ListenAndServe(":8000", handlers.RecoveryHandler()(handlers.CORS(rules)(r)))
    // now:
    r := gmb.NewRouter()
    r.UseOver(handlers.RecoveryHandler())
    r.UseOver(handlers.CORS(rules))
    // or inline:
    r.UseOver(handlers.RecoveryHandler(), handlers.CORS(rules))
    
    // or you can use third party packages like `github.com/rs/cors` as middleware:
    r := gmb.NewRouter()
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://foo.com", "http://foo.com:8080"},
        AllowCredentials: true,
    })
    r.UseOver(c.Handler)
    ```

* mount routes by http method name:
    ```go
    // was:
    r := mux.NewRouter()
    r.HandleFunc("/hello", handler.Hello).Methods("GET")
    r.HandleFunc("/create", handler.Create).Methods("POST")
    r.HandleFunc("/edit", handler.Edit).Methods("PUT")
    // now:
    r := gmb.NewRouter()
    r.GET("/hello", handler.Hello)
    r.POST("/create", handler.Create)
    r.PUT("/edit", handler.Edit)
    ```

* possible use regex contractions in uri definition:
    ```go
    // was:
    r.GET("/user/{user_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}", handler.User) // UUID pattern
    // now:
    r.GET("/user/{user_id:@uuid@}", handler.User) // @uuid@ is a predefined contraction
    
    // if you'd like to define your own contraction do it like:
    gmb.RegisterRegex("md5", "[a-f0-9]{32}")
    r.GET("/invite/{code:@md5@}", handler.Invite)
    ```

    predefined default contractions:
    * `@uuid@`:     `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`
    * `@num@`:      `[0-9]+`
    * `@alpha@`:    `[a-zA-Z]+`,
	* `@alphanum@`: `[a-zA-Z0-9]+`,
