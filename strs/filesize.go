package strs

import (
	"fmt"
)

const bUnits = "EPTGMK"

// Size 将文件大小转换为人类可读的格式
func Size(bytes int64) string {
	for i, c := range bUnits {
		base := 1 << ((len(bUnits) - i) * 10)
		if float64(bytes) >= float64(base) {
			return fmt.Sprintf("%6.2f%cB", float64(bytes)/float64(base), c)
		}
	}
	return fmt.Sprintf("%dB", int(bytes))
}

func Percent(cur, total int64) string {
	return fmt.Sprintf("%.2f%%", float64(cur)/float64(total)*100)
}
