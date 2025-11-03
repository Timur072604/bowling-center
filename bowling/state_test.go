package bowling

import (
	"testing"
)

func TestNewState(t *testing.T) {
	t.Run("Создание состояния с 3 дорожками", func(t *testing.T) {
		numLanes := 3
		state := NewState(numLanes)

		if state == nil {
			t.Fatal("state не должен быть nil")
		}

		if len(state.Lanes) != numLanes {
			t.Errorf("ожидалось %d дорожек, получено %d", numLanes, len(state.Lanes))
		}

		for i := 1; i <= numLanes; i++ {
			if _, ok := state.Lanes[i]; !ok {
				t.Errorf("дорожка с ID %d не найдена", i)
			}
			if state.Lanes[i] != nil {
				t.Errorf("дорожка %d должна быть свободна (nil), но она занята", i)
			}
		}

		if state.WaitingQueue == nil {
			t.Error("WaitingQueue не должна быть nil")
		}
		if len(state.WaitingQueue) != 0 {
			t.Error("WaitingQueue должна быть пустой при создании")
		}

		if state.FinishedGames == nil {
			t.Error("FinishedGames не должна быть nil")
		}
		if len(state.FinishedGames) != 0 {
			t.Error("FinishedGames должна быть пустой при создании")
		}
	})
}
