// petrinet is a simple petri net execution library
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
	TokenIds           [][]int                `json:"token-identifier"`    // State
	Variables          map[string]interface{} `json:"variables"`           // variablen die mit dem Prozess mitlaufen
	EnabledTransitions []int                  `json:"enabled_transitions"` // list of transitions which can be fired
}

var tokenID int // token id counter

func (net *Net) Init() {
	tokenID = 0
	// tokenIDs vergeben
	net.TokenIds = make([][]int, len(net.InputMatrix[0]))

	// token ids for initial state
	for i, tokens := range net.State {
		for n := 0; n < tokens; n++ {
			net.TokenIds[i] = append(net.TokenIds[i], net.nextTokenID())
		}
	}

	net.EnabledTransitions = net.evaluateNextPossibleTransitions();
}

func (net *Net) nextTokenID() int {
	tokenID++
	return tokenID
}

func (net *Net) FireWithTokenId(transition int, tokenID int) error {

	for p, _ := range net.InputMatrix[transition] {

		if len(net.TokenIds[p]) > 0 { // nur wenn es elemente hat prüfen
 			if net.TokenIds[p][0] != tokenID {
				for index, tokenIDVal := range net.TokenIds[p] {
					// tokenID finden und TokenID an erste stelle setzen
					if (tokenIDVal == tokenID) {
						a := net.TokenIds[p][index]
						b := net.TokenIds[p][0]
						net.TokenIds[p][index] = b
						net.TokenIds[p][0] = a
						return net.Fire(transition)
					}
				}
			}else{
				return net.Fire(transition)
			}
		}
	}

	return errors.New("TokenID not found")
}

// fires an enabled transition.

func (f *Net) Fire(transition int) error {
	var err error
	var mutex = &sync.Mutex{}

	if f.TransitionEnabled(transition) {
		mutex.Lock()

		f.EnabledTransitions = f.fastfire(transition)

		mutex.Unlock()
		return err
	} else {
		err = errors.New(fmt.Sprintf("Transition %v not enabled", transition))
		return err
	}
}

func (net *Net) TransitionEnabled(t int) bool {
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

	for place, step := range net.InputMatrix[transition] {
		net.State[place] = net.State[place] - step

		// id tokens entfernen (an erster stelle)
		for n := 0; n < step; n++ {
			// pop

			net.TokenIds[place] = net.TokenIds[place][1:]

		}
	}
	for place, step := range net.OutputMatrix[transition] {
		net.State[place] = net.State[place] + step

		// id tokens erzeugen (an letzter stelle)
		for n := 0; n < step; n++ {
			// push
			nextID := net.nextTokenID()

			net.TokenIds[place] = append(net.TokenIds[place], nextID)

		}

	}

	return net.evaluateNextPossibleTransitions();
}