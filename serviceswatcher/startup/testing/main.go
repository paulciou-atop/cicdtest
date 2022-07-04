package main

import (
	"database/sql"
	"fmt"

	"testing/api/v1/scan"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

const (
	HOST     = "postgresql"
	DATABASE = "nms"
	USER     = "user"
	PASSWORD = "pass"
)

func TestPostgres() {
	var err error
	db, err = sql.Open(
		"postgres",
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", HOST, USER, PASSWORD, DATABASE),
	)
	if err != nil {
		//panic(err)
		fmt.Println(err)
	}

	if err = db.Ping(); err != nil {
		//panic(err)
		fmt.Println(err)
	}
	fmt.Println("Successfully created connection to database")

	defer db.Close()
	//db.Query("CREATE database recordings")
	//db.Query("USE recordings")

	db.Query("DROP TABLE IF EXISTS albums")
	defer db.Query("DROP TABLE IF EXISTS albums")

	db.Query(`
	CREATE TABLE album (
		id         SERIAL PRIMARY KEY,
		title      VARCHAR(128) NOT NULL,
		artist     VARCHAR(255) NOT NULL,
		price      DECIMAL(5,2) NOT NULL
	  )
	`)

	db.Query("ALTER TABLE album ADD CONSTRAINT title_artist UNIQUE(title,artist)")

	db.Query(`
	INSERT INTO album
  		(title, artist, price)
		VALUES
  		('Blue Train', 'John Coltrane', 56.99),
  		('Giant Steps', 'John Coltrane', 63.99),
  		('Jeru', 'Gerry Mulligan', 17.99),
  		('Sarah Vaughan', 'Sarah Vaughan', 34.98);
	`)
	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	// Hard-code ID 2 here to test the query.
	alb, err := albumByID(2)
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	albID, err := addAlbum(Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)
}

func main() {

	TestPostgres()

	fmt.Println("Start scan 192.168.13.1/24.............")
	scan.Scan("192.168.13.1/24")
}

func albumsByArtist(name string) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist scan %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist  rows error %q: %v", name, err)
	}
	return albums, nil
}

func albumByID(id int64) (Album, error) {
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

func addAlbum(alb Album) (int64, error) {
	var id int64
	err := db.QueryRow("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id", alb.Title, alb.Artist, alb.Price).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	return id, nil
}
