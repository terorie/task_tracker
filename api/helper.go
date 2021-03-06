package api

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Request struct {
	Ctx *fasthttp.RequestCtx
}

func (r *Request) OkJson(object interface{}) {

	resp, err := json.Marshal(object)
	handleErr(err, r)

	r.Ctx.Response.Header.Set("Content-Type", "application/json")
	_, err = r.Ctx.Write(resp)
	handleErr(err, r)
}

func (r *Request) Json(object interface{}, code int) {

	resp, err := json.Marshal(object)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"code": code,
		}).Error("Error during json encoding of object")
	}

	r.Ctx.Response.SetStatusCode(code)
	r.Ctx.Response.Header.Set("Content-Type", "application/json")
	_, err = r.Ctx.Write(resp)
	if err != nil {
		panic(err) //todo handle differently
	}

}

func (r *Request) GetJson(x interface{}) bool {

	err := json.Unmarshal(r.Ctx.Request.Body(), x)
	handleErr(err, r)

	return err == nil
}
