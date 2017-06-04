function err(msg, f) {
    layer.alert(msg, {icon: 2, closeBtn: 0, title: false}, f);
}

function succ(msg, f) {
    layer.msg(msg, {
        icon: 1,
        time: 1000
    }, f);
}

function handle_json(json, f) {
    if (json.msg.length > 0) {
        err(json.msg);
    } else {
        succ('恭喜，操作成功：）', f);
    }
}

function login() {
    $.post("/auth/login", {
        "username": $("#username").val(),
        "password": $("#password").val()
    }, function(json){
        handle_json(json, function(){
            location.href=$("#callback").val();
        });
    });
}

function article_delete(id) {
	layer.confirm("确定要删除？", {
        icon: 3
    }, function(index) {
        layer.close(index);
        $.ajax({
            url: "/backend/article/" + id,
            type: "DELETE",
            success: function(json) {
                handle_json(json, function() {
                    location.reload();
                });
            }
        });
    });
}

function article_update(){
    var status = $("#status").val();
    var url = $.trim($("#url").val());
    $.post("/backend/articles", {
        "title": $("#title").val(),
        "url": url,
        "tags": $("#tags").val(),
        "content": $("#content").val(),
        "status": status,
        "id": $("#hid_article_id").val()
    }, function(json){
        handle_json(json, function(){
            if(status == "0"){
                location.href = "/backend/draft/" + url + ".html";
            }else{
                location.href = "/article/" + url + ".html";
            }
        });
    });
}

function about_update() {
    $.ajax({
        url: "/backend/about",
        type: "PUT",
        data: {
            "content": $("#content").val()
        },
        success: function(json) {
            handle_json(json, function() {
                location.href="/about";
            });
        }
    });
}

function upload() {
    $("#image").click();
}
