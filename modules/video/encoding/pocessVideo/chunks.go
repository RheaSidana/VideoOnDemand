package videoProcessing

import (
	"fmt"
	"strconv"
)

func createVideoChunks(
	loc, vidName, extension string,
	numberOf int,
) []string {
	var chunks []string
	for i := 1; i <= numberOf; i++ {
		chunk := fmt.Sprint(
			loc + "/" + vidName +
				strconv.Itoa(i) +
				extension,
		)
		chunks = append(
			chunks,
			chunk,
		)
	}
	return chunks
}
