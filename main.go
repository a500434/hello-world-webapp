package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	version = "v0.0.3"
)

var (
	port = "8000"
)

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
}

func main() {
	log.Println("✨ Starting helloworld application from 2026-06-10 11:08")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		log.Printf("%s is showing /\n", ip)
		message := ""
		user := r.Header.Get("X-Auth-Request-User")
		if user == "" {
			message = "No user header found"
		} else {
			message = fmt.Sprintf("Hello %s", user)
		}

		internalIP, err := getIP()
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}

		message = message + "\nServer internal IP is: " + internalIP.String()

		w.Write([]byte(message))
		return
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		log.Printf("%s is showing health\n", ip)
		w.WriteHeader(200)
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		log.Printf("%s is showing version\n", ip)
		fmt.Fprintf(w, version)
	})

	s := http.Server{Addr: ":" + port}
	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")

	s.Shutdown(context.Background())
}

func getIP() (net.IP, error) {
	var ip net.IP
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Something went wrong getting IP")
		return ip, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println("Something went wrong getting IP")
			return ip, err
		}
		for _, addr := range addrs {

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
		}
	}

	return ip, nil
}

func getClientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}
//KICK IT
