package graphviz

import (
	"fmt"
	"io"
	"os"
)

type LinkMode int

const (
	Directed LinkMode = iota
	Undirected
)

func (mode LinkMode) Name() string {
	if mode == Directed {
		return "digraph"
	}
	return "graph"
}

func (mode LinkMode) String() string {
	if mode == Directed {
		return "->"
	}
	return "--"
}

type Entity struct {
	Name string
	Attr string
}

func NewEntity(name string, attr string) *Entity {
	return &Entity{Name: name, Attr: attr}
}

func (e *Entity) copyAttr(e2 Entity) {
	e.Attr = e2.Attr
}

type Target struct {
	Name string
	Attr string
}

func PlainTarget(name string) Target {
	return Target{Name: name}
}

type Graphviz struct {
	name string
	mode LinkMode

	entities map[string]*Entity
	pairs    map[string][]Target
}

func New(name string, mode LinkMode) *Graphviz {
	return &Graphviz{
		name:     name,
		mode:     mode,
		entities: make(map[string]*Entity),
		pairs:    make(map[string][]Target),
	}
}

func (g *Graphviz) tryInsertEntity(name string, entity *Entity) *Entity {
	origin, ok := g.entities[name]
	if ok {
		if entity != nil {
			origin.copyAttr(*entity)
		}
		return origin
	} else {
		if entity == nil {
			entity = new(Entity)
			entity.Name = name
		}
		g.entities[name] = entity
		return entity
	}
}

func (g *Graphviz) FindEntity(name string) *Entity {
	if e, ok := g.entities[name]; ok {
		return e
	}
	return nil
}

func (g *Graphviz) insertPair(from string, target Target) {
	pair := g.pairs[from]
	pair = append(pair, target)
	g.pairs[from] = pair
}

func (g *Graphviz) AddPlain(from, to, attr string) {
	g.tryInsertEntity(from, nil)
	g.tryInsertEntity(to, nil)
	g.insertPair(from, Target{Name: to, Attr: attr})
}

func (g *Graphviz) Add(from, to *Entity, attr string) {
	g.tryInsertEntity(from.Name, from)
	g.tryInsertEntity(to.Name, to)
	g.insertPair(from.Name, Target{to.Name, attr})
}

func (g *Graphviz) Write(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%s %s {\n", g.mode.Name(), g.name); err != nil {
		return err
	}
	for from, tos := range g.pairs {
		for _, to := range tos {
			if _, err := fmt.Fprintf(w, "\t%s %v %s%s;\n", from, g.mode, to.Name, to.Attr); err != nil {
				return err
			}
		}
	}
	for _, e := range g.entities {
		if e.Attr != "" {
			if _, err := fmt.Fprintf(w, "\t%s %s;\n", e.Name, e.Attr); err != nil {
				return err
			}
		}
	}
	if _, err := fmt.Fprintf(w, "}\n"); err != nil {
		return err
	}
	return nil
}

func (g *Graphviz) WriteFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return g.Write(file)
}
