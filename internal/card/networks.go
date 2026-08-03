package card

import (
	"cmp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const MaxIINLength = 8

type CardIssuer string

const (
	Visa                    CardIssuer = "Visa"
	Mastercard                         = "Mastercard"
	AmericanExpress                    = "American Express"
	Discover                           = "Discover"
	DinersClubInternational            = "Diners Club International"
	JCB                                = "JCB"
	UnionPay                           = "UnionPay"
)

type IINRange struct {
	start string
	end   string
}

type NetworkCardRange struct {
	start  int
	end    int
	issuer CardIssuer
}

type CardNetwork struct {
	issuer        CardIssuer
	iinRanges     []IINRange
	numberLengths []int
}

type NetworkStorage interface {
	FindIssuer(iin string) (CardIssuer, bool)
	GetNumberLengths(issuer CardIssuer) []int
}

type networkStorageImpl struct {
	networkRanges []NetworkCardRange
	networkMap    map[CardIssuer]*CardNetwork
}

func NewNetworkStorage() (NetworkStorage, error) {
	networks := []*CardNetwork{
		{
			issuer:        Visa,
			iinRanges:     []IINRange{{"4", "4"}},
			numberLengths: []int{13, 16, 19},
		},
		{
			issuer:        Mastercard,
			iinRanges:     []IINRange{{"2221", "2720"}, {"51", "55"}},
			numberLengths: []int{16},
		},
		{
			issuer:        AmericanExpress,
			iinRanges:     []IINRange{{"34", "34"}, {"37", "37"}},
			numberLengths: []int{15},
		},
		{
			issuer: Discover,
			iinRanges: []IINRange{
				{"6011", "6011"},
				{"622126", "622925"},
				{"644", "649"},
				{"65", "65"},
			},
			numberLengths: []int{16, 17, 18, 19},
		},
		{
			issuer: DinersClubInternational,
			iinRanges: []IINRange{
				{"30", "30"},
				{"36", "36"},
				{"38", "38"},
				{"39", "39"},
			},
			numberLengths: []int{14, 15, 16, 17, 18, 19},
		},
		{
			issuer:        JCB,
			iinRanges:     []IINRange{{"3528", "3589"}},
			numberLengths: []int{16, 17, 18, 19},
		},
		{
			issuer:        UnionPay,
			iinRanges:     []IINRange{{"62", "62"}},
			numberLengths: []int{16, 17, 18, 19},
		},
	}

	s := &networkStorageImpl{
		networkRanges: []NetworkCardRange{},
		networkMap:    make(map[CardIssuer]*CardNetwork),
	}

	for _, network := range networks {
		s.networkMap[network.issuer] = network
		for _, iinRange := range network.iinRanges {
			start, end, err := s.normalizeIINRange(iinRange)
			if err != nil {
				return nil, err
			}
			s.networkRanges = append(s.networkRanges, NetworkCardRange{
				issuer: network.issuer,
				start:  start,
				end:    end,
			})
		}
	}

	slices.SortFunc(s.networkRanges, func(a, b NetworkCardRange) int {
		return cmp.Compare(a.start, b.start)
	})

	return s, nil
}

func (s *networkStorageImpl) normalizeIINRange(iinRange IINRange) (start, end int, err error) {
	start, err = s.normalizeIIN(iinRange.start, false)
	if err != nil {
		return
	}
	end, err = s.normalizeIIN(iinRange.end, true)
	return
}

func (s *networkStorageImpl) normalizeIIN(iin string, isEnd bool) (int, error) {
	var padCh string
	if isEnd {
		padCh = "9"
	} else {
		padCh = "0"
	}

	iin += strings.Repeat(padCh, MaxIINLength-len(iin))
	num, err := strconv.Atoi(iin)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (s *networkStorageImpl) FindIssuer(iin string) (issuer CardIssuer, found bool) {
	iinNumber, err := s.normalizeIIN(iin, false)
	if err != nil {
		found = false
		return
	}

	i := sort.Search(len(s.networkRanges), func(i int) bool {
		return s.networkRanges[i].start >= iinNumber
	})

	if i == len(s.networkRanges) {
		found = false
		return
	}

	if s.networkRanges[i].end < iinNumber {
		found = false
		return
	}

	return s.networkRanges[i].issuer, true
}

func (s *networkStorageImpl) GetNumberLengths(issuer CardIssuer) []int {
	network, ok := s.networkMap[issuer]
	if !ok {
		return []int{}
	}
	return slices.Clone(network.numberLengths)
}
