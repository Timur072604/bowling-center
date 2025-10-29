package bowling

import (
	"math/rand"
	"sync"
	"time"
)

type DurationConfig struct {
	Base    time.Duration
	Variant time.Duration
}

type Config struct {
	NumLanes          int
	NumClients        int
	MaxClientWaitTime time.Duration
	ClientArrival     DurationConfig
	GameDuration      DurationConfig
}

type Center struct {
	config    Config
	clientsCh chan Client
	resultsCh chan GameResult
	wg        sync.WaitGroup
	rng       *rand.Rand
	state     *State
}

func New(cfg Config, state *State) *Center {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	return &Center{
		config:    cfg,
		clientsCh: make(chan Client, cfg.NumClients),
		resultsCh: make(chan GameResult, cfg.NumClients),
		rng:       rng,
		state:     state,
	}
}

func (c *Center) Run() {
	c.wg.Add(1)
	go c.runManager()

	go c.generateClients()

	for i := 0; i < c.config.NumClients; i++ {
		<-c.resultsCh
	}

	c.wg.Wait()
}

func (c *Center) runManager() {
	defer c.wg.Done()

	freeLanes := make(chan int, c.config.NumLanes)
	for i := 1; i <= c.config.NumLanes; i++ {
		freeLanes <- i
	}

	var waitingQueue []Client
	var timeoutCh <-chan time.Time

	for {
		if len(waitingQueue) > 0 {
			headClient := waitingQueue[0]
			remainingWaitTime := c.config.MaxClientWaitTime - time.Since(headClient.ArrivalTime)
			if remainingWaitTime < 0 {
				remainingWaitTime = 0
			}
			timeoutCh = time.After(remainingWaitTime)
		} else {
			timeoutCh = nil
		}

		select {
		case client, ok := <-c.clientsCh:
			if !ok {
				c.clientsCh = nil
			} else {
				waitingQueue = append(waitingQueue, client)
				c.state.mu.Lock()
				c.state.WaitingQueue = append(c.state.WaitingQueue, client)
				c.state.mu.Unlock()
			}
		case laneID := <-freeLanes:
			if len(waitingQueue) > 0 {
				nextClient := waitingQueue[0]
				waitingQueue = waitingQueue[1:]

				c.state.mu.Lock()
				c.state.WaitingQueue = waitingQueue
				c.state.Lanes[laneID] = &nextClient
				c.state.mu.Unlock()

				c.wg.Add(1)
				go playOnLane(nextClient, laneID, freeLanes, c.resultsCh, &c.wg, c.rng, c.state, c.config)
			} else {
				freeLanes <- laneID
			}
		case <-timeoutCh:
			if len(waitingQueue) > 0 {
				leavingClient := waitingQueue[0]
				waitingQueue = waitingQueue[1:]

				result := GameResult{Client: leavingClient, Status: StatusLeft}

				c.state.mu.Lock()
				c.state.WaitingQueue = waitingQueue
				c.state.FinishedGames = append(c.state.FinishedGames, result)
				c.state.mu.Unlock()

				c.resultsCh <- result
			}
		}

		if c.clientsCh == nil && len(waitingQueue) == 0 {
			break
		}
	}
}

func (c *Center) generateClients() {
	defer close(c.clientsCh)
	for i := 1; i <= c.config.NumClients; i++ {
		client := Client{
			ID:          i,
			GameString:  GenerateRandomGame(c.rng),
			ArrivalTime: time.Now(),
		}
		c.clientsCh <- client

		sleepTime := c.config.ClientArrival.Base + time.Duration(c.rng.Intn(int(c.config.ClientArrival.Variant)))
		time.Sleep(sleepTime)
	}
}

func playOnLane(client Client, laneID int, freeLanes chan<- int, resultsCh chan<- GameResult, wg *sync.WaitGroup, rng *rand.Rand, state *State, cfg Config) {
	defer wg.Done()

	gameTime := cfg.GameDuration.Base + time.Duration(rng.Intn(int(cfg.GameDuration.Variant)))
	time.Sleep(gameTime)

	finalScore, err := Score(client.GameString)
	result := GameResult{Client: client, Score: finalScore, Err: err, Status: StatusPlayed}

	state.mu.Lock()
	state.Lanes[laneID] = nil
	state.FinishedGames = append(state.FinishedGames, result)
	state.mu.Unlock()

	resultsCh <- result
	freeLanes <- laneID
}
