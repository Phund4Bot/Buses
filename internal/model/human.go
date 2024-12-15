package model

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const (
	minMinutes = 7
	maxMinutes = 15
)

type Human struct {
	ID          string
	WaitingTime time.Duration
}

func GeneratePeople(count int) []Human {
	people := make([]Human, 0, count)

	for i := 0; i < count; i++ {
		people = append(people, Human{
			ID:          uuid.NewString(),
			WaitingTime: time.Duration(rand.Intn(maxMinutes - minMinutes) + minMinutes) * time.Minute,
		})
	}

	return people
}
