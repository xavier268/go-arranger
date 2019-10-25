package arrangers

import (
	"fmt"
	"testing"
)

func TestDistancerGradient(t *testing.T) {
	eps := 0.000001
	x, y, xx, yy := 0.5, 0.1, 0.9, 0.2
	tgt := 0.33

	// Compare computed and estimated gradient
	dx, dy, dxx, dyy := dlosstarget(x, y, xx, yy, tgt)
	ex, ey, exx, eyy := estimateddlosstarget(x, y, xx, yy, tgt, eps)
	ex, ey, exx, eyy = dx-ex, dy-ey, dxx-exx, dyy-eyy
	if ex > eps*10 || ey > eps*10 || exx > eps*10 || eyy > eps*10 {
		fmt.Println("Errors between computed and estimated gradients")
		fmt.Println(ex, ey, exx, eyy, " epsilon = ", eps)
		t.Fatal("Gradient error is too large !")
	}
}

func estimateddlosstarget(x, y, xx, yy, tgt, eps float64) (dx, dy, dxx, dyy float64) {
	return (losstarget(x+eps, y, xx, yy, tgt) - losstarget(x, y, xx, yy, tgt)) / eps,
		(losstarget(x, y+eps, xx, yy, tgt) - losstarget(x, y, xx, yy, tgt)) / eps,
		(losstarget(x, y, xx+eps, yy, tgt) - losstarget(x, y, xx, yy, tgt)) / eps,
		(losstarget(x, y, xx, yy+eps, tgt) - losstarget(x, y, xx, yy, tgt)) / eps

}
