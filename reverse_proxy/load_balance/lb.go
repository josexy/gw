package load_balance

type LbType int

const (
	LbRandom LbType = iota
	LbRoundRobin
	LbWeightRoundRobin
	LbConsistentHash
	LbIPHash
)

type LoadBalance interface {
	Add(...string) error
	Get(string) (string, error)
}
