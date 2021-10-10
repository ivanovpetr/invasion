package simulator

import "regexp"

const cityConstraint = "^[a-zA-Z1-9\\-\\_]+$"

var (
	cityRegex = regexp.MustCompile(cityConstraint)
)

func isValidCityName(name string) bool {
	return cityRegex.Match([]byte(name))
}

func isValidDirection(name string) bool {
	switch name {
	case directionNorth:
		return true
	case directionSouth:
		return true
	case directionEast:
		return true
	case directionWest:
		return true
	default:
		return false
	}
}
