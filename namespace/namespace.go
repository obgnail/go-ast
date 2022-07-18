package namespace

import "fmt"

var (
	ErrGet = fmt.Errorf("error getter")
)

type Secondary struct {
	getterMap map[string]map[string]interface{}
}

func NewSecondary() *Secondary {
	hub := &Secondary{getterMap: make(map[string]map[string]interface{})}
	return hub
}

func (hub *Secondary) Set(ns1, ns2 string, value interface{}) {
	m, ok := hub.getterMap[ns1]
	if !ok {
		m = make(map[string]interface{})
		hub.getterMap[ns1] = m
	}
	m[ns2] = value
}

func (hub *Secondary) Get(ns1 string, ns2 ...string) (interface{}, bool) {
	switch len(ns2) {
	case 0:
		if m1, ok := hub.getterMap[ns1]; ok {
			return m1, true
		}
	case 1:
		if m1, ok := hub.getterMap[ns1]; ok {
			if m2, ok := m1[ns2[0]]; ok {
				return m2, true
			}
		}
	}
	return nil, false
}

func (hub *Secondary) AppendMap(ns1 string, m map[string]interface{}) {
	for ns2, value := range m {
		hub.Set(ns1, ns2, value)
	}
}

func (hub *Secondary) All() map[string]map[string]interface{} {
	return hub.getterMap
}

type Tertiary struct {
	getterMap map[string]map[string]map[string]interface{}
}

func NewTertiary() *Tertiary {
	hub := &Tertiary{getterMap: make(map[string]map[string]map[string]interface{})}
	return hub
}

func (hub *Tertiary) Set(ns1, ns2, ns3 string, value interface{}) {
	m, ok := hub.getterMap[ns1]
	if !ok {
		m = make(map[string]map[string]interface{})
		hub.getterMap[ns1] = m
	}
	m2, ok := m[ns2]
	if !ok {
		m2 = make(map[string]interface{})
		hub.getterMap[ns1][ns2] = m2
	}

	m2[ns3] = value
}

func (hub *Tertiary) Get(ns1 string, other ...string) (interface{}, bool) {
	switch len(other) {
	case 0:
		if m1, ok := hub.getterMap[ns1]; ok {
			return m1, true
		}
	case 1:
		if m1, ok := hub.getterMap[ns1]; ok {
			if m2, ok := m1[other[0]]; ok {
				return m2, true
			}
		}
	case 2:
		ns2, ns3 := other[0], other[1]
		if m1, ok := hub.getterMap[ns1]; ok {
			if m2, ok := m1[ns2]; ok {
				if m3, ok := m2[ns3]; ok {
					return m3, true
				}
			}
		}
	}
	return nil, false
}

func (hub *Tertiary) AppendSecondary(ns1 string, s *Secondary) {
	for ns2, kv := range s.getterMap {
		for ns3, value := range kv {
			hub.Set(ns1, ns2, ns3, value)
		}
	}
}

func (hub *Tertiary) All() map[string]map[string]map[string]interface{} {
	return hub.getterMap
}
