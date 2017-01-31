package admin

import (
	"io/ioutil"

	"github.com/VirgilSecurity/virgild/config"
	"github.com/valyala/fasthttp"
)

type AdminHandlers struct {
	Index fasthttp.RequestHandler
}

func Init(conf *config.App) *AdminHandlers {
	return &AdminHandlers{
		Index: mainPage("./templates/index.html"),
	}
}

func mainPage(path string) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		f, err := ioutil.ReadFile(path)
		if err != nil {
			ctx.Error("", fasthttp.StatusInternalServerError)
			return
		}
		ctx.SetContentType("text/html; charset=UTF-8")
		ctx.Write(f)
	}
}
