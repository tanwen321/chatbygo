package gchat

import (
	//	"fmt"
	"io/ioutil"
	//	"net"
	"net/http"
	//	"os"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	Whome    = "/Users/tanwen/Desktop/G_work/gchatv2/priv"
	uini     = Whome + "/../data/udb.ini"
	gini     = Whome + "/../data/gdb.ini"
	Time_out = 60
)

type user struct {
	mid     int
	mname   string
	uimage  string
	passwd  string
	group   []int
	friends []int
	unkname string
	uword   string
	ucook   *http.Cookie
	uws     *wsConnection
	//	ucon    net.Conn
	uctime time.Time
}

type udb map[int]user

type group struct {
	gid    int
	gname  string
	gimage string
	gword  string
	gctime time.Time
}

type gdb map[int]group

var groupdb = make(gdb)
var userdb = make(udb)

var pngtmp = make(map[string]string)

func db_start() {
	userdb = init_user_db()
	groupdb = init_group_db()
}

func init_user_db() udb {
	var db = make(udb)
	admin := user{mid: 10000, mname: "tanwen", uimage: rand_img(), passwd: "admin@123", group: []int{20000}, unkname: "管理员", uword: "this is a game."}
	fb, err := ioutil.ReadFile(uini)
	if err != nil {
		return udb{10000: admin}
	}
	lines := strings.Split(string(fb), "\n")
	for _, line := range lines {
		u := strings.Split(line, " ")
		if len(u) >= 7 {
			id, err := strconv.Atoi(u[0])
			if err == nil {
				u1 := user{mid: id, mname: u[1], uimage: rand_img(), passwd: u[2], group: s_to_list(u[3]), friends: s_to_list(u[4]), unkname: u[5], uword: u[6]}
				db[id] = u1
			}
		}
	}
	if len(db) == 0 {
		return udb{10000: admin}
	}
	return db
}

func init_group_db() gdb {
	var db = make(gdb)
	gadmin := group{gid: 20000, gname: "默认小组", gimage: rand_img(), gword: "iPhone12系列刘海变小实锤了"}
	fb, err := ioutil.ReadFile(gini)
	if err != nil {
		return gdb{20000: gadmin}
	}
	lines := strings.Split(string(fb), "\n")
	for _, line := range lines {
		g := strings.Split(line, " ")
		if len(g) >= 3 {
			id, err := strconv.Atoi(g[0])
			if err == nil {
				g1 := group{gid: id, gname: g[1], gimage: rand_img(), gword: g[2]}
				db[id] = g1
			}
		}
	}
	if len(db) == 0 {
		return gdb{20000: gadmin}
	}
	return db
}

func s_to_list(s string) []int {
	var list []int
	ids := strings.Split(s, ",")
	for _, id := range ids {
		i, err := strconv.Atoi(id)
		if err == nil {
			list = append(list, i)
		}
	}
	if len(list) == 0 {
		return []int{1}
	}
	return list
}

func list_to_s(l []int) string {
	var s string
	j := len(l) - 1
	for i := 0; i < j; i++ {
		s += (strconv.Itoa(l[i]) + ",")
	}
	return s + strconv.Itoa(l[j])
}

func rand_img() string {
	finfo, err := ioutil.ReadDir(Whome + "/touxiang")
	if err != nil {
		return "/touxiang/icon01.png"
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(finfo))
	if finfo[n].Name() == ".DS_Store" {
		return "/touxiang/tx001.jpg"
	}
	return "/touxiang/" + finfo[n].Name()
}
