package letter

import (
	"sync"
)

// FreqMap records the frequency of each rune in a given text.
type FreqMap map[rune]int

// Frequency counts the frequency of each rune in a given text and returns this
// data as a FreqMap.
func Frequency(text string) FreqMap {
	frequencies := FreqMap{}
	for _, r := range text {
		frequencies[r]++
	}
	return frequencies
}

// ConcurrentFrequency counts the frequency of each rune in the given strings,
// by making use of concurrency.
func ConcurrentFrequency(texts []string) FreqMap {

	var wg sync.WaitGroup
	FreqMapChan := make(chan FreqMap)

	for _, text := range texts {
		wg.Add(1)
		go func(text string) {
			defer wg.Done()
			FM := Frequency(text)
			FreqMapChan <- FM
		}(text)

	}

	go func() {
		wg.Wait()
		close(FreqMapChan)
	}()

	finalFM := FreqMap{}
	for fm := range FreqMapChan {
		finalFM = finalFM.AggregateWithFM(fm)
	}

	return finalFM

}

func (mainFreqMap FreqMap) AggregateWithFM(FM_ToAdd FreqMap) FreqMap {
	for key, value := range FM_ToAdd {
		if _, exists := mainFreqMap[key]; exists {
			mainFreqMap[key] += value
		} else {
			mainFreqMap[key] = value
		}
	}
	return mainFreqMap
}

// TODO:
// 1. make a function to aggregate the results of different FreqMap
// 		it has to make sure to aggregate all the key value pairs, and make sure all the new keys are also added together
// 2. start the Frequency functions in parallel, using a WaitGroup and a channel to communicate and pass the FreqMap of each routine
// 3. then only return the final value once the wg is Done and the channel is closed 