package arrangers

import (
	"fmt"
	"math"
	"testing"
)

func TestDLossTarget(t *testing.T) {
	compareTargetDistGrad(0.5, 0.1, 0.2, 0.3, 0.3333, t)
	compareTargetDistGrad(0.5, 5, 12, 0.3, 0.3333, t)
	compareTargetDistGrad(0.5, -2, 0.2, 0.3, 0.3333, t)
	compareTargetDistGrad(0.5, 0.1, 0.2, -99, 14, t)
}

// Compare computed and estimated gradient
func compareTargetDistGrad(x, y, xx, yy, tgt float64, t *testing.T) {
	eps := 0.000001
	dx, dy, dxx, dyy := dlosstarget(x, y, xx, yy, tgt)
	ex, ey, exx, eyy := (losstarget(x+eps, y, xx, yy, tgt)-losstarget(x, y, xx, yy, tgt))/eps,
		(losstarget(x, y+eps, xx, yy, tgt)-losstarget(x, y, xx, yy, tgt))/eps,
		(losstarget(x, y, xx+eps, yy, tgt)-losstarget(x, y, xx, yy, tgt))/eps,
		(losstarget(x, y, xx, yy+eps, tgt)-losstarget(x, y, xx, yy, tgt))/eps
	ex, ey, exx, eyy = math.Abs(dx-ex), math.Abs(dy-ey), math.Abs(dxx-exx), math.Abs(dyy-eyy)
	if ex > eps*100 || ey > eps*100 || exx > eps*100 || eyy > eps*100 {
		fmt.Println("Errors between computed and estimated gradients is above ", eps*100)
		fmt.Println(ex, ey, exx, eyy)
		t.Fatal("Gradient error of dlosstarget is too large !")
	}
}
