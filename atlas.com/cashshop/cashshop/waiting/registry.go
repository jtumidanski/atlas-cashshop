package waiting

import (
	"atlas-cashshop/model"
	"errors"
	"sync"
)

type registry struct {
	characters map[uint32]Model
	lock       sync.RWMutex
}

var r *registry
var once sync.Once

func GetRegistry() *registry {
	once.Do(func() {
		r = &registry{
			characters: make(map[uint32]Model, 0),
			lock:       sync.RWMutex{},
		}
	})
	return r
}

func (r *registry) Add(worldId byte, channelId byte, characterId uint32) error {
	r.lock.Lock()
	if _, ok := r.characters[characterId]; ok {
		r.lock.Unlock()
		return errors.New("character already waiting")
	}
	r.characters[characterId] = Model{worldId: worldId, channelId: channelId, characterId: characterId}
	r.lock.Unlock()
	return nil
}

func (r *registry) AddApproval(characterId uint32) error {
	r.lock.Lock()
	if val, ok := r.characters[characterId]; ok {
		r.characters[characterId] = val.AddApproval()
		r.lock.Unlock()
		return nil
	}
	r.lock.Unlock()
	return errors.New("character is not waiting for approval")
}

func (r *registry) ProcessAllApproved(minimum uint32, operator model.SliceOperator[Model]) error {
	r.lock.Lock()
	results := make([]Model, 0)
	for characterId, m := range r.characters {
		if m.Approvals() >= minimum {
			results = append(results, m)
			delete(r.characters, characterId)
		}
	}
	r.lock.Unlock()
	return operator(results)
}

func (r *registry) ProcessIfApproved(characterId uint32, minimum uint32, operator model.Operator[Model]) error {
	r.lock.Lock()
	if val, ok := r.characters[characterId]; ok {
		if val.Approvals() >= minimum {
			delete(r.characters, characterId)
			r.lock.Unlock()
			return operator(val)
		}
		r.lock.Unlock()
		return nil
	}
	r.lock.Unlock()
	return errors.New("character is not waiting for approval")
}

func (r *registry) Remove(characterId uint32) (Model, error) {
	r.lock.Lock()
	if val, ok := r.characters[characterId]; ok {
		delete(r.characters, characterId)
		r.lock.Unlock()
		return val, nil
	}
	r.lock.Unlock()
	return Model{}, errors.New("character is not waiting for approval")
}
