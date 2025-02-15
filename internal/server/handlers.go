package server

import (
	"errors"
	"wordofwisdom/internal/pow"
	"wordofwisdom/internal/protocol/requests"
	"wordofwisdom/internal/protocol/responses"
)

type serverHandlers struct {
	challengeDifficulty int64
}

func NewServerHandlers(challengeDifficulty int64) *serverHandlers {
	return &serverHandlers{
		challengeDifficulty: challengeDifficulty,
	}
}

func (h *serverHandlers) handleRequestWisdom(svrCtx ServerContext) error {
	challenge := pow.GenerateChallenge(h.challengeDifficulty)
	challengeResponse := responses.ChallengeResponse{
		Data:           challenge.Data,
		Timestamp:      uint64(challenge.Timestamp),
		Difficulty:     uint64(challenge.Difficulty),
		ExpectedPrefix: challenge.ExpectedPrefix,
	}
	svrCtx.SendSuccessMessage(responses.RES_CODE_CHALLENGE, &challengeResponse)

	message, err := svrCtx.WaitMessage()
	if err != nil {
		return err
	}

	if message.Opcode != requests.OPCODE_REQUEST_CHALLENGE_PROOF {
		return errors.New("invalid opcode")
	}

	challengeProofRequest := requests.ChallengeProofRequest{}
	if err := challengeProofRequest.Decode(message.Data); err != nil {
		return err
	}

	if !challenge.Verify(challengeProofRequest.Nonce) {
		svrCtx.SendError(requests.OPCODE_REQUEST_CHALLENGE_PROOF)
		return nil
	}

	quote := GetRandomQuote()
	svrCtx.SendSuccessMessage(responses.RES_CODE_WISDOM, &responses.WisdomResponse{Quote: quote})
	return nil
}
