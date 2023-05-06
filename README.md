# 👾 InvaSim 👾

The alien invasion simulator.

InvaSim reads map files that describe cities and the roads between them. It then positions N aliens randomly in the map and simulates how they move and fight across the map. It finally prints to standard output how the world looks like when the invasion ends (all aliens have been destroyed or each of them has moved at least 10,000 times).

# Running InvaSim

```
$> invasim <path_to_map_file> <num_aliens>
```

Where `<path_to_map_file>` is the path to the map file describing the world and `<num_aliens>` is the number of aliens that will be unleashed in the invasion.
