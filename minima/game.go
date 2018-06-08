// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"bufio"
	"encoding/json"
	"io"
	"math/rand"

	"github.com/mccoyst/min-game/ai"
	"github.com/mccoyst/min-game/animal"
	"github.com/mccoyst/min-game/geom"
	"github.com/mccoyst/min-game/item"
	"github.com/mccoyst/min-game/ui"
	"github.com/mccoyst/min-game/world"
)

var TileSize = world.TileSize

type Game struct {
	wo         *world.World
	cam        ui.Camera
	base       Base
	Astro      *Player
	Herbivores []animal.Herbivores
	Treasure   []item.Treasure
}

// ReadGame returns a *Game, read from the given
// reader.
func ReadGame(r io.Reader) (*Game, error) {
	g := new(Game)
	in := bufio.NewReader(r)
	var err error
	if g.wo, err = world.Read(in); err != nil {
		return nil, err
	}
	g.cam = ui.Camera{Torus: g.wo.Pixels, Dims: ScreenDims}
	crashSite := geom.Pt(float64(g.wo.X0), float64(g.wo.Y0)).Mul(TileSize)
	g.Astro = NewPlayer(g.wo, crashSite)
	g.base = NewBase(crashSite)

	if err := json.NewDecoder(in).Decode(&g); err != nil {
		panic(err)
	}
	g.CenterOnTile(g.wo.Tile(g.Astro.body.Center()))
	return g, nil
}

func (*Game) Transparent() bool {
	return false
}

// CenterOnTile centers the display on a given tile.
func (g *Game) CenterOnTile(x, y int) {
	pt := geom.Pt(float64(x), float64(y))
	halfTile := TileSize.Div(geom.Pt(2, 2))
	g.cam.Center(pt.Mul(TileSize).Add(halfTile))
}

func (g *Game) Draw(d ui.Drawer) {
	pt := ScreenDims.Div(TileSize)
	w, h := int(pt.X), int(pt.Y)
	x0, y0 := g.wo.Tile(g.cam.Pt)

	for x := x0; x <= x0+w; x++ {
		for y := y0; y <= y0+h; y++ {
			l := g.wo.At(x, y)
			g.cam.Draw(d, ui.Sprite{
				Name:   l.Terrain.Name,
				Bounds: geom.Rectangle{geom.Pt(0, 0), TileSize},
				Shade:  shade(l),
			}, geom.Pt(float64(x), float64(y)).Mul(TileSize))
		}
	}

	g.base.Draw(d, g.cam)
	for _, t := range g.Treasure {
		if t.Item == nil {
			continue
		}
		g.cam.Draw(d, ui.Sprite{
			Name:   "Present",
			Bounds: geom.Rect(0, 0, t.Box.Dx(), t.Box.Dy()),
			Shade:  shade(g.wo.At(g.wo.Tile(t.Box.Center()))),
		}, t.Box.Min)
	}
	g.Astro.Draw(d, g.cam)
	for i := range g.Herbivores {
		g.Herbivores[i].Draw(d, g.cam)
	}

	g.Astro.drawO2(d)

	if !*debug {
		return
	}
	d.SetFont(DialogFont, 8)
	d.SetColor(White)
	sz := d.TextSize(g.Astro.info)
	d.Draw(g.Astro.info, geom.Pt(0, ScreenDims.Y-sz.Y))
}

// Shade returns the shade value for a location.
func shade(l *world.Loc) float32 {
	const minSh = 0.15
	const slope = (1 - minSh) / world.MaxElevation
	return slope*float32(l.Elevation-l.Depth) + minSh
}

func (g *Game) Handle(stk *ui.ScreenStack, ev ui.Event) error {
	k, ok := ev.(ui.Key)
	if !ok || !k.Down {
		return nil
	}

	switch {
	case k.Button == ui.Menu:
		stk.Push(NewPauseScreen(g.Astro))

	case k.Button == ui.Action:
		it, box := g.GetTreasure(g.Astro.body.Box)
		if it == nil && g.wo.Pixels.Overlaps(g.Astro.body.Box, g.base.Box) {
			stk.Push(NewBaseScreen(g.Astro, &g.base))
			break
		} else if it == nil {
			break
		}
		scr := NewNormalMessage("Bravo! You got the " + it.Name + "!")
		if !g.Astro.PutPack(it) {
			scr = NewNormalMessage("You don't have room for that in your pack.")
			g.Treasure = append(g.Treasure, item.Treasure{it, box})
		}
		stk.Push(scr)

	case k.Button == ui.Hands && g.Astro.Held != nil:
		pt := g.Astro.HeldLoc()
		box := geom.Rectangle{pt, pt.Add(TileSize)}
		g.Treasure = append(g.Treasure, item.Treasure{g.Astro.Held, box})
		g.Astro.Held = nil

	case k.Button == ui.Hands:
		it, _ := g.GetTreasure(g.Astro.body.Box)
		if it == nil {
			break
		}
		scr := NewNormalMessage("Ahh, you decided to hold onto the " + it.Name + "!")
		if it.Name == item.Scrap {
			g.Astro.Scrap++
			scr = NewNormalMessage("Bravo! You got the " + it.Name + "!")
		} else {
			g.Astro.Held = it
		}
		stk.Push(scr)
	}
	return nil
}

func (g *Game) GetTreasure(b geom.Rectangle) (*item.Item, geom.Rectangle) {
	for i, t := range g.Treasure {
		if !g.wo.Pixels.Overlaps(b, t.Box) {
			continue
		}
		g.Treasure[i] = g.Treasure[len(g.Treasure)-1]
		g.Treasure = g.Treasure[:len(g.Treasure)-1]
		return t.Item, t.Box
	}
	return nil, geom.Rectangle{}
}

func (g *Game) Update(stk *ui.ScreenStack) error {
	const speed = 4 // px

	if g.Astro.o2 == 0 && !*debug {
		if et := g.Astro.FindEtele(); et == nil {
			stk.Push(NewGameOverScreen())
		} else {
			et.Uses--
			g.Astro.body.Vel = geom.Pt(0, 0)
			dims := geom.Pt(g.Astro.body.Box.Dx(), g.Astro.body.Box.Dy())
			g.Astro.body.Box.Min = g.base.Box.Min
			g.Astro.body.Box.Max = g.base.Box.Min.Add(dims)
			g.Astro.RefillO2()
		}
	}

	g.Astro.body.Vel = geom.Pt(0, 0)
	if stk.Buttons&ui.Left != 0 {
		g.Astro.body.Vel.X -= speed
	}
	if stk.Buttons&ui.Right != 0 {
		g.Astro.body.Vel.X += speed
	}
	if stk.Buttons&ui.Down != 0 {
		g.Astro.body.Vel.Y += speed
	}
	if stk.Buttons&ui.Up != 0 {
		g.Astro.body.Vel.Y -= speed
	}
	g.Astro.Move(g.wo)
	g.cam.Center(g.Astro.body.Box.Center())

	for i := range g.Herbivores {
		ai.UpdateBoids(stk.NFrames, g.Herbivores[i], &g.Astro.body, g.wo)
		g.Herbivores[i].Move(g.wo)
	}

	return nil
}

func randPoint(xmax, ymax float64) geom.Point {
	return geom.Pt(rand.Float64()*xmax, rand.Float64()*ymax)
}
