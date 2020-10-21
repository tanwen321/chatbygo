package gchat

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func index(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "text/html")
	f, err := ioutil.ReadFile(Whome + "/index.html")
	if err != nil {
		fmt.Fprintf(w, "页面找不到了")
		return
	}
	w.Write(f)
}

func login(w http.ResponseWriter, req *http.Request) {
	//	w.Header().Set("content-type", "text/html")
	req.ParseForm()
	p_u, _ := req.Form["user"]
	p_p, _ := req.Form["passwd"]

	if len(p_u) == 0 || len(p_p) == 0 {
		fmt.Fprintf(w, "请输入用户名和密码")
		return
	}

	id := chnti[p_u[0]]
	db_user, ok := userdb[id]
	if !ok {
		fmt.Fprintf(w, "用户名错误")
		return
	}
	if onuser[id] {
		fmt.Fprintf(w, "用户:%v已经登陆", p_u[0])
		return
	}
	if db_user.passwd != p_p[0] {
		fmt.Fprintf(w, "密码错误")
		return
	}
	//	c := md5.Sum([]byte(p_u))
	cookie := &http.Cookie{
		Name:   "mmssid",
		Value:  p_u[0],
		MaxAge: 300,
	}
	http.SetCookie(w, cookie)
	db_user.ucook = cookie
	userdb[id] = db_user
	//	onuch <- db_user
	fmt.Fprintf(w, "ok")
}
