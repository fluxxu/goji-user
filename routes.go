package user

import (
	"github.com/fluxxu/goji-rbac"
	"github.com/fluxxu/util"
	"github.com/zenazn/goji/web"
	"net/http"
)

//create root user and admin role
func routeInit(c web.C, w http.ResponseWriter, r *http.Request) {
	var count int
	err := opts.Dbx.Get(&count, "SELECT COUNT(*) FROM user")
	if err != nil {
		util.Response(w).Error(err.Error())
		return
	}

	if count != 0 {
		util.Response(w).Error("not allowed", 405)
		return
	}

	type req struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"display_name"`
	}

	var body req
	if err = util.Request(r).DecodeBody(&body); err != nil {
		util.Response(w).Error("decode body:"+err.Error(), 400)
		return
	}

	u := NewUser()
	u.Id = 1
	u.Email = body.Email
	u.SetPassword(body.Password)
	u.DisplayName = body.DisplayName

	if err = u.Insert(); err != nil {
		ve, ok := err.(util.ValidationErrorInterface)
		if ok {
			util.Response(w).Error("validation error", 400, map[string]interface{}{
				"errors": ve.ValidationErrors(),
			})
			return
		}
		util.Response(w).Error("can not save: " + err.Error())
		return
	}

	util.Response(w).Send(200, u)
}

func routePost(c web.C, w http.ResponseWriter, r *http.Request) {
	req := util.Request(r)
	res := util.Response(w)
	type reqBody struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"display_name"`
	}
	var body reqBody
	err := req.DecodeBody(&body)
	if err != nil {
		res.Error("can not parse body: " + err.Error())
		return
	}

	u := NewUser()
	u.Email = body.Email
	u.Password = body.Password
	u.DisplayName = body.DisplayName
	if err = u.Insert(); err != nil {
		if ve, ok := err.(util.ValidationErrorInterface); ok {
			res.Error(err.Error(), 400, ve.ToResponseData())
			return
		}
		res.Error(err.Error())
		return
	}

	res.Send(200, u)
}

func routeList(c web.C, w http.ResponseWriter, r *http.Request) {
	res := util.Response(w)
	p := NewUserProvider()
	p.ParseRequest(r)
	data := []User{}
	resData, err := p.Read(&data)
	if err != nil {
		res.Error(err.Error())
		return
	}
	res.Send(200, resData)
}
