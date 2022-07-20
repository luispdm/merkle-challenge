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
	fDesc = "Given the known root hash ':hash', " +
		"when a request is made to '/piece/:hash/:pieceIndex', " +
		"and the piece hash is calculated, " +
		"and a request is made to ':pieceIndex' sibling, " +
		"the piece hash matches the first element of the proofs array of its sibling"
	sDesc = "Given the known root hash ':hash', " +
		"when a request is made to '/piece/:hash/:pieceIndex', " +
		"and a request is made to ':pieceIndex' sibling, " +
		"the uncles of both pieces (elements of the proofs array except for the first one) match"
)

var _ = Describe("Sibling tests", func() {
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

	getPiece := func(pI int) *entities.Piece {
		var p entities.Piece
		pPath := fmt.Sprintf("%s/%s/%d", manager.Piece, cfg.IconsHash, pI)
		sC, err := m.Get(manager.NewRequest(pPath), &p)
		ExpectWithOffset(1, err).ShouldNot(HaveOccurred(),
			fmt.Sprintf(errMsg, "getting", "piece", pI))
		ExpectWithOffset(1, sC).To(Equal(http.StatusOK),
			fmt.Sprintf(errMsg, "getting", "piece", pI))
		return &p
	}

	decodeString := func(s string) []byte {
		content, err := b.DecodeString(s)
		ExpectWithOffset(1, err).ShouldNot(HaveOccurred(),
			fmt.Sprintf(errMsg, "decoding", "base64-encoded", "string"))
		ExpectWithOffset(1, content).ShouldNot(BeNil(),
			fmt.Sprintf(errMsg, "decoding", "base64-encoded", "string"))
		return content
	}

	BeforeEach(func() {
		setVars()
	})

	It(fDesc, func() {
		for pI := cfg.FirstPieceIndex; pI <= cfg.LastPieceIndex-1; pI++ {
			content := decodeString(getPiece(pI).Content)

			pHash := h.EncodeToString(merkle.GetSHA256(content))

			sibling := getPiece(merkle.GetSiblingIndex(pI))

			Expect(pHash).
				To(Equal(sibling.Proofs[0]),
					fmt.Sprintf("Piece hash '%s' doesn't match sibling's first element "+
						"of the proofs array '%s'", pHash, sibling.Proofs[0]))
		}
	})

	It(sDesc, func() {
		for pI := cfg.FirstPieceIndex; pI <= cfg.LastPieceIndex-2; pI += 2 {
			p := getPiece(pI)
			s := getPiece(merkle.GetSiblingIndex(pI))

			pUncles := p.Proofs[1:]
			sUncles := s.Proofs[1:]

			Expect(pUncles).
				To(Equal(sUncles),
					fmt.Sprintf("Uncles of piece at index '%d' don't match sibling's uncles", pI))
		}
	})
})
