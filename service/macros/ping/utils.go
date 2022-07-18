package ping

import "sort"

func computeAvgOfPing(pings []uint16) uint16 {
	result := uint16(0)
	totalMS := pings[:]
	sort.Slice(totalMS, func(i, j int) bool { return totalMS[i] < totalMS[j] })
	mediumMS := totalMS[len(totalMS)/2]
	threshold := 300
	realCount := uint16(0)
	for _, delay := range totalMS {
		if -threshold < int(delay)-int(mediumMS) && int(delay)-int(mediumMS) < threshold {
			realCount += 1
		}
	}
	for _, delay := range totalMS {
		if -threshold < int(delay)-int(mediumMS) && int(delay)-int(mediumMS) < threshold {
			result += delay / realCount
		}
	}
	return result
}
