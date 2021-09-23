package types

// class FullBlock(Streamable):
type FullBlockChia struct {
	// # All the information required to validate a block
	// FinishedSubSlots: List[EndOfSubSlotBundle]  # If first sb
	FinishedSubSlots []*EndOfSubSlotBundle
	// RewardChainBlock: RewardChainBlock  # Reward chain trunk data
	RewardChainBlock *RewardChainBlock
	// challengeChainSpProof: Optional[VDFProof]  # If not first sp in sub-slot
	ChallengeChainSpProof *VDFProof
	// challengeChainIpProof: VDFProof
	ChallengeChainIpProof *VDFProof
	// rewardChainSpProof: Optional[VDFProof]  # If not first sp in sub-slot
	RewardChainSpProof *VDFProof
	// rewardChainIpProof: VDFProof
	RewardChainIpProof *VDFProof
	// infusedChallengeChainIpProof: Optional[VDFProof]  # Iff deficit < 4
	InfusedChallengeChainIpProof *VDFProof
	// foliage: Foliage  # Reward chain foliage data
	Foliage *Foliage
	// foliageTransactionBlock: Optional[FoliageTransactionBlock]  # Reward chain foliage data (tx block)
	FoliageTransactionBlock *FoliageTransactionBlock
	// transactionsInfo: Optional[TransactionsInfo]  # Reward chain foliage data (tx block additional)
	TransactionsInfo *TransactionsInfo
	// transactionsGenerator: Optional[SerializedProgram]  # Program that generates transactions
	TransactionsGenerator *SerializedProgram
	// transactionsGeneratorRefList: List[
	// uint32
	// ]  # List of block heights of previous generators referenced in this block
	TransactionsGeneratorRefList []uint64
}

type SerializedProgram struct {
}

func (f *FullBlockChia) PrevHeaderHash() *HashData {
	return f.Foliage.PrevBlockHash
}

func (f *FullBlockChia) Height() uint64 {
	return f.RewardChainBlock.Height
}

func (f *FullBlockChia) Weight() BigInt {
	return f.RewardChainBlock.Weight
}

func (f *FullBlockChia) TotalIters() BigInt {
	return f.RewardChainBlock.TotalIters
}

func (f *FullBlockChia) HeaderHash() *HashData {
	return nil
}

func (f *FullBlockChia) IsTransactionBlock() bool {
	return f.RewardChainBlock != nil
}

func (f *FullBlockChia) GetIncludedRewardCoins() []*Coin {
	if !f.IsTransactionBlock() {
		return nil
	}
	if f.TransactionsInfo == nil {
		panic("TransactionsInfo is nil")
	}
	// todo list to set
	return f.TransactionsInfo.RewardClaimsInCorporated
}

func (f *FullBlockChia) IsFullyCompactified() bool {
	for _, subSlot := range f.FinishedSubSlots {
		if subSlot.Proofs.challengeChainSlotProof.WitnessType != 0 || !subSlot.Proofs.challengeChainSlotProof.NormalizedToIdentity {
			return false
		}
		if subSlot.Proofs.infusedChallengeChainSlotProof != nil && (subSlot.Proofs.infusedChallengeChainSlotProof.WitnessType != 0 ||
			!subSlot.Proofs.infusedChallengeChainSlotProof.NormalizedToIdentity) {
			return false
		}
	}
	if f.ChallengeChainSpProof != nil && (f.ChallengeChainSpProof.WitnessType != 0 || !f.ChallengeChainSpProof.NormalizedToIdentity) {
		return false
	}
	if f.ChallengeChainIpProof.WitnessType != 0 || !f.ChallengeChainIpProof.NormalizedToIdentity {
		return false
	}

	return true
}
