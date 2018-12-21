package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Stat struct {
	File string `json:"file"`
	Size int64  `json:"size"`
}

func main() {

	p := flag.Int("p", 12345, "specify port")

	flag.Parse()

	initWebServer(*p)
}

func initWebServer(port int) {

	fmt.Printf("starting local server and listening on port %d...\n\n", port)

	defineRoutes()

	if port == 443 {
		if err := http.ListenAndServeTLS(":12345", "./server.crt", "./server.key", nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	} else {
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}
}

func defineRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		msg := `server listening on port 12345 at localhost.
call http://localhost:12345/files [POST] with body {"name": "foo", "path": "~/folder/"}
call http://localhost:12345/stats [GET]
`
		io.WriteString(w, msg)
	})

	http.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		if s, err := generateRequest(body); err != nil {
			fmt.Printf("whoops... %s", err.Error())
		} else {
			fmt.Printf("body recieved %s\n", string(body))
			GetStatisticsStore().Add(s.File, s.Size)
		}

		io.WriteString(w, "ok\n")
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		io.WriteString(w, GetStatisticsStore().Print())
	})
}

func generateRequest(reqBytes []byte) (*Stat, error) {

	s := Stat{}

	err := json.Unmarshal(reqBytes, &s)

	if err != nil {
		return nil, err
	}
	return &s, nil
}
