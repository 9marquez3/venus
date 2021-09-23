package types

type EndOfSubSlotBundle struct {
	// challengeChain: ChallengeChainSubSlot
	ChallengeChain *ChallengeChainSubSlot
	// infusedChallengeChain: Optional[InfusedChallengeChainSubSlot]
	InfusedChallengeChain *InfusedChallengeChainSubSlot
	// rewardChain: RewardChainSubSlot
	RewardChain *RewardChainSubSlot
	// proofs: SubSlotProofs
	Proofs *SubSlotProofs
}

type ChallengeChainSubSlot struct {
	// ChallengeChainEndOfSlotVdf: VDFInfo
	ChallengeChainEndOfSlotVdf *VDFInfo
	// infusedChallengeChainSubSlotHash: Optional[bytes32]  # Only at the end of a slot
	InfusedChallengeChainSubSlotHash *HashData
	// subepochSummaryHash: Optional[bytes32]  # Only once per sub-epoch, and one sub-epoch delayed
	SubepochSummaryHash *HashData
	// newSubSlotIters: Optional[uint64]  # Only at the end of epoch, sub-epoch, and slot
	NewSubSlotIters *uint64
	// newDifficulty: Optional[uint64]  # Only at the end of epoch, sub-epoch, and slot
	NewDifficulty *uint64
}

type InfusedChallengeChainSubSlot struct {
	// InfusedChallengeChainEndOfSlotVdf: VDFInfo
	InfusedChallengeChainEndOfSlotVdf *VDFInfo
}

// RewardChainSubSlot
// class RewardChainSubSlot(Streamable):
type RewardChainSubSlot struct {
	// endOfSlotVdf: VDFInfo
	EndOfSlotVdf *VDFInfo
	// challengeChainSubSlotHash: bytes32
	ChallengeChainSubSlotHash *HashData
	// infusedChallengeChainSubSlotHash: Optional[bytes32]
	InfusedChallengeChainSubSlotHash *HashData
	// deficit: uint8  # 16 or less. usually zero
	Deficit uint8
}

// SubSlotProofs
// class SubSlotProofs(Streamable):
type SubSlotProofs struct {
	// challengeChainSlotProof: VDFProof
	challengeChainSlotProof *VDFProof
	// infusedChallengeChainSlotProof: Optional[VDFProof]
	infusedChallengeChainSlotProof *VDFProof
	// rewardChainSlotProof: VDFProof
	rewardChainSlotProof *VDFProof
}
