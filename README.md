# gorilla-mux-booster
The simple wrapper to improve features of [github.com/gorilla/mux](https://github.com/gorilla/mux)

### Features
* add support for idiomatic right way to add CORS and other middlewares over the routes.
    
    Old way:
    ```go
    r := mux.NewRouter()
    ...
    http.ListenAndServe(":8000", handlers.RecoveryHandler()(handlers.CORS(rules)(r)))
    ```
    New way:
    ```go
    r := gmb.NewRouter()
    r.UseOver(handlers.RecoveryHandler())
    r.UseOver(handlers.CORS(rules))
    // OR
    // r.UseOver(
    //   handlers.RecoveryHandler(),
    //   handlers.CORS(rules),
    // )
    http.ListenAndServe(":8000", r)
    ```

    Or you can use [github.com/rs/cors](http://github.com/rs/cors) as middleware:
    ```go
    r := gmb.NewRouter()
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://foo.com", "http://foo.com:8080"},
        AllowCredentials: true,
    })
    r.UseOver(c.Handler)
    http.ListenAndServe(":8000", r)
    ```

* add http methods as methods of route.
    Now you could mount handlers like below:
    ```go
    r := gmb.NewRouter()
    r.GET("/hello", handler.Hello)
    r.POST("/create", handler.Create)
    r.PUT("/edit", handler.Edit)
    http.ListenAndServe(":8000", r)
    ```

* possible use macros in uri definition
    ```go
    // was:
    r.GET("/user/{user_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}", handler.User)
    // become:
    r.GET("/user/{user_id:@uuid@}", handler.User)
    ```

    supported macros:
    * `@uuid@`: `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`
	* `@num@`:  `[0-9]+`