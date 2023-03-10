package fiatshamir

import (
	"github.com/xingzwh/gnark/frontend"
	"github.com/xingzwh/gnark/std/hash"
)

type Settings struct {
	Transcript     *Transcript
	Prefix         string
	BaseChallenges []frontend.Variable
	Hash           hash.Hash
}

func WithTranscript(transcript *Transcript, prefix string, baseChallenges ...frontend.Variable) Settings {
	return Settings{
		Transcript:     transcript,
		Prefix:         prefix,
		BaseChallenges: baseChallenges,
	}
}

func WithHash(hash hash.Hash, baseChallenges ...frontend.Variable) Settings {
	return Settings{
		BaseChallenges: baseChallenges,
		Hash:           hash,
	}
}
