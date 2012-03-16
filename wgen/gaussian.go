package main

import "math"

// Gaussian2d is a 2-dimensional normal distribution.
type Gaussian2d struct{
	Mx, My float64 	// mean
	Sx, Sy float64	// standard deviation
	Cov float64	// covariance

	// pre-computed values for the PDF
	covcov, cov2, sxsy, sxsx, sysy, coeff, expcoeff float64
}

// NewGaussian2d returns a new 2-dimensional Gaussian2d with
// the given mean, standard deviation, covariance
//
// A bunch of values used during PDF computation are
// pre-computed for efficiency, so if any of the exported
// fields are changed then Precompute() must be called
// to re-compute these fields.
func NewGaussian2d(mx, my, sx, sy, cov float64) *Gaussian2d {
	g := &Gaussian2d{
		Mx: mx,
		My: my,
		Sx: sx,
		Sy: sy,
		Cov: cov,
	}
	g.Precompute()
	return g
}

// precompute pre-computes a bunch of information
// that is used during PDF computation.  It must
// be called if any of the exported fields of
// the gaussion are changed after NewGaussian2d
// is called.
func (g *Gaussian2d) Precompute() {
	sxsy := g.Sx*g.Sy
	covcov := g.Cov*g.Cov
	g.covcov = covcov
	g.cov2 = g.Cov*2
	g.sxsy = sxsy
	g.sxsx = g.Sx*g.Sx
	g.sysy = g.Sy*g.Sy
	g.coeff = 1 / (2 * math.Pi * sxsy * math.Sqrt(1-covcov))
	g.expcoeff = -0.5 * (1-covcov)
}

// PDF returns the probability density.
func (g *Gaussian2d) PDF(x, y float64) float64 {
	devx := x - g.Mx
	devy := y - g.My
	vl := devx*devx / g.sxsx
	vl += devy*devy / g.sysy
	vl -= g.cov2 * devx * devy / g.sxsy
	return g.coeff * math.Exp(g.expcoeff*vl)
}
