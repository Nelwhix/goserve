package main

import (
	"flag"
	"fmt"
	"github.com/Nelwhix/goserve/server"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	figure "github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

func main() {
	ascii := figure.NewFigure("GOSERVE", "basic", true)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s - Static file serving and directory listing \n", ascii.String())
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright "+strconv.Itoa(time.Now().Year())+"\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage Information:")
		flag.PrintDefaults()
	}

	root := flag.String("root", ".", "Root Directory to serve")
	port := flag.Int64("p", 3000, "Port to start server on")
	flag.Parse()

	sig := make(chan os.Signal, 1)
	errCh := make(chan error)
	signal.Notify(sig, syscall.SIGINT)

	cyan := color.New(color.BgCyan).SprintFunc()
	fmt.Fprintf(os.Stdout, "\n\n"+cyan("INFO")+" Accepting connections at http://localhost:%v \n\n", *port)
	go func() {
		server.Serve(*port, *root, errCh)
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
