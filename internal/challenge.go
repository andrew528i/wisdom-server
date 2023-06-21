// The internal package is a collection of functions and types used by the Wisdom Server internally.
// It should not be imported or used directly by external packages.

package internal

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"time"

	"andrew528i/wisdom_server/internal/config"
	"andrew528i/wisdom_server/proto"
)

// Challenge represents a single challenge sent by the server to the client.
// The client must return a solution to the challenge before the deadline or the challenge will be considered invalid.
type Challenge struct {
	data      []byte    // random byte string of length ChallengeSize
	deadline  time.Time // the time at which the challenge expires
	nonce     uint64    // a unique nonce for the challenge
	solution  uint64    // the solution to the challenge
	signature []byte    // the cryptographic signature for the challenge
}

// NewChallenge generates a new challenge with a random byte string and a deadline from config.
// The challenge is signed using the secret provided as an argument.
func NewChallenge(nonce uint64, secret []byte) (*Challenge, error) {
	data := make([]byte, config.ChallengeSize)
	if _, err := rand.Read(data); err != nil {
		return nil, err
	}

	// set deadline and signature
	deadline := time.Now().Add(config.SignatureTimeout)
	signature := make([]byte, config.SignatureSize)
	challenge := &Challenge{
		data:      data,
		deadline:  deadline,
		nonce:     nonce,
		signature: signature,
	}

	// sign challenge
	if err := challenge.sign(secret); err != nil {
		return nil, err
	}

	return challenge, nil
}

// ChallengeFromMessage converts a proto.Challenge message to a Challenge struct.
func ChallengeFromMessage(challenge *proto.Challenge) *Challenge {
	data := make([]byte, config.ChallengeSize)
	copy(data, challenge.GetData())

	signature := make([]byte, config.SignatureSize)
	copy(signature, challenge.GetSignature())

	deadline := time.Unix(challenge.GetDeadline(), 0)

	return &Challenge{
		data:      data,
		deadline:  deadline,
		nonce:     challenge.GetNonce(),
		solution:  challenge.GetSolution(),
		signature: signature,
	}
}

// Copy returns a copy of the Challenge struct.
func (c Challenge) Copy() *Challenge {
	data := make([]byte, config.ChallengeSize)
	copy(data, c.data)

	signature := make([]byte, config.SignatureSize)
	copy(signature, c.signature)

	return &Challenge{
		data:      data,
		deadline:  c.deadline,
		nonce:     c.nonce,
		solution:  c.solution,
		signature: signature,
	}
}

// Solve solves the challenge with the provided difficulty.
func (c *Challenge) Solve(difficulty int) {
	// initialize target difficulty
	target := big.NewInt(1)
	target.Lsh(target, uint(32*8-difficulty*4))

	// increment solution until hash meets difficulty target
	for {
		hash := big.NewInt(0)
		hash.SetBytes(c.Hash())

		if hash.Cmp(target) == -1 {
			break
		}

		c.solution++
	}
}

// sign signs the challenge with the provided secret.
func (c *Challenge) sign(secret []byte) error {
	var emptySignature [32]byte

	// check that signature isn't already set
	if hex.EncodeToString(c.signature) != hex.EncodeToString(emptySignature[:]) {
		return ErrSigSet
	}

	hash := sha256.New()
	hash.Write(c.Hash())
	hash.Write(secret)

	c.signature = hash.Sum(nil)

	return nil
}

// Hash returns the SHA256 hash of the challenge data, nonce, solution, and deadline.
func (c Challenge) Hash() []byte {
	deadline := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(deadline, c.deadline.Unix())

	nonce := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(nonce, c.nonce)

	solution := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(solution, c.solution)

	hash := sha256.New()
	hash.Write(c.data[:])
	hash.Write(deadline)
	hash.Write(nonce)
	hash.Write(solution)

	return hash.Sum(nil)
}

// Check checks that the challenge is valid, given the current nonce, secret, and difficulty.
func (c Challenge) Check(difficulty int, currentNonce uint64, secret []byte) bool {
	// check deadline
	if c.deadline.Sub(time.Now()) < 0 {
		return false
	}

	// check nonce
	if currentNonce-c.nonce > config.NonceMaxDelta {
		return false
	}

	// check signature
	signature := make([]byte, config.SignatureSize)
	copy(signature, c.signature)

	solution := c.solution
	c.solution = 0
	c.signature = make([]byte, config.SignatureSize)
	if err := c.sign(secret); err != nil {
		return false
	}

	// in case user changed data, deadline or nonce (not counting solution which he of course should change)
	if hex.EncodeToString(c.signature) != hex.EncodeToString(signature) {
		return false
	}

	c.solution = solution

	// check solution
	target := big.NewInt(1)
	target.Lsh(target, uint(32*8-difficulty*4))
	hash := big.NewInt(0)
	hash.SetBytes(c.Hash())

	return hash.Cmp(target) == -1
}

func (c Challenge) Message() *proto.Challenge {
	data := make([]byte, config.ChallengeSize)
	copy(data, c.data)

	signature := make([]byte, config.SignatureSize)
	copy(signature, c.signature)

	return &proto.Challenge{
		Data:      data,
		Deadline:  c.deadline.Unix(),
		Nonce:     c.nonce,
		Solution:  c.solution,
		Signature: signature,
	}
}
