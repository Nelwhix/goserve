package server

import (
	"context"
	"fmt"
	"github.com/Nelwhix/goserve/utils"
	"github.com/Nelwhix/goserve/watcher"
	"github.com/gorilla/handlers"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Serve(port int64, root string, errChan chan error) {
	eventCh := make(chan string)
	go watcher.StartWatcher(root, eventCh)

	router := http.NewServeMux()
	router.HandleFunc("GET /events", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "eventCh", eventCh)
		r = r.WithContext(ctx)

		streamEvents(w, r)
	})
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "root", root)
		r = r.WithContext(ctx)

		serveFile(w, r)
	})

	server := http.Server{
		Addr:        ":" + strconv.FormatInt(port, 10),
		Handler:     handlers.LoggingHandler(os.Stdout, router),
		ReadTimeout: 5 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		errChan <- err
		return
	}
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	root := r.Context().Value("root").(string)
	wd, err := os.Getwd()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "INTERNAL SERVER ERROR")
		os.Exit(1)
	}

	var fs http.Handler

	if root == "." {
		fs = http.FileServer(http.Dir(wd))
	} else {
		fs = http.FileServer(http.Dir(root))
	}

	fs.ServeHTTP(w, r)
}

func streamEvents(w http.ResponseWriter, r *http.Request) {
	eventCh := r.Context().Value("eventCh").(chan string)
	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/event-stream")
	w.Header().Add("Access-Control-Allow-Origin", "*")

	for raw := range eventCh {
		fmt.Println("File changed")
		event, err := utils.FormatSSE(raw)

		if err != nil {
			fmt.Println(err)
			break
		}

		_, err = fmt.Fprint(w, event)
		if err != nil {
			fmt.Println(err)
		}

		flusher.Flush()
		time.Sleep(5 * time.Second)
	}
}
