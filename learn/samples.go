package learn

// TODO (abresk) write tests for samples

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const (
	hashSeperator = "#"
)

// Sample holds the sample data, value is just used for regression annotation
type Sample struct {
	Vector      []float64
	Output      []float64
	VectorHash  string
	OutputHash  string
	Label       string
	ClassNumber int
	Value       float64
}

func splitSamples(set *Set, ratio float64) (Set, Set) {
	normalizedRatio := int(ratio * 100.0)
	firstSet := Set{
		Samples:      make([]Sample, 0),
		ClassToLabel: set.ClassToLabel,
	}
	secondSet := Set{
		Samples:      make([]Sample, 0),
		ClassToLabel: set.ClassToLabel,
	}
	for i := range set.Samples {
		if rand.Intn(100) <= normalizedRatio {
			firstSet.Samples = append(firstSet.Samples, set.Samples[i])
		} else {
			secondSet.Samples = append(secondSet.Samples, set.Samples[i])
		}
	}
	return firstSet, secondSet
}

func (s *Sample) updateHashes() {
	text := ""
	for k, v := range s.Vector {
		text += fmt.Sprintf("%v:%v;", k, v)
	}
	s.VectorHash = calculateHash(text)
	text = ""
	for k, v := range s.Label {
		text += fmt.Sprintf("%v:%v;", k, v)
	}
	s.OutputHash = calculateHash(text)
}

func (s *Sample) getVectorHash() string {
	return s.VectorHash
}

func (s *Sample) getOutputHash() string {
	return s.OutputHash
}

func (s *Sample) getOverallHash() string {
	if s.OutputHash == "" || s.VectorHash == "" {
		s.updateHashes()
	}
	return s.VectorHash + hashSeperator + s.OutputHash
}

func calculateHash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func problemToMap(problem string) (map[int]float64, string, error) {
	sliced := strings.Split(problem, " ")
	m := make(map[int]float64)
	label := sliced[0]
	features := sliced[1:len(sliced)]
	for feature := range features {
		if features[feature] == "" {
			continue
		}
		splitted := strings.Split(features[feature], ":")
		idx, errIdx := strconv.Atoi(splitted[0])
		value, errVal := strconv.ParseFloat(splitted[1], 64)
		if errIdx == nil && errVal == nil {
			m[idx] = value
		}
	}
	return m, label, nil
}

// this function returns the highest index found
func scanSamples(path string) int {
	highest := 0
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("error while opening file")
		os.Exit(-1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m, _, err := problemToMap(scanner.Text())
		if err != nil {
			fmt.Printf("error while scanning files: %v", err)
			os.Exit(-1)
		}
		for k := range m {
			if k > highest {
				highest = k
			}
		}
	}
	return highest
}
