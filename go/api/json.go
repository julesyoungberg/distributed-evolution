package api

import (
	"encoding/json"

	"github.com/MaxHalford/eaopt"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

/**
 * Since Shapes contains pointers, custom json marshalling is needed.
 * This file contains the JSON-encodable versions of Shapes and eaopt.Individual,
 * and handles the logic associated with them.
 */
type JSONShapes struct {
	Bounds util.Vector `json:"bounds"`
	Type   string      `json:"type"`

	Circles   []Circle   `json:"circles"`
	Polygons  []Polygon  `json:"polygon"`
	Triangles []Triangle `json:"triangle"`
}

type JSONIndividual struct {
	Genome  JSONShapes `json:"genome"`
	Fitness float64    `json:"fitness"`
	ID      string     `json:"ID"`
}

// convert a Shapes instance to a json encodable version
func (s Shapes) ToJSON() JSONShapes {
	j := JSONShapes{
		Bounds: s.Bounds,
		Type:   s.Type,
	}

	nMembers := len(s.Members)

	switch s.Type {
	case "circles":
		j.Circles = make([]Circle, nMembers)
		for i, m := range s.Members {
			j.Circles[i] = m.(Circle)
		}
	case "polygons":
		j.Polygons = make([]Polygon, nMembers)
		for i, m := range s.Members {
			j.Polygons[i] = m.(Polygon)
		}
	case "triangles":
		j.Triangles = make([]Triangle, nMembers)
		for i, m := range s.Members {
			j.Triangles[i] = m.(Triangle)
		}
	}

	return j
}

// convert a decoded json shapes intance to a Shapes instance
func (j JSONShapes) ToShapes() Shapes {
	s := Shapes{
		Bounds: j.Bounds,
		Type:   j.Type,
	}

	switch s.Type {
	case "circles":
		s.Members = make([]Shape, len(j.Circles))
		for i, m := range j.Circles {
			s.Members[i] = m
		}
	case "polygons":
		s.Members = make([]Shape, len(j.Polygons))
		for i, m := range j.Polygons {
			s.Members[i] = m
		}
	case "triangles":
		s.Members = make([]Shape, len(j.Triangles))
		for i, m := range j.Triangles {
			s.Members[i] = m
		}
	}

	return s
}

func getJSONIndividual(i eaopt.Individual) JSONIndividual {
	j := JSONIndividual{
		Fitness: i.Fitness,
		ID:      i.ID,
	}

	if i.Genome == nil {
		j.Genome = JSONShapes{}
	} else {
		j.Genome = i.Genome.(Shapes).ToJSON()
	}

	return j
}

func getJSONIndividuals(pop []eaopt.Individual) []JSONIndividual {
	population := make([]JSONIndividual, len(pop))
	for i, m := range pop {
		population[i] = getJSONIndividual(m)
	}

	return population
}

func individualFromJSON(j JSONIndividual) eaopt.Individual {
	return eaopt.Individual{
		Genome:  j.Genome.ToShapes(),
		Fitness: j.Fitness,
		ID:      j.ID,
	}
}

func (t *Task) MarshalJSON() ([]byte, error) {
	type Alias Task

	return json.Marshal(&struct {
		BestFit    JSONIndividual   `json:"bestFit"`
		Population []JSONIndividual `json:"population"`
		*Alias
	}{
		BestFit:    getJSONIndividual(t.BestFit),
		Population: getJSONIndividuals(t.Population),
		Alias:      (*Alias)(t),
	})
}

func (t *Task) UnmarshalJSON(data []byte) error {
	type Alias Task

	aux := &struct {
		BestFit    JSONIndividual   `json:"bestFit"`
		Population []JSONIndividual `json:"population"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	population := make([]eaopt.Individual, len(aux.Population))
	for i, m := range aux.Population {
		population[i] = individualFromJSON(m)
	}

	t.BestFit = individualFromJSON(aux.BestFit)
	t.Population = population

	return nil
}
