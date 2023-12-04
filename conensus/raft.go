package raft

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// LogEntry represents a log entry in the Raft algorithm.
type LogEntry struct {
	Term    int
	Command interface{}
}

// State represents the state of a Raft node.
type State int

const (
	Follower State = iota
	Candidate
	Leader
)

// Raft represents a node in the Raft distributed consensus algorithm.
type Raft struct {
	mu            sync.Mutex
	ID            int
	currentTerm   int
	votedFor      int
	log           []LogEntry
	commitIndex   int
	lastApplied   int
	nextIndex     map[int]int
	matchIndex    map[int]int
	state         State
	electionTimer *time.Timer
	applyCh       chan ApplyMsg
}

// ApplyMsg represents a message to apply to the state machine.
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

// NewRaft creates a new Raft node with the given ID.
func NewRaft(id int, applyCh chan ApplyMsg) *Raft {
	return &Raft{
		ID:            id,
		state:         Follower,
		electionTimer: time.NewTimer(randomElectionTimeout()),
		applyCh:       applyCh,
	}
}

// Start starts the main event loop of the Raft node.
func (r *Raft) Start() {
	for {
		select {
		case <-r.electionTimer.C:
			r.startElection()
			// Handle other events and RPCs here
		}
	}
}

// startElection transitions the node to the Candidate state and initiates an election.
func (r *Raft) startElection() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.state = Candidate
	r.currentTerm++
	r.votedFor = r.ID
	r.electionTimer.Reset(randomElectionTimeout())

	// Send RequestVote RPCs to other nodes

	// Handle votes from other nodes
}

// ClientRequest handles a client request to set a value in the Raft cluster.
func (r *Raft) ClientRequest(command interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != Leader {
		fmt.Printf("Node %d is not the leader. Redirect the request to the leader.\n", r.ID)
		// Redirect the request to the leader
		return
	}

	// Append the command to the log
	entry := LogEntry{Term: r.currentTerm, Command: command}
	r.log = append(r.log, entry)

	// Replicate the log entry to other nodes
	// ...

	// Notify the client that the command is applied
	r.applyCh <- ApplyMsg{
		CommandValid: true,
		Command:      command,
		CommandIndex: len(r.log),
	}
}

func randomElectionTimeout() time.Duration {
	rand.Seed(time.Now().UnixNano())
	return time.Duration(rand.Intn(300)+300) * time.Millisecond
}
