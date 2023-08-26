package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

const (
	listenAddr      = "127.0.0.1:8080"
	shutdownTimeout = 5 * time.Second
)

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var arseniy = &person{
	Name: "Arseniy",
	Age:  18,
}

func personHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		response, err := json.Marshal(arseniy)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		_, err = writer.Write(response)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	case "POST":
		d := json.NewDecoder(request.Body)
		newPerson := &person{}
		err := d.Decode(newPerson)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		arseniy = newPerson
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		_, err := fmt.Fprintf(writer, "I can't do that.")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := runServer(ctx); err != nil {
		log.Fatal(err)
	}
}

func runServer(ctx context.Context) error {
	server := http.Server{
		Addr: listenAddr,
	}

	http.HandleFunc("/person/", personHandler)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Listen and serve: %v", err)
		}
	}()

	log.Printf("Listening on %s", listenAddr)
	<-ctx.Done()

	log.Println("Shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	longShutdown := make(chan bool)

	go func() {
		time.Sleep(3 * time.Second)
		longShutdown <- true
	}()

	select {
	case <-shutdownCtx.Done():
		return fmt.Errorf("server shutdown: %w", ctx.Err())
	case <-longShutdown:
		log.Println("Finished")
	}

	return nil
}
