# Invasion
![Invasion](./assets/invasion.png)

This extremely useful tool will help you with the most important task. Everyone have to deal with the danger of alien invasion. The main problem in this case estimation of possible damage to your planet.
 
Invasion allow you to simulate an alien invasion using the map of your planet. You just have to translate it in our own map format. You will be able to simulate a plenty of possible invasion scenarios.

## Map format
In order to have convenient format for describing planets we invented our own data format.

It MUST be a file with utf8 encoding.

Every line MUST contain one city name and from 1 to 4 directions.

City name MUST contain only alphanumeric characters, underscores or dashes

City names MUST be unique

Direction MUST be one of {south, east, west, north} and have format city {mapDirection}={city name}
For example 
```
south=City17
```

Direction MUST point only to cities which exist in the provided map file.

Every direction attached to a city MUST be unique for that city.
For example the following line is invalid as it contains two `north` directions.
```
City1 north=City2 north=City3
```

Different directions attached to the city CAN'T lead to the same city
For example the following line is invalid. Because two different directions point to the same city.
```
City1 north=City2 south=City2
```

Directions MAY be specified in any order.
For example the following lines mean the same
```
City1 north=City2 south=City3
 ```
```
City1 south=City3 north=City2
```

### Map example
```Paris south=Boston west=Los_Angeles east=Moscow
Boston north=New_York west=Berlin east=London
New_York north=London south=Moscow
Salehard4 north=Berlin south=Los_Angeles west=Paris east=New_York
Moscow north=London west=Salehard4 east=Berlin
Los_Angeles north=Paris south=London west=Salehard4 east=Boston
London north=Salehard4 south=Berlin  east=Paris
Berlin south=Boston west=Salehard4 east=Moscow
```

## Usage
Build invasion using 
```
make build
```
And then launch simulation with
```
.dest/invasion simulate path/to/map
```
By default, your planet will be invaded with 15 aliens, but you can change it with `-n` flag
```
.dest/invasion simulate path/to/map --n=40
```

