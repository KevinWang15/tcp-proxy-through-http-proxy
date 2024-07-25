package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	localPort := flag.String("local-port", "12345", "Local port to listen on")
	proxyURL := flag.String("proxy-url", "https://user:pass@proxy.com", "HTTPS proxy URL")
	targetURL := flag.String("target-url", "https://a.com", "Target HTTPS URL")
	flag.Parse()

	parsedProxyURL, err := url.Parse(*proxyURL)
	if err != nil {
		log.Fatalf("Failed to parse proxy URL: %v", err)
	}

	parsedTargetURL, err := url.Parse(*targetURL)
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}

	// Create a transport that uses the proxy
	transport := &http.Transport{
		Proxy: http.ProxyURL(parsedProxyURL),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Note: This skips certificate verification. Use with caution in production.
		},
	}

	// Create an HTTP client that uses the custom transport
	client := &http.Client{Transport: transport}

	// Create a handler function for incoming requests
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Construct the URL for the target server
		targetReqURL := *parsedTargetURL
		targetReqURL.Path = r.URL.Path
		targetReqURL.RawQuery = r.URL.RawQuery

		// Create a new request to send to the target server
		targetReq, err := http.NewRequest(r.Method, targetReqURL.String(), r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Copy headers from the original request
		for name, values := range r.Header {
			for _, value := range values {
				targetReq.Header.Add(name, value)
			}
		}

		// Send the request to the target server through the proxy
		resp, err := client.Do(targetReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Copy headers from the response
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		// Set the status code
		w.WriteHeader(resp.StatusCode)

		// Stream the response body
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Printf("Error streaming response: %v", err)
		}
	}

	// Start the HTTP server
	addr := fmt.Sprintf(":%s", *localPort)
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, http.HandlerFunc(handler)))
}
