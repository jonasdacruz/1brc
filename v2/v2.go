package v2

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"
)

type slice struct {
	start int64
	end   int64
}

type station struct {
	min, max, sum float64
	count         int
}

func ProcessFile(fp string) {
	defer func(t time.Time) {
		fmt.Println("process file taken: ", time.Since(t))
	}(time.Now())

	f, fs, err := openAndSize(fp)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	slices := calcFileSlices(f, int(fs), 10)

	_ = slices
	chsta := make(chan map[string]*station, len(slices))

	for _, s := range slices {
		go processSlice(f, s, chsta)
	}

	processStations(slices, chsta)
}

func openAndSize(fp string) (*os.File, int64, error) {
	defer func(t time.Time) {
		fmt.Println("open and get size taken: ", time.Since(t))
	}(time.Now())

	f, err := os.Open(fp)
	if err != nil {
		return nil, 0, err
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, 0, err
	}

	return f, fi.Size(), nil
}

func calcFileSlices(f *os.File, fs int, num int) []slice {
	defer func(t time.Time) {
		fmt.Println("calc file slices taken: ", time.Since(t))
	}(time.Now())

	ss := int64(fs / num)

	fmt.Printf("fs: %v, ss: %v\n", fs, ss)

	sliceBuffer := make([]byte, ss)
	offset := int64(0)
	nextOffset := offset

	slices := make([]slice, num)
	c := 0
	for {
		//		time.Sleep(1 * time.Second)
		_, err := f.ReadAt(sliceBuffer, offset)
		if err != nil {
			break
		}

		nextOffset = offset + int64(bytes.LastIndexByte(sliceBuffer, '\n')) + 1
		nextOffset = min(nextOffset, int64(fs))
		slices[c] = slice{start: offset, end: nextOffset}
		c++
		offset = nextOffset

		if nextOffset >= int64(fs) {
			break
		}
	}

	// fmt.Printf("slices: %v (%v)\n", slices, len(slices))
	return slices
}

func processSlice(f *os.File, s slice, chsta chan<- map[string]*station) {

	// load file slice in memory
	slice := make([]byte, s.end-s.start)
	_, err := f.ReadAt(slice, s.start)
	if err != nil {
		fmt.Println("Error reading file slice:", err)
		return
	}

	stations := make(map[string]*station)

	var el int
	lines := 0
	for {
		lines++
		el = bytes.IndexByte(slice, '\n')
		if el < 0 {
			break
		}
		stations = processLine(stations, slice[:el])
		slice = slice[el+1:]
	}

	chsta <- stations
}

func processLine(stations map[string]*station, line []byte) map[string]*station {
	b, e, _ := bytes.Cut(line, []byte{';'})

	sn := string(b)
	st, _ := strconv.ParseFloat(string(e), 64) // thanks to Alex

	tmpStation := make(map[string]*station)
	tmpStation[sn] = &station{min: st, max: st, sum: st, count: 1}

	return mergeStations(stations, tmpStation)
}

func mergeStations(stations map[string]*station, s map[string]*station) map[string]*station {
	for k, v := range s {
		if st, ok := stations[k]; ok {
			st.min = min(st.min, v.min)
			st.max = max(st.max, v.max)
			st.sum += v.sum
			st.count += v.count
		} else {
			stations[k] = v
		}
	}
	return stations
}

func processStations(slices []slice, chsta <-chan map[string]*station) {
	stations := make(map[string]*station)

	for i := 0; i < len(slices); i++ {
		stations = mergeStations(stations, <-chsta)
	}
}
