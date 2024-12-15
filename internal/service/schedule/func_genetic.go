package schedule

import (
	"fmt"
	"math/rand"
	"time"

	"kurs/internal/service/logger"
	"kurs/internal/utils"
)

type genetic struct {
	numGenerations int
	populationSize int
	mutationRate   float64
	population     [][]int
}

func NewGenetic(numGenerations, populationSize int, mutationRate float64) Genetic {
	return &genetic{
		numGenerations: numGenerations,
		populationSize: populationSize,
		mutationRate:   mutationRate,
	}
}

func (g *genetic) initializePopulation() {
	g.population = make([][]int, g.populationSize)
	for i := 0; i < g.populationSize; i++ {
		individual := make([]int, 8)
		for j := 0; j < 8; j++ {
			individual[j] = rand.Intn(2) + 1
		}
		g.population[i] = individual
	}
}

func (g *genetic) fitness(driverTypes []int) (int, error) {
	schedule, err := Initialize(driverTypes)
	if err != nil {
		return -1, fmt.Errorf("cannot initialize schedule: %v", err)
	}
	return schedule.RunSimulation(), nil
}

func (g *genetic) selectParents(fitnessScores []int) ([]int, []int) {
	totalFitness := 0.0
	probabilities := make([]float64, len(fitnessScores))

	for _, score := range fitnessScores {
		totalFitness += 1.0 / float64(1+score)
	}

	for i, score := range fitnessScores {
		probabilities[i] = (1.0 / float64(1+score)) / totalFitness
	}

	parent1 := g.population[selectRandomIndex(probabilities)]
	parent2 := g.population[selectRandomIndex(probabilities)]
	return parent1, parent2
}

func (g *genetic) crossover(parent1, parent2 []int) ([]int, []int) {
	point := rand.Intn(len(parent1)-1) + 1
	child1 := append(parent1[:point], parent2[point:]...)
	child2 := append(parent2[:point], parent1[point:]...)
	return child1, child2
}

func (g *genetic) mutate(individual []int) []int {
	if rand.Float64() < g.mutationRate {
		index := rand.Intn(len(individual))
		if individual[index] == 1 {
			individual[index] = 2
		} else {
			individual[index] = 1
		}
	}
	return individual
}

func (g *genetic) RunGenetic() []int {
	logger.Info(time.Now(), "Starting genetic algorithm")
	defer utils.LogElapsed("Finished genetic algorithm")()

	var (
		err          error
		bestSolution []int
		bestResult  int
	)

	g.initializePopulation()
	for generation := 0; generation < g.numGenerations; generation++ {

		fitnessScores := make([]int, g.populationSize)
		for i, individual := range g.population {
			fitnessScores[i], err = g.fitness(individual)
			if err != nil {
				logger.Errorf(time.Now(), "Cannot fitness: %v", err)
			}
		}

		newPopulation := [][]int{}
		for len(newPopulation) < g.populationSize {
			parent1, parent2 := g.selectParents(fitnessScores)
			child1, child2 := g.crossover(parent1, parent2)
			newPopulation = append(newPopulation, g.mutate(child1))
			if len(newPopulation) < g.populationSize {
				newPopulation = append(newPopulation, g.mutate(child2))
			}
		}
		g.population = newPopulation

		for i, score := range fitnessScores {
			if score > bestResult {
				bestResult = score
				bestSolution = g.population[i]
			}
		}

		logger.Infof(time.Now(), "Generation %d: Solution = %v, Result=%v", generation+1, bestSolution, bestResult)
	}

	return bestSolution
}
