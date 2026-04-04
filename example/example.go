package example

import (
	"math/rand"
	"time"
)

// This is an unrelated type that should be ignored by the enumify generator.
type Example struct {
	Name   string
	Day    Day
	Status Status
	Date   time.Time
	Tags   []string
}

// This is an unrelated function that should be ignored by the enumify generator.
func New() (*Example, error) {
	adjectives := exampleNames[0]
	nouns := exampleNames[1]
	name := adjectives[rand.Intn(len(adjectives))] + " " + nouns[rand.Intn(len(nouns))]

	day := Day(rand.Intn(int(Sunday) + 1))
	status := Status(rand.Intn(int(StatusCancelled) + 1))

	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	n := rand.Intn(6) // 0–5 tags
	tags := make([]string, n)
	for i := range tags {
		tags[i] = exampleTags[rand.Intn(len(exampleTags))]
	}

	return &Example{
		Name:   name,
		Day:    day,
		Status: status,
		Date:   date,
		Tags:   tags,
	}, nil
}

// This is an unrelated variable that should be ignored by the enumify generator.
var exampleTags = []string{
	"low", "medium", "high",
	"red", "green", "blue",
	"primary", "secondary", "success", "danger", "warning", "info",
	"foo", "bar", "baz",
}

// This is an unrelated 2D array that should be ignored by the enumify generator.
var exampleNames = [][]string{
	{"curious", "ancient", "vibrant", "subtle", "brittle", "serene", "chaotic", "luminous", "hollow", "nimble", "terse", "ornate"},
	{"mountain", "river", "castle", "lighthouse", "violin", "compass", "telescope", "orchard", "glacier", "harbor", "cathedral", "parchment"},
}
