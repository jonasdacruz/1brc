package v1

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

type station struct {
	min, max, sum int32
	count         int64
}

var totals = make(map[string]station)

func customParseIntWithPoint(b []byte) int32 {
	var res int32

	if b[0] == '-' {
		b = b[1:]
		defer func() {
			res = -res
		}()
	}

	max := 5
	for _, c := range b {
		max--
		if max < 0 {
			break
		}
		if c == '.' {
			continue
		}
		res = res + int32(c-'0')
	}

	return res
}

func processTotals(sta []byte, tmp []byte) {
	temp := customParseIntWithPoint(tmp)
	s := totals[string(sta)]
	if s.count == 0 {
		totals[string(sta)] = station{
			min:   temp,
			max:   temp,
			sum:   temp,
			count: 1,
		}
	} else {
		s.min = min(s.min, temp)
		s.max = max(s.max, temp)
		s.sum += temp
		s.count++
	}
}

func processChunk(chunk []byte) {
	defer func(s time.Time) {
		fmt.Println("process a chunk taken: ", time.Since(s))
	}(time.Now())

	lines := 0
	for {
		newline := bytes.LastIndexByte(chunk, '\n')
		//fmt.Println("newline: ", newline)
		if newline < 0 {
			break
		}

		chunk = chunk[:newline+1]

		sta, tmp, sep := bytes.Cut(chunk, []byte(";"))
		// fmt.Printf("sta: %s, tmp: %s, sep: %v\n", sta, tmp, sep)
		if !sep {
			break
		}

		processTotals(sta, tmp)

		chunk = chunk[:len(chunk)-len(sta)-len(tmp)-1]
		lines++
	}

	fmt.Println("lines: ", lines)
}

func ProcessFile(filePath string) {
	// Open the file
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	fs := fi.Size()

	// Create a buffer to keep chunks that are read
	chunk := make([]byte, 0, 1024*1024*10)
	buf := make([]byte, 1024*10024*100)
	numChunks := 0

	for {
		// Read a chunk
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading file:", err)
			return
		}
		if n == 0 {
			break
		}

		chunk = append(chunk, buf[:n]...)
		processChunk(chunk)
		chunk = chunk[:0]

		numChunks++
		fmt.Printf("progress: %d%%\n", int64(len(buf)*numChunks)*100/fs)
	}

	fmt.Printf("totals: %+v\n", totals)
}
