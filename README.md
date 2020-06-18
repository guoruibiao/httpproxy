# httpproxy

a simple http proxy for helping me request at office network.

![demo](/demo.jpg)

## supported methods

- [x] `GET`
- [x] `POST`
- [x] `HEAD`
- [x] `DELETE`
- [ ] `OPTIONS`

## usage

```
[work@host90 tools]$ ./httpproxy -port 80
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /proxy                    --> main.Post (3 handlers)
[GIN-debug] Listening and serving HTTP on :80
```

The key parameter is `url`, which is the target URL you want to proxy.
And the `METHOD` will keep in the same with the origin request.Any other parameters 
will be send to `url` with no modified.

There can be many features to be added into this repo, so maybe i can do it better 
in the future.

And now, just enjoy it.
