package har

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

func Load(filePaths []string, maxConcurrency int) (map[string]Har, error) {
	if maxConcurrency == 0 {
		maxConcurrency = 5
	}
	var wg sync.WaitGroup
	entriesSyncMap := &sync.Map{}

	fileChan := make(chan string, len(filePaths))

	for i := 0; i < maxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range fileChan {
				entry, err := readJSONFile(filePath)
				if err != nil {
					log.Printf("Error reading file %s: %v", filePath, err)
					continue
				}
				entriesSyncMap.Store(filePath, entry)
			}
		}()
	}

	for _, filePath := range filePaths {
		fileChan <- filePath
	}
	close(fileChan)

	wg.Wait()

	entries := make(map[string]Har)
	entriesSyncMap.Range(func(key, value interface{}) bool {
		entries[key.(string)] = value.(Har)
		return true
	})

	return entries, nil
}

func readJSONFile(filePath string) (Har, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return Har{}, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var h Har
	if err := json.Unmarshal(content, &h); err != nil {
		return Har{}, fmt.Errorf("failed to unmarshal JSON from file %s: %w", filePath, err)
	}

	return h, nil
}
