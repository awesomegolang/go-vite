package chain

import (
	"bytes"
	"github.com/vitelabs/go-vite/vm/contracts/abi"
	"math/big"

	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/vm/quota"
	"github.com/vitelabs/go-vite/vm_context"
)

func (c *chain) GetContractGidByAccountBlock(block *ledger.AccountBlock) (*types.Gid, error) {
	if block == nil {
		return nil, nil
	}

	return c.GetContractGid(&block.AccountAddress)
}

// TODO cache
func (c *chain) GetContractGid(addr *types.Address) (*types.Gid, error) {
	if addr == nil {
		return nil, nil
	}

	if bytes.Equal(addr.Bytes(), abi.AddressRegister.Bytes()) ||
		bytes.Equal(addr.Bytes(), abi.AddressVote.Bytes()) ||
		bytes.Equal(addr.Bytes(), abi.AddressPledge.Bytes()) ||
		bytes.Equal(addr.Bytes(), abi.AddressConsensusGroup.Bytes()) ||
		bytes.Equal(addr.Bytes(), abi.AddressMintage.Bytes()) {
		return &types.DELEGATE_GID, nil
	}

	account, getAccountErr := c.chainDb.Account.GetAccountByAddress(addr)
	if getAccountErr != nil {

		c.log.Error("Query account failed. Error is "+getAccountErr.Error(), "method", "GetContractGid")

		return nil, getAccountErr
	}
	if account == nil {
		return nil, nil
	}

	gid, err := c.chainDb.Ac.GetContractGid(account.AccountId)

	if err != nil {
		c.log.Error("GetContractGid failed, error is "+err.Error(), "method", "GetContractGid")
		return nil, err
	}

	return gid, nil
}

func (c *chain) GetPledgeQuotas(snapshotHash types.Hash, beneficialList []types.Address) (map[types.Address]uint64, error) {
	pledgeDb, err := vm_context.NewVmContext(c, &snapshotHash, nil, nil)
	if err != nil {
		c.log.Error("NewVmContext failed, error is "+err.Error(), "method", "GetPledgeQuotas")
		return nil, err
	}
	quotas := make(map[types.Address]uint64)
	for _, addr := range beneficialList {
		balanceDb, err := vm_context.NewVmContext(c, &snapshotHash, nil, &addr)
		if err != nil {
			c.log.Error("NewVmContext failed, error is "+err.Error(), "method", "GetPledgeQuotas")
			return nil, err
		}
		pledgeAmount := abi.GetPledgeBeneficialAmount(pledgeDb, addr)
		quotas[addr] = quota.GetPledgeQuota(balanceDb, addr, pledgeAmount)
	}
	return quotas, nil
}
func (c *chain) GetPledgeQuota(snapshotHash types.Hash, beneficial types.Address) (uint64, error) {
	vmContext, err := vm_context.NewVmContext(c, &snapshotHash, nil, &beneficial)
	if err != nil {
		c.log.Error("NewVmContext failed, error is "+err.Error(), "method", "GetPledgeQuota")
		return 0, err
	}
	pledgeAmount := abi.GetPledgeBeneficialAmount(vmContext, beneficial)
	return quota.GetPledgeQuota(vmContext, beneficial, pledgeAmount), nil
}

func (c *chain) GetRegisterList(snapshotHash types.Hash, gid types.Gid) ([]*types.Registration, error) {
	vmContext, err := vm_context.NewVmContext(c, &snapshotHash, nil, nil)
	if err != nil {
		c.log.Error("NewVmContext failed, error is "+err.Error(), "method", "GetCandidateList")
		return nil, err
	}
	return abi.GetCandidateList(vmContext, gid, nil), nil
}

func (c *chain) GetVoteMap(snapshotHash types.Hash, gid types.Gid) ([]*types.VoteInfo, error) {
	vmContext, err := vm_context.NewVmContext(c, &snapshotHash, nil, nil)
	if err != nil {
		c.log.Error("NewVmContext failed, error is "+err.Error(), "method", "GetVoteList")
		return nil, err
	}
	return abi.GetVoteList(vmContext, gid, nil), nil
}

func (c *chain) GetPledgeAmount(snapshotHash types.Hash, beneficial types.Address) (*big.Int, error) {
	vmContext, err := vm_context.NewVmContext(c, &snapshotHash, nil, nil)
	if err != nil {
		c.log.Error("NewVmContext failed, error is "+err.Error(), "method", "GetPledgeBeneficialAmount")
		return nil, err
	}
	return abi.GetPledgeBeneficialAmount(vmContext, beneficial), nil
}

func (c *chain) GetConsensusGroupList(snapshotHash types.Hash) ([]*types.ConsensusGroupInfo, error) {
	vmContext, err := vm_context.NewVmContext(c, &snapshotHash, nil, nil)
	if err != nil {
		c.log.Error("NewVmContext failed, error is "+err.Error(), "method", "GetActiveConsensusGroupList")
		return nil, err
	}
	return abi.GetActiveConsensusGroupList(vmContext, nil), nil
}

func (c *chain) GetBalanceList(snapshotHash types.Hash, tokenTypeId types.TokenTypeId, addressList []types.Address) (map[types.Address]*big.Int, error) {
	vmContext, err := vm_context.NewVmContext(c, &snapshotHash, nil, nil)
	if err != nil {
		c.log.Error("NewVmContext failed, error is "+err.Error(), "method", "GetBalanceList")
		return nil, err
	}

	balanceList := make(map[types.Address]*big.Int)
	for _, addr := range addressList {
		balanceList[addr] = vmContext.GetBalance(&addr, &tokenTypeId)
	}
	return balanceList, nil
}

func (c *chain) GetTokenInfoById(tokenId *types.TokenTypeId) (*types.TokenInfo, error) {
	vmContext, err := vm_context.NewVmContext(c, nil, nil, &abi.AddressMintage)
	if err != nil {
		c.log.Error("NewVmContext failed, error is "+err.Error(), "method", "GetTokenInfoById")
		return nil, err
	}
	return abi.GetTokenById(vmContext, *tokenId), nil
}
