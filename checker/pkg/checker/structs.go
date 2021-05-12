package checker

type InfoMessage struct {
	ServiceName   string `json:"serviceName"`
	FlagVariants  uint64 `json:"flagVariants"`
	NoiseVariants uint64 `json:"noiseVariants"`
	HavocVariants uint64 `json:"havocVariants"`
}

type Result string

var (
	ResultOk      = Result("OK")
	ResultError   = Result("INTERNAL_ERROR")
	ResultMumble  = Result("MUMBLE")
	ResultOffline = Result("OFFLINE")
)

type ResultMessage struct {
	Result  Result `json:"result"`
	Message string `json:"message"`
}

type TaskMessageMethod string

var (
	TaskMessageMethodPutFlag  = TaskMessageMethod("putflag")
	TaskMessageMethodGetFlag  = TaskMessageMethod("getflag")
	TaskMessageMethodPutNoise = TaskMessageMethod("putnoise")
	TaskMessageMethodGetNoise = TaskMessageMethod("getnoise")
	TaskMessageMethodHavoc    = TaskMessageMethod("havoc")
)

type TaskMessage struct {
	TaskId         uint64            `json:"taskId"`
	Method         TaskMessageMethod `json:"method"`
	Address        string            `json:"address"`
	TeamId         uint64            `json:"teamId"`
	TeamName       string            `json:"teamName"`
	CurrentRoundId uint64            `json:"currentRoundId"`
	RelatedRoundId uint64            `json:"relatedRoundId"`
	Flag           string            `json:"flag"`
	VariantId      uint64            `json:"variantId"`
	Timeout        uint64            `json:"timeout"`
	RoundLength    uint64            `json:"roundLength"`
	TaskChainId    string            `json:"taskChainId"`
}
