What they've done here is build 4 different linear feedback shift registers with multiple feedback bits. They're asking us to find the cycle length of each, and then the LCM of those lengths.

I figured this much out by looking at the graph of the connections and recognizing it from my hardware engineering days.

I've checked in the dot file that shows this.

I "solved" it by special-casing the end states of each LFSR, printing out the cycle times for each, and calculating the LCM of them by hand.
It made me feel dirty.