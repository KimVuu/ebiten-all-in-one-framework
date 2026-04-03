package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	projectdata "github.com/kimyechan/ebiten-aio-framework/projects/dice-rogue/data"
)

var (
	gameDataOnce sync.Once
	gameDataErr  error

	characterCatalog []UnitState
	encounterCatalog map[string]EncounterDefinition
	actMapNodes      map[string]EncounterNode
)

type characterCatalogFile struct {
	Characters []characterData `json:"characters"`
}

type encounterCatalogFile struct {
	Encounters []encounterData `json:"encounters"`
}

type mapNodeFile struct {
	Nodes []mapNodeData `json:"nodes"`
}

type characterData struct {
	ID    string    `json:"id"`
	Name  string    `json:"name"`
	Role  string    `json:"role"`
	MaxHP int       `json:"max_hp"`
	Dice  []dieData `json:"dice"`
}

type encounterData struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Kind    string      `json:"kind"`
	Enemies []enemyData `json:"enemies"`
}

type enemyData struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	MaxHP    int           `json:"max_hp"`
	Patterns []patternData `json:"patterns"`
}

type mapNodeData struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Kind        string   `json:"kind"`
	EncounterID string   `json:"encounter_id"`
	NextIDs     []string `json:"next_ids"`
}

type dieData struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Kind     string     `json:"kind"`
	EffectID string     `json:"effect_id,omitempty"`
	Faces    []faceData `json:"faces"`
}

type faceData struct {
	Kind  string `json:"kind"`
	Label string `json:"label"`
	Value int    `json:"value,omitempty"`
}

type patternData struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Attacks []int  `json:"attacks,omitempty"`
	Defense int    `json:"defense,omitempty"`
}

func init() {
	mustLoadGameData()
}

func mustLoadGameData() {
	gameDataOnce.Do(func() {
		gameDataErr = loadGameData()
	})
	if gameDataErr != nil {
		panic(gameDataErr)
	}
}

func loadGameData() error {
	characters, err := loadCharacterCatalog()
	if err != nil {
		return err
	}
	encounters, err := loadEncounterCatalog()
	if err != nil {
		return err
	}
	nodes, err := loadMapNodes()
	if err != nil {
		return err
	}
	start, ok := nodes["start"]
	if !ok || len(start.NextIDs) == 0 {
		return fmt.Errorf("load game data: start node is missing or has no exits")
	}

	characterCatalog = characters
	encounterCatalog = encounters
	actMapNodes = nodes
	return nil
}

func loadCharacterCatalog() ([]UnitState, error) {
	var file characterCatalogFile
	if err := loadJSONDataFile("characters.json", &file); err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	units := make([]UnitState, 0, len(file.Characters))
	for _, entry := range file.Characters {
		if entry.ID == "" {
			return nil, fmt.Errorf("load characters.json: character id is required")
		}
		if _, exists := seen[entry.ID]; exists {
			return nil, fmt.Errorf("load characters.json: duplicate character id %q", entry.ID)
		}
		seen[entry.ID] = struct{}{}
		if entry.MaxHP <= 0 {
			return nil, fmt.Errorf("load characters.json: character %q max_hp must be positive", entry.ID)
		}
		unit := UnitState{
			ID:       entry.ID,
			Name:     entry.Name,
			Role:     entry.Role,
			MaxHP:    entry.MaxHP,
			HP:       entry.MaxHP,
			Dice:     make([]DieSpec, 0, len(entry.Dice)),
			Counters: map[string]int{},
		}
		for _, dieEntry := range entry.Dice {
			die, err := toDieSpec(entry.ID, dieEntry)
			if err != nil {
				return nil, fmt.Errorf("load characters.json: %w", err)
			}
			unit.Dice = append(unit.Dice, die)
		}
		units = append(units, unit)
	}
	return units, nil
}

func loadEncounterCatalog() (map[string]EncounterDefinition, error) {
	var file encounterCatalogFile
	if err := loadJSONDataFile("encounters.json", &file); err != nil {
		return nil, err
	}
	catalog := make(map[string]EncounterDefinition, len(file.Encounters))
	for _, entry := range file.Encounters {
		if entry.ID == "" {
			return nil, fmt.Errorf("load encounters.json: encounter id is required")
		}
		kind, err := parseEncounterKind(entry.Kind)
		if err != nil {
			return nil, fmt.Errorf("load encounters.json: encounter %q: %w", entry.ID, err)
		}
		enemies := make([]UnitState, 0, len(entry.Enemies))
		for _, enemyEntry := range entry.Enemies {
			if enemyEntry.ID == "" {
				return nil, fmt.Errorf("load encounters.json: encounter %q has enemy without id", entry.ID)
			}
			if enemyEntry.MaxHP <= 0 {
				return nil, fmt.Errorf("load encounters.json: enemy %q max_hp must be positive", enemyEntry.ID)
			}
			patterns := make([]EncounterPattern, 0, len(enemyEntry.Patterns))
			for _, patternEntry := range enemyEntry.Patterns {
				patterns = append(patterns, EncounterPattern{
					ID:      patternEntry.ID,
					Label:   patternEntry.Label,
					Attacks: append([]int(nil), patternEntry.Attacks...),
					Defense: patternEntry.Defense,
				})
			}
			enemies = append(enemies, UnitState{
				ID:       enemyEntry.ID,
				Name:     enemyEntry.Name,
				Role:     "적",
				MaxHP:    enemyEntry.MaxHP,
				HP:       enemyEntry.MaxHP,
				Patterns: patterns,
				Counters: map[string]int{},
			})
		}
		catalog[entry.ID] = EncounterDefinition{
			ID:      entry.ID,
			Name:    entry.Name,
			Kind:    kind,
			Enemies: enemies,
		}
	}
	return catalog, nil
}

func loadMapNodes() (map[string]EncounterNode, error) {
	var file mapNodeFile
	if err := loadJSONDataFile("map_nodes.json", &file); err != nil {
		return nil, err
	}
	nodes := make(map[string]EncounterNode, len(file.Nodes))
	for _, entry := range file.Nodes {
		if entry.ID == "" {
			return nil, fmt.Errorf("load map_nodes.json: node id is required")
		}
		kind, err := parseNodeKind(entry.Kind)
		if err != nil {
			return nil, fmt.Errorf("load map_nodes.json: node %q: %w", entry.ID, err)
		}
		nodes[entry.ID] = EncounterNode{
			ID:          entry.ID,
			Name:        entry.Name,
			Kind:        kind,
			EncounterID: entry.EncounterID,
			NextIDs:     append([]string(nil), entry.NextIDs...),
		}
	}
	return nodes, nil
}

func loadJSONDataFile(name string, target any) error {
	data, err := readGameDataFile(name)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("parse %s: %w", name, err)
	}
	return nil
}

func readGameDataFile(name string) ([]byte, error) {
	candidates := candidateDataPaths(name)
	var readErr error
	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err == nil {
			return data, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			readErr = err
		}
	}
	data, err := fs.ReadFile(projectdata.FS, name)
	if err == nil {
		return data, nil
	}
	if readErr != nil {
		return nil, fmt.Errorf("read %s: %w", name, readErr)
	}
	return nil, fmt.Errorf("read %s: file not found in %v and embedded fallback failed: %w", name, candidates, err)
}

func candidateDataPaths(name string) []string {
	seen := map[string]struct{}{}
	paths := make([]string, 0, 12)
	addCandidates := func(base string) {
		if base == "" {
			return
		}
		for _, dir := range walkUpDirs(base, 4) {
			path := filepath.Clean(filepath.Join(dir, "data", name))
			if _, exists := seen[path]; exists {
				continue
			}
			seen[path] = struct{}{}
			paths = append(paths, path)
		}
	}

	if cwd, err := os.Getwd(); err == nil {
		addCandidates(cwd)
	}
	if executable, err := os.Executable(); err == nil {
		addCandidates(filepath.Dir(executable))
	}
	return paths
}

func walkUpDirs(start string, levels int) []string {
	dirs := make([]string, 0, levels+1)
	current := filepath.Clean(start)
	for i := 0; i <= levels; i++ {
		dirs = append(dirs, current)
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return dirs
}

func toDieSpec(ownerID string, entry dieData) (DieSpec, error) {
	kind, err := parseDieKind(entry.Kind)
	if err != nil {
		return DieSpec{}, fmt.Errorf("character %q die %q: %w", ownerID, entry.ID, err)
	}
	faces := make([]DieFace, 0, len(entry.Faces))
	for _, faceEntry := range entry.Faces {
		faceKind, err := parseFaceKind(faceEntry.Kind)
		if err != nil {
			return DieSpec{}, fmt.Errorf("character %q die %q: %w", ownerID, entry.ID, err)
		}
		faces = append(faces, DieFace{
			Label: faceEntry.Label,
			Kind:  faceKind,
			Value: faceEntry.Value,
		})
	}
	return DieSpec{
		ID:              entry.ID,
		OwnerID:         ownerID,
		Name:            entry.Name,
		Kind:            kind,
		Faces:           faces,
		EffectID:        entry.EffectID,
		EnabledInBattle: true,
	}, nil
}

func parseDieKind(raw string) (DieKind, error) {
	kind := DieKind(raw)
	switch kind {
	case DieKindAttack, DieKindDefense, DieKindSkill:
		return kind, nil
	default:
		return "", fmt.Errorf("unknown die kind %q", raw)
	}
}

func parseFaceKind(raw string) (FaceKind, error) {
	kind := FaceKind(raw)
	switch kind {
	case FaceKindValue, FaceKindSuccess, FaceKindFailure, FaceKindCritical, FaceKindEscape:
		return kind, nil
	default:
		return "", fmt.Errorf("unknown face kind %q", raw)
	}
}

func parseEncounterKind(raw string) (EncounterKind, error) {
	kind := EncounterKind(raw)
	switch kind {
	case EncounterKindNormal, EncounterKindElite, EncounterKindBoss:
		return kind, nil
	default:
		return "", fmt.Errorf("unknown encounter kind %q", raw)
	}
}

func parseNodeKind(raw string) (NodeKind, error) {
	kind := NodeKind(raw)
	switch kind {
	case NodeKindStart, NodeKindNormal, NodeKindElite, NodeKindBoss, NodeKindRest:
		return kind, nil
	default:
		return "", fmt.Errorf("unknown node kind %q", raw)
	}
}

func newCharacterState(id string) (UnitState, bool) {
	mustLoadGameData()
	for _, unit := range characterCatalog {
		if unit.ID == id {
			return cloneUnits([]UnitState{unit})[0], true
		}
	}
	return UnitState{}, false
}

func characterChoices() []UnitState {
	mustLoadGameData()
	return cloneUnits(characterCatalog)
}

func encounterByID(id string) (EncounterDefinition, bool) {
	mustLoadGameData()
	encounter, ok := encounterCatalog[id]
	if !ok {
		return EncounterDefinition{}, false
	}
	encounter.Enemies = cloneUnits(encounter.Enemies)
	return encounter, true
}

func mapNodeByID(id string) (EncounterNode, bool) {
	mustLoadGameData()
	node, ok := actMapNodes[id]
	return node, ok
}
