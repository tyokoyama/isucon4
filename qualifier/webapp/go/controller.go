package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	// "github.com/gorilla/sessions"

)

var index = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <link rel="stylesheet" href="/stylesheets/bootstrap.min.css">
    <link rel="stylesheet" href="/stylesheets/bootflat.min.css">
    <link rel="stylesheet" href="/stylesheets/isucon-bank.css">
    <title>isucon4</title>
  </head>
  <body>
    <div class="container">
      <h1 id="topbar">
        <a href="/"><img src="/images/isucon-bank.png" alt="いすこん銀行 オンラインバンキングサービス"></a>
      </h1>
<div id="be-careful-phising" class="panel panel-danger">
  <div class="panel-heading">
    <span class="hikaru-mozi">偽画面にご注意ください！</span>
  </div>
  <div class="panel-body">
    <p>偽のログイン画面を表示しお客様の情報を盗み取ろうとする犯罪が多発しています。</p>
    <p>ログイン直後にダウンロード中や、見知らぬウィンドウが開いた場合、<br>すでにウィルスに感染している場合がございます。即座に取引を中止してください。</p>
    <p>また、残高照会のみなど、必要のない場面で乱数表の入力を求められても、<br>絶対に入力しないでください。</p>
  </div>
</div>

<div class="page-header">
  <h1>ログイン</h1>
</div>

{{ if .Flash }}
  <div id="notice-message" class="alert alert-danger" role="alert">{{ .Flash }}</div>
{{ end }}

<div class="container">
  <form class="form-horizontal" role="form" action="/login" method="POST">
    <div class="form-group">
      <label for="input-username" class="col-sm-3 control-label">お客様ご契約ID</label>
      <div class="col-sm-9">
        <input id="input-username" type="text" class="form-control" placeholder="半角英数字" name="login">
      </div>
    </div>
    <div class="form-group">
      <label for="input-password" class="col-sm-3 control-label">パスワード</label>
      <div class="col-sm-9">
        <input type="password" class="form-control" id="input-password" name="password" placeholder="半角英数字・記号（２文字以上）">
      </div>
    </div>
    <div class="form-group">
      <div class="col-sm-offset-3 col-sm-9">
        <button type="submit" class="btn btn-primary btn-lg btn-block">ログイン</button>
      </div>
    </div>
  </form>
</div>
    </div>

  </body>
</html>
	`

var mypage = `
<html>
  <head>
    <meta charset="UTF-8">
    <link rel="stylesheet" href="/stylesheets/bootstrap.min.css">
    <link rel="stylesheet" href="/stylesheets/bootflat.min.css">
    <link rel="stylesheet" href="/stylesheets/isucon-bank.css">
    <title>isucon4</title>
  </head>
  <body>
    <div class="container">
      <h1 id="topbar">
        <a href="/"><img src="/images/isucon-bank.png" alt="いすこん銀行 オンラインバンキングサービス"></a>
      </h1>
<div class="alert alert-success" role="alert">
  ログインに成功しました。<br>
  未読のお知らせが０件、残っています。
</div>

<dl class="dl-horizontal">
  <dt>前回ログイン</dt>
  <dd id="last-logined-at">{{ .LastLogin.CreatedAt.Format "2006-01-02 15:04:05" }}</dd>
  <dt>最終ログインIPアドレス</dt>
  <dd id="last-logined-ip">{{ .LastLogin.IP }}</dd>
</dl>

<div class="panel panel-default">
  <div class="panel-heading">
    お客様ご契約ID：{{ .LastLogin.Login }} 様の代表口座
  </div>
  <div class="panel-body">
    <div class="row">
      <div class="col-sm-4">
        普通預金<br>
        <small>東京支店　1111111111</small><br>
      </div>
      <div class="col-sm-4">
        <p id="zandaka" class="text-right">
          ―――円
        </p>
      </div>

      <div class="col-sm-4">
        <p>
          <a class="btn btn-success btn-block">入出金明細を表示</a>
          <a class="btn btn-default btn-block">振込・振替はこちらから</a>
        </p>
      </div>

      <div class="col-sm-12">
        <a class="btn btn-link btn-block">定期預金・住宅ローンのお申込みはこちら</a>
      </div>
    </div>
  </div>
</div>
    </div>

  </body>
</html>
`

var tIndex = template.Must(template.New("index").Parse(index))
var tMyPage = template.Must(template.New("mypage").Parse(mypage))

func IndexController(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "isucon_go_session")

	if err := tIndex.Execute(w, session.Values["notice"]); err != nil {
		http.Error(w, "500 page Error", http.StatusInternalServerError)
	}
}

func LoginController(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "isucon_go_session")
	user, err := attemptLogin(r)

	notice := ""
	if err != nil || user == nil {
		switch err {
		case ErrBannedIP:
			notice = "You're banned."
		case ErrLockedUser:
			notice = "This account is locked."
		default:
			notice = "Wrong username or password"
		}

		fmt.Println(err)
		session.Values["notice"] = notice
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	session.Values["user_id"] = strconv.Itoa(user.ID)
	session.Save(r, w)
	http.Redirect(w, r, "/mypage", http.StatusFound)
}

func MyPageController(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "isucon_go_session")

	currentUser := getCurrentUser(session.Values["user_id"])

	if currentUser == nil {
		session.Values["notice"] = "You must be logged in"
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	currentUser.getLastLogin()
	if err := tMyPage.Execute(w, currentUser); err != nil {
		http.Error(w, "500 page Error", http.StatusInternalServerError)
	}
}

func ReportController(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(map[string][]string {
		"banned_ips":   bannedIPs(),
		"locked_users": lockedUsers(),
		})

	if err != nil {
		http.Error(w, "500 page Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", string(b))
}
