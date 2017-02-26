package main

type registry struct {
	slaveMap map[string]*slave
}

func (r *registry) addSlave(s *slave) {
	r.slaveMap[s.id] = s
}

func (r *registry) removeSlave(slaveID string) {
	slave, ok := r.slaveMap[slaveID]
	if !ok || slave == nil {
		return
	}
	slave.destroy()
	delete(r.slaveMap, slaveID)
}
