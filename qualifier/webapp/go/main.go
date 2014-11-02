package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
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
	r := mux.NewRouter()

	r.Handle("/images/{rest}", http.StripPrefix("/images/", http.FileServer(http.Dir("../public/images"))))
	r.Handle("/stylesheets/{rest}", http.StripPrefix("/stylesheets/", http.FileServer(http.Dir("../public/stylesheets"))))

	r.HandleFunc("/", IndexController)
	r.HandleFunc("/login", LoginController).Methods("POST")
	r.HandleFunc("/mypage", MyPageController)
	r.HandleFunc("/report", ReportController)

	r.NotFoundHandler = http.HandlerFunc(NotFound)

	http.ListenAndServe(":8080", LoggingServeMux(r))
}

func LoggingServeMux(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
}
