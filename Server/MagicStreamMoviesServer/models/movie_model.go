package model

import "go.mongodb.org/mongo-driver/v2/bson"

// Жаны
type Genre struct {
	GenreID   int    `bson:"genre_id" json:"genre_id" validate:"required"`
	GenreName string `bson:"genre_name" json:"genre_name" validate:"required,min=2,max=500"`
}

// Классификация
type Ranking struct {
	RankingValue int    `bson:"ranking_value" json:"ranking_value" validate:"required"`
	RankingName  string `bson:"ranking_name" json:"ranking_name" validate:"oneof=Excellent Good Okay bad Terrible"`
}

// Фильмы
type Movie struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ImdbID      string        `bson:"imdb_id" json:"imdb_id" validate:"required"`
	Title       string        `bson:"title" json:"title" validate:"required,min=2,max=500"`
	
	// ИСПРАВЛЕНО: valdate -> validate. 
	PosterPath  string        `bson:"poster_path" json:"poster_path" validate:"required,url"` 
	
	// ИСПРАВЛЕНО: valdate -> validate
	YouTubeId   string        `bson:"youtube_id" json:"youtube_id" validate:"required"`

	Genre       []Genre       `bson:"genre" json:"genre" validate:"required,dive"`
	AdminReview string        `bson:"admin_review" json:"admin_review"`
	Ranking     Ranking       `bson:"ranking" json:"ranking" validate:"required"`
}