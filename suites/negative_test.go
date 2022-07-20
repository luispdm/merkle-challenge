package suites

import (
	"fmt"
	"merkle-challenge/internal/client"
	"merkle-challenge/internal/entities"
	"merkle-challenge/internal/manager"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Malformed requests to '/piece' endpoint", func() {
	var (
		m *manager.Manager
	)
	setVars := func() {
		var err error
		m, err = manager.New(client.NewHTTP(cfg.ServerURL, cfg.ClientTimeout))
		ExpectWithOffset(1, err).ShouldNot(HaveOccurred(),
			fmt.Sprintf(errMsg, "creating", "manager", "object"))
		ExpectWithOffset(1, m).ShouldNot(BeNil(),
			fmt.Sprintf(errMsg, "creating", "manager", "object"))
	}

	BeforeEach(func() {
		setVars()
	})

	DescribeTable("Given the known root hash ':hash'", func(pieceIndex string, expSC int) {
		piecePath := fmt.Sprintf("%s/%s/%s", manager.Piece, cfg.IconsHash, pieceIndex)

		sC, err := m.Get(manager.NewRequest(piecePath), &entities.Piece{})
		Expect(sC).To(Equal(expSC),
			fmt.Sprintf("Expecting request to '/piece/%s/%s' to return '%d', got '%d'",
				cfg.IconsHash, pieceIndex, expSC, sC))
		Expect(err).ShouldNot(BeNil(),
			fmt.Sprintf("Expecting request to '/piece/%s/%s' to return an error but no error found",
				cfg.IconsHash, pieceIndex))
	},
		EntryDescription("when a request is made to '/piece/:hash/%s', the server returns %d"),
		Entry(nil, "-1", 404),
		Entry(nil, "10", 404),
		Entry(nil, "abc", 400),
		Entry(nil, "🫠", 400),
	)
})
