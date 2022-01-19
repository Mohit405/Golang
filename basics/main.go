package main

import (
	"html/template"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var templates *template.Template
var client *redis.Client
var store = sessions.NewCookieStore([]byte("t0p-se3ret"))

func main() {
	r := mux.NewRouter()
	templates = template.Must(template.ParseGlob("Page/*.html"))
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	r.HandleFunc("/", Authorization(GetComments)).Methods("get")
	r.HandleFunc("/", Authorization(PostComments)).Methods("post")
	r.HandleFunc("/login", LoginGetHandler).Methods("get")
	r.HandleFunc("/login", LoginPostHandler).Methods("post")
	r.HandleFunc("/register", RegisterGetHandler).Methods("get")
	r.HandleFunc("/register", RegisterPostHandler).Methods("post")
	r.HandleFunc("/test", TestLogin).Methods("Get")

	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.Handle("/", r)
	http.ListenAndServe(":8080", r)
}
//MIddleware part
func Authorization(handler http.HandlerFunc)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		session,_:=store.Get(r,"session")
		_,ok:=session.Values["username"]
		if !ok{
			http.Redirect(w,r,"/login",302)
			return
		}
		handler.ServeHTTP(w,r)
	}
}
func GetComments(w http.ResponseWriter, r *http.Request) {
	
	comments, err := client.LRange("comments", 0, 10).Result()
	if err != nil {
		return
	}
	templates.ExecuteTemplate(w, "index.html", comments)
}
func PostComments(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	comment := r.PostForm.Get("comment")
	client.LPush("comments", comment)
	http.Redirect(w, r, "/", 302)
}

func LoginGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}
func LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password:=r.PostForm.Get("password")
	hash,err:=client.Get("user:"+username).Bytes()
	if err!=nil{
		return
	}
	err=bcrypt.CompareHashAndPassword(hash,[]byte(password))
	if err!=nil{
		return
	}
	session, _ := store.Get(r, "session")
	session.Values["username"] = username
	session.Save(r, w)
	http.Redirect(w,r,"/",302)

}

//lets test the login
func TestLogin(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	value, ok := session.Values["username"]
	if !ok {
		return
	}
	username, ok := value.(string)
	if !ok {
		return
	}
	w.Write([]byte(username))
}

func RegisterGetHandler(w http.ResponseWriter,r* http.Request){
	templates.ExecuteTemplate(w,"register.html",nil)
}

func RegisterPostHandler(w http.ResponseWriter,r* http.Request){
	r.ParseForm()
	username:=r.PostForm.Get("username")
	password:=r.PostForm.Get("password")
	cost:=bcrypt.DefaultCost
	hash,err:=bcrypt.GenerateFromPassword([]byte(password),cost)
	if err!=nil{
		return
	}
	client.Set("user:"+username,hash,0)
	http.Redirect(w,r,"/login",302)
}