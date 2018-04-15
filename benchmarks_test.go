package tizzy

import (
	"archive/zip"
	"log"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var locationNames = loadLocationNames()

func BenchmarkLoadLocationValue(b *testing.B) {
	sample := prepareSample(b.N, locationNames)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LoadLocationValue(sample[i])
	}
}

func BenchmarkLoadLocation(b *testing.B) {
	sample := prepareSample(b.N, locationNames)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LoadLocation(sample[i])
	}
}

func BenchmarkTimeLoadLocation(b *testing.B) {
	sample := prepareSample(b.N, locationNames)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		time.LoadLocation(sample[i])
	}
}

func loadLocationNames() []string {
	zoneinfo := runtime.GOROOT() + "/lib/time/zoneinfo.zip"

	r, err := zip.OpenReader(zoneinfo)
	if err != nil {
		log.Fatalf("Could not open zoneinfo.zip for reading: path '%s' error: '%v'", zoneinfo, err)
	}
	defer r.Close()

	var locs = []string{}

	for _, f := range r.File {
		if f == nil {
			log.Fatalf("zoneinfo.zip contained nil file pointer")
		}

		if f.FileHeader.FileInfo().IsDir() {
			continue
		}

		locs = append(locs, f.Name)
	}
	return locs
}

func prepareSample(n int, data []string) []string {
	dn := len(data)
	sample := make([]string, n)
	for i := 0; i < n; i++ {
		sample[i] = data[rand.Intn(dn)]
	}

	return sample
}
