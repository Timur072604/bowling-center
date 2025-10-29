package bowling

import "sync"

type StateSnapshot struct {
	Lanes         map[int]*Client
	WaitingQueue  []Client
	FinishedGames []GameResult
}

type State struct {
	mu            sync.RWMutex
	Lanes         map[int]*Client
	WaitingQueue  []Client
	FinishedGames []GameResult
}

func NewState(numLanes int) *State {
	lanes := make(map[int]*Client, numLanes)
	for i := 1; i <= numLanes; i++ {
		lanes[i] = nil
	}

	return &State{
		Lanes:         lanes,
		WaitingQueue:  make([]Client, 0),
		FinishedGames: make([]GameResult, 0),
	}
}

func (s *State) GetSnapshot() StateSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	lanesCopy := make(map[int]*Client, len(s.Lanes))
	for k, v := range s.Lanes {
		lanesCopy[k] = v
	}

	queueCopy := make([]Client, len(s.WaitingQueue))
	copy(queueCopy, s.WaitingQueue)

	finishedCopy := make([]GameResult, len(s.FinishedGames))
	copy(finishedCopy, s.FinishedGames)

	return StateSnapshot{
		Lanes:         lanesCopy,
		WaitingQueue:  queueCopy,
		FinishedGames: finishedCopy,
	}
}
