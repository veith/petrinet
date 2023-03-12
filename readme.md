# a petrinet library

PACKAGE DOCUMENTATION

package petrinet
    import "./"

    petrinet is a simple petri net execution library



TYPES

type Net struct {
    InputMatrix        [][]int                `json:"-"`                   // Input Matrix
    OutputMatrix       [][]int                `json:"-"`                   // Output Matrix
    ConditionMatrix    [][]string             `json:"-"`                   // Condition Matrix
    State              []int                  `json:"-"`                   // State
    Variables          map[string]interface{} `json:"variables"`           // variablen die mit dem Prozess mitlaufen
    EnabledTransitions []int                  `json:"enabled_transitions"` // list of transitions which can be fired
}

func (f *Net) Fire(transition int) error
    fires an enabled transition.

func (net *Net) Init()


