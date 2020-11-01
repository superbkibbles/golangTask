package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"github.com/speps/go-hashids"
	"log"
	"net/http"
	"time"
)

const PORT = ":8080"

var (
	DB *sql.DB
	dbErr error
	c *cache.Cache
)

type Urls struct {
	ID string `json:"id,omitempty"`
	LongUrl string `json:"longUrl,omitempty"`
	ShortUrl string `json:"shortUrl,omitempty"`
}

func init() {
	c = cache.New(5*time.Minute, 10*time.Minute)
}

func CreateEndpoint(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(202)
	//w.Write([]byte("Hello"))
	//return
	var url Urls

	_ = json.NewDecoder(r.Body).Decode(&url)
	var exist bool
	rows, err := DB.Query("select exists(select * from urls WHERE long_url = ?)", url.LongUrl)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&exist)
	}
	if !exist {
		hashId := hashids.NewData()
		h, _:= hashids.NewWithData(hashId)
		now := time.Now()
		url.ID, _ = h.Encode([]int{int(now.Unix())})
		url.ShortUrl = "http://localhost:8080/" + url.ID
		// Storing key values in the database
		DB.Query(`insert into urls(id, long_url, short_url) VALUES(?,?,?)`, url.ID, url.LongUrl, url.ShortUrl)
	} else {
		var rowUrl Urls
		rows, err = DB.Query("select * from urls WHERE long_url = ?", url.LongUrl)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte(err.Error()))
			return
		}
		defer rows.Close()

		for rows.Next() {
			rows.Scan(&rowUrl.ID, &rowUrl.LongUrl, &rowUrl.ShortUrl)
		}
		url = rowUrl
	}

	json.NewEncoder(w).Encode(url)
}

func ShowEndpoint(w http.ResponseWriter, r *http.Request) {
	var url Urls
	param := r.URL.Query().Get("shortUrl")

	rows, err := DB.Query(`select * from urls WHERE short_url = ?`, param)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}

	for rows.Next() {
		rows.Scan(&url.ID, &url.LongUrl, &url.ShortUrl)
	}
	json.NewEncoder(w).Encode(url)
}

func RootEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	rows, err := DB.Query(`select * from urls where id = ?`, id)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	var url Urls
	for rows.Next() {
		rows.Scan(&url.ID, &url.LongUrl, &url.ShortUrl)
	}
	_, err = http.Get(url.LongUrl)

	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}

	// save redirect in temp memory
	//go saveInCache(url)

	http.Redirect(w, r, url.LongUrl, http.StatusFound)
}

//func saveInCache(url Urls) {
//	var urls []interface{}
//	res, found := c.Get("visitedSites")
//	if !found {
//	//	 add mew one
//		urls = append(urls, url)
//		c.Set("visitedSites", &urls, cache.DefaultExpiration)
//
//		return
//	}
//
//	for _, links := range res {
//
//	}
//	urls = append(urls, res)
//	c.Set("visitedSites", &urls, cache.DefaultExpiration)
//
//	fmt.Println(urls)
//	// add new one to slice
//}


func main() {
	// DELETE MY PATH and password Later on
	DB , dbErr = ConnectDb("root:A3201888118a@/shortner")
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	defer DB.Close()

	r := mux.NewRouter()
	r.HandleFunc("/create", CreateEndpoint).Methods("PUT")
	r.HandleFunc("/show/", ShowEndpoint).Methods("GET")
	r.HandleFunc("/{id}", RootEndpoint).Methods("GET")

	if err := http.ListenAndServe(PORT, r); err != nil {
		log.Fatalf("Error while listening to port %s", PORT)
	}
}

func ConnectDb(path string) (*sql.DB, error) {
	db, err := sql.Open("mysql", path)

	if err != nil {
		return nil, err
	}

	return db, nil
}