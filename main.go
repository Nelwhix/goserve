package main

import (
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
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright " + strconv.Itoa(time.Now().Local().Year()) + "\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage Information:")
		flag.PrintDefaults()
	}

	port := flag.Int64("port", 3000, "Port to start server on")
	flag.Parse()

	sig := make(chan os.Signal, 1)
	errCh := make(chan error)
	signal.Notify(sig, syscall.SIGINT)

	go func() {
		serve(*port, errCh)
	}()
	
	for {
		select {
		case _ = <-sig:
			signal.Stop(sig)
			fmt.Println("Gracefully shutting down..")
		case err := <-errCh:
			fmt.Errorf("Error: %w", err)
			os.Exit(1)
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
	fmt.Fprintf(os.Stdout, "Server starting on port: %v", port)
}