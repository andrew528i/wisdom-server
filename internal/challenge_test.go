package internal

import (
	"crypto/rand"
	"testing"
	"time"

	"andrew528i/wisdom_server/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestChallenge_Solve_Check(t *testing.T) {
	secret := make([]byte, config.SecretSize)
	_, err := rand.Read(secret)
	assert.NoError(t, err)

	evilSecret := make([]byte, config.SecretSize)
	_, err = rand.Read(secret)
	assert.NoError(t, err)

	// generate challenge itself and solve it to perform different checks for Check method
	currentNonce := uint64(1000)
	challenge, err := NewChallenge(currentNonce, secret)
	assert.NoError(t, err)

	challenge.Solve(config.Difficulty)

	for _, tc := range []struct {
		description  string
		difficulty   int
		currentNonce uint64
		secret       []byte
		success      bool
		before       func(challenge *Challenge)
	}{
		{
			"all is ok",
			config.Difficulty,
			currentNonce + 1,
			secret,
			true,
			nil,
		}, {
			"swap secret",
			config.Difficulty,
			currentNonce + 1,
			evilSecret,
			false,
			nil,
		}, {
			"exceed nonce",
			config.Difficulty,
			currentNonce + config.NonceMaxDelta + 1,
			secret,
			false,
			nil,
		}, {
			"increase difficulty",
			config.Difficulty + 1,
			currentNonce + 1,
			secret,
			false,
			nil,
		}, {
			"make solution invalid",
			config.Difficulty,
			currentNonce + 1,
			secret,
			false,
			func(challenge *Challenge) { challenge.solution-- },
		}, {
			"exceed signature deadline",
			config.Difficulty,
			currentNonce + 1,
			secret,
			false,
			func(*Challenge) {
				time.Sleep(config.SignatureTimeout)
			},
		}, {
			"spoof data",
			config.Difficulty,
			currentNonce + 1,
			secret,
			false,
			func(challenge *Challenge) {
				_, err := rand.Read(challenge.data)
				assert.NoError(t, err)
			},
		}, {
			"spoof nonce",
			config.Difficulty,
			currentNonce + 1,
			secret,
			false,
			func(challenge *Challenge) {
				challenge.nonce++
			},
		}, {
			"spoof deadline",
			config.Difficulty,
			currentNonce + 1,
			secret,
			false,
			func(challenge *Challenge) {
				challenge.deadline.Add(time.Millisecond)
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			challengeCopy := challenge.Copy()

			if tc.before != nil {
				tc.before(challengeCopy)
			}

			result := challengeCopy.Check(tc.difficulty, tc.currentNonce, tc.secret)
			assert.Equal(t, result, tc.success)
		})
	}
}
