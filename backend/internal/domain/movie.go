package domain

import "time"

type Movie struct {
	ID          string `json:"id" bson:"_id"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Duration    int    `json:"duration" bson:"duration"`
	PosterURL   string `json:"poster_url" bson:"poster_url"`
	Genre       string `json:"genre" bson:"genre"`
}

type Showtime struct {
	ID        string    `json:"id" bson:"_id"`
	MovieID   string    `json:"movie_id" bson:"movie_id"`
	Screen    string    `json:"screen" bson:"screen"`
	StartTime time.Time `json:"start_time" bson:"start_time"`
	EndTime   time.Time `json:"end_time" bson:"end_time"`
	Price     float64   `json:"price" bson:"price"`
}
