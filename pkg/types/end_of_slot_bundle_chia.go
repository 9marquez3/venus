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
	// end_of_slot_vdf: VDFInfo
	// challenge_chain_sub_slot_hash: bytes32
	// infused_challenge_chain_sub_slot_hash: Optional[bytes32]
	// deficit: uint8  # 16 or less. usually zero
}
// @dataclass(frozen=True)
// @streamable
// class SubSlotProofs(Streamable):
// challenge_chain_slot_proof: VDFProof
// infused_challenge_chain_slot_proof: Optional[VDFProof]
// reward_chain_slot_proof: VDFProof
