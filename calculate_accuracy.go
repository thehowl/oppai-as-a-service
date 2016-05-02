package main

// taken from ripple-cron-go
func calculateAccuracy(count300, count100, count50, countgeki, countkatu, countmiss, playMode int) float64 {
	var accuracy float64
	switch playMode {
	case 1:
		// Please note this is not what is written on the wiki.
		// However, what was written on the wiki didn't make any sense at all.
		totalPoints := (count100*50 + count300*100)
		maxHits := (countmiss + count100 + count300)
		accuracy = float64(totalPoints) / float64(maxHits*100)
	case 2:
		fruits := count300 + count100 + count50
		totalFruits := fruits + countmiss + countkatu
		accuracy = float64(fruits) / float64(totalFruits)
	case 3:
		totalPoints := (count50*50 + count100*100 + countkatu*200 + count300*300 + countgeki*300)
		maxHits := (countmiss + count50 + count100 + count300 + countgeki + countkatu)
		accuracy = float64(totalPoints) / float64(maxHits*300)
	default:
		totalPoints := (count50*50 + count100*100 + count300*300)
		maxHits := (countmiss + count50 + count100 + count300)
		accuracy = float64(totalPoints) / float64(maxHits*300)
	}
	return accuracy * 100
}
