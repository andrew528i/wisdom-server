package config

import "time"

const (
	// ChallengeSize determines size in bytes of challenge random data
	ChallengeSize = 32

	// SecretSize is amount of bytes allocated for secret
	SecretSize = 32

	// SignatureSize is size of SHA256 based signature
	SignatureSize = 32

	// SignatureTimeout is the amount of time signed challenge is valid after generation by server
	SignatureTimeout = 2 * time.Second

	// Difficulty can be thought as zero count in hex representation of Challenge.Hash()
	// Important: Difficulty is correlated with SignatureTimeout so keep it in mind before changing them
	Difficulty = 5

	// NonceMaxDelta is maximum allowed delta = currentNonce - challengeNonce to pass Challenge.Check
	NonceMaxDelta = 100
)

// These are used only in defining cli arguments with flag package
const (
	DefaultHost = "0.0.0.0"
	DefaultPort = 9000
)

// Quotes are better to be fetched from database rather than config so keep this in mind before production
var Quotes = []string{
	"The fool doth think he is wise, but the wise man knows himself to be a fool.",
	"It is better to remain silent at the risk of being thought a fool, than to talk and remove all doubt of it.",
	"Knowing yourself is the beginning of all wisdom.",
	"The only true wisdom is in knowing you know nothing.",
	"The saddest aspect of life right now is that science gathers knowledge faster than society gathers wisdom.",
	"Count your age by friends, not years. Count your life by smiles, not tears.",
	"It is the mark of an educated mind to be able to entertain a thought without accepting it.",
	"The secret of life, though, is to fall seven times and to get up eight times.",
}
