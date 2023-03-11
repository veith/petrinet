package petrinet_test

import (
	"github.com/veith/petrinet"
	"testing"

	"fmt"
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
		// ConditionMatrix: [][]string{{"true"}, {}, {}, {}, {}, {}, {}, {}},
		State:     []int{1, 0, 0, 0, 0, 0, 0, 0, 0},
		Variables: map[string]interface{}{"a": 4, "b": 8},
	}
	return f
}

func TestPetriNet_FireBadCondition(t *testing.T) {
	flow := makeExampleNet()
	flow.Variables = map[string]interface{}{"a": 2, "b": 8}
	flow.ConditionMatrix = [][]string{{"a > 11", "b == 8", "a != b", "true"}, {}, {}, {}, {}, {}, {}, {}}

	flow.Init()

	err := flow.Fire(0)

	if err == nil {
		t.Error("should not fire")
	}

	flow.UpdateVariable("a", 12)

	errr := flow.Fire(0)

	if errr != nil {
		t.Error("should have no error ")
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

func BenchmarkNet_Fire(b *testing.B) {
	flow := makeExampleNet()
	flow.ConditionMatrix = [][]string{{"a > 11", "b == 8", "a != b", "true"}, {}, {}, {}, {}, {}, {}, {}}
	flow.OutputMatrix[0][0] = 1
	flow.Init()

	for i := 0; i < b.N; i++ {
		flow.Fire(0)
	}
}
func BenchmarkNet_FireWithoutConditions(b *testing.B) {
	flow := makeExampleNet()
	flow.OutputMatrix[0][0] = 1
	flow.Init()

	for i := 0; i < b.N; i++ {
		flow.Fire(0)
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

func TestPetriNet_TokenID(t *testing.T) {
	flow := makeExampleNet()
	flow.State = []int{3, 0, 0, 0, 0, 0, 0, 0, 0}
	flow.Init()
	if flow.TokenIds[0][0] != 1 {
		t.Error("token id counter wrong ", flow.TokenIds)
	}

	fmt.Println(flow.State)
	err := flow.Fire(0)
	err = flow.Fire(2)
	err = flow.Fire(1)
	err = flow.Fire(3)
	err = flow.Fire(4)
	flow.Fire(0)
	flow.Fire(0)

	err = flow.Fire(2)
	err = flow.Fire(1)
	err = flow.Fire(3)
	err = flow.Fire(4)
	err = flow.Fire(2)
	err = flow.Fire(1)
	err = flow.Fire(3)
	err = flow.Fire(4)

	if flow.TokenIds[6][0] != 9 {
		t.Error("token id counter wrong ", flow.TokenIds)
	}

	if len(flow.TokenIds[0]) != 0 {
		t.Error("token id counter wrong should be empty haves", len(flow.TokenIds[0]))
	}

	if err != nil {
		t.Error(err)
	}

	flow.FireWithTokenId(5, 17)

	if flow.TokenIds[6][0] != 9 {
		t.Error("token id counter wrong ", flow.TokenIds)
	}

	if flow.TokenIds[7][0] != 22 {
		t.Error("token id counter wrong ", flow.TokenIds)
	}

	err = flow.FireWithTokenId(5, 2)

	if err == nil {
		t.Error("should return error tokenid not found")
	}
	err = flow.FireWithTokenId(5, 21)

	if flow.TokenIds[6][0] != 9 {
		t.Error("token id counter wrong ", flow.TokenIds)
	}
	flow.Fire(6)
	flow.Fire(6)
	flow.Fire(5)
	flow.Fire(6)

	if len(flow.EnabledTransitions) != 0 {
		t.Error("should have no enabled transitions, got ", len(flow.EnabledTransitions))
	}

}
