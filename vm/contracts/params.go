package contracts

import (
	"github.com/vitelabs/go-vite/vm/util"
	"math/big"
)

const (
	registerGas               uint64 = 62200
	updateRegistrationGas     uint64 = 62200
	cancelRegisterGas         uint64 = 83200
	rewardGas                 uint64 = 83200
	calcRewardGasPerPage      uint64 = 200
	maxRewardCount            uint64 = 150000000
	voteGas                   uint64 = 62000
	cancelVoteGas             uint64 = 62000
	pledgeGas                 uint64 = 21000
	cancelPledgeGas           uint64 = 103400
	createConsensusGroupGas   uint64 = 62200
	cancelConsensusGroupGas   uint64 = 83200
	reCreateConsensusGroupGas uint64 = 62200
	mintageGas                uint64 = 83200
	mintageCancelPledgeGas    uint64 = 83200

	minPledgeHeight uint64 = 3600 * 24 * 3 // Minimum pledge height

	cgNodeCountMin                   uint8  = 3       // Minimum node count of consensus group
	cgNodeCountMax                   uint8  = 101     // Maximum node count of consensus group
	cgIntervalMin                    int64  = 1       // Minimum interval of consensus group in second
	cgIntervalMax                    int64  = 10 * 60 // Maximum interval of consensus group in second
	cgPerCountMin                    int64  = 1
	cgPerCountMax                    int64  = 10 * 60
	cgPerIntervalMin                 int64  = 1
	cgPerIntervalMax                 int64  = 10 * 60
	createConsensusGroupPledgeHeight uint64 = 3600 * 24 * 3

	rewardHeightLimit uint64 = 60 * 30 // Get snapshot block reward of 30 minutes before current
	dbPageSize        uint64 = 10000   // Batch get snapshot blocks from vm database to calc snapshot block reward

	tokenNameLengthMax   int    = 40 // Maximum length of a token name(include)
	tokenSymbolLengthMax int    = 10 // Maximum length of a token symbol(include)
	mintagePledgeHeight  uint64 = 3600 * 24 * 30 * 3
)

var (
	viteTotalSupply                  = new(big.Int).Mul(big.NewInt(1e9), util.AttovPerVite)
	rewardPerBlock                   = new(big.Int).Div(viteTotalSupply, big.NewInt(1051200000)) // Reward pre snapshot block, rewardPreBlock * blockNumPerYear / viteTotalSupply = 3%
	pledgeAmountMin                  = new(big.Int).Mul(big.NewInt(10), util.AttovPerVite)
	mintageFee                       = new(big.Int).Mul(big.NewInt(1e3), util.AttovPerVite) // Mintage cost choice 1, destroy ViteToken
	mintagePledgeAmount              = new(big.Int).Mul(big.NewInt(1e5), util.AttovPerVite) // Mintage cost choice 2, pledge ViteToken for 3 month
	createConsensusGroupPledgeAmount = new(big.Int).Mul(big.NewInt(1000), util.AttovPerVite)
)