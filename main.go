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

	fmt.Fprintf(os.Stdout, "Server starting on port: %v \n", *port)
	go func() {
		serve(*port, errCh)
	}()
	
	for {
		select {
		case <-sig:
			signal.Stop(sig)
			fmt.Println("Gracefully shutting down..")
			os.Exit(1)
		case err := <-errCh:
			log.Fatalf("Error starting server: %v", err)
		}
	}
}

func serve(port int64, errChan chan error) {
	wd, err := os.Getwd()

	if err != nil {
		errChan <- err
		return
	}
	fs := http.FileServer(http.Dir(wd))

	http.Handle("/", fs)
	err = http.ListenAndServe(":" + strconv.FormatInt(port, 10), nil)

	if err != nil {
		errChan <- err
		return
	}
}