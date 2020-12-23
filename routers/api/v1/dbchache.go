package v1

import (
	"github.com/tanwen321/chatbygo/models"
	"github.com/tanwen321/chatbygo/pkg/logging"
	"github.com/tanwen321/chatbygo/pkg/setting"

	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type nowuser struct {
	user models.User

	// ucook  *http.Cookie
	uws *wsConnection
	uch chan struct{}
	// uctime time.Time
}

type nowgroup struct {
	group models.Group

	// ucook  *http.Cookie
	uws *wsConnection
	// uctime time.Time
}

// 在线用户token-》用户id
var onutot = make(map[string]int)

// 在线群组-》用户
var ongtou = make(map[int][]int)

// 在线用户-》群组
var onutog = make(map[int][]int)

// 在线用户-》好友
var onutof = make(map[int][]int)

//在线用户
var onuser = make(map[int]nowuser)

//在线群组
var ongroup = make(map[int]nowgroup)

//在线用户通道
var onuch = make(chan nowuser)

//下线用户接受
var offuch = make(chan int)

//头像缓存
var imagelist []os.FileInfo

func DbCache() {
	models.ClearToken()
	go online_status()
}

func online_status() {
	for {
		select {
		case u := <-onuch:
			id := u.user.Id
			us := models.GetUser(id)
			u.user = us
			onuser[id] = u
			ugroup := models.GetGroupidbyUid(id)
			onutog[id] = ugroup
			ufriend := models.GetFrendidbyUid(id)
			onutof[id] = ufriend
			for _, gid := range ugroup {
				_, ok := ongroup[gid]
				if !ok {
					gp := models.GetGroup(gid)
					ongroup[gid] = nowgroup{group: gp}
				}
				og, ok := ongtou[gid]
				if ok {
					ongtou[gid] = addpara(og, id)
				} else {
					ongtou[gid] = []int{id}
				}
			}
			close(u.uch)
		case id := <-offuch:
			_, ok := onuser[id]
			if ok {
				delete(onuser, id)
				delete(onutog, id)
				delete(onutof, id)
				models.OfflineUser(id)
				ugroup := models.GetGroupidbyUid(id)
				for _, gid := range ugroup {
					og, ok := ongtou[gid]
					if ok {
						ongtou[gid] = removepara(og, id)
					}
				}
			}
		}
	}
}

func removepara(sli []int, n int) []int {
	for i, u := range sli {
		if u == n {
			return append(sli[:i], sli[i+1:]...)
		}
	}
	return sli
}

func addpara(sli []int, n int) []int {
	for _, u := range sli {
		if u == n {
			return sli
		}
	}
	return append(sli, n)
}

func randimage() (s string) {
	if len(imagelist) == 0 {
		readimage()
	}
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(imagelist))
	s = imagelist[i].Name()
	return
}

func readimage() {
	dir := setting.WEB_HOME + "/touxiang"
	var err error
	imagelist, err = ioutil.ReadDir(dir)
	if err != nil {
		logging.Error("read touxiang dir error:%s", err)
	}
}
