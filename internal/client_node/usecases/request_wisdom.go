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
	"wordofwisdom/pkg/server_sdk"
)

var (
	ErrUnexpectedServerResponse = errors.New("unexpected server response")
)

func RequestWisdom(ctx *client_context.ClientContext) error {
	if err := ctx.Sdk.SendMessage(true, requests.OPCODE_REQUEST_WISDOM, nil); err != nil {
		return err
	}

	msg, err := ctx.Sdk.PopMessage()
	if err != nil {
		if errors.Is(err, server_sdk.ErrPopMessageTimeout) {
			ctx.Sdk.CloseConnection()
			log.Println("Closing connection due to message pop timeout.")
			return err
		}
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

	msg, err = ctx.Sdk.PopMessage()
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

	fmt.Printf("Wisdom received; [CHALLENGE TIME: %.4f seconds]\n", elapsed.Seconds())
	return nil
}
