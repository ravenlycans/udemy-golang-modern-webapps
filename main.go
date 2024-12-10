package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

const portNumber = 8080

// Home is the http handler for the "/" route.
func Home(w http.ResponseWriter, r *http.Request) {
	n, err := fmt.Fprintf(w, "This is the Home Page")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Home: wrote %d bytes to %s\n", n, r.RemoteAddr)
}

// About is the http handler for the "/about" route.
func About(w http.ResponseWriter, r *http.Request) {
	var sum int
	var err error
	var errMessage []string
	var bw int
	var tbw int
	var x, y int

	// Let's tell the client that we're going to send partial html.
	w.Header().Set("Content-Type", "text/html")

	// Let's extract the querystring from the request.
	x, err = strconv.Atoi(r.URL.Query().Get("x"))
	if err != nil {
		errMessage = append(errMessage, "Cannot convert querystring x to int")
		x = 0
	}
	y, err = strconv.Atoi(r.URL.Query().Get("y"))
	if err != nil {
		errMessage = append(errMessage, "Cannot convert querystring y to int")
		y = 0
	}

	sum, err = addValues(x, y)
	bw, err = fmt.Fprintf(w, "<p>This is the About Page</p>")
	tbw += bw
	bw, err = fmt.Fprintf(w, "<p>Adding values of querystring variables x and y, remember they need to be numbers.</p>")
	tbw += bw
	bw, err = fmt.Fprintf(w, "<p>And the sum of %d+%d is %d</p>", x, y, sum)
	tbw += bw

	if len(errMessage) > 0 {
		bw, err = fmt.Fprintf(w, "<p>YARR, there be errors, following is a list of them:\n\n")
		tbw += bw
		bw, err = fmt.Fprintf(w, "<ol>")
		tbw += bw
		for _, msg := range errMessage {
			bw, err = fmt.Fprintf(w, "<li>%s</li>\n", msg)
			tbw += bw
		}
		bw, err = fmt.Fprintf(w, "</ol></p>")
		tbw += bw
	}

	// Check for any error from above code, and if there is one, log it and throw a http status code 500.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("About: wrote %d bytes to %s before encounting error: %s\n", tbw, r.RemoteAddr, err.Error())
		return
	}

	fmt.Printf("About: wrote %d bytes to %s\n", tbw, r.RemoteAddr)
}

// FavIcon serves the favicon.ico in the server root.
func FavIcon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	http.ServeFile(w, r, "favicon.ico")
	fmt.Printf("FavIcon: wrote %d bytes to %s\n", w.Header().Get("Content-Length"), r.RemoteAddr)
}

func Divide(w http.ResponseWriter, r *http.Request) {
	var result float32
	var err error
	var errMessage []string
	var bw int
	var x, y float64
	var res string

	res = "<p>This is the Divide Page</p>"
	res = res + "<p>Enter two numbers in the querystring (x,y) and we'll divide them.</p>"

	// Let's set the content-type to html.
	w.Header().Set("Content-Type", "text/html")

	x, err = strconv.ParseFloat(r.URL.Query().Get("x"), 32)
	if err != nil {
		errMessage = append(errMessage, "Cannot convert querystring x to float")
		x = 0
	}

	y, err = strconv.ParseFloat(r.URL.Query().Get("y"), 32)
	if err != nil {
		errMessage = append(errMessage, "Cannot convert querystring y to float")
		y = -1
	}

	if y != -1 && x != 0 {
		// Okay we have sane values, let's divide them.
		result, err = divideValues(float32(x), float32(y))
		if err != nil {
			res = res + fmt.Sprintf("<p>Error dividing %f by %f - Error Message: %s</p>", x, y, err.Error())
		} else {
			res = res + fmt.Sprintf("<p>The result of dividing %f by %f is %f</p>", x, y, result)
		}
	} else {
		res = res + fmt.Sprintf("<p>Please input some valid floating point numbers in the querystring as x and y (?x=1.2&y=2.3)</p>")
	}

	if len(errMessage) > 0 {
		res = res + "<p>YARR, there be errors, following is a list of them:"
		res = res + "<ol>"
		for _, msg := range errMessage {
			res = res + fmt.Sprintf("<li>%s</li>\n", msg)
		}
		res = res + "</ol></p>"
	}

	// Finally let's write out the res string.
	bw, err = fmt.Fprintf(w, res)

	fmt.Printf("Divide: wrote %d bytes to %s\n", bw, r.RemoteAddr)
}

// addValues adds to integers and returns the sum.
func addValues(x, y int) (int, error) {
	if reflect.TypeOf(x).Kind() != reflect.Int || reflect.TypeOf(y).Kind() != reflect.Int {
		return 0, fmt.Errorf("AddValues: x and y must be integers")
	}

	return x + y, nil
}

// divideValues divides to float's and returns the result, it checks if it is called with the right parameters.
func divideValues(x, y float32) (float32, error) {
	if reflect.TypeOf(x).Kind() != reflect.Float32 || reflect.TypeOf(y).Kind() != reflect.Float32 {
		return 0, fmt.Errorf("DivideValues: x and y must be floats")
	}
	if y <= 0 {
		return 0, fmt.Errorf("DivideValues: y must be greater than 0")
	}

	result := x / y
	return result, nil
}

// main is the application entrypoint
func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)
	http.HandleFunc("/divide", Divide)
	http.HandleFunc("/favicon.ico", FavIcon)

	fmt.Printf("Server is listening on port %d\n", portNumber)

	_ = http.ListenAndServe(fmt.Sprintf(":%d", portNumber), nil)
}
