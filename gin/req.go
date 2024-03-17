package gin

import (
	"mime/multipart"

	"github.com/spf13/cast"
)

type IRequest interface {
	DefaultQueryInt(key string, def int) (int, bool)
	DefaultQueryInt64(key string, def int64) (int64, bool)
	DefaultQueryFloat32(key string, def float32) (float32, bool)
	DefaultQueryFloat64(key string, def float64) (float64, bool)
	DefaultQueryBool(key string, def bool) (bool, bool)
	DefaultQueryString(key string, def string) (string, bool)
	DefaultQueryStringSlice(key string, def []string) ([]string, bool)
	DefaultQuery(key string) any

	DefaultParamInt(key string, def int) (int, bool)
	DefaultParamInt64(key string, def int64) (int64, bool)
	DefaultParamFloat32(key string, def float32) (float32, bool)
	DefaultParamFloat64(key string, def float64) (float64, bool)
	DefaultParamBool(key string, def bool) (bool, bool)
	DefaultParamString(key string, def string) (string, bool)
	DefaultParam(key string) any

	DefaultFormInt(key string, def int) (int, bool)
	DefaultFormInt64(key string, def int64) (int64, bool)
	DefaultFormFloat32(key string, def float32) (float32, bool)
	DefaultFormFloat64(key string, def float64) (float64, bool)
	DefaultFormBool(key string, def bool) (bool, bool)
	DefaultFormStringSlice(key string, def []string) ([]string, bool)
	DefaultFormFile(key string) (*multipart.FileHeader, error)
	DefaultForm(key string) any

	BindJson(obj any) error
	BindXml(obj any) error
	BindRaw(obj any) error

	Uri() string
	Method() string
	Host() string
	ClientIp() string

	Headers() map[string]string
	Header(key string) (string, bool)

	Cookies() map[string]string
	Cookie(key string) (string, bool)
}

func (ctx *Context) QueryAll() map[string][]string {
	ctx.initQueryCache()
	return map[string][]string(ctx.queryCache)
}

func (ctx *Context) DefaultQueryInt(key string, def int) (int, bool) {
	if v, ok := ctx.QueryAll()[key]; ok {
		if len(v) > 0 {
			return cast.ToInt(v[0]), true
		}
	}

	return def, false
}

func (ctx *Context) DefaultQueryInt64(key string, def int64) (int64, bool) {
	if v, ok := ctx.QueryAll()[key]; ok {
		if len(v) > 0 {
			return cast.ToInt64(v[0]), true
		}
	}

	return def, false
}

func (ctx *Context) DefaultQueryFloat32(key string, def float32) (float32, bool) {
	if v, ok := ctx.QueryAll()[key]; ok {
		if len(v) > 0 {
			return cast.ToFloat32(v[0]), true
		}
	}

	return def, false
}

func (ctx *Context) DefaultQueryFloat64(key string, def float64) (float64, bool) {
	if v, ok := ctx.QueryAll()[key]; ok {
		if len(v) > 0 {
			return cast.ToFloat64(v[0]), true
		}
	}

	return def, false
}

func (ctx *Context) DefaultQueryBool(key string, def bool) (bool, bool) {
	if v, ok := ctx.QueryAll()[key]; ok {
		if len(v) > 0 {
			return cast.ToBool(v[0]), true
		}
	}

	return def, false
}

func (ctx *Context) DefaultQueryString(key string, def string) (string, bool) {
	if v, ok := ctx.QueryAll()[key]; ok {
		if len(v) > 0 {
			return cast.ToString(v[0]), true
		}
	}

	return def, false
}

func (ctx *Context) DefaultQueryStringSlice(key string, def []string) ([]string, bool) {
	if v, ok := ctx.QueryAll()[key]; ok {
		if len(v) > 0 {
			return cast.ToStringSlice(v[0]), true
		}
	}

	return def, false
}

func (ctx *Context) ParamVal(key string) any {
	if v, ok := ctx.Params.Get(key); ok {
		return v
	}

	return nil
}

func (ctx *Context) DefaultParamInt(key string, def int) (int, bool) {
	if v := ctx.ParamVal(key); v != nil {
		return cast.ToInt(v), true
	}

	return def, false
}

func (ctx *Context) DefaultParamInt64(key string, def int64) (int64, bool) {
	if v := ctx.ParamVal(key); v != nil {
		return cast.ToInt64(v), true
	}

	return def, false
}

func (ctx *Context) DefaultParamFloat32(key string, def float32) (float32, bool) {
	if v := ctx.ParamVal(key); v != nil {
		return cast.ToFloat32(v), true
	}

	return def, false
}

func (ctx *Context) DefaultParamFloat64(key string, def float64) (float64, bool) {
	if v := ctx.ParamVal(key); v != nil {
		return cast.ToFloat64(v), true
	}

	return def, false
}

func (ctx *Context) DefaultParamBool(key string, def bool) (bool, bool) {
	if v := ctx.ParamVal(key); v != nil {
		return cast.ToBool(v), true
	}

	return def, false
}

func (ctx *Context) DefaultParamString(key string, def string) (string, bool) {
	if v := ctx.ParamVal(key); v != nil {
		return cast.ToString(v), true
	}

	return def, false
}

func (ctx *Context) FormAll() map[string][]string {
	ctx.initFormCache()

	return map[string][]string(ctx.formCache)
}

func (ctx *Context) DefaultFormInt(key string, def int) (int, bool) {
	if v, ok := ctx.FormAll()[key]; ok {
		if len(v) > 0 {
			return cast.ToInt(v[0]), true
		}
	}

	return def, false
}

func (ctx *Context) DefaultFormInt64(key string, def int64) (int64, bool) {
	if v, ok := ctx.FormAll()[key]; ok {
		if len(v) > 0 {
			return cast.ToInt64(v[0]), true
		}
	}

	return def, false
}

func (ctx *Context) DefaultFormFloat32(key string, def float32) (float32, bool) {
	if v, ok := ctx.FormAll()[key]; ok {
		if len(v) > 0 {
			return cast.ToFloat32(v[0]), true
		}
	}

	return def, false
}

func (ctx *Context) DefaultFormFloat64(key string, def float64) (float64, bool) {
	if v, ok := ctx.FormAll()[key]; ok {
		if len(v) > 0 {
			return cast.ToFloat64(v[0]), true
		}
	}

	return def, false
}

func (ctx *Context) DefaultFormBool(key string, def bool) (bool, bool) {
	if v, ok := ctx.FormAll()[key]; ok {
		if len(v) > 0 {
			return cast.ToBool(v[0]), true
		}
	}

	return def, false
}

func (ctx *Context) DefaultFormStringSlice(key string, def []string) ([]string, bool) {
	if v, ok := ctx.FormAll()[key]; ok {
		if len(v) > 0 {
			return cast.ToStringSlice(v[0]), true
		}
	}

	return def, false
}

func (ctx *Context) DefaultForm(key string) any {
	if v, ok := ctx.FormAll()[key]; ok {
		if len(v) > 0 {
			return v[0]
		}
	}

	return nil
}
