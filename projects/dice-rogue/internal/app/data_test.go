package app

import "testing"

func TestGameDataCatalogsLoadFromJSON(t *testing.T) {
	choices := characterChoices()
	if got, want := len(choices), 6; got != want {
		t.Fatalf("character count mismatch: got %d want %d", got, want)
	}

	encounter, ok := encounterByID("boss-1")
	if !ok {
		t.Fatalf("expected boss encounter from json catalog")
	}
	if got, want := encounter.Kind, EncounterKindBoss; got != want {
		t.Fatalf("boss encounter kind mismatch: got %q want %q", got, want)
	}

	start, ok := mapNodeByID("start")
	if !ok {
		t.Fatalf("expected start node from json catalog")
	}
	if got, want := len(start.NextIDs), 2; got != want {
		t.Fatalf("start next node count mismatch: got %d want %d", got, want)
	}
}

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
