package util

import "log"

func Pretty(list [][]string) {
	if len(list) == 0 {
		return
	}

	rows, cols := len(list), len(list[0])
	lens := make([][]int, rows)
	for i := 0; i < rows; i++ {
		lens[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			lens[i][j] = len(list[i][j])
		}
	}

	data := make([]int, cols)
	for j := 0; j < cols; j++ {
		for i := 0; i < rows; i++ {
			if data[j] < lens[i][i] {
				data[j] = lens[i][j]
			}
		}
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			log.Print(list[i][j])
			pad := data[j] - lens[i][j] + 2
			for k := 0; k < pad; k++ {
				log.Print(" ")
			}
		}
		log.Print("\n")
	}
}
