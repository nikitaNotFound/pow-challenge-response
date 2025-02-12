package protocol

const (
	RES_CODE_CHALLENGE uint32 = 1
	RES_CODE_WISDOM    uint32 = 2
)

type ChallengeResponse struct {
	Data       string
	Timestamp  int64
	Difficulty int
}

func (cr *ChallengeResponse) EncodeChallengeResponse() ([]byte, error) {
	buff := make([]byte, 4)
}

type WisdomResponse struct {
	Quote string
}
