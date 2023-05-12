# Design ideas behind InvaSim

This doc describes the main ideas and considerations that shaped the implementation of InvaSim.

## Assumptions

Natural languages are far more flexible than formal ones. They are also ambiguous, and this ambiguity emanates from their flexibility. While the original requirements doc is quite clear in terms of what needs to be implemented, there are details that leave some room for interpretation.

The following assumptions have an impact on how the program works:

- There can't be more than one alien in a city in the initial state. This implies that the number of aliens can never be greater than the number of cities in the given world.

- Fights can involve more than two aliens. Aliens move randomly in the world, if two or more of them end up in the same place, all of them will be destroyed, along with the city.

- Map files are parsed strictly. The parser doesn't trim leading or trailing spaces or interprets several spaces as one. If the file is not formatted exactly as expected, an error will be produced.

## Implementation, data structures and access patterns

It is usual for simulation software to be implemented as a loop where each iteration represents the events that take place in a given point in time. This is how InvaSim is implemented. Each iteration, the program simulates the movement of the different aliens invading the world. It then processes any battles that can take place and updates the state to reflect their results.

This state contains two main pieces of data: how the world evolves from its initial state and the positions of the different aliens in that world. To decide which data structures are better suited to hold that information, it is essential to reason about the way the data will be accessed.

### World

The world is an undirected, potentially cyclic, graph. The most common data structures used to represent graphs are the adjacency matrix and the adjacency list. From the two, the adjacency list is chosen in this case as it allows checking adjacency in constant time (`O(1)`). The adjacency check is the most relevant operation as it is used both to find adjacent cities an alien can travel to from a given origin and to find both ends of roads that need to be destroyed when a city is destroyed.

To improve the check, roads to adjacent cities are stored as a map instead of the usual array. This allows finding the corresponding road (interesting when looking for specific road ends) also in constant time. The difference in performance is minimal, though, since the maximum degree of any given vertex in our graph is 4 and traversing a 4-element array is not an expensive operation.

Thus, the `World` is a map from city name keys to `Road` maps, which in turn are maps from `Direction`s to destination city names.

### Alien tracking

To keep track of the position of each alien on the map, another data structure is used. This information could have been embedded in the world representation. Aside from clearly separating concerns, having a separate data structure allows iteration over the aliens that still exist rather than iterating over the cities in the world looking for aliens to move or destroy. There is a performance gain in doing so, because the number of aliens will always be less or equal than the number of cities, and it also decreases faster.

The `Tracker` is the type that holds the information about the placement of each alien. It is a map from the alien name to the city name the alien is currently at. This city will be used as a key into the world map when looking for places an alien can move to.

### Auxiliary structure: visited cities

A third data structure is used to efficiently look up which cities and aliens are involved in battles during each iteration of the simulation. As opposed to the world map or the alien tracker, which are mutated from an initial state as the simulation progresses, this structure is scoped to each iteration.

When aliens move to their new destinations during a given iteration, the cities they end up at are collected in a `VisitedCities` map. Keys in the map are city names and values are slices of alien names, which represent the aliens that are currently in that city. This is exactly the data required to know what cities and aliens are to be destroyed.

## Additional considerations

Aside from the implementation details provided above, there are other design decisions that may be interesting to highlight:

- The program is structured in a modular way, where each package has clear responsibilities. The `main` function is only charged with parsing the parameters in the command line and setting up dependencies.

- At first I was inclined to use the non-deterministic iteration order exposed by maps in Go. After reading the [spec](https://go.dev/ref/spec#For_statements) a second time, however, I ended up using the `rand` package as a more traditional source of randomness. The spec states that

  > The iteration order over maps is not specified and is not guaranteed to be the same from one iteration to the next.

  The behaviour is not specified in either way, so counting on it being random feels as much of an error as counting on it being deterministic.

- Functions have side effects. The simulation progresses by mutating some initial state. Functions and methods receive the data structures storing that state and update them in place. I tend to prefer pure functions that don't mutate input parameters. In this case, however, it made sense to make the trade-off for performance reasons. Updating the state in place avoids the potentially expensive operations of creating new data structures and copying the required elements over.

- The project doesn't include end to end tests. They didn't seem to add a lot of value in this case because, as mentioned before, `main` doesn't contain any logic related with the simulation itself. The functions used by `main` are already covered by unit tests. That means the only code e2e tests would cover that is not covered yet is parameter parsing, which is done via the `flag` package, and producing the expected error messages when parameters don't have the expected values, which is not critical.
