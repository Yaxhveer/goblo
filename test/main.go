package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func serveBackend(name string, port string) {
	mu := http.NewServeMux()
	mu.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d := 3000+rand.Intn(2000)
		time.Sleep(time.Duration(d)*time.Millisecond)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Backend server name: %v\n", name)
		fmt.Fprintf(w, "Backend service time: %vms\n", d)
		fmt.Fprintf(w, "Response header: %v\n", r.Header)
	}))
	mu.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Backend server name: %v\n", name)
		fmt.Fprintf(w, "Server working Fine.\n")
	}))
	http.ListenAndServe(port, mu)
}

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(4)

	go func() {
		serveBackend("web1", ":3331")
		wg.Done()
	}()

	go func() {
		serveBackend("web2", ":3332")
		wg.Done()
	}()

	go func() {
		serveBackend("web3", ":3333")
		wg.Done()
	}()

	go func() {
		serveBackend("web4", ":3334")
		wg.Done()
	}()

	go func() {
		serveBackend("web5", ":3335")
		wg.Done()
	}()

	go func() {
		serveBackend("web6", ":3336")
		wg.Done()
	}()

	go func() {
		serveBackend("web7", ":3337")
		wg.Done()
	}()

	wg.Wait()
}
