package main

import (
	"github.com/vgheri/pi-go/Godeps/_workspace/src/github.com/gorilla/mux"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/pigo/{digits:[0-9]+}", PigoHandler).Methods("GET")
	http.Handle("/", r)
	log.Printf("Server started and listening on port %d.", 3141)
	log.Fatal(http.ListenAndServe(":3141", nil))
}

// PigoHandler handles requests
func PigoHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	vars := mux.Vars(r)
	param := vars["digits"]
	n, _ := strconv.Atoi(param)
	pi := calculatePi(n)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(pi.String()))
	duration := time.Since(start)
	log.Printf("\t%s\t%s",
		r.RequestURI,
		duration)
	return
}

func calculatePi(n int) *big.Int {
	digits := big.NewInt(int64(n))
	unity := big.NewInt(0)
	unity.Exp(big.NewInt(10), digits, nil)
	pi := big.NewInt(0)
	four := big.NewInt(4)
	pi.Mul(four, pi.Sub(pi.Mul(four, arccot(5, unity)), arccot(239, unity)))
	return pi
}

// This is the same as the first Machin-based Pi
// program, except that it uses the "big" package's
// infinite-sized integers to get as many digits as we
// want.  It still computes the formula:
// pi := 4 * (4 * arccot(5) - arccot(239))

// We start out by defining a high-precision arc cotangent
// function.  This one returns the response as an integer
// - normally it would be a floating point number.  Here,
// the integer is multiplied by the "unity" that we pass in.
// If unity is 10, for example, and the answer should be "0.5",
// then the answer will come out as 5

func arccot(x int64, unity *big.Int) *big.Int {
	bigx := big.NewInt(x)
	xsquared := big.NewInt(x * x)
	sum := big.NewInt(0)
	sum.Div(unity, bigx)
	xpower := big.NewInt(0)
	xpower.Set(sum)
	n := int64(3)
	zero := big.NewInt(0)
	sign := false

	term := big.NewInt(0)
	for {
		xpower.Div(xpower, xsquared)
		term.Div(xpower, big.NewInt(n))
		if term.Cmp(zero) == 0 {
			break
		}
		if sign {
			sum.Add(sum, term)
		} else {
			sum.Sub(sum, term)
		}
		sign = !sign
		n += 2
	}
	return sum
}
