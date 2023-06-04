package main

import (
	"log"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/fatih/color"
)

func main() {
	flag.Usage = func () {
		fmt.Fprintf(flag.CommandLine.Output(), 
		"Serve - Static file serving and directory listing \n")
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright " + strconv.Itoa(time.Now().Year()) + "\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage Information:")
		flag.PrintDefaults()
	}

	port := flag.Int64("p", 3000, "Port to start server on")
	flag.Parse()

	sig := make(chan os.Signal, 1)
	errCh := make(chan error)
	signal.Notify(sig, syscall.SIGINT)

	cyan := color.New(color.BgCyan).SprintFunc()
	fmt.Fprintf(os.Stdout, "\n\n" + cyan("INFO") + " Accepting connections at http://localhost:%v \n\n", *port)
	go func() {
		serve(*port, errCh)
	}()
	
	for {
		select {
		case <-sig:
			signal.Stop(sig)
			fmt.Fprintf(os.Stdout, "\n\n%s Gracefully shutting down. Please wait...", cyan("INFO"))
			os.Exit(1)
		case err := <-errCh:
			log.Fatalf("Error starting server: %v", err)
		}
	}
}

func serve(port int64, errChan chan error) {
	http.HandleFunc("/", serveFile)
	err := http.ListenAndServe(":" + strconv.FormatInt(port, 10), nil)
	if err != nil {
		errChan <- err
		return
	}
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	wd, err := os.Getwd()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "INTERNAL SERVER ERROR")
		os.Exit(1)
	}

	fs := http.FileServer(http.Dir(wd))

}

func logRequest(r *http.Request) {
	now := time.Now()
	date := fmt.Sprintf("%s/%s/%s", strconv.Itoa(now.Day()), now.Month(), strconv.Itoa(now.Year()))
	time := fmt.Sprintf("%s:%s:%s", strconv.Itoa(now.Hour()), strconv.Itoa(now.Minute()), strconv.Itoa(now.Second()))
	fmt.Fprintf(os.Stdout, "%s\t%s\t%s\t%s\t%s\t%s\n", r.Proto, date, time, r.Host, r.Method, r.RequestURI)
}