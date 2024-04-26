# Algorithm

## Overview
The solution to the water jug problem implemented here uses a breadth-first search (BFS) algorithm. BFS is chosen because it ensures the discovery of the shortest path to the solution, which is the minimum number of steps required to measure exactly the target amount of water using two jugs with specified capacities.

## Detailed Explanation

### Initial Setup
The BFS algorithm starts with an initial state where both jugs are empty. This state is represented by `(x=0, y=0)` and enqueued as the starting point of our search. A `visited` map is used to keep track of visited states to prevent processing the same state multiple times.

### State Representation
Each state of the jugs is represented as a tuple `(x, y)`, where `x` and `y` are the current amounts of water in each jug, respectively.

### State Transition
From any given state, the following transitions are possible:

1. Fill either jug to its maximum capacity.
2. Empty either jug completely.
3. Pour water from one jug to the other until one jug is either full or the other is empty.

### Goal Check
At each step, the algorithm checks whether the target amount has been measured in either jug. If so, the current sequence of actions leading to this state forms a valid solution.

## BFS Algorithm Steps

1. Initialize the queue with the starting state `(0, 0)`.
2. Loop until the queue is empty:
    - Dequeue the front state.
    - Generate all possible states from the current state based on the state transitions.
    - For each generated state:
        - Check if it meets the goal criteria (either jug equals the target amount).
        - If it does, trace back the steps taken to reach this state and return the solution.
        - If not, and if the state has not been visited, mark it as visited and enqueue it.
3. If the queue is exhausted without finding a solution, return that no solution is possible.

## Considerations

**Greatest Common Divisor (GCD)**: Before starting BFS, the algorithm checks if the GCD of the two jugs' capacities is a divisor of the target volume. If not, it is proven by the theory of Diophantine equations that no solution exists.

## Optimization: Bidirectional Search

### Overview

To enhance the efficiency of the solution, especially for large inputs, the algorithm utilizes a bidirectional search strategy. This advanced technique involves simultaneously running two breadth-first search (BFS) operations:

1. **Forward Search**: Starts from the initial state with both jugs empty.
2. **Reverse Search**: Begins from a set of goal states where the target amount is achieved in one or both jugs.

### Execution
The bidirectional search aims to meet in the middle, significantly reducing the number of states that need to be explored compared to a traditional BFS approach. This optimization can dramatically decrease both the time complexity and space complexity of the solution by shortening the paths that both searches need to explore.

### Implementation Details

- **Initialization**: Both searches are initialized with their respective starting states:
    - Forward search starts with `(0, 0)`, indicating both jugs are empty.
    - Reverse search starts from states where either jug contains the target amount, e.g., `(target, 0)` or `(0, target)`, and other feasible combinations based on the jug capacities and the target.

- **State Space Exploration**:

- Each search independently explores possible transitions (filling, emptying, transferring water).
- States from both searches are stored in separate visited sets to track progress and prevent reprocessing of the same state in each 
respective search.

- **Meeting Point**:

- The algorithm continuously checks for intersections between the sets of states visited by the two searches.
- Once a common state is found, the algorithm concatenates the paths from the initial state to the meeting state (from the forward search) and from the meeting state to the goal state (from the reverse search), creating a complete path from start to finish.

## Benefits
This approach halves the search space required by each BFS operation, as each only needs to explore up to the middle of the total path length instead of the entire path. It's particularly effective in cases with large state spaces, making it feasible to find solutions more quickly and with less memory usage.