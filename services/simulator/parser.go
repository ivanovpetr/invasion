package simulator

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/scanner"
	"unicode"
)

const (
	directionSouth = "south"
	directionNorth = "north"
	directionWest  = "west"
	directionEast  = "east"
)

type parsedCity struct {
	line                     int
	north, east, south, west string
}

func (pc *parsedCity) directionExists(name string) bool {
	return pc.north == name || pc.east == name || pc.south == name || pc.west == name
}

func (pc *parsedCity) isDirectionSet(name string) bool {
	switch name {
	case directionSouth:
		return pc.south != ""
	case directionNorth:
		return pc.north != ""
	case directionEast:
		return pc.east != ""
	case directionWest:
		return pc.west != ""
	default:
		return false
	}
}

func (pc *parsedCity) setDirection(direction, value string) {
	switch direction {
	case directionSouth:
		pc.south = value
	case directionNorth:
		pc.north = value
	case directionEast:
		pc.east = value
	case directionWest:
		pc.west = value
	}
}

func (pc *parsedCity) getDirections() []mapDirection {
	directions := make([]mapDirection, 0, 4)
	if pc.south != "" {
		directions = append(directions, mapDirection{
			directionType:  directionSouth,
			directionValue: pc.south,
		})
	}
	if pc.north != "" {
		directions = append(directions, mapDirection{
			directionType:  directionNorth,
			directionValue: pc.north,
		})
	}
	if pc.west != "" {
		directions = append(directions, mapDirection{
			directionType:  directionWest,
			directionValue: pc.west,
		})
	}
	if pc.east != "" {
		directions = append(directions, mapDirection{
			directionType:  directionEast,
			directionValue: pc.east,
		})
	}
	return directions
}

type expectation byte

const (
	expectCity expectation = iota
	expectDirectionType
	expectDirectionValue
	expectEqualSign
)

type parserError struct {
	message  string
	position scanner.Position
}

func (p parserError) Error() string {
	return fmt.Sprintf("%s:%d:%d: %s", p.position.Filename, p.position.Line, p.position.Column, p.message)
}

func newParserError(position scanner.Position, msg string) error {
	return parserError{
		message:  msg,
		position: position,
	}
}

type parser struct {
	parsed             bool
	currentExpectation expectation
	currentCity        string
	currentDirection   string
	s                  scanner.Scanner
	parsedCities       map[string]*parsedCity
}

func newParser(src io.Reader, filename string) *parser {
	p := &parser{
		parsedCities: map[string]*parsedCity{},
	}
	p.s.Init(src)
	p.s.Filename = filename
	p.s.Whitespace ^= 1 << '\n'
	p.s.IsIdentRune = func(ch rune, i int) bool {
		return ch != '=' && (ch > '!' && ch < '~' || unicode.IsLetter(ch))
	}
	return p
}

func (p *parser) currentCityHasAtLeastOneDirection() bool {
	return len(p.parsedCities[p.currentCity].getDirections()) != 0
}

func (p *parser) parse() error {
	if p.parsed {
		return errors.New("already parsed")
	}

	for tok := p.s.Scan(); tok != scanner.EOF; tok = p.s.Scan() {
		switch tok {
		case '\n':
			if !p.currentCityHasAtLeastOneDirection() {
				return newParserError(p.s.Pos(), "unexpected newline, city must contain at least one direction")
			}
			p.currentExpectation = expectCity
		default:
			err := p.handleToken(p.s.TokenText())
			if err != nil {
				return err
			}
		}
	}

	if p.currentExpectation == expectEqualSign || p.currentExpectation == expectDirectionValue {
		return newParserError(p.s.Position, "Unexpected EOF")
	}

	if !p.currentCityHasAtLeastOneDirection() {
		return newParserError(p.s.Pos(), "unexpected EOF, city must contain at least one direction")
	}

	p.parsed = true
	return nil
}

func (p *parser) handleToken(token string) error {
	switch p.currentExpectation {
	case expectCity:
		err := p.handleCityToken(token)
		if err != nil {
			return err
		}
	case expectDirectionType:
		err := p.handleDirectionType(token)
		if err != nil {
			return err
		}
	case expectEqualSign:
		err := p.handleEqualSign(token)
		if err != nil {
			return err
		}
	case expectDirectionValue:
		err := p.handleDirectionValue(token)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *parser) handleCityToken(token string) error {
	// check for existence
	if _, ok := p.parsedCities[token]; ok {
		return newParserError(p.s.Pos(), fmt.Sprintf("got city duplication for %s previously declared on line %d", token, p.parsedCities[token].line))

	}
	// validate city
	if !isValidCityName(token) {
		return newParserError(p.s.Pos(), fmt.Sprintf("expected a valid city name, got %s", token))
	}
	// write new city
	p.parsedCities[token] = &parsedCity{
		line: p.s.Pos().Line,
	}
	p.currentCity = token
	p.currentExpectation = expectDirectionType
	return nil
}

func (p *parser) handleDirectionValue(token string) error {
	// validate city
	// micro optimization: try to check mapDirection value against parsed city and avoid usage of regexp
	if _, ok := p.parsedCities[token]; !ok {
		if !isValidCityName(token) {
			return newParserError(p.s.Pos(), fmt.Sprintf("expected a valid city name as a mapDirection value, got %s", token))
		}
	}
	// check against mapDirection value duplication
	if p.parsedCities[p.currentCity].directionExists(token) {
		// mapDirection value duplication
		return newParserError(p.s.Pos(), fmt.Sprintf("got mapDirection value duplication %s for city %s", token, p.currentCity))
	}
	// write mapDirection to current city current mapDirection
	p.parsedCities[p.currentCity].setDirection(p.currentDirection, token)
	p.currentExpectation = expectDirectionType

	return nil
}

func (p *parser) handleDirectionType(token string) error {
	// validate mapDirection
	if !isValidDirection(token) {
		// unexpected mapDirection type
		return newParserError(p.s.Pos(), fmt.Sprintf("got unexpected mapDirection type %s, expected one of south,north,west,east", token))
	}
	// check for duplication
	if p.parsedCities[p.currentCity].isDirectionSet(token) {
		// mapDirection type duplication
		return newParserError(p.s.Pos(), fmt.Sprintf("got mapDirection type duplication %s for city %s", token, p.currentCity))
	}

	// write current mapDirection
	p.currentDirection = token
	p.currentExpectation = expectEqualSign
	return nil
}

func (p *parser) handleEqualSign(token string) error {
	if token != "=" {
		return newParserError(p.s.Pos(), fmt.Sprintf("unexpected token %s, expected =", token))
	}
	p.currentExpectation = expectDirectionValue
	return nil
}

func (p *parser) checkDirectionValuesExistence() error {
	if !p.parsed {
		return errors.New("cannot check direction value for unparsed file")
	}
	for n, c := range p.parsedCities {
		for _, d := range c.getDirections() {
			if _, ok := p.parsedCities[d.directionValue]; !ok {
				return fmt.Errorf("city %s on line %d has direction %s which points to non existent city %s", n, c.line, d.directionType, d.directionValue)
			}
		}
	}
	return nil
}

func (p *parser) buildSimulation() (*Simulation, error) {
	// check parsed flag
	if !p.parsed {
		return nil, errors.New("cannot build simulation for unparsed file")
	}
	// go through parsed cities
	// fulfill simulation
	s := &Simulation{initialMap: map[string]*city{}}
	for name, pc := range p.parsedCities {

		s.initialMap[name] = &city{
			name:        name,
			isDestroyed: false,
			directions:  pc.getDirections(),
		}
	}
	return s, nil
}

func CreateSimulationFromPath(path string) (*Simulation, error) {
	mapFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer mapFile.Close()

	return createSimulation(mapFile, filepath.Base(path))
}

func createSimulation(src io.Reader, filename string) (*Simulation, error) {
	p := newParser(src, filename)
	err := p.parse()
	if err != nil {
		return nil, err
	}
	err = p.checkDirectionValuesExistence()
	if err != nil {
		return nil, err
	}
	return p.buildSimulation()
}
