package app

import "testing"

func TestHumanGuardUsesThreeDefenseDice(t *testing.T) {
	choices := characterChoices()

	var guard *UnitState
	for i := range choices {
		if choices[i].ID == "human-guard" {
			guard = &choices[i]
			break
		}
	}
	if guard == nil {
		t.Fatalf("expected human guard in character catalog")
	}

	defenseDice := 0
	skillDice := 0
	for _, die := range guard.Dice {
		switch die.Kind {
		case DieKindDefense:
			defenseDice++
		case DieKindSkill:
			skillDice++
		}
	}

	if got, want := defenseDice, 3; got != want {
		t.Fatalf("human guard defense dice mismatch: got %d want %d", got, want)
	}
	if got, want := skillDice, 1; got != want {
		t.Fatalf("human guard skill dice mismatch: got %d want %d", got, want)
	}
}
