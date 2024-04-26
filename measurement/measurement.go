// Copyright (c) 2024 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package measurement

import "github.com/tiagomelo/golang-waterjug-api/measurement/models"

// state represents the state of the water jugs, including their amounts and the previous state.
type state struct {
	x      int
	y      int
	action string
	status string
	prev   *state
}

// gcd computes the greatest common divisor (GCD) of two integers, a and b.
// It uses Euclid's algorithm, which is an efficient method for computing
// the GCD. The algorithm repeatedly replaces the larger number by its remainder
// when divided by the smaller number until one of the numbers is zero.
// The non-zero number at this point is the GCD of a and b.
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// bfs performs a breadth-first search (BFS) to find the minimum steps required
// to measure exactly the target amount of water using two jugs with capacities xMax and yMax.
func bfs(xMax, yMax, target int) *state {
	if gcd(xMax, yMax) != gcd(xMax, target) {
		return nil
	}
	visited := make(map[[2]int]bool)
	queue := []*state{{x: 0,
		y:      0,
		action: "Start",
		prev:   nil,
	}}
	visited[[2]int{0, 0}] = true
	for len(queue) > 0 {
		currentState := queue[0]
		queue = queue[1:]
		possibleStates := []*state{
			{x: xMax, y: currentState.y, action: "Fill bucket X", prev: currentState},
			{x: currentState.x, y: yMax, action: "Fill bucket Y", prev: currentState},
			{x: 0, y: currentState.y, action: "Empty bucket X", prev: currentState},
			{x: currentState.x, y: 0, action: "Empty bucket Y", prev: currentState},
			{x: currentState.x - min(currentState.x, yMax-currentState.y), y: currentState.y + min(currentState.x, yMax-currentState.y), action: "Transfer from bucket X to Y", prev: currentState},
			{x: currentState.x + min(currentState.y, xMax-currentState.x), y: currentState.y - min(currentState.y, xMax-currentState.x), action: "Transfer from bucket Y to X", prev: currentState},
		}
		for _, nextState := range possibleStates {
			if nextState.x == target || nextState.y == target {
				nextState.status = "Solved"
				return nextState // found the solution.
			}
			if !visited[[2]int{nextState.x, nextState.y}] {
				visited[[2]int{nextState.x, nextState.y}] = true
				queue = append(queue, nextState)
			}
		}
	}
	return nil
}

// Measure calculates the solution to the water jug problem.
func Measure(xMax, yMax, target int) *models.Solution {
	s := bfs(xMax, yMax, target)
	if s != nil {
		solution := &models.Solution{}
		states := []*state{}
		for stateStep := s; stateStep != nil; stateStep = stateStep.prev {
			states = append(states, stateStep)
		}
		var stepNumber int
		for i := len(states) - 2; i >= 0; i-- { // Start from len(steps) - 2 to skip the initial state
			stepNumber++
			step := &models.Step{
				Number:  stepNumber,
				BucketX: states[i].x,
				BucketY: states[i].y,
				Action:  states[i].action,
				Status:  states[i].status,
			}
			solution.Steps = append(solution.Steps, step)
		}
		return solution
	}
	return nil
}
