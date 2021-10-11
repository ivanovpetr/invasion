package simulator

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
)

const (
	invasionDuration         = 10000
	cityDestructionThreshold = 2
)

type planetMap map[string]*city

// getRandomCity returns a random city on the map. If map is empty returns an empty string
func (p planetMap) getRandomCity() string {
	for name := range p {
		return name
	}
	return ""
}

// mapDirection one of the fourth directions which city can have
type mapDirection struct {
	directionType  string
	directionValue string
}

// city is the main part of a simulation
type city struct {
	name string
	// specifies either city is eligible for visit
	isDestroyed bool

	// City directions, can be a city name or empty string
	directions []mapDirection

	// Contains the identifiers of all aliens which are located in the city
	aliens []int64
}

// shuffleDirections shuffle directions slice order
func (c *city) shuffleDirections() {
	for i := range c.directions {
		j := rand.Intn(i + 1)
		c.directions[i], c.directions[j] = c.directions[j], c.directions[i]
	}
}

// addAlien adds the alien to the city. If the alien is already in the city, does nothing
func (c *city) addAlien(alienID int64) {
	for _, id := range c.aliens {
		if id == alienID {
			return
		}
	}
	c.aliens = append(c.aliens, alienID)
}

// removeAlien removes the alien from the city. If the alien is not in the city, does nothing.
func (c *city) removeAlien(alienID int64) {
	for i, a := range c.aliens {
		if a == alienID {
			c.aliens[i] = c.aliens[len(c.aliens)-1]
			c.aliens = c.aliens[:len(c.aliens)-1]
		}
	}
}

// battleMessage returns a battle message based on the city name and aliens in it
func (c *city) battleMessage() string {
	builder := strings.Builder{}
	builder.WriteString("Aliens: ")
	for i, a := range c.aliens {
		builder.WriteString("👾")
		builder.WriteString(strconv.Itoa(int(a)))
		if i != len(c.aliens)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString(fmt.Sprintf(" have met in the city of %s. ⚔ Battle destroyed the city.", c.name))

	return builder.String()
}

// alien is an earth invader which is moving from one city to another
type alien struct {
	// City where the alien is located
	city string

	// specifies either alien is alive and can move or dead
	isDead bool
}

// Simulation allows running invasion scenarios with different number of aliens
type Simulation struct {
	initialMap planetMap
}

// getMapCopy returns copy af a simulation map
func (s *Simulation) getMapCopy() planetMap {
	cp := make(planetMap)
	for name, c := range s.initialMap {
		cp[name] = c
	}

	return cp
}

// SimulationResult represents a final result of a simulation, contains resulted aliens and logs of simulation.
type SimulationResult struct {
	ResultMap planetMap
	Aliens    []alien
	Logs      []string
}

// PrintResultMap prints out result state of a map in the standard map format
func (sr *SimulationResult) PrintResultMap(out io.Writer) error {
	for _, c := range sr.ResultMap {
		if c.isDestroyed {
			continue
		}
		output := strings.Builder{}
		output.WriteString(c.name)
		for _, d := range c.directions {
			if sr.ResultMap[d.directionValue].isDestroyed {
				continue
			}
			output.WriteString(fmt.Sprintf(" %s=%s", d.directionType, d.directionValue))
		}
		output.WriteByte('\n')
		_, err := out.Write([]byte(output.String()))
		if err != nil {
			return fmt.Errorf("failed to print out result map: %w", err)
		}
	}

	return nil
}

// Run runs simulation with provided number of aliens, returns result of a simulation
// with battle logs, aliens and final state of a map
func (s *Simulation) Run(numberOfAliens int64) *SimulationResult {
	// All aliens take their actions simultaneously. There are three main simulation phases Spawn, Battle, Moving.
	// Spawn. Create aliens and put every of them in a random city
	simulationMap := s.getMapCopy()
	aliens := make([]alien, numberOfAliens)
	logs := []string{fmt.Sprintf("Simulate invasion with %d aliens", numberOfAliens)}

	for id := range aliens {
		aliens[id].city = simulationMap.getRandomCity()
		simulationMap[aliens[id].city].addAlien(int64(id))
	}

	for i := 0; i < invasionDuration; i++ {
		// Battle stage. Try to begin a battle in every city.
		for _, c := range simulationMap {
			if c.isDestroyed {
				continue
			}
			if len(c.aliens) >= cityDestructionThreshold {
				// destroy the city and the aliens
				c.isDestroyed = true
				for _, id := range c.aliens {
					aliens[id].isDead = true
				}
				logs = append(logs, c.battleMessage())
			}
		}

		// Moving. Move every alien to a new destination
		moves := 0
		aliveAliens := 0
		for id, a := range aliens {
			// if alien is dead, don't move him
			if a.isDead {
				continue
			}
			// increment number of alive aliens for the current iterations
			aliveAliens++
			// shuffle directions order in order to start with a random one
			simulationMap[a.city].shuffleDirections()
			for _, direction := range simulationMap[a.city].directions {
				// if the mapDirection leads to a not destroyed city move alien to it
				if !simulationMap[a.city].isDestroyed {
					// remove alien from its current city
					simulationMap[a.city].removeAlien(int64(id))
					// add alien to the new location
					simulationMap[direction.directionValue].addAlien(int64(id))
					// update alien location
					aliens[id].city = direction.directionValue
					// increment moves counter
					moves++
					break
				}
			}
		}
		// if all aliens are dead or locked without ability to move then end the simulation.
		if aliveAliens == 0 {
			logs = append(logs, fmt.Sprintf("All aliens are dead, simulations is over on turn number %d", i))
			break
		} else if moves == 0 {
			logs = append(logs, fmt.Sprintf("All aliens are either dead or locked, simulations is over on turn number %d", i))
			break
		}

		if i == invasionDuration-1 {
			logs = append(logs, fmt.Sprintf("%d turns are finished. Simulation is over", invasionDuration))
		}
	}

	return &SimulationResult{
		ResultMap: simulationMap,
		Aliens:    aliens,
		Logs:      logs,
	}
}
