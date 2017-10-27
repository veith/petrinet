package petrinet_test

import (
	"testing"
	"github.com/veith/petrinet"
)

func makeExampleNet() petrinet.Net {
	f := petrinet.Net{
		InputMatrix: [][]int{
			{1, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 1, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 1, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 1, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 1, 0},
		},
		OutputMatrix: [][]int{
			{0, 1, 1, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 1, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 1, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 1, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 1, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 1},
		},

		State:           []int{1, 0, 0, 0, 0, 0, 0, 0, 0},
		Variables:       map[string]interface{}{"a": 4, "b": 8},
	}
	return f
}

func TestPetriNet_FireBadCondition(t *testing.T) {
	flow := makeExampleNet()
	flow.Variables = map[string]interface{}{"a": 2, "b": 8}
	flow.ConditionMatrix =  [][]string{{"a > 11", "b == 8", "a != b", "true"}, {}, {}, {}, {}, {}, {}, {}}

	flow.Init()
	err := flow.Fire(0)

	if err == nil {
		t.Error("should not fire")
	}

}

func TestPetriNet_WithoutConditionsMatrix(t *testing.T) {
	flow := makeExampleNet()

	flow.Init()
	err := flow.Fire(0)

	if err != nil {
		t.Error(err)
	}
	// example should have 2 possible transitions
	if len(flow.EnabledTransitions) != 2 {
		t.Error("Expected 2, got ", len(flow.EnabledTransitions))
	}

}

func TestPetriNet_Fire(t *testing.T) {
	flow := makeExampleNet()
	flow.ConditionMatrix = [][]string{{"true"}, {}, {}, {}, {}, {}, {}, {}}
	flow.Init()
	err := flow.Fire(0)

	if err != nil {
		t.Error(err)
	}
	// example should have 2 possible transitions
	if len(flow.EnabledTransitions) != 2 {
		t.Error("Expected 2, got ", len(flow.EnabledTransitions))
	}

}
