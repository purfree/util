package math

import "math"

/* 保留小数位
 * @param n 小数位数，精度
 */
func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Round(f*n10) / n10
}
