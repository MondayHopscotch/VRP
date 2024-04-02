# Vehicle Routing Problem

This application provides a simple CLI that ingests a list of loads and outputs a set of driver routes to complete all load pickup/dropoffs while attempting to minimize the cost of doing so.

## Instructions

From a terminal, run `go run cmd/vrp/main.go <path to load file>`

   * This tool was built with `golang 1.21.6`, please ensure either this or a newer version of golang is installed prior to running this command

## Assumptions

Example input file
```
loadNumber pickup dropoff
1 (-50.1,80.0) (90.1,12.2)
2 (-24.5,-19.2) (98.5,1,8)
3 (0.3,8.9) (40.9,55.0)
4 (5.3,-61.1) (77.8,-5.4)

```

There are a few assumptions made based on time constraints:

* File parsing assumes a header row is present in all test files. The first line of each file is dropped and will not be parsed as a load entry
* Load numbers are expected to be sequential starting from `1`

## Notes

This section intends to capture some of my thought process while building this tool

* Nearest Neighbor
    * The current heuristic assumes that each load added to a given route is the final load, meaning that the return-trip to the hub is part of the calculation when selecting the next load to pickup
        * This provides safety because if we can't find another load for a given driver, their return to hub is already factored into their current route. No recalculation needed.
    * I attempted to make certain aspects of the Nearest Neighbor more effective when selecting the next loads to add to the routes, however many of them actually yielded _worse_ performance against the test set
        * Selecting the next load to add based only on the lowest cost addition to next load dropoff (not accounting for the return to hub time) yielded ~1-2% slower mean cost to the solution
    * Getting the proper number of drivers proved underwhelming with NN. I automatically add drivers as loads are unable to be completed by the current set of drivers
        * If additional drivers were needed, I attempted to recalculate with varying the number of starting drivers as this intuitively feels like it would yield more optimal routes. However, my NN implementation seemed to yield
          the same results, even when increasing the starting driver count incrementally. For readability and CPU consumption, I reduced this to a single recalculation.
 * Some test scenarios proved difficult to write a concise unit test to cover, however the provided test samples were able to be used to target certain scenarios such as underestimating the number of drivers needed.

## Future Enhancements

There are many potential ways to improve the current performance. This section will briefly explore ideas of how that can be achieved that fell outside of the time constraints

* More complex search heuristics
    * The current solution follows a Nearest Neighbor heuristic for choosing which load to assign to each route. This is somewhat short sighted and can miss efficiencies
    that may be available when looking one, or even more, loads ahead. This can lead to efficiencies at the cost of code complexity and potential computation time.
* Raw performance
    * Despite some of the recalculations with the current heuristic not yielding improved routes/cost, I'm certain there is a heuristic where recalculating with a different set of
      initial state will yield improved routing. In this scenario, these recalculations can be executed asynchronously to improve overall execution time.


## References

* [Study vehicle routing problem using Nearest Neighbor Algorithm](https://iopscience.iop.org/article/10.1088/1742-6596/2421/1/012027/pdf#:~:text=Vehicle%20routing%20problem%20(VRP)%20has,get%20the%20most%20optimal%20results.)
* [nptelhrd Lecture 29](https://www.youtube.com/watch?v=A1wsIFDKqBk)
* [VRP Wikipedia](https://en.wikipedia.org/wiki/Vehicle_routing_problem)
