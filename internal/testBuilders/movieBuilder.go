package testBuilders

import (
	"intern/models"
	"time"
)

type MovieBuilder struct {
	movie models.Movie
}

func NewMovieBuilder() *MovieBuilder {
	return &MovieBuilder{}
}

func (b *MovieBuilder) WithID(id int) *MovieBuilder {
	b.movie.ID = id
	return b
}

func (b *MovieBuilder) WithTitle(title string) *MovieBuilder {
	b.movie.Title = title
	return b
}

func (b *MovieBuilder) WithDesc(desc string) *MovieBuilder {
	b.movie.Description = desc
	return b
}

func (b *MovieBuilder) WithRelease(release time.Time) *MovieBuilder {
	b.movie.ReleaseDate = release
	return b
}

func (b *MovieBuilder) WithRating(rating int) *MovieBuilder {
	b.movie.Rating = rating
	return b
}

func (b *MovieBuilder) Build() models.Movie {
	return b.movie
}
