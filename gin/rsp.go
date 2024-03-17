package gin

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/imajinyun/goframe/gin/internal/json"
)

type IResponse interface {
	ToJson(obj any) IResponse
	ToJsonp(obj any) IResponse
	ToXml(obj any) IResponse
	ToHtml(file string, obj any) IResponse
	ToText(format string, values ...any) IResponse
	ToRedirect(path string) IResponse
	ToSetHeader(key string, val string) IResponse
	ToSetCookie(key string,
		val string,
		maxAge int,
		path string,
		domain string,
		secure bool,
		httpOnly bool,
	) IResponse
	ToSetStatus(code int) IResponse
	ToSetOkStatus() IResponse
}

func (ctx *Context) ToJson(obj any) IResponse {
	byt, err := json.Marshal(obj)
	if err != nil {
		return ctx.ToSetStatus(http.StatusInternalServerError)
	}
	ctx.ToSetHeader("Content-Type", "application/json")
	ctx.Writer.Write(byt)

	return ctx
}

func (ctx *Context) ToJsonp(obj any) IResponse {
	fn := ctx.Query("callback")
	ctx.ToSetHeader("Content-Type", "application/javascript")
	callback := template.JSEscapeString(fn)

	_, err := ctx.Writer.Write([]byte(callback))
	if err != nil {
		return ctx
	}

	_, err = ctx.Writer.Write([]byte("("))
	if err != nil {
		return ctx
	}

	res, err := json.Marshal(obj)
	if err != nil {
		return ctx
	}
	_, err = ctx.Writer.Write(res)
	if err != nil {
		return ctx
	}
	// 输出右括号
	_, err = ctx.Writer.Write([]byte(")"))
	if err != nil {
		return ctx
	}

	return ctx
}

func (ctx *Context) ToXml(obj any) IResponse {
	byt, err := xml.Marshal(obj)
	if err != nil {
		return ctx.ToSetStatus(http.StatusInternalServerError)
	}
	ctx.ToSetHeader("Content-Type", "application/html")
	ctx.Writer.Write(byt)

	return ctx
}

func (ctx *Context) ToHtml(file string, obj any) IResponse {
	tpl, err := template.New("output").ParseFiles(file)
	if err != nil {
		return ctx
	}

	if err := tpl.Execute(ctx.Writer, obj); err != nil {
		return ctx
	}

	ctx.ToSetHeader("Content-Type", "application/html")

	return ctx
}

func (ctx *Context) ToText(format string, values ...any) IResponse {
	out := fmt.Sprintf(format, values...)
	ctx.ToSetHeader("Content-Type", "application/text")
	ctx.Writer.Write([]byte(out))

	return ctx
}

func (ctx *Context) ToRedirect(path string) IResponse {
	http.Redirect(ctx.Writer, ctx.Request, path, http.StatusMovedPermanently)

	return ctx
}

func (ctx *Context) ToSetHeader(key string, val string) IResponse {
	ctx.Writer.Header().Add(key, val)

	return ctx
}

func (ctx *Context) ToSetCookie(
	key string,
	val string,
	maxAge int,
	path string,
	domain string,
	secure bool,
	httpOnly bool,
) IResponse {
	if path == "" {
		path = "/"
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     key,
		Value:    url.QueryEscape(val),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: 1,
		Secure:   secure,
		HttpOnly: httpOnly,
	})

	return ctx
}

func (ctx *Context) ToSetStatus(code int) IResponse {
	ctx.Writer.WriteHeader(code)

	return ctx
}

func (ctx *Context) ToSetOkStatus() IResponse {
	ctx.Writer.WriteHeader(http.StatusOK)

	return ctx
}
