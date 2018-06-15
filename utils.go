package utils

import (
	"sort"
	"os"
	"strings"
	"strconv"
)

type ByNumericalFilename []os.FileInfo

func (nf ByNumericalFilename) Len() int      { return len(nf) }
func (nf ByNumericalFilename) Swap(i, j int) { nf[i], nf[j] = nf[j], nf[i] }
func (nf ByNumericalFilename) Less(i, j int) bool {

	pathA := nf[i].Name()
	pathB := nf[j].Name()

	subA := pathA[0:strings.LastIndex(pathA, ".")]
	subB := pathB[0:strings.LastIndex(pathB, ".")]

	if subA == subB {
		pathA_ := nf[i].Name()
		pathB_ := nf[j].Name()
		a, err1 := strconv.ParseInt(pathA_[strings.LastIndex(pathA_, ".")+1:len(pathA_)], 10, 64)
		b, err2 := strconv.ParseInt(pathB_[strings.LastIndex(pathB_, ".")+1:len(pathB_)], 10, 64)

		if err1 == nil && err2 == nil {
			return a < b
		}
	}
	return pathA < pathB
}

func SortFileInfo(files []os.FileInfo) {
	sort.Sort(ByNumericalFilename(files))
}
