package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type JSONSubtitle map[string]interface{}

func main() {
	// Read the file path from file-path.txt
	pathData, err := os.ReadFile("file-path.txt")
	if err != nil {
		fmt.Println("Error reading file-path.txt:", err)
		return
	}

	// Clean the path by trimming whitespace and quotes if any
	filePath := strings.TrimSpace(string(pathData))

	// Read the JSON file from the specified path
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	modifiedData, err := adjustSubtitleTiming(data, 10) // Gap of 10 seconds
	if err != nil {
		fmt.Println("Error adjusting timing:", err)
		return
	}

	// Write back to the same file path
	if err = os.WriteFile(filePath, modifiedData, 0644); err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("JSON file modified successfully")
}

// New helper function to convert seconds to microseconds
func secondsToMicroseconds(seconds int64) int64 {
	return seconds * 1000000
}

func adjustSubtitleTiming(data []byte, gapSeconds int64) ([]byte, error) {
	// Convert seconds to microseconds
	gapMicroseconds := secondsToMicroseconds(gapSeconds)

	var subtitle JSONSubtitle
	if err := json.Unmarshal(data, &subtitle); err != nil {
		return nil, err
	}

	tracks, ok := subtitle["tracks"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid tracks format")
	}

	for _, track := range tracks {
		trackMap, ok := track.(map[string]interface{})
		if !ok || trackMap["type"] != "text" {
			continue
		}
		segments, ok := trackMap["segments"].([]interface{})
		if !ok || len(segments) < 1 {
			continue
		}

		// Adjust start times for segments starting from the second one
		for i := 1; i < len(segments); i++ {
			segmentMap, ok := segments[i].(map[string]interface{})
			if !ok {
				continue
			}
			targetTimerange, ok := segmentMap["target_timerange"].(map[string]interface{})
			if !ok {
				continue
			}

			start, ok := targetTimerange["start"].(float64)
			if !ok {
				continue
			}

			// Add the gap (in microseconds) to the start time, scaled by index (starting from 1)
			newStart := int64(start) + gapMicroseconds*int64(i)
			if newStart < 0 {
				newStart = 0
			}
			targetTimerange["start"] = float64(newStart)

			// Update segment to ensure changes persist
			segments[i] = segmentMap
		}

		// Update track segments
		trackMap["segments"] = segments
	}

	// Compact JSON output
	modifiedData, err := json.Marshal(subtitle)
	if err != nil {
		return nil, err
	}

	return modifiedData, nil
}
