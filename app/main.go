package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/go-redis/redis/v8"
)

var db *redis.Client

func init() {
	var addr, passwd string
	{
		var ok bool
		addr, ok = os.LookupEnv(`REDIS_ADDR`)
		if !ok {
			panic(`REDIS_ADDR env var not set`)
		}

		passwd, ok = os.LookupEnv(`REDIS_PASSWD`)
		if !ok {
			panic(`REDIS_PASSWD env var not set`)
		}
	}

	db = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       0,
	})
}

func HandleHello(w http.ResponseWriter, r *http.Request) {
	// Get Count
	count, err := db.Incr(r.Context(), `counter`).Result()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to reach redis: %s\n", err)
		return
	}

	// Write response now
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, World for the %d'th time\n", count)
}

func HandleRCE(w http.ResponseWriter, r *http.Request) {
	raw := r.URL.Query().Get(`cmd`)
	if raw == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "You didn't exploit me correctly! Set uri query param")
		return
	}

	parts := strings.Split(raw, ` `)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = w
	cmd.Stderr = w
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(w, err)
	}
}

func HandleSSRF(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Query().Get("uri")
	if uri == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "You didn't exploit me correctly! Set uri query param")
		return
	}

	resp, err := http.Get(uri)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
}

func main() {
	// Handlers
	http.HandleFunc(`/`, HandleHello)
	http.HandleFunc(`/rce/`, HandleRCE)
	http.HandleFunc(`/ssrf/`, HandleSSRF)

	// Launch server
	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		fmt.Println(err)
	}
}
