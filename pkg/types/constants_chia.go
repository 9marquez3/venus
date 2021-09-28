package types

type ConsensusConstants struct {
	// SlotBlocksTarget: uint32  # How many blocks to target per sub-slot
	SlotBlocksTarget uint64
	// MinBlocksPerChallengeBlock: uint8  # How many blocks must be created per slot (to make challenge sb)
	MinBlocksPerChallengeBlock uint8
	// # Max number of blocks that can be infused into a sub-slot.
	// # Note: this must be less than SubEpochBlocks/2, and > SlotBlocksTarget
	// MaxSubSlotBlocks: uint32
	MaxSubSlotBlocks uint64
	// NumSpsSubSlot: uint32  # The number of signage points per sub-slot (including the 0th sp at the sub-slot start)
	NumSpsSubSlot uint64
	// SubSlotItersStarting: uint64  # The sub_slot_iters for the first epoch
	SubSlotItersStarting uint64
	// DifficultyConstantFactor: uint128  # Multiplied by the difficulty to get iterations
	DifficultyConstantFactor BigInt
	// DifficultyStarting: uint64  # The difficulty for the first epoch
	DifficultyStarting uint64
	// # The maximum factor by which difficulty and sub_slot_iters can change per epoch
	// DifficultyChangeMaxFactor: uint32
	DifficultyChangeMaxFactor uint64
	// SubEpochBlocks: uint32  # The number of blocks per sub-epoch
	SubEpochBlocks uint64
	// EpochBlocks: uint32  # The number of blocks per sub-epoch, must be a multiple of SubEpochBlocks
	EpochBlocks uint64
	// SignificantBits: int  # The number of bits to look at in difficulty and min iters. The rest are zeroed
	SignificantBits int64
	// DiscriminantSizeBits: int  # Max is 1024 (based on ClassGroupElement int size)
	DiscriminantSizeBits int64
	// NumberZeroBitsPlotFilter: int  # H(plot id + challenge hash + signage point) must start with these many zeroes
	NumberZeroBitsPlotFilter int64
	// MinPlotSize: int
	MinPlotSize int64
	// MaxPlotSize: int
	MaxPlotSize int64
	// SubSlotTimeTarget: int  # The target number of seconds per sub-slot
	SubSlotTimeTarget int64
	// NumSpIntervalsExtra: int  # The difference between signage point and infusion point (plus required_iters)
	NumSpIntervalsExtra int64
	// MaxFutureTime: int  # The next block can have a timestamp of at most these many seconds more
	MaxFutureTime int64
	// NumberOfTimestamps: int  # Than the average of the last NumberOfTimestamps blocks
	NumberOfTimestamps int64
	// # Used as the initial cc rc challenges, as well as first block back pointers, and first SES back pointer
	// # We override this value based on the chain being run (testnet0, testnet1, mainnet, etc)
	// GenesisChallenge: bytes32
	GenesisChallenge [32]byte
	// # Forks of chia should change this value to provide replay attack protection
	// AggSigMeAdditionalData: bytes
	AggSigMeAdditionalData []byte
	// GenesisPreFarmPoolPuzzleHash: bytes32  # The block at height must pay out to this pool puzzle hash
	GenesisPreFarmPoolPuzzleHash [32]byte
	// GenesisPreFarmFarmerPuzzleHash: bytes32  # The block at height must pay out to this farmer puzzle hash
	GenesisPreFarmFarmerPuzzleHash [32]byte
	// MaxVdfWitnessSize: int  # The maximum number of classgroup elements within an n-wesolowski proof
	MaxVdfWitnessSize int64
	// # Size of mempool = 10x the size of block
	// MempoolBlockBuffer: int
	MempoolBlockBuffer int64
	// # Max coin amount uint(1 << 64). This allows coin amounts to fit in 64 bits. This is around 18M chia.
	// MaxCoinAmount: int
	MaxCoinAmount int64
	// # Max block cost in clvm cost units
	// MaxBlockCostClvm: int
	MaxBlockCostClvm int64
	// # Cost per byte of generator program
	// CostPerByte: int
	CostPerByte int64
	// WeightProofThreshold: uint8
	WeightProofThreshold uint8
	// WeightProofRecentBlocks: uint32
	WeightProofRecentBlocks uint64
	// MaxBlockCountPerRequests: uint32
	MaxBlockCountPerRequests uint64
	// RustConditionChecker: uint64
	RustConditionChecker uint64
	// BlocksCacheSize: uint32
	BlocksCacheSize uint64
	// NetworkType: int
	NetworkType int64
	// MaxGeneratorSize: uint32
	MaxGeneratorSize uint64
	// MaxGeneratorRefListSize: uint32
	MaxGeneratorRefListSize uint64
	// PoolSubSlotIters: uint64
	PoolSubSlotIters uint64
}
