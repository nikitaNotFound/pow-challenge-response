package usecases

import (
	"errors"
	"fmt"
	"log"
	"time"
	"wordofwisdom/internal/client_node/client_context"
	"wordofwisdom/internal/pow"
	"wordofwisdom/pkg/protocol/requests"
	"wordofwisdom/pkg/protocol/responses"
)

var (
	ErrUnexpectedServerResponse = errors.New("unexpected server response")
)

func RequestWisdom(ctx *client_context.ClientContext) error {
	if err := ctx.Sdk.SendMessage(true, requests.OPCODE_REQUEST_WISDOM, nil); err != nil {
		return err
	}

	msg, err := ctx.Sdk.WaitMessage()
	if err != nil {
		return err
	}
	if msg.Opcode != responses.RES_CODE_CHALLENGE {
		return ErrUnexpectedServerResponse
	}

	challengeRes := responses.ChallengeResponse{}
	if err := challengeRes.Decode(msg.Data); err != nil {
		return err
	}

	challenge := pow.Challenge{
		Data:           challengeRes.Data,
		Timestamp:      challengeRes.Timestamp,
		Difficulty:     challengeRes.Difficulty,
		ExpectedPrefix: challengeRes.ExpectedPrefix,
	}

	log.Printf("Challenge received [DIFFICULTY: %d]. Solving challenge...", challenge.Difficulty)
	started := time.Now()
	proof, err := challenge.Solve()
	if err != nil {
		return err
	}
	elapsed := time.Since(started)

	proofRequest := requests.ChallengeProofRequest{Nonce: proof}
	if err := ctx.Sdk.SendMessage(true, requests.OPCODE_REQUEST_CHALLENGE_PROOF, proofRequest); err != nil {
		return err
	}

	msg, err = ctx.Sdk.WaitMessage()
	if err != nil {
		return err
	}

	if msg.Opcode != responses.RES_CODE_WISDOM {
		return ErrUnexpectedServerResponse
	}

	wisdomRes := responses.WisdomResponse{}
	if err := wisdomRes.Decode(msg.Data); err != nil {
		return err
	}

	fmt.Printf("Wisdom received: %s ; [CHALLENGE TIME: %.4f seconds]\n", wisdomRes.Quote, elapsed.Seconds())
	return nil
}
