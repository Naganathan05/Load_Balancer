package main

import (
	redisDB "Load_Balancer_Server/services"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type LoadBalancer struct {
	Port string
}

func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, r *http.Request) {
	optimalServer, err := redisDB.GetOptimalServer("www.google.com")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Optimal server : %s\n", optimalServer)
	serverURL, err := url.Parse(optimalServer)
	proxy := httputil.NewSingleHostReverseProxy(serverURL)
	proxy.ServeHTTP(rw, r)
}

var lb *LoadBalancer

func InitLoadBalancer() {
	lb = &LoadBalancer{Port: "5000"}
}

func redirectRequest(rw http.ResponseWriter, r *http.Request) {
	lb.serveProxy(rw, r)
}

func main() {
	fmt.Println("Starting Load Balancer...")

	redisDB.InitRedisClient()

	InitLoadBalancer()

	http.HandleFunc("/", redirectRequest)

	err := http.ListenAndServe(":"+lb.Port, nil)
	if err != nil {
		panic(err)
	}
}
