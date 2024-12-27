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
	var domainKey string

	if r.Header.Get("X-Domain-Name") != "" {
		domainKey = r.Header.Get("X-Domain-Name") + "_Backend"
		fmt.Printf("Domain : %s Backend", domainKey)
	} else if r.Host != "" {
		domainKey = r.Host + "_Frontend"
		fmt.Printf("Domain : %s Frontend", domainKey)
	} else {
		http.Error(rw, "Domain not provided", http.StatusBadRequest)
		return
	}

	// Fetch the optimal server
	optimalServer, err := redisDB.GetOptimalServer(domainKey)
	if err != nil {
		http.Error(rw, "Failed to fetch optimal server", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Optimal server: %s\n", optimalServer)

	// Ensure the URL includes the scheme
	if !urlHasScheme(optimalServer) {
		optimalServer = "http://" + optimalServer
	}

	// Parse the server URL
	serverURL, err := url.Parse(optimalServer)
	if err != nil {
		http.Error(rw, "Invalid server URL", http.StatusInternalServerError)
		return
	}

	// Create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(serverURL)
	proxy.ServeHTTP(rw, r)
	newVal, err := redisDB.UpdateActiveCount(domainKey, optimalServer, 1)
	if err != nil {
		http.Error(rw, "Failed to Update active count", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Host server : %s\nUpdated Value: %f", r.Host, newVal)
	defer redisDB.UpdateActiveCount(domainKey, optimalServer, -1)
}

// Helper function to check if a URL includes a scheme
func urlHasScheme(urlStr string) bool {
	return len(urlStr) > 0 && (urlStr[:7] == "http://" || urlStr[:8] == "https://")
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
