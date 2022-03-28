package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-aybukeertekin/domain/library"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-aybukeertekin/domain/model"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-aybukeertekin/domain/repository"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-aybukeertekin/domain/utils"
	"github.com/Picus-Security-Golang-Backend-Bootcamp/homework-4-aybukeertekin/infrastructure"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var libraryOp *library.Library

func CORSOptions() {
	handlers.AllowedOrigins([]string{"https://www.library.com"})
	handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	handlers.AllowedMethods([]string{"POST", "GET", "PUT", "PATCH"})
}

func init() {
	//Create books
	infrastructure.Read("./books.csv")
	host := "localhost"
	port := 5432
	user := "postgres"
	password := "admin"
	dbname := "library_db"
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db := infrastructure.NewPostgreDB(psqlInfo)
	bookRepository := repository.NewBookRepository(db)
	bookRepository.Migration()
	authorRepository := repository.NewAuthorRepository(db)
	authorRepository.Migration()
	libraryOp = library.NewLibrary(bookRepository, authorRepository)
}

func main() {
	router := mux.NewRouter()
	CORSOptions()
	router.Use(loggingMiddleware)
	router.Use(authenticationMiddleware)

	createReqHandler := authenticationMiddleware(http.HandlerFunc(CreateRequestHandler))
	subrouter := router.PathPrefix("/books").Subrouter()
	subrouter.HandleFunc("/{bookId:[0-9]+}/buy", BuyRequestHandler).Methods("PUT")
	subrouter.HandleFunc("/list", ListRequestHandler).Methods("GET")
	subrouter.Handle("/", createReqHandler).Methods("POST")
	subrouter.HandleFunc("/{bookId}:[0-9]+}", BuyRequestHandler).Methods("DELETE")

	srv := &http.Server{
		Addr:         "localhost:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	ShutdownServer(srv, time.Second*10)
}

func DeleteRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookId := vars["bookId"]
	//response := ""
	if bookId == "" {
		//response = "bookId cannot be null"
	} else {
		id, err1 := strconv.Atoi(bookId)
		if id <= 0 || err1 != nil {
			//response = "id should be an integer larger than 0"
		} else {
			err := libraryOp.DeleteBook(id)
			if err == nil {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

func ListRequestHandler(w http.ResponseWriter, r *http.Request) {
	books := libraryOp.ListBooks()
	w.Header().Set("Content-type", "application/json")
	resp, _ := json.Marshal(books)
	w.Write(resp)
}

func BuyRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	count := r.URL.Query().Get("count")
	bookId := vars["bookId"]
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	response := ""
	number, err1 := strconv.Atoi(count)
	id, err2 := strconv.Atoi(bookId)
	if err1 != nil {
		response = "count should be an integer larger than 0"
	} else if err2 != nil {
		response = "bookId should be an integer larger than 0"
	} else {
		err, book := libraryOp.BuyBooks(id, number)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			byteResponse, _ := json.Marshal(book)
			w.Write(byteResponse)
		}
	}
	if response != "" {
		resp := ApiResponse{
			Data: response,
		}
		byteResponse, _ := json.Marshal(resp)
		w.Write(byteResponse)
	}
}

func CreateRequestHandler(w http.ResponseWriter, r *http.Request) {
	var book model.Book
	err := utils.DecodeJSONBody(w, r, &book)
	if err != nil {
		var mr *utils.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = libraryOp.CreateBook(&book)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else {
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		byteResponse, _ := json.Marshal(book)
		w.Write(byteResponse)
	}
}

type ApiResponse struct {
	Data interface{} `json:"data"`
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.URL.Query())
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if strings.HasPrefix(r.URL.Path, "/books") {
			if token != "" {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Token not found", http.StatusUnauthorized)
			}
		} else {
			next.ServeHTTP(w, r)
		}

	})
}

func ShutdownServer(srv *http.Server, timeout time.Duration) {
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	//Block until we receive our signal.
	<-c
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
