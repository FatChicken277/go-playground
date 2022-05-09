package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	port = os.Getenv("PORT")

	notFoundError = "object not found"
)

// Movie represents a movie model
type Movie struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

// Director represents a director model
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(movies)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var movie Movie

	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(100000000))
	movies = append(movies, movie)

	json.NewEncoder(w).Encode(movie)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, movie := range movies {
		if movie.ID == params["id"] {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}

	http.Error(w, notFoundError, http.StatusNotFound)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for idx, movie := range movies {
		if movie.ID == params["id"] {
			movies = append(movies[:idx], movies[idx+1:]...)
			fmt.Fprintf(w, "successful delete")
			return
		}
	}

	http.Error(w, notFoundError, http.StatusNotFound)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var newMovie Movie

	json.NewDecoder(r.Body).Decode(&newMovie)

	for _, movie := range movies {
		if movie.ID == params["id"] {
			movie.Title = newMovie.Title
			movie.Director = newMovie.Director

			json.NewEncoder(w).Encode(movie)
			return
		}
	}

	http.Error(w, notFoundError, http.StatusNotFound)
}

func main() {
	r := mux.NewRouter()

	movies = append(movies,
		Movie{
			ID:    "1",
			Title: "Movie1",
			Director: &Director{
				Firstname: "John",
				Lastname:  "Doe",
			},
		},
		Movie{
			ID:    "2",
			Title: "Movie2",
			Director: &Director{
				Firstname: "John",
				Lastname:  "Dee",
			},
		},
	)

	r.HandleFunc("/movies", getMovies).Methods(http.MethodGet)
	r.HandleFunc("/movies", createMovie).Methods(http.MethodPost)

	r.HandleFunc("/movies/{id}", getMovie).Methods(http.MethodGet)
	r.HandleFunc("/movies/{id}", updateMovie).Methods(http.MethodPut)
	r.HandleFunc("/movies/{id}", deleteMovie).Methods(http.MethodDelete)

	fmt.Println("Serving at :" + port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
