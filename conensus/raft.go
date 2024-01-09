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
	peers         []Raft // Assuming there's a Node type for peer communication
	leaderID      int
}

// ApplyMsg represents a message to apply to the state machine.
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

// RequestVoteArgs represents the arguments for the RequestVote RPC.
type RequestVoteArgs struct {
	Term         int
	CandidateID  int
	LastLogIndex int
	LastLogTerm  int
}

// RequestVoteReply represents the reply for the RequestVote RPC.
type RequestVoteReply struct {
	Term        int
	VoteGranted bool
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

	if r.state == Leader {
		// If the node is already the leader, reset the election timer and return.
		r.electionTimer.Reset(randomElectionTimeout())
		return
	}

	r.state = Candidate
	r.currentTerm++
	r.votedFor = r.ID
	r.electionTimer.Reset(randomElectionTimeout())

	// Request votes from other nodes
	args := RequestVoteArgs{
		Term:         r.currentTerm,
		CandidateID:  r.ID,
		LastLogIndex: len(r.log) - 1,
		LastLogTerm:  r.log[len(r.log)-1].Term,
	}

	var votesReceived int
	voteCh := make(chan bool, len(r.peers))

	for _, peer := range r.peers {
		go func(peerID int) {
			var reply RequestVoteReply
			if r.sendRequestVote(peerID, &args, &reply) {
				r.mu.Lock()
				defer r.mu.Unlock()

				if reply.VoteGranted {
					votesReceived++
					if votesReceived >= len(r.peers)/2+1 {
						// Received majority of votes, become the leader
						r.state = Leader
						r.leaderID = r.ID

						// Initialize leader-specific data structures and start sending heartbeats
						// r.initLeaderState()

						fmt.Printf("Node %d became the leader (Term %d)\n", r.ID, r.currentTerm)
					}
				} else if reply.Term > r.currentTerm {
					// Detected a higher term, transition to Follower state
					r.state = Follower
					r.currentTerm = reply.Term
					r.votedFor = -1
				}
			}
			voteCh <- reply.VoteGranted
		}(peer.ID)
	}

	// Wait for the votes to be collected
	votesNeeded := len(r.peers)/2 + 1
	votesReceived = 1 // Count the vote from the candidate itself
	for i := 0; i < len(r.peers)-1; i++ {
		if <-voteCh {
			votesReceived++
		}
		if votesReceived >= votesNeeded {
			// Received majority of votes, become the leader
			r.state = Leader
			r.leaderID = r.ID

			// Initialize leader-specific data structures and start sending heartbeats
			// r.initLeaderState()

			fmt.Printf("Node %d became the leader (Term %d)\n", r.ID, r.currentTerm)
			return
		}
	}

	// Did not receive majority of votes, reset to Follower state
	r.state = Follower
	r.currentTerm++
	r.votedFor = -1
	fmt.Printf("Node %d did not receive majority of votes, back to Follower state (Term %d)\n", r.ID, r.currentTerm)
}

// sendRequestVote sends a RequestVote RPC to another node.
func (r *Raft) sendRequestVote(peerID int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	// Implement the logic to send a RequestVote RPC to the given peer and
	// handle the reply. Return true if the RPC was successful, false otherwise.
	// You can use Go's net/rpc package for this.
	// Example:
	// if err := r.peers[peerID].Call("Node.RequestVote", args, reply); err != nil {
	//     return false
	// }
	return true
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
