package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/zenazn/goji"
	"net/http"
	"strconv"
)

var db *sql.DB
var (
	UserLockThreshold int
	IPBanThreshold    int
)

func init() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		getEnv("ISU4_DB_USER", "root"),
		getEnv("ISU4_DB_PASSWORD", ""),
		getEnv("ISU4_DB_HOST", "localhost"),
		getEnv("ISU4_DB_PORT", "3306"),
		getEnv("ISU4_DB_NAME", "isu4_qualifier"),
	)

	var err error

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	UserLockThreshold, err = strconv.Atoi(getEnv("ISU4_USER_LOCK_THRESHOLD", "3"))
	if err != nil {
		panic(err)
	}

	IPBanThreshold, err = strconv.Atoi(getEnv("ISU4_IP_BAN_THRESHOLD", "10"))
	if err != nil {
		panic(err)
	}
}

var store = sessions.NewCookieStore([]byte("secret-isucon"))

func main() {
	flag.Set("bind", ":8080")

	goji.Get("/", IndexController)
	goji.Post("/login", LoginController)
	goji.Get("/mypage", MyPageController)
	goji.Get("/report", ReportController)
	goji.Get("/images/*", http.StripPrefix("/images/", http.FileServer(http.Dir("../public/images"))))
	goji.Get("/stylesheets/*", http.StripPrefix("/stylesheets/", http.FileServer(http.Dir("../public/stylesheets"))))
	goji.NotFound(NotFound)

	goji.Serve()

}

func NotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
}
