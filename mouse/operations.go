package mouse

import (
	"math"
	"remy/config"

	"github.com/go-vgo/robotgo"
)

func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(x*(-1)))
}

func cap(x float64) float64 {
	DPI := float64(config.GetConfig().DPI)
	return sigmoid(x)*DPI - (DPI / 2)
}

func cutoff(x float64) float64 {
	sensitivity := float64(config.GetConfig().Sensitivity)
	if math.Abs(x) < sensitivity {
		return 0
	}

	return x
}

func MoveByDelta(x float64, y float64) {
	xDeltaPercentage := cap(cutoff(x))
	yDeltaPercentage := cap(cutoff(y))

	xMoveOffset, yMoveOffset := robotgo.MoveScale(int(xDeltaPercentage), int(yDeltaPercentage), 0)

	robotgo.MoveRelative(-xMoveOffset, -yMoveOffset)
}

func ClickButton(button MouseButton) {
	robotgo.Click(string(button))
}
