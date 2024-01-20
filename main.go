package main

import (
	"leHttpServ/serv"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	srv := serv.NewSrv()

	mux.HandleFunc("/vote", srv.Vote)
	mux.HandleFunc("/stats", srv.Stats)

	http.ListenAndServe(":8080", mux)
}
