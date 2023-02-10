package utils

/* From https://www.arduino.cc/reference/en/language/functions/math/map/ */
func Map(x, originMin, originMax, targetMin, targetMax float64) float64 {
	return (x-originMin)*(targetMax-targetMin)/(originMax-originMin) + targetMin
}

func Map2(x, in_min, in_max, out_min, out_max float64) float64 {
	return (x-in_min)*(out_max-out_min)/(in_max-in_min) + out_min
}
