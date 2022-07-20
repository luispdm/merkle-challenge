package suites

import (
	"fmt"
	"net/http"

	"merkle-challenge/internal/client"
	"merkle-challenge/internal/encoder"
	"merkle-challenge/internal/entities"
	"merkle-challenge/internal/manager"
	"merkle-challenge/internal/merkle"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	desc = "Given the known root hash ':hash', " +
		"when a request is made to '/piece/:hash/:pieceIndex' and " +
		"the root hash is reconstructed, " +
		"the reconstructed hash matches with ':hash'"
)

var _ = Describe("Verify reconstructed hash with known root hash", func() {
	var (
		m *manager.Manager
		b encoder.Decoder
		h encoder.EncodeDecoder
	)

	setVars := func() {
		var err error
		m, err = manager.New(client.NewHTTP(cfg.ServerURL, cfg.ClientTimeout))
		ExpectWithOffset(1, err).ShouldNot(HaveOccurred(),
			fmt.Sprintf(errMsg, "creating", "manager", "object"))
		ExpectWithOffset(1, m).ShouldNot(BeNil(),
			fmt.Sprintf(errMsg, "creating", "manager", "object"))
		b = encoder.NewB64()
		h = encoder.NewHex()
	}

	decodeString := func(s string, d encoder.Decoder) []byte {
		content, err := d.DecodeString(s)
		Expect(err).ShouldNot(HaveOccurred(),
			fmt.Sprintf(errMsg, "decoding", "encoded", "string"))
		Expect(content).ShouldNot(BeNil(),
			fmt.Sprintf(errMsg, "decoding", "encoded", "string"))
		return content
	}

	getPiece := func(pI int) *entities.Piece {
		piecePath := fmt.Sprintf("%s/%s/%d", manager.Piece, cfg.IconsHash, pI)
		var p entities.Piece
		Expect(m.Get(manager.NewRequest(piecePath), &p)).
			To(Equal(http.StatusOK), fmt.Sprintf(errMsg, "getting", "piece", pI))
		return &p
	}

	BeforeEach(func() {
		setVars()
	})

	It(desc, func() {
		for pI := cfg.FirstPieceIndex; pI <= cfg.LastPieceIndex; pI++ {
			p := getPiece(pI)
			content := decodeString(p.Content, b)

			curHash := merkle.GetSHA256(content)

			for proofIndex := 0; proofIndex < len(p.Proofs); proofIndex++ {
				decProof := decodeString(p.Proofs[proofIndex], h)

				if merkle.IsRelativeOnTheRight(pI, proofIndex+1) {
					curHash = merkle.GetSHA256(merkle.ConcatByteSlices(curHash, decProof))
				} else {
					curHash = merkle.GetSHA256(merkle.ConcatByteSlices(decProof, curHash))
				}
			}

			encHash := h.EncodeToString(curHash)
			Expect(encHash).
				To(Equal(cfg.IconsHash),
					fmt.Sprintf("Reconstructed root hash '%s' from piece at index '%d' doesn't match"+
						" known root hash '%s'",
						encHash, pI, cfg.IconsHash))
		}
	})
})
