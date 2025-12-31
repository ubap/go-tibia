package assets

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadItemsJson reads the JSON file and populates the global 'things' slice.
func LoadItemsJson(path string) error {
	// Read the file from disk
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to open items file: %w", err)
	}

	// Parse JSON into a temporary slice
	var loadedItems []ItemType
	if err := json.Unmarshal(fileBytes, &loadedItems); err != nil {
		return fmt.Errorf("failed to parse items json: %w", err)
	}

	// initialize the Global Registry (Allocate memory)
	// We add +1 because IDs are 0-indexed in the slice.
	maxId := getMaxId(loadedItems)
	initialize(maxId + 1)

	// 5. Populate the Global Registry
	count := 0
	for _, item := range loadedItems {
		things[item.ID] = item
		count++
	}

	fmt.Printf("[Data] Loaded %d items from %s. Max ID: %d\n", count, path, maxId)
	return nil
}

func getMaxId(items []ItemType) int {
	maxID := 0
	for _, item := range items {
		if int(item.ID) > maxID {
			maxID = int(item.ID)
		}
	}
	return maxID
}
