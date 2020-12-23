//枚举，代表三个页面
var Pages = {
    //登录页面
    login: 0,
    //联系人
    contacts: 1,
    //私聊界面
    privateChat: 2,
    //群聊界面
    groupChat: 3
};

var Msgtype = {
    //普通消息
    firtype:1,
    grptype:2,
    ontype:3,
    offtype:4
}

//标记所处的页面和当前聊天的好友
var currentPage = {
    page: Pages.login,
    currentUser: null,
    currentChatFriend: null,
    currentChatGroup: null
};

var lockReconnect = false;  //避免ws重复连接
var ws = null;          // 判断当前浏览器是否支持WebSocket
var wsUrl = "ws://" + location.host + "/api/v1/mmchat_socket";
var hei;
var wstoken = ""
var chatService;

//为了避免controller层侵入到service层，用controller层的代码覆盖service层的预置方法
// chatService.onNewPrivateMessageReceive = onNewPrivateMessageReceive;
// chatService.onPrivateHistoryLoad = onPrivateHistoryLoad;
// chatService.onNewGroupMessageReceive = onNewGroupMessageReceive;
// chatService.onGroupHistoryLoad = onGroupHistoryLoad;
// chatService.onFriendListChange = onFriendListChange;
// chatService.onGroupListChange = onGroupListChange;

//登录
// function login() {
//     var username = $("#username").val();
//     var password = $("#password").val();
//     var formdata = new FormData();
//     formdata.append('user', username);
//     formdata.append('passwd', password);
//     var xhr = new XMLHttpRequest();
//     xhr.timeout = 3000;
//     xhr.responseType = "text"
//     xhr.open('POST', 'http://10.68.61.9/login.html', true);
//     xhr.onload = function(e) {
//         if (this.status == 200 || this.status == 304) {
//             //初始化GoEasy和本地好友列表
//             chatService.connectIM();
//             $("#login-box").hide();
//             //切换到联系人
//             switchToContacts();
//         } else {
//         //显示错误信息
//             $("#login-error-box").show();
//         }
//     }
// }

function login() {
    var username = $("#username").val();
    var password = $("#password").val();
    auth={user:username, passwd:password};
    var resultStatus = false;
    $.ajaxSettings.async = false;
    $.ajax({
        type:'POST',
        traditional:true,
        // data:JSON.stringify(auth),
        // contentType :'application/json',
        // dataType:'json',
        data:auth,
        url:'login.html',
        success :function(data) {
            // console.log(data);
            if (data.msg=="ok") {
                wstoken = data.data.token
                resultStatus=true
            } else {
                $("#login-error-msg").html(data.msg)
            }
        },
        error :function(e) {
            alert("error"); 
        } 
    });
    if (resultStatus) {
        //初始化GoEasy和本地好友列表
        createWebSocket(wsUrl);
        $("#login-box").hide();
        //切换到联系人
//        switchToContacts();
    } else {
        //显示错误信息
        $("#login-error-box").show();
    }
}


//切换到联系人界面
function switchToContacts() {
//    chatService = data
    //设置当前页面为好友列表页面
    currentPage.page = Pages.contacts;
    //修改页面的title
    $("title").text("联系人");
    $('.login').hide()

    //初始化当前用户
    currentPage.currentUser = chatService.currentUser[0];
//    currentPage.currentUser = currentUser;
    $("#current-user-avatar").attr("src", currentPage.currentUser.avatar);
    $("#current-user-name").text(currentPage.currentUser.name);

    //隐藏聊天窗口
    $("#chat-box").hide();
    $("#group-chat-box").hide();

    //获取好友数据
    var friendsMap = chatService.friends;
    //绘制好友列表
    renderFriendList(friendsMap);
    //显示好友列表
    $("#friends-box").show();

    //获取群数据
    var groupsMap = chatService.groups;
    //绘制群列表
    renderGroupList(groupsMap);
    //显示群列表
    $("#group-box").show()
}

//绘制好友列表
function renderFriendList(friends) {
    var friendListDiv = $("#friend-list");
    friendListDiv.empty();
    for (var key in friends) {
       if (friends[key] != undefined){
        var item = $("#friend-item-template").clone();//好友信息的模版
        item.remove('#friend-item-template');
        var friend = friends[key];
        if (friend.unReadMessage == undefined) {
            var narry=new Array(new Object());
            friend.unReadMessage=0;
            friend.Message=narry;
        }
        //更新好友uuid
        item.attr("id", friend.uuid);
        //设置好友头像
        var friendAvatarImg = item.find("img");
        friendAvatarImg.attr("src", friend.avatar);
        friendAvatarImg.attr("id", "avatar_" + key);

        //设置好友名称
        item.find(".friend-name").html(friend.name);
        //更新未读消息数
        var unReadMessage = friend.unReadMessage;
        item.find(".message-count").text(unReadMessage);

        //显示或隐藏未读消息数
        var messageBadge = item.find(".friend-item-message-badge");
        if (unReadMessage == 0) {
            //没有未读消息，隐藏消息数量
            messageBadge.hide();
        } else {
            //有未读消息，展现未读消息量
            messageBadge.show();
        }

        //添加点击事件，点击好友条目后，进入聊天窗口
        (function (key) {
            var friend = friends[key]
            item.click(function () {
                switchToPrivateChat(friend);
            });
        })(key)

        //将一条好友添加到好友的列表中
        friendListDiv.append(item);

        // if (friend.online) {
        $("#avatar_" + key).removeClass("friend-avatar-desaturate");
        // }
    }
    }
}

//绘制群列表
function renderGroupList(groups) {
    var groupListDiv = $("#group-list");
    groupListDiv.empty();
    for (var key in groups) {
        var item = $("#group-item-template").clone();//群信息的模版
        item.remove('#group-item-template');
        var group = groups[key];
        if (group.unReadMessage == undefined) {
            var narry=new Array(new Object());
            group.unReadMessage=0
            group.Message=narry;
        }
        //设置群id
        item.attr("id", key);
        //设置群头像
        var groupAvatarImg = item.find("img");
        groupAvatarImg.attr("src", group.avatar);
        groupAvatarImg.attr("id", "avatar_" + key);
        //设置群名称
        item.find(".group-name").html(group.name+" (群组)");
        //更新未读消息数
        var unReadMessage = group.unReadMessage;
        item.find(".message-count").text(unReadMessage);

        //显示或隐藏未读消息数
        var messageBadge = item.find(".friend-item-message-badge");
        if (unReadMessage == 0) {
            //没有未读消息，隐藏消息数量
            messageBadge.hide();
        } else {
            //有未读消息，展现未读消息量
            messageBadge.show();
        }

        //添加点击事件，点击好友条目后，进入聊天窗口
        (function (key) {
            var group = groups[key]
            item.click(function () {
                switchToGroupChat(group);
            });
        })(key)

        //将一条好友添加到好友的列表中
        groupListDiv.append(item);
    }
}

//切换到私聊界面
function switchToPrivateChat(friend) {
    //设置当前页面为聊天页面和聊天的好友
    friend.unReadMessage=0
    currentPage.page = Pages.privateChat;
    currentPage.currentChatFriend = friend;
    $("title").text("私聊界面");

    //先清空之前的聊天内容
    var chatBoxContent = $("#chat-box-content");
    chatBoxContent.empty();

    //隐藏好友列表
    $("#friends-box").hide();
    //更新当前聊天的好友名称
    $(".current-friend-name").text(friend.name);

    //显示聊天窗口
    $("#chat-box").show();

    $("#sendMessageButton").off('click').on("click", function () {
        sendPrivateChatMessage(friend.uuid);
    });

    //绘制聊天消息
    var messages = friend.Message;
    if (messages.length != 0) {
        renderPrivateChatMessage(messages,true)
    }
}

//切换到群聊界面
function switchToGroupChat(group) {
    //设置当前页面为聊天页面和聊天的好友
    group.unReadMessage=0
    currentPage.page = Pages.groupChat;
    currentPage.currentChatGroup = group;
    $("title").text("群聊界面");

    //先清空之前的聊天内容
    var chatBoxContent = $("#group-chat-box-content");
    chatBoxContent.empty();

    //隐藏好友列表
    $("#friends-box").hide();
    //更新当前聊天的好友名称
    $(".current-friend-name").text(group.name);

    //显示聊天窗口
    $("#group-chat-box").show();

    $("#groupSendMessageButton").off('click').on("click", function () {
        sendGroupChatMessage(group.uuid);
    });


//    绘制界面聊天消息
    var messages = group.Message
    renderGroupChatMessage(messages,true)
}

//私聊回到联系人
function privateChatBackToContacts() {
//    chatService.resetFriendUnReadMessage(currentPage.currentChatFriend);
    switchToContacts()
}

// 发送私聊消息
function sendPrivateChatMessage() {
    //获取content并赋值
    var messageInputBox = $("#send-input-box");

    var content = messageInputBox.val();
    if (content != '' && content.trim().length > 0) {
        // 发送消息
//        chatService.sendPrivateMessage(currentPage.currentChatFriend.uuid, content);
        sendPrivateMessage(currentPage.currentChatFriend.uuid, content)
        //发送消息后输入框清空
        messageInputBox.val("");
    }
}

function sendPrivateMessage(uuid, text) {
    data = JSON.stringify({"id":uuid,"type":Msgtype.firtype,"mdata":text});
    ws.send(data)
}

//加载私聊历史消息
function loadPrivateHistory() {
    var messages = chatService.getPrivateMessages(currentPage.currentChatFriend.uuid);
    let earliestMessageTimeStamp = Date.now();
    let earliestMessage = messages[0];
    if (earliestMessage) {
        earliestMessageTimeStamp = earliestMessage.timestamp;
    }
    this.chatService.loadPrivateHistoryMessage(currentPage.currentChatFriend.uuid, earliestMessageTimeStamp)
}

//监听私聊消息加载
function onPrivateHistoryLoad(friendId, messages) {
    if (messages.length == 0) {
        $('#top').html('已经没有更多的历史消息');
        $('#top').css({color: 'gray', textDecoration: 'none'});
        return
    }
    var chatMessages = chatService.getPrivateMessages(friendId)
    renderPrivateChatMessage(chatMessages)
}

//绘制界面私聊消息
function renderPrivateChatMessage(privateMessages, scrollToBottom) {
    var chatBoxContent = $("#chat-box-content");
    chatBoxContent.empty();
    privateMessages.forEach(function (message) {
        var messageTemplate;
        //判断这条消息是谁发的
        if (message.senderId == undefined){
        } else {
        if (message.senderId === chatService.currentUser[0].uuid) {
            //自己发送的消息展示在右边
            messageTemplate = $("#chat-box-self-message-template").clone();
            messageTemplate.remove('#chat-box-self-message-template');
            //更新头像
            messageTemplate.find("img").attr("src", chatService.currentUser[0].avatar);
        } else {
            //如果该为好友发送的消息展示在左边
            messageTemplate = $("#chat-box-friend-message-template").clone();
            messageTemplate.find("img").attr("src", currentPage.currentChatFriend.avatar);
        }
        messageTemplate.find(".chat-message").text(message.payload);
        //显示一条消息到页面上
        chatBoxContent.append(messageTemplate);
    }
    });
    //将滚动条拉到最下
    scrollToBottom && $('#private-box').scrollTop($('#private-box')[0].scrollHeight);
}

//群聊回到联系人
function groupChatBackToContacts() {
//    chatService.resetGroupUnReadMessage(currentPage.currentChatGroup);
    switchToContacts()
}

//发送群聊消息
function sendGroupChatMessage() {
    //获取content并赋值
    var messageInputBox = $("#group-send-input-box");

    var content = messageInputBox.val();
    if (content != '' && content.trim().length > 0) {
        // 发送消息
//        chatService.sendGroupMessage(currentPage.currentChatGroup.uuid, content);
        sendGroupMessage(currentPage.currentChatGroup.uuid, content)
        //发送消息后输入框清空
        messageInputBox.val("");
    }
}

function sendGroupMessage(uuid, text) {
    data = JSON.stringify({"id":uuid,"type":Msgtype.grptype,"mdata":text});
    ws.send(data)
}

//加载群聊历史消息
function loadGroupHistory() {
    var messages = chatService.getGroupMessages(currentPage.currentChatGroup.uuid);
    let earliestMessageTimeStamp = Date.now();
    let earliestMessage = messages[0];
    if (earliestMessage) {
        earliestMessageTimeStamp = earliestMessage.timestamp;
    }
    this.chatService.loadGroupHistoryMessage(currentPage.currentChatGroup.uuid, earliestMessageTimeStamp)
}

//监听群聊历史消息加载
function onGroupHistoryLoad(groupId, messages) {
    if (messages.length == 0) {
        $('#group-top').html('已经没有更多的历史消息');
        $('#group-top').css({color: 'gray', textDecoration: 'none'});
        return
    }
    ;
    var chatMessage = chatService.getGroupMessages(groupId)
    renderGroupChatMessage(chatMessage)
}

//绘制群聊界面消息
function renderGroupChatMessage(groupMessages, scrollToBottom) {
    var currentUser = chatService.currentUser[0];
    var chatBoxContent = $("#group-chat-box-content");
    chatBoxContent.empty();
    groupMessages.forEach(function (message) {
        var messageTemplate;
        //判断这条消息是谁发的
        if (message.senderId == undefined) {
        } else {
        if (message.senderId === currentUser.uuid) {
            //自己发送的消息展示在右边
            messageTemplate = $("#chat-box-self-message-template").clone();
            messageTemplate.remove('#chat-box-self-message-template');
            //更新头像
            messageTemplate.find("img").attr("src", currentUser.avatar);
        } else {
            //如果该为好友发送的消息展示在左边
            messageTemplate = $("#chat-box-friend-message-template").clone();
            messageTemplate.remove('#chat-box-self-message-template');

            var friend = chatService.friends[message.senderId]
            messageTemplate.find("img").attr("src", friend.avatar);
        }
        messageTemplate.find(".chat-message").text(message.payload);
        //显示一条消息到页面上
        chatBoxContent.append(messageTemplate);
    }
    });
    groupMessages = [];
    //将滚动条拉到最下
    scrollToBottom && $('#group-box').scrollTop($('#group-box')[0].scrollHeight);
}

//监听接收私聊消息
function onNewPrivateMessageReceive(buf) {
    let friendId = buf.id
    friend = chatService.friends[friendId];
    // if (friend == undefined) {
    //     var fob = new Object();
    //     var narry=new Array(new Object());
    //     friend=fob;
    //     friend.Message=narry;
    // }
    gmsglen = friend.Message.length
    friend.Message[gmsglen] = buf

    //如果当前窗口是在好友列表页面，只显示未读消息数
    if (currentPage.page == Pages.contacts) {
        friend.unReadMessage = friend.unReadMessage+1
        renderFriendList(chatService.friends)
    } else {
        if (friendId == currentPage.currentChatFriend.uuid) {
            friend.unReadMessage = 0
            renderPrivateChatMessage(friend.Message, true)
        }
    }
}

//监听接收群聊消息
function onNewGroupMessageReceive(buf) {
    let groupId = buf.id
    let group = chatService.groups[groupId];
    let gmsglen = group.Message.length
    group.Message[gmsglen] = buf

    //如果当前窗口是在好友列表页面，只显示未读消息数
    if (currentPage.page == Pages.contacts) {
        group.unReadMessage = group.unReadMessage + 1
        var groupItem = $("#" + groupId);
        groupItem.find(".message-count").text(group.unReadMessage);
        groupItem.find(".friend-item-message-badge").show();
    } else {
        if (groupId == currentPage.currentChatGroup.uuid) {
            renderGroupChatMessage(group.Message, true)
        }
    }
}

//更新好友上线
function onFriendonline(buf) {
    var a = new Object();
    chatService.friends[buf.id] = a
    chatService.friends[buf.id].uuid=buf.id;
    chatService.friends[buf.id].name=buf.name;
    chatService.friends[buf.id].avatar=buf.avatar;
    renderFriendList(chatService.friends)
}

//更新好友下线
function onFriendoffline(buf) {
    if (chatService.friends[buf.id] != undefined){
        chatService.friends[buf.id] = undefined
    }
    renderFriendList(chatService.friends);
}

//更新群列表
function onGroupListChange(groups) {
    renderGroupList(groups)
}

//显示群成员
function showGroupMember() {
    $('#group-member-layer').show();
    var members = chatService.groups[currentPage.currentChatGroup.uuid];
    $('.group-member-amount').html("成员(" + Object.keys(members).length + ')')
    let str = "";
    for (var key in members) {
        str += '<img src="' + members[key].avatar + '"/>'
    }
    $('.layer-container').html(str)
}

//隐藏群成员
function hideGroupMember() {
    $('#group-member-layer').hide();
}


//ws
function createWebSocket(url) {
    try {
        if ('WebSocket' in window) {
            ws = new WebSocket(url + "?token="+ wstoken);
        } else if ('MozWebSocket' in window) {
            ws = new MozWebSocket(url + "?token="+ wstoken);
        } else {
            alert("您的浏览器不支持websocket协议,建议使用新版谷歌、火狐等浏览器，请勿使用IE10以下浏览器，360浏览器请使用极速模式，不要使用兼容模式！");
        }
        initEventHandle();
    } catch (e) {
        //x]alert("line98");
        reconnect(url);
        console.log(e);
    }
}

function initEventHandle() {
    ws.onclose = function () {
        console.log("llws连接关闭!" + new Date().toUTCString());
        reconnect(wsUrl);
    };
    ws.onerror = function () {
        console.log("llws连接错误!");
        reconnect(wsUrl);
    };
    ws.onopen = function () {
        // ws.send(wstoken)
        heartCheck.reset().start();
        console.log("llws连接成功!" + new Date().toUTCString());
    };
    ws.onmessage = function (event) {
        heartCheck.reset().start();
        var eventData = event.data;
        handMsg(eventData);
    };
}


window.onbeforeunload = function () {
    ws.close();
}

function reconnect(url) {
    if (lockReconnect) return;
    lockReconnect = true;
    setTimeout(function () {     //没连接上会一直重连，设置延迟避免请求过多
    createWebSocket(url);
    lockReconnect = false;
    }, 2000);
}

//心跳检测
var heartCheck = {
    //timeout: 540000,        //9分钟发一次心跳
    //timeout: 3600,        //1分钟发一次心跳
    timeout: 10800,        //3分钟发一次心跳
    timeoutObj: null,
    serverTimeoutObj: null,
    reset: function () {
        clearTimeout(this.timeoutObj);
        clearTimeout(this.serverTimeoutObj);
        return this;
    },
    start: function () {
        var self = this;
        this.timeoutObj = setTimeout(function () {
            ws.send("ping");
            self.serverTimeoutObj = setTimeout(function () {
                ws.close();     
            },self.timeout)
        },this.timeout)
    }
}

function handMsg(edata){
    if(edata != "pong"){
        var buf=JSON.parse(edata);
        if(buf.init == undefined){
            if(buf.type == Msgtype.firtype){
                onNewPrivateMessageReceive(buf)
            }
            else if (buf.type == Msgtype.grptype){
                onNewGroupMessageReceive(buf)
            }
            else if (buf.type == Msgtype.ontype){
                onFriendonline(buf)
            }
            else if (buf.type == Msgtype.offtype){
                onFriendoffline(buf)
            }
            else {
            }                
        }
        else{
            chatService = buf.init
            switchToContacts()
        }}
    }

