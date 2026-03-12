package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/vector-ops/goships/types"
)

func SaveMapState(prefix string, m map[types.Position]types.Cell, unplacedShips []string) {
	_, err := os.Stat("logs")
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("logs", 0o755)
			if err != nil {
				WriteError(err)
				return
			}
		} else {
			return
		}
	}
	timestamp := time.Now().Unix()
	filepath := filepath.Join("logs", fmt.Sprintf("%s_map_%d.json", prefix, timestamp))
	file, err := os.Create(filepath)
	if err != nil {
		WriteError(err)
		return
	}
	defer file.Close()

	type SerializableCell struct {
		Type     string `json:"cell_type"`
		ShipType string `json:"ship_type"`
		Color    int16  `json:"cell_color"`
		Content  string `json:"cell_content"`
		Hit      bool   `json:"hit"`
	}

	type MapEntry struct {
		Key   types.Position   `json:"key"`
		Value SerializableCell `json:"value"`
	}

	state := make(map[string]interface{})
	state["unplacedShips"] = unplacedShips

	var entries []MapEntry
	for k, v := range m {
		entries = append(entries, MapEntry{Key: k, Value: SerializableCell{
			Type:     v.Type.String(),
			ShipType: v.ShipType.String(),
			Color:    v.Color,
			Content:  string(v.Content),
			Hit:      v.Hit,
		}})
	}

	state["map"] = entries

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(state)
	if err != nil {
		WriteError(err)
	}
}

func WriteTurns(turns map[string]int) {
	filepath := filepath.Join("logs", "turns.log")
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		WriteError(err)
		return
	}
	defer file.Close()

	for k, v := range turns {
		_, err = fmt.Fprintf(file, "%s: %d\n", k, v)
		if err != nil {
			WriteError(err)
			return
		}
	}
}

func WriteError(e error) {
	if e == nil {
		return
	}

	timestamp := time.Now().Unix()
	filepath := filepath.Join("logs", fmt.Sprintf("error_%d.log", timestamp))
	file, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	errorBytes := []byte(e.Error())
	_, err = file.Write(errorBytes)
	if err != nil {
		panic(err)
	}
}

func RemoveFilesByPattern(pattern string) error {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, file := range files {
		os.Remove(file)
	}
	return nil
}
