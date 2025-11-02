 package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Movie struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Year       int    `json:"year"`
	ActorCount int    `json:"actor_count"`
}

func main() {
	db := initDB()
	defer db.Close()

	router := gin.Default()
	router.GET("/movies", func(c *gin.Context) { getMovies(c, db) })
	router.Run(":8090")
}

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./movies.db")
	if err != nil {
		panic(err)
	}
	
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS movies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			year INTEGER NOT NULL
		)`)
	if err != nil {
		panic(err)
	}
	
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS actors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			movie_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			FOREIGN KEY (movie_id) REFERENCES movies(id)
		)`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("DELETE FROM actors")
	_, err = db.Exec("DELETE FROM movies")
	
	_, err = db.Exec(`INSERT INTO movies (title, year) VALUES 
		('Inception', 2010), 
		('The Dark Knight', 2008), 
		('Interstellar', 2014),
		('The Matrix', 1999), 
		('Avatar', 2009)`)
	if err != nil {
		panic(err)
	}
	
	_, err = db.Exec(`INSERT INTO actors (movie_id, name) VALUES 
		(1, 'Leonardo DiCaprio'), 
		(1, 'Joseph Gordon-Levitt'), 
		(1, 'Ellen Page'),
		(2, 'Christian Bale'), 
		(2, 'Heath Ledger'), 
		(2, 'Aaron Eckhart'),
		(3, 'Matthew McConaughey'), 
		(3, 'Anne Hathaway'),
		(4, 'Keanu Reeves'), 
		(4, 'Laurence Fishburne'),
		(5, 'Sam Worthington'), 
		(5, 'Zoe Saldana')`)
	if err != nil {
		panic(err)
	}
	
	return db
}

func getMovies(c *gin.Context, db *sql.DB) {
	start := time.Now()

	yearMin := c.Query("year_min")
	yearMax := c.Query("year_max")
	limit := c.Query("limit")
	offset := c.Query("offset")

	query := `
		SELECT m.id, m.title, m.year, COUNT(a.id) as actor_count
		FROM movies m
		LEFT JOIN actors a ON m.id = a.movie_id
		WHERE 1=1
	`
	var args []interface{}

	if yearMin != "" {
		query += " AND m.year >= ?"
		yearMinInt, _ := strconv.Atoi(yearMin)
		args = append(args, yearMinInt)
	}

	if yearMax != "" {
		query += " AND m.year <= ?"
		yearMaxInt, _ := strconv.Atoi(yearMax)
		args = append(args, yearMaxInt)
	}

	query += " GROUP BY m.id, m.title, m.year"
	query += " ORDER BY m.year DESC"

	if limit != "" {
		query += " LIMIT ?"
		limitInt, _ := strconv.Atoi(limit)
		args = append(args, limitInt)
	}

	if offset != "" {
		query += " OFFSET ?"
		offsetInt, _ := strconv.Atoi(offset)
		args = append(args, offsetInt)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		err := rows.Scan(&m.ID, &m.Title, &m.Year, &m.ActorCount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		movies = append(movies, m)
	}

	queryTime := time.Since(start)
	c.Header("X-Query-Time", queryTime.String())
	c.JSON(http.StatusOK, movies)
}