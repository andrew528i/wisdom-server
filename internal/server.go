package internal

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"
	"math/rand"
	"sync"
	"time"

	"andrew528i/wisdom_server/internal/config"
	"andrew528i/wisdom_server/proto"
	log "github.com/sirupsen/logrus"
)

var _ proto.WisdomBookServer = (*WisdomBookServer)(nil)

type WisdomBookServer struct {
	proto.UnimplementedWisdomBookServer
	mut sync.Mutex

	nonce  uint64
	secret []byte
}

func NewWisdomBookServer() (*WisdomBookServer, error) {
	nonce := uint64(1)
	secret := make([]byte, config.SecretSize)

	_, err := cryptorand.Read(secret)
	if err != nil {
		return nil, err
	}

	return &WisdomBookServer{
		nonce:  nonce,
		secret: secret,
	}, nil
}

func (s *WisdomBookServer) GetChallenge(_ context.Context, _ *proto.Empty) (*proto.Challenge, error) {
	challenge, err := NewChallenge(s.nonce, s.secret)
	if err != nil {
		return nil, err
	}

	s.mut.Lock()
	defer s.mut.Unlock()
	s.nonce++

	log.WithField("method", "GetChallenge").
		WithField("difficulty", config.Difficulty).
		WithField("nonce", challenge.nonce).
		WithField("signature", hex.EncodeToString(challenge.signature)).
		Info("generated new challenge")

	return challenge.Message(), nil
}

func (s *WisdomBookServer) GetQuote(_ context.Context, protoChallenge *proto.Challenge) (*proto.Quote, error) {
	challenge := ChallengeFromMessage(protoChallenge)

	s.mut.Lock()
	defer s.mut.Unlock()

	checkResult := challenge.Check(config.Difficulty, s.nonce, s.secret)

	log.WithField("method", "GetQuote").
		WithField("challenge_hash", hex.EncodeToString(challenge.Hash())).
		WithField("challenge_nonce", challenge.nonce).
		WithField("challenge_ms_until_deadline", challenge.deadline.Sub(time.Now()).Milliseconds()).
		WithField("challenge_solution", challenge.solution).
		WithField("challenge_check_result", checkResult).
		Info("received challenge")

	// if there will be more than one rpc method it's better to create middleware to check challenge
	if !checkResult {
		return nil, ErrChallengeInvalid
	}

	randomQuote := config.Quotes[rand.Intn(len(config.Quotes))]

	log.WithField("method", "GetQuote").
		WithField("quote", randomQuote).
		Info("client deserved to get quote")

	return &proto.Quote{Quote: randomQuote}, nil
}
