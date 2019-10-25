package arrangers

// Distancer rearrange a Grapher to ensure distance
// between linked nodes are as close as possible to target,
// and distance between non linked nodes are 'sufficient'.
// Graph should be Normalized first.
// Target distances are typically between 0 and 0.5.
type Distancer struct {
	dl float64 // target distance for linked nodes
}

// Loss for staying inside the target perimeter
func lossCenter(x, y float64) float64 {
	r := (x - 0.5) * (y - 0.5)
	return r * r
}

// corresponding gradient
func dlossCenter(x, y float64) (dx, dy float64) {
	r := (x - 0.5) * (y - 0.5)
	return 2 * r * (y - 0.5), 2 * r * (x - 0.5)
}

// loss to provide attraction between two nodes
func lossattract(x, y, xx, yy float64) float64 {
	return (x-xx)*(x-xx) + (y-yy)*(y-yy)
}

// gradient for loassattarct
func dlossattract(x, y, xx, yy float64) (dx, dy, dxx, dyy float64) {
	return 2 * (x - xx), 2 * (y - yy), 2 * (xx - x), 2 * (yy - y)
}

// loss for a fixed distance, tgt, typically << 1
func losstarget(x, y, xx, yy, tgt float64) float64 {
	d := lossattract(x, y, xx, yy) / tgt
	return d + 1/d
}

// gradient for losstarget
func dlosstarget(x, y, xx, yy, tgt float64) (dx, dy, dxx, dyy float64) {
	d := lossattract(x, y, xx, yy)
	d = (1/tgt - tgt/(d*d))
	dx, dy, dxx, dyy = dlossattract(x, y, xx, yy)
	return d * dx, d * dy, d * dxx, d * dyy
}
