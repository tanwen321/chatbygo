package models

type User struct {
	Id       int    `json:"user_id" gorm:"index"`
	Username string `json:"user"`
	Password string `json:"pass"`
	Urool    int    `json:"urool"`
	Nickname string `json:"nickname"`
	Uimage   string `json:"uimage"`
	Usinge   string `json:"usinge"`
	Status   int    `json:"status"`
}

type Group struct {
	Id       int    `json:"group_id" gorm:"index"`
	Nickname string `json:"nickname"`
	Gimage   string `json:"gimage"`
	Usinge   string `json:"usinge"`
	Status   int    `json:"status"`
}

type UserToken struct {
	Userid int `gorm:"index"`
	Utoken string
	Status int `json:"status"`
}

type UserGroup struct {
	Userid  int
	Groupid int
	Status  int `json:"status"`
}

type UserFriend struct {
	Userid  int
	Frendid int
	Status  int `json:"state"`
}

func GetUser(id int) (user User) {
	db.First(&user, id)
	return
}

func GetGroup(id int) (group Group) {
	db.First(&group, id)
	return
}

func CheckAuth(username, password string) (int, bool) {
	var user User
	db.Select("id").Where(User{Username: username, Password: password}).Find(&user)

	if user.Id > 0 {
		var ut UserToken
		db.Where("userid = ?", user.Id).Find(&ut)
		if ut.Userid == 0 {
			return user.Id, true
		}
		return ut.Userid, false
	}
	return 0, false
}

func OnlineUser(userid int, token string) bool {
	var ut UserToken

	db.Where("userid = ?", userid).Find(&ut)
	if ut.Userid == 0 {
		db.Save(&UserToken{Userid: userid, Utoken: token, Status: 1})
		return true
	}
	return false
}

func GetUserBytoken(s string) (int, bool) {
	var ut UserToken
	db.Where("Utoken = ?", s).Find(&ut)
	if ut.Userid > 0 {
		return ut.Userid, true
	}
	return 0, false
}

func GetGroupidbyUid(uid int) (gids []int) {
	var ugs []UserGroup
	db.Select("groupid").Where("Userid = ?", uid).Find(&ugs)
	for _, ug := range ugs {
		gids = append(gids, ug.Groupid)
	}
	return
}

func GetFrendidbyUid(uid int) (fids []int) {
	var ufs []UserFriend
	db.Select("frendid").Where("Userid = ?", uid).Find(&ufs)
	for _, uf := range ufs {
		fids = append(fids, uf.Frendid)
	}
	return
}

func OfflineUser(uid int) {
	var ut UserToken
	db.Where("userid = ?", uid).Find(&ut)
	if ut.Userid > 0 {
		db.Where("userid = ?", uid).Delete(UserToken{})
	}
}

func CheckAdmin(uid int) bool {
	var user User
	db.First(&user, uid)
	if user.Id > 0 && user.Urool == 0 {
		return true
	}
	return false
}

func ExistUserByName(name string) int {
	var user User
	db.Select("id").Where("username = ?", name).Find(&user)
	return user.Id
}

func ExistGroupByName(name string) int {
	var group Group
	db.Select("id").Where("nickname = ?", name).Find(&group)
	return group.Id
}

func AddUser(data map[string]interface{}) {
	s := data["touxiang"].(string)
	image := "/touxiang/" + s
	db.Create(&User{
		Username: data["name"].(string),
		Password: data["pass"].(string),
		Nickname: data["nikename"].(string),
		Usinge:   data["words"].(string),
		Urool:    data["rool"].(int),
		Uimage:   image,
		Status:   data["status"].(int),
	})
}

func AddGroup(data map[string]interface{}) {
	s := data["touxiang"].(string)
	image := "/touxiang/" + s
	db.Create(&Group{
		Nickname: data["name"].(string),
		Usinge:   data["words"].(string),
		Gimage:   image,
		Status:   data["status"].(int),
	})
}

func AddUserToGroup(uid, gid int) bool {
	var ug UserGroup
	db.Select("userid").Where(UserGroup{Userid: uid, Groupid: gid}).First(&ug)
	if ug.Userid > 0 {
		return false
	}
	db.Create(&UserGroup{
		Userid:  uid,
		Groupid: gid,
		Status:  1,
	})
	return true
}

func AddUserToFriend(uid, fid int) bool {
	var uf UserFriend
	db.Select("userid").Where(UserFriend{Userid: uid, Frendid: fid}).First(&uf)
	if uf.Userid > 0 {
		return false
	}
	db.Create(&UserFriend{
		Userid:  uid,
		Frendid: fid,
		Status:  1,
	})
	db.Create(&UserFriend{
		Userid:  fid,
		Frendid: uid,
		Status:  1,
	})
	return true
}

func ClearToken() {
	db.Where("userid > ?", 0).Delete(UserToken{})
}
