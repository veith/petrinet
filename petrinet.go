package petrinet

import (
	"sync"
	"github.com/Knetic/govaluate"
	"errors"
	"fmt"
)

type Net struct {
	InputMatrix        [][]int                `json:"-"`                   // Input Matrix
	OutputMatrix       [][]int                `json:"-"`                   // Output Matrix
	ConditionMatrix    [][]string             `json:"-"`                   // Condition Matrix
	State              []int                  `json:"-"`                   // State
	Variables          map[string]interface{} `json:"variables"`           // variablen die mit dem Prozess mitlaufen
	EnabledTransitions []int                  `json:"enabled_transitions"` // list of transitions which can be fired
}

func (net *Net) Init() {
	net.EnabledTransitions = net.evaluateNextPossibleTransitions();
}

/**
 * fire a transition.
 */
func (f *Net) Fire(transition int) error {
	var err error
	var mutex = &sync.Mutex{}
	if f.isTransitionEnabled(transition) {
		mutex.Lock()
		f.EnabledTransitions = f.fastfire(transition)
		mutex.Unlock()
		return err
	} else {
		err = errors.New(fmt.Sprintf("Transition %v not enabled", transition))
		return err
	}
}

func (net *Net) isTransitionEnabled(t int) bool {
	for _, b := range net.EnabledTransitions {
		if b == t {
			return true
		}
	}
	return false
}

// prüfe ob Transition ungeachtet der arc bedingungen gefeuert werden könnte
func (net *Net) fastCheck(transition int) bool {
	for place, p := range net.InputMatrix[transition] {
		if p != 0 && net.State[place]-p < 0 {
			return false
		}
	}
	return true
}

// finde die möglichen nächsten Transitionen und löse automatische Transitionen direkt aus
// Timed Transitionen werden gestartet
func (net *Net) evaluateNextPossibleTransitions() []int {
	var possibleTransitions []int

	// mögliche transitionen finden
	for t := 0; t < len(net.InputMatrix); t++ {
		if net.fastCheck(t) {
			possibleTransitions = append(possibleTransitions, t)
		}
	}

	var lockedTransitions []int
	// conditions prüfen
	for t := 0; t < len(possibleTransitions); t++ {

		// sobald eine Bedingung auf einem arc zu einer Transition nicht erfüllt ist, ist die Transition nicht mehr feuerbar
		// conditionen sind als string eingetragen
		if !net.proveConditions(t) {
			lockedTransitions = append(lockedTransitions, t)
			possibleTransitions = removeFromIntFromArray(possibleTransitions, t)
		}
	}
	return possibleTransitions
}

func removeFromIntFromArray(l []int, item int) []int {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func (net *Net) proveConditions(transitionIndex int) bool {
	if len(net.ConditionMatrix) > 0 {
		for _, condition := range net.ConditionMatrix[transitionIndex] {
			expression, err := govaluate.NewEvaluableExpression(condition);
			result, err := expression.Evaluate(net.Variables);
			if err != nil || !result.(bool) {
				return false
			}
		}
	}
	return true
}

// fire ohne Check
func (net *Net) fastfire(transition int) []int {
	for i, step := range net.InputMatrix[transition] {
		net.State[i] = net.State[i] - step
	}
	for i, step := range net.OutputMatrix[transition] {
		net.State[i] = net.State[i] + step
	}
	return net.evaluateNextPossibleTransitions();
}
