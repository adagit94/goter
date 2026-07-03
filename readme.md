# Goter

Router implemented in Go. Agnostic of specific http implementation. There is faster, static variant, and dynamic, but slower one.

## Examples
Little wrapper is needed for usage in context of specific http implementation like fasthttp.

### fasthttp
    package server

    import (
      "github.com/adagit94/goter"
      "github.com/valyala/fasthttp"
    )

    type reqHandler func(ctx *fasthttp.RequestCtx, params goter.IParams)

    func main() {
      router := goter.CreateRouter[reqHandler]()

      router.Route("a/:b").Get(func(ctx *fasthttp.RequestCtx, params goter.IParams) {
        b := params.Path("b")
      })

      fasthttp.ListenAndServe("localhost:4000", func(ctx *fasthttp.RequestCtx) {
        handler, params := router.Select(string(ctx.Path()), string(ctx.Request.URI().QueryString()), string(ctx.Method()))
        handler(ctx, params)
      })
    }