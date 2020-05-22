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
    ```

    Or you can use third party packeages like [github.com/rs/cors](http://github.com/rs/cors) as middleware:
    ```go
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

* possible use predefined regex macros in uri definition:
    ```go
    // was:
    r.GET("/user/{user_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}", handler.User)
    // now:
    r.GET("/user/{user_id:@uuid@}", handler.User)
    
    // if you'd like to define your own macros do it like:
    gmb.Macros("alphabet", "[A-Z0-9]")
    r.GET("/invite/{code:@alphabet@}", handler.Invite)
    ```

    default macros:
    * `@uuid@`: `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`
    * `@num@`:  `[0-9]+`
