# Introduction #

wgen is a stand-alone program used to generate worlds for minima.  It takes a small number of parameters (including the size of the world and optionally the random seed) and outputs a world file to standard output.  The world can be read by both the game proper or by other tools which, in turn, may output additional information, such as enemy or item locations, along with the world.  This document describes the wgen program an the techniques that it uses to generate random landscape.

wgen proceeds in two distinct phases.  The first phase generates the elevations of each location and the second phase adds terrain types.  These two phases are described in the following sections.

# Elevation Generation #

Worlds in minima are comprised of a grid of locations that are connected in a torus shape.  This means that when you walk beyond the edge you wrap around it.  When generating the elevations, they must be _tileable_, e.g.  the elevations at the bottom edge of the world must fit well with the elevations at the top edge of the world, and likewise for the left and right edges.

Each location's elevation is initialized to half of the maximum world height.  wgen then `grows' the landscape using a mixture of random 2D normal distributions, a.k.a. Gaussians (http://en.wikipedia.org/wiki/Normal_distribution).  The center point for the mean of each Gaussian is selected uniformly from all locations in the world.  The height of the Gaussian is chosen from a normal distribution; negative heights will create valleys and positive heights will create mountains.  The density of each Gaussian is added to the elevation of all locations within approximately 2 standard deviations of the mean, taking into account wrapping if these locations are beyond an edge of the world.

The elevation generation is controlled by a variety of parameters that are documented in the code (see wgen.go).  After the elevations have been generated, each location's elevation is clamped to fall between 0 and the maximum world elevation.  Next the terrain is generated.

# Terrain Generation #

The terrain generation phase assigns a terrain type to each location.  Initially there are only two terrain types: mountains are above a given threshold and everything else is grass.  After initializing the terrain, other features are added such as oceans, lakes, forests, deserts, and glaciers.  These features are added by either flooding and growing.  Flooding adds water or other liquids by filling up local minima to certain heights, and growing begins with a set of seed locations and grows out randomly from there.  These two techniques are described in the following sections.

## Flooding ##

As mentioned previously, flooding is used to add pools of liquid to the world.  For efficiency, flooding does not work on a location-by-location basis, instead it make use of a data structure called a topographical graph (topo graph).  Nodes in the topo graph are called contours and each contour represents a connected group of locations with the same _height_, where the height of a location is its elevation minus its depth.  Two contours _c0_ and _c1_ are connected by an edge if and only if there are two locations _l0_ and _l1_ such that _l0_ is a member of _c0_ and _l1_ is a member of _c1_ and _l0_ and _l1_ are adjacent.  Typically the topo graph is much smaller than the graph representing the locations in the world.  For a 500x500 world, there are 25,000 locations, however there is usually only a few thousand contours.

The flooding algorithm starts by finding all local minima in the topo graph, it then selects random minima and considers flooding them to different random heights.  Each minimum is considered only a single time so the algorithm is guaranteed to terminate.  The contours that would be submerged by a flood are found by using a depth-first search in the topo graph starting from the initial minimum and pruning adjacent contours that are too high to be submerged.  When considering a flood, the algorithm compares the number of locations that would be submerged with a set of parameters used to differentiate between placing large oceans and placing smaller lakes.  The algorithm terminates when it has either considered each local minimum or when it has placed at least a given amount of new liquid; the latter case is the most common.

## Growing ##

After water has been added to the world, the terrain generation process adds forests, deserts, and glaciers.  These features are added by a growing technique.  The growing algorithm also works on with a topo graph instead of locations,  not for efficiency but to give the grown features a less regular look.  The algorithm begins by finding all contours with an acceptable terrain type upon which the new terrain can be grown (the acceptable terrains are given as a parameter and differ for each different feature).  A given number of seed contours are selected from the acceptable set and are converted to the terrain type that is being grown.  After the seeds are chosen, the algorithm iteratively selects a random contour to convert from among those locations adjacent to contours that have already been converted.  The algorithm terminates when at least a specified number of locations have been converted.

Grown features can be given different looks by changing the number of seeds and the number of locations to be converted.  If there are few seeds then the grown feature will be more rare like deserts and glaciers, however, if there are more seeds then the feature will be more common like forests.

## Rivers ##

The river placement works by generating costs for each location of the map using Perlin noise.  Next, a set of starting locations are chosen from the set of locations that are mountains or water.  Each river is a shortest path (with respect to the Perlin costs) from a starting point to an ocean.  Rivers that are too short are then discarded, this is repeated until a desired number of river locations are added.