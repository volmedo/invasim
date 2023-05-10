# 👾 InvaSim 👾

The alien invasion simulator.

InvaSim reads a map file that describes a world by declaring cities and the roads between them. It then positions N aliens randomly in the map and simulates how they move and fight across the world.

The simulation ends either when all aliens have been destroyed or when each of them has moved at least 10,000 times. Once the simulation finishes, InvaSim will print to standard output how the world looks like after the invasion.

## Building

You will need a working Go installation (v1.20 or later) to build InvaSim. Check the [official docs](https://go.dev/doc/install) if you need help doing that.

The repo has a `Makefile` you can use if you have `make` installed that makes building the project a breeze:

```
$> make build
```

The previous command will get you an `./build/invasim` binary ready to use.

Alternatively, you can issue the build command yourself. That can be as easy as issuing:

```
$> go build ./cmd/invasim
```

## Running

InvaSim is a CLI tool. Run it in the terminal as:

```
$> invasim -map <path_to_map_file> -aliens <num_aliens>
```

where `<path_to_map_file>` is the path to the map file describing the world and `<num_aliens>` is the number of aliens that will be unleashed in the invasion.

> **Note**
>
> If you used `make build` previously to build the binary, remember that it will be at `./build/invasim`.

## Map file format

Worlds to be invaded are described by means of map files. These are regular text files that consist on a series of lines, where each line contains the declaration of a city along with the cities that can be reached from it taking roads in different directions. Each of these lines has the format `<city_name> [<road> [<road>]...]`, where `<city_name>` is a string. `<road>` is a pair `<direction>=<destination_city_name>`. `<direction>` can only be one of `"east"`, `"north"`, `"south"` and `"west"`.

If you are the kind of person that enjoys formal definitions, the map format can be expressed in EBNF notation as:

```ebnf
map file = city line , { city line } ;
city line = city name , {" " , road} ;
city name = ( alpha | digit ) , { alpha | digit } ;
road = direction , "=" , city name ;
direction = "east" | "north" | "south" | "west" ;
```

## Design and implementation

If you are interested in how InvaSim has been implemented and want to know more, check the [DESIGN](./DESIGN.md) doc.
