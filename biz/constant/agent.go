package constant

const (
	SlotAgent    = "slot"
	OncallAgent  = "oncall"
	NeedClassify = "need_classify"
)

type NodeState int

const (
	NodeStateInit NodeState = iota
	NodeStateRestart
	NodeStateContinue
	NodeStateFinish
)

type InputType int

const (
	InputTypeText InputType = iota
	InputTypeCard
)

type AgentCode int

const (
	AgentCodeSuccess AgentCode = iota
	AgentCodeFailure
	AgentCodeNeedClassify
	AgentCodeGateway
)
