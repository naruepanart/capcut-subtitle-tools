package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
)

func main() {
	data, err := os.ReadFile("subtitles.srt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	modifiedData := adjustSubtitleTiming(data, 10)

	if err = os.WriteFile("subtitles-mod.srt", modifiedData, 0644); err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("SRT file modified successfully")
}

func adjustSubtitleTiming(data []byte, gapSeconds int) []byte {
	var timeRegex = regexp.MustCompile(`(\d{2}:\d{2}:\d{2},\d{3}) --> (\d{2}:\d{2}:\d{2},\d{3})`)
	var buffer bytes.Buffer
	buffer.Grow(len(data) + len(data)/4)

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(data, len(data))

	var cumulativeGapMs int64
	gapMs := int64(gapSeconds) * 1000

	for scanner.Scan() {
		line := scanner.Bytes()

		if matches := timeRegex.FindSubmatch(line); matches != nil {
			if len(matches) != 3 {
				buffer.Write(line)
				buffer.WriteByte('\n')
				continue
			}
			startTime := matches[1]
			endTime := matches[2]

			newStartMs := parseTimeToMs(startTime) + cumulativeGapMs
			if newStartMs < 0 {
				newStartMs = 0
			}
			durationMs := parseTimeToMs(endTime) - parseTimeToMs(startTime)
			if durationMs < 0 {
				buffer.Write(line)
				buffer.WriteByte('\n')
				continue
			}
			newEndMs := newStartMs + durationMs

			writeTime(&buffer, newStartMs)
			buffer.WriteString(" --> ")
			writeTime(&buffer, newEndMs)
			buffer.WriteByte('\n')

			cumulativeGapMs += gapMs
		} else {
			buffer.Write(line)
			buffer.WriteByte('\n')
		}
	}

	return buffer.Bytes()
}

func parseTimeToMs(time []byte) int64 {
	hour := int64((time[0]-'0')*10 + (time[1] - '0'))
	minute := int64((time[3]-'0')*10 + (time[4] - '0'))
	second := int64((time[6]-'0')*10 + (time[7] - '0'))
	milli := int64((time[9]-'0')*100 + (time[10]-'0')*10 + (time[11] - '0'))
	return (hour*3600+minute*60+second)*1000 + milli
}

func writeTime(buf *bytes.Buffer, ms int64) {
	if ms < 0 {
		buf.WriteString("00:00:00,000")
		return
	}

	hours := ms / (3600 * 1000)
	ms %= 3600 * 1000
	minutes := ms / (60 * 1000)
	ms %= 60 * 1000
	seconds := ms / 1000
	millis := ms % 1000

	buf.WriteByte(byte('0' + hours/10))
	buf.WriteByte(byte('0' + hours%10))
	buf.WriteByte(':')
	buf.WriteByte(byte('0' + minutes/10))
	buf.WriteByte(byte('0' + minutes%10))
	buf.WriteByte(':')
	buf.WriteByte(byte('0' + seconds/10))
	buf.WriteByte(byte('0' + seconds%10))
	buf.WriteByte(',')
	buf.WriteByte(byte('0' + millis/100))
	buf.WriteByte(byte('0' + (millis/10)%10))
	buf.WriteByte(byte('0' + millis%10))
}
