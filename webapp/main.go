package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// CLI flag to enable the error outcome.
	enableError = flag.Bool("e", false, "Enable error outcome in /flip endpoint")
	enableDelay = flag.Bool("d", false, "Enable delay in /flip endpoint")
)

// coinFlipTotal is a counter vector that records each coin flip response.
// It uses two labels:
// - http_code: the HTTP status code returned ("200" for success, "500" for error)
// - result: the outcome of the flip ("head", "tails", or "error")
var coinFlipTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "coin_flip_total",
		Help: "Total number of coin flip responses, partitioned by HTTP status code and outcome",
	},
	[]string{"http_code", "result"},
)

func init() {
	// Register the coinFlipTotal metric with Prometheus.
	prometheus.MustRegister(coinFlipTotal)
}

// flipHandler implements the /flip endpoint.
func flipHandler(w http.ResponseWriter, r *http.Request) {
	var outcome string
	var delay time.Duration

	if *enableError {
		// Randomly choose among "head", "tails", or "error".
		switch rand.Intn(3) {
		case 0:
			outcome = "head"
		case 1:
			outcome = "tails"
		default:
			outcome = "error"
		}
	} else {
		// Choose between "head" and "tails" only.
		if rand.Intn(2) == 0 {
			outcome = "head"
		} else {
			outcome = "tails"
		}
	}

	if *enableDelay {
		delay = time.Duration(rand.Intn(5))
		time.Sleep(delay * time.Second)
	}

	w.Header().Set("Content-Type", "application/json")

	if outcome == "error" {
		coinFlipTotal.With(prometheus.Labels{"http_code": "500", "result": "error"}).Inc()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
	} else {
		coinFlipTotal.With(prometheus.Labels{"http_code": "200", "result": outcome}).Inc()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"result": outcome})
	}
}

func main() {
	flag.Parse()

	http.HandleFunc("/flip", flipHandler)
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
