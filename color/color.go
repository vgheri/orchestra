package color

import (
	"fmt"
	"math/rand"
	"runtime"
)

type Color string

var reset Color = "\033[0m"
var Red Color = "\033[31m"
var Green Color = "\033[32m"
var Yellow Color = "\033[33m"
var Blue Color = "\033[34m"
var Purple Color = "\033[35m"
var Cyan Color = "\033[36m"
var Gray Color = "\033[37m"
var LightGreen Color = "\033[92m"
var LightYellow Color = "\033[93m"
var LightBlue Color = "\033[94m"
var LightMagenta Color = "\033[95m"
var LightCyan Color = "\033[96m"
var White Color = "\033[97m"
var colors []Color

func init() {
	if runtime.GOOS == "windows" {
		reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
		LightGreen = ""
		LightYellow = ""
		LightBlue = ""
		LightMagenta = ""
		LightCyan = ""
	}
	colors = append(colors, Green, Yellow, Blue, Purple, Cyan, Gray, LightGreen, LightYellow, LightBlue, LightMagenta, LightCyan)
}

// RandomColor picks a random color from the list, and removes the color from it.
// Not thread safe, but not a problem as it's called from a for loop.
func RandomColor() Color {
	r := rand.Int31n(int32(len(colors)))
	color := colors[r]
	colors = append(colors[:r], colors[r+1:]...)
	return color
}

func Colorise(c Color, s string) string {
	return fmt.Sprintf("%s%s%s", c, s, reset)
}

func PrintFailure(s string) string {
	return fmt.Sprintf("%s%s%s", Red, s, reset)
}
