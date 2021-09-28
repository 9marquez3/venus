package types

// BlockRecord
//     This class is not included or hashed into the blockchain, but it is kept in memory as a more
//     efficient way to maintain data about the blockchain. This allows us to validate future blocks,
//     difficulty adjustments, etc, without saving the whole header block in memory.
// class BlockRecord(Streamable):
type BlockRecord struct {
	// header_hash: bytes32
	HeaderHash *HashData
	// prev_hash: bytes32  # Header hash of the previous block
	PrevHash *HashData
	// height: uint32
	Height uint32
	// weight: uint128  # Total cumulative difficulty of all ancestor blocks since genesis
	Weight BigInt
	// total_iters: uint128  # Total number of VDF iterations since genesis, including this block
	TotalIters BigInt
	// signage_point_index: uint8
	SignagePointIndex uint8
	// challenge_vdf_output: ClassGroupElement  # This is the intermediary VDF output at ip_iters in challenge chain
	ChallengeVdfOutput *ClassGroupElement
	// infused_challenge_vdf_output: Optional[
	//     ClassGroupElement
	// ]  # This is the intermediary VDF output at ip_iters in infused cc, iff deficit <= 3
	InfusedChallengeVdfOutput *ClassGroupElement
	// reward_infusion_new_challenge: bytes32  # The reward chain infusion output, input to next VDF
	RewardInfusionNewChallenge []byte
	// challenge_block_info_hash: bytes32  # Hash of challenge chain data, used to validate end of slots in the future
	ChallengeBlockInfoHash *HashData
	// sub_slot_iters: uint64  # Current network sub_slot_iters parameter
	SubSlotIters uint64
	// pool_puzzle_hash: bytes32  # Need to keep track of these because Coins are created in a future block
	PoolPuzzleHash *HashData
	// farmer_puzzle_hash: bytes32
	FarmerPuzzleHash *HashData
	// required_iters: uint64  # The number of iters required for this proof of space
	RequiredIters uint64
	// deficit: uint8  # A deficit of 16 is an overflow block after an infusion. Deficit of 15 is a challenge block
	Deficit uint8
	// overflow: bool
	Overflow bool
	// prev_transaction_block_height: uint32
	PrevTransactionBlockHeight uint32
	// # Transaction block (present iff isTransactionBlock)
	// timestamp: Optional[uint64]
	Timestamp uint64
	// prev_transaction_block_hash: Optional[bytes32]  # Header hash of the previous transaction block
	PrevTransactionBlockHash *HashData
	// fees: Optional[uint64]
	Fees uint64
	// reward_claims_incorporated: Optional[List[Coin]]
	RewardClaimsIncorporated []*Coin
	// # Slot (present iff this is the first SB in sub slot)
	// finished_challenge_slot_hashes: Optional[List[bytes32]]
	FinishedChallengeSlotHashes []*HashData
	// finished_infused_challenge_slot_hashes: Optional[List[bytes32]]
	FinishedInfusedChallengeSlotHashes []*HashData
	// finished_reward_slot_hashes: Optional[List[bytes32]]
	FinishedRewardSlotHashes []*HashData
	// # Sub-epoch (present iff this is the first SB after sub-epoch)
	// sub_epoch_summary_included: Optional[SubEpochSummary]
	SubEpochSummaryIncluded *SubEpochSummary
}

// IsTransactionBlock
// @property
// def isTransactionBlock(self) -> bool:
//     return self.timestamp is not None
func (b *BlockRecord) IsTransactionBlock() bool {
	return b.Timestamp != 0
}

// FirstInSubSlot
// @property
// def firstInSubSlot(self) -> bool:
//     return self.finished_challenge_slot_hashes is not None
func (b *BlockRecord) FirstInSubSlot() bool {
	return b.FinishedChallengeSlotHashes != nil
}

// IsChallengeBlock
// def isChallengeBlock(self, constants: ConsensusConstants) -> bool:
//     return self.deficit == constants.MIN_BLOCKS_PER_CHALLENGE_BLOCK - 1
func (b *BlockRecord) IsChallengeBlock(constants *ConsensusConstants) bool {
	return b.Deficit == constants.MinBlocksPerChallengeBlock-1
}

// SpSubSlotTotalIters
// def spSubSlotTotalIters(self, constants: ConsensusConstants) -> uint128:
//     if self.overflow:
//         return uint128(self.total_iters - self.IpIters(constants) - self.sub_slot_iters)
//     else:
//         return uint128(self.total_iters - self.IpIters(constants))
func (b *BlockRecord) SpSubSlotTotalIters(constants *ConsensusConstants) BigInt {
	if b.Overflow {
		return BigSub(BigSub(b.TotalIters, NewInt(b.IpIters(constants))), NewInt(b.SubSlotIters))
	}
	return BigSub(b.TotalIters, NewInt(b.IpIters(constants)))
}

// IpSubSlotTotalIters
// def ipSubSlotTotalIters(self, constants: ConsensusConstants) -> uint128:
//     return uint128(self.total_iters - self.IpIters(constants))
func (b *BlockRecord) IpSubSlotTotalIters(constants *ConsensusConstants) BigInt {
	return BigSub(b.TotalIters, NewInt(b.IpIters(constants)))
}

// SpIters
// def spIters(self, constants: ConsensusConstants) -> uint64:
//     return CalculateSpIters(constants, self.sub_slot_iters, self.signage_point_index)
func (b *BlockRecord) SpIters(constants *ConsensusConstants) uint64 {
	return CalculateSpIters(constants, b.SubSlotIters, uint64(b.SignagePointIndex))
}

// IpIters
// def ip_iters(self, constants: ConsensusConstants) -> uint64:
//     return calculate_ip_iters( constants,self.sub_slot_iters,self.signage_point_index,self.required_iters)
func (b *BlockRecord) IpIters(constants *ConsensusConstants) uint64 {
	return CalculateSpIters(constants, b.SubSlotIters, uint64(b.SignagePointIndex))
}

// SpTotalIters
// def spTotalIters(self, constants: ConsensusConstants):
//     return self.sp_sub_slot_total_iters(constants) + self.sp_iters(constants)
func (b *BlockRecord) SpTotalIters(constants *ConsensusConstants) BigInt {
	return BigAdd(b.SpSubSlotTotalIters(constants), NewInt(b.SpIters(constants)))
}
