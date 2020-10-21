package gchat

import (
	//	"fmt"
	"log"
	"net/http"
)

// 在线群组-》用户
var ongtou = make(map[int][]user)

//在线用户
var onuser = make(map[int]bool)

// 用户name -》id
var chnti = make(map[string]int)

//在线用户通道
var onuch = make(chan user)

//下线用户接受
var offuch = make(chan user)

func Start_Server() {
	//	ws := NewWsServer()
	db_start()
	go online_status()
	//	go sync_db()
	http.HandleFunc("/", index)
	http.HandleFunc("/index.html", index)
	http.HandleFunc("/login.html", login)
	http.HandleFunc("/mmchat_socket", wsHandler)
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(Whome+"/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(Whome+"/css"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(Whome+"/images"))))
	http.Handle("/touxiang/", http.StripPrefix("/touxiang/", http.FileServer(http.Dir(Whome+"/touxiang"))))
	log.Fatal(http.ListenAndServe("10.68.61.6:8000", nil))
}

func online_status() {
	for id, u := range userdb {
		chnti[u.mname] = id
	}
	for {
		select {
		case u := <-onuch:
			//			fmt.Println("88888-----online:", u.mname)
			onuser[u.mid] = true
			for _, gid := range u.group {
				og, ok := ongtou[gid]
				if ok {
					ongtou[gid] = addpara(og, u)
				} else {
					ongtou[gid] = []user{u}
				}
			}
		case u := <-offuch:
			//			fmt.Println("666666-----offline:", u.mname)
			onuser[u.mid] = false
			for _, gid := range u.group {
				og, ok := ongtou[gid]
				if ok {
					ongtou[gid] = removepara(og, u)
				}
			}
		}
	}
}

func removepara(sli []user, n user) []user {
	for i, u := range sli {
		if u.mid == n.mid {
			return append(sli[:i], sli[i+1:]...)
		}
	}
	return sli
}

func addpara(sli []user, n user) []user {
	for _, u := range sli {
		if u.mid == n.mid {
			return sli
		}
	}
	return append(sli, n)
}
