# Vehicle Routing Problem

This application provides a simple CLI that injests a list of loads and outputs a set of driver routes to complete all load pickup/dropoffs while attempting to minimize deadhead.

## Instructions

To run the cli, either:

* From a terminal, run `go run cmd/vrp/main.go <path to load file>`
    * This tool was built with `golang 1.21.6`, please ensure either this or a newer version of golang is installed prior to running this command
* Download the binary from the Releases page and execute it from the command line with `./vrp <path to load file>`

## Assumptions

There are a few assumptions made based on time constraints:

* File parsing assumes a header row is present in all test files. The first line of each file is dropped and will not be parsed as a load entry

## Future Enhancements

There are many potential ways to improve the current performance. This section will briefly explore ideas of how that can be achieved that fell outside of the time constraints

* Raw performance
    *
* More complex search heuristics
    * The current solution follows a Nearest Neighbor heuristic for choosing which load to assign to each route. This is somewhat short sighted and can miss efficiencies
    that may be available when looking one, or even more, loads ahead. This can lead to efficiencies at the cost of code complexity and potential computation time.

## Notes

This section intends to capture some of my thought process while building this tool

* There is a brute force solver that was used as my test code to familiarize myself with the problem while building out the data structures. This remains in the code, but is incomplete.
* Nearest Neighbor tweaks
    * The current heuristic assumes that each load added to a given route is the final load, meaning that the return-trip to the hub is part of the calculation when selecting the next load to pickup
        * This provides safety because if we can't find another load for a given driver, their return to hub is already factored into their current route. No recalculation needed.
    * I attempted to make certain aspects of the Nearest Neighbor more effective at selecting the next loads to add to the routes, however many of them actually yielded _worse_ performance against the test set
        * Selecting the next load to add based on the lowest cost addition to next load dropoff yielded ~1-2% slower mean cost to the solution

## References

* [Study vehicle routing problem using Nearest Neighbor Algorithm](https://iopscience.iop.org/article/10.1088/1742-6596/2421/1/012027/pdf#:~:text=Vehicle%20routing%20problem%20(VRP)%20has,get%20the%20most%20optimal%20results.)
* [nptelhrd Lecture 29](https://www.youtube.com/watch?v=A1wsIFDKqBk)
