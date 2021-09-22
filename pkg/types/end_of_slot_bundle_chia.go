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
	InfusedChallengeChainSubSlotHash []byte
	// subepochSummaryHash: Optional[bytes32]  # Only once per sub-epoch, and one sub-epoch delayed
	SubepochSummaryHash [32]byte
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
	ChallengeChainSubSlotHash [32]byte
	// infusedChallengeChainSubSlotHash: Optional[bytes32]
	InfusedChallengeChainSubSlotHash []byte
	// deficit: uint8  # 16 or less. usually zero
	Deficit uint8
}
// SubSlotProofs
// class SubSlotProofs(Streamable):
type SubSlotProofs struct {
	// challengeChainSlotProof: VDFProof
	challengeChainSlotProof *VDFInfo
	// infusedChallengeChainSlotProof: Optional[VDFProof]
	infusedChallengeChainSlotProof *VDFInfo
	// rewardChainSlotProof: VDFProof
	rewardChainSlotProof *VDFInfo
}

