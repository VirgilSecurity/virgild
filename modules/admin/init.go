package admin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/VirgilSecurity/virgild/config"
	"github.com/valyala/fasthttp"
)

type AdminHandlers struct {
	Index        fasthttp.RequestHandler
	GetConfig    fasthttp.RequestHandler
	UpdateConfig fasthttp.RequestHandler
	CardInfo     fasthttp.RequestHandler
}

func Init(conf *config.App) *AdminHandlers {
	return &AdminHandlers{
		Index:        mainPage("./templates/index.html"),
		GetConfig:    getConf(conf.Common.ConfigUpdate),
		CardInfo:     getVirgilDCardInfo(conf.Site.VirgilD),
		UpdateConfig: updateConf(conf.Common.ConfigUpdate),
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

func getConf(updater *config.Updater) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		b, err := json.Marshal(updater.Config())
		if err != nil {
			ctx.Error("", fasthttp.StatusInternalServerError)
			return
		}
		ctx.SetContentType("application/json")
		ctx.Write(b)
	}
}

func updateConf(updater *config.Updater) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")
		conf := new(config.Config)
		err := json.Unmarshal(ctx.PostBody(), conf)
		if err != nil {
			ctx.Error("{'message':'JSON invalid'}", fasthttp.StatusBadRequest)
			return
		}
		err = updater.Update(*conf)
		if err != nil {
			ctx.Error(fmt.Sprintf("{'message':'%v'}", err), fasthttp.StatusBadRequest)
		}
	}
}

func getVirgilDCardInfo(info config.VirgilDCard) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		b, err := json.Marshal(info)
		if err != nil {
			ctx.Error("", fasthttp.StatusInternalServerError)
			return
		}
		ctx.SetContentType("application/json")
		ctx.Write(b)
	}
}
