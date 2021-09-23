package types

import (
	"bytes"
	"encoding/hex"
	"fmt"

	bls "github.com/cnc-project/cnc-bls"
)

var ToHashBytes32 = bls.HashDigest256FromBytes

type Signature struct {
	Data [96]byte
}

func (s *Signature) Bytes() []byte {
	return s.Data[:]
}

func (s *Signature) Bytes96() [96]byte {
	return s.Data
}

func NewSignFromBytes(b []byte) *Signature {
	var d [96]byte
	copy(d[:], b)
	return &Signature{d}
}

type HashData struct {
	Data bls.HashDigest256
}

func (h *HashData) GetHashDate() bls.HashDigest256 {
	return h.Data
}

func (h *HashData) IsZero() bool {
	return h.Data.IsZero()
}

func (h *HashData) Bytes() []byte {
	return h.Data.Bytes()
}

func NewHashDataFromBytes(b []byte) *HashData {
	return &HashData{ToHashBytes32(b)}
}

type RewardChainBlockUnfinished struct {
	// total_iters: uint128
	TotalIters BigInt
	// signage_point_index: uint8
	SignagePointIndex uint8
	// pos_ss_cc_challenge_hash: bytes32
	PosSsCcChallengeHash *HashData
	// proof_of_space: ProofOfSpace
	ProofOfSpace *ProofOfSpace
	// challenge_chain_sp_vdf: Optional[VDFInfo]  # Not present for first sp in slot
	ChallengeChainSpVdf *VDFInfo
	// challenge_chain_sp_signature: G2Element
	ChallengeChainSpSignature *Signature
	// reward_chain_sp_vdf: Optional[VDFInfo]  # Not present for first sp in slot
	RewardChainSpVdf *VDFInfo
	// reward_chain_sp_signature: G2Element
	RewardChainSpSignature *Signature
}

type RewardChainBlock struct {
	// weight: uint128
	Weight BigInt
	// height: uint32
	Height uint64
	// total_iters: uint128
	TotalIters BigInt
	// signage_point_index: uint8
	SignagePointIndex uint8
	// pos_ss_cc_challenge_hash: bytes32
	PosSsCcChallengeHash *HashData
	// proof_of_space: ProofOfSpace
	ProofOfSpace *ProofOfSpace
	// challenge_chain_sp_vdf: Optional[VDFInfo]  # Not present for first sp in slot
	ChallengeChainSpVdf *VDFInfo
	// challenge_chain_sp_signature: G2Element
	ChallengeChainSpSignature *Signature
	// challenge_chain_ip_vdf: VDFInfo
	ChallengeChainIpVdf *VDFInfo
	// reward_chain_sp_vdf: Optional[VDFInfo]  # Not present for first sp in slot
	RewardChainSpVdf *VDFInfo
	// reward_chain_sp_signature: G2Element
	RewardChainSpSignature *Signature
	// reward_chain_ip_vdf: VDFInfo
	RewardChainIpVdf *VDFInfo
	// infused_challenge_chain_ip_vdf: Optional[VDFInfo]  # Iff deficit < 16
	InfusedChallengeChainIpVdf *VDFInfo
	// is_transaction_block: bool
	IsTransactionBlock bool
}

func (r RewardChainBlock) GetUnfinished() RewardChainBlockUnfinished {
	return RewardChainBlockUnfinished{
		TotalIters:                r.TotalIters,
		SignagePointIndex:         r.SignagePointIndex,
		PosSsCcChallengeHash:      r.PosSsCcChallengeHash,
		ProofOfSpace:              r.ProofOfSpace,
		ChallengeChainSpVdf:       r.ChallengeChainIpVdf,
		ChallengeChainSpSignature: r.ChallengeChainSpSignature,
		RewardChainSpVdf:          r.RewardChainIpVdf,
		RewardChainSpSignature:    r.RewardChainSpSignature,
	}
}

type ProofOfSpace struct {
	// challenge: bytes32
	Challenge [32]byte
	// pool_public_key: Optional[G1Element]  # Only one of these two should be present
	PoolPublicKey *bls.PublicKey
	// pool_contract_puzzle_hash: Optional[bytes32]
	PoolContractPuzzleHash *HashData
	// plot_public_key: G1Element
	PlotPublicKey *bls.PublicKey
	// size: uint8
	Size uint8
	// proof: bytes
	Proof []byte
}

func (p ProofOfSpace) GetPlotId() []byte {
	if p.PoolPublicKey == nil {
		return p.CalculatePlotIdPh(p.PoolContractPuzzleHash, p.PlotPublicKey)
	}
	return p.CalculatePlotIdPk(p.PoolPublicKey, p.PlotPublicKey)
}

func (p ProofOfSpace) VerifyAndGetQualityString(constants *ConsensusConstants, originalChallengeHash, signagePoint *HashData) ([]byte, error) {

	if p.PoolPublicKey == nil && p.PoolContractPuzzleHash.IsZero() {
		return nil, fmt.Errorf("fail 1")
	}
	if p.PoolPublicKey != nil && !p.PoolContractPuzzleHash.IsZero() {
		return nil, fmt.Errorf("fail 2")
	}
	if int64(p.Size) < constants.MinPlotSize {
		return nil, fmt.Errorf("fail 3")
	}
	if int64(p.Size) > constants.MaxPlotSize {
		return nil, fmt.Errorf("fail 4")
	}

	plotId := p.GetPlotId()
	plotId32 := NewHashDataFromBytes(plotId)

	challenge := p.CalculatePosChallenge(plotId32, originalChallengeHash, signagePoint)

	if !bytes.Equal(challenge, p.Challenge[:]) {
		return nil, fmt.Errorf("new challenge %s is not Equal challenge %s",
			hex.EncodeToString(challenge), hex.EncodeToString(p.Challenge[:]))
	}

	if !p.PassesPlotFilter(constants, plotId32, originalChallengeHash, signagePoint) {
		return nil, fmt.Errorf("fail 5")
	}
	return p.GetQualityString(plotId32), nil
}

func (p ProofOfSpace) GetQualityString(plotId *HashData) []byte {
	// todo C++ chiapos
	return nil
}

func (p ProofOfSpace) PassesPlotFilter(constants *ConsensusConstants, plotId, challengeHash, signagePoint *HashData) bool {
	input := p.CalculatePlotFilterInput(plotId, challengeHash, signagePoint)
	for i := int64(0); i < constants.NumberZeroBitsPlotFilter; i++ {
		if input[i] != 0 {
			return false
		}
	}
	return true
}

func (p ProofOfSpace) CalculatePlotFilterInput(plotId, challengeHash, signagePoint *HashData) []byte {
	return bls.CalculatePlotFilterInput(plotId.GetHashDate(), challengeHash.GetHashDate(), signagePoint.GetHashDate())
}

func (p ProofOfSpace) CalculatePosChallenge(plotId, challengeHash, signagePoint *HashData) []byte {
	return bls.CalculatePosChallenge(plotId.GetHashDate(), challengeHash.GetHashDate(), signagePoint.GetHashDate())
}

func (p ProofOfSpace) CalculatePlotIdPk(poolContractPuzzleHash, plotPublicKey *bls.PublicKey) []byte {
	return bls.CalculatePlotIdPk(*poolContractPuzzleHash, *plotPublicKey)
}

func (p ProofOfSpace) CalculatePlotIdPh(poolContractPuzzleHash *HashData, plotPublicKey *bls.PublicKey) []byte {
	return bls.CalculatePlotIdPh(poolContractPuzzleHash.GetHashDate(), *plotPublicKey)
}

func (p ProofOfSpace) GeneratePlotPublicKey(localPk, farmerPk *bls.PublicKey, includeTaproot bool) *bls.PublicKey {
	publicKey := bls.GeneratePlotPublicKey(*localPk, *farmerPk, includeTaproot)
	return &publicKey
}

func (p ProofOfSpace) GenerateTaprootSk(localPk, farmerPk *bls.PublicKey) *bls.PrivateKey {
	privateKey := bls.GenerateTaprootSk(*localPk, *farmerPk)
	return &privateKey
}

// VDFProof
// class VDFProof(Streamable):
type VDFProof struct {
	// witnessType: uint8
	WitnessType uint8
	// witness: bytes
	Witness []byte
	// normalizedToIdentity: bool
	NormalizedToIdentity bool
}

type VDFInfo struct {
	// challenge: bytes32  # Used to generate the discriminant (VDF group)
	Challenge [32]byte
	// number_of_iterations: uint64
	NumberOfIterations uint64
	// output: ClassgroupElement {data: bytes100}
	Output ClassGroupElement
}

type ClassGroupElement struct {
	data [100]byte
}

func (c ClassGroupElement) FromBytes(d []byte) ClassGroupElement {
	var data ClassGroupElement
	copy(data.data[:], d)
	return data
}

func (c ClassGroupElement) GetDefaultElement() ClassGroupElement {
	return c.FromBytes([]byte{0x08})
}

func (c ClassGroupElement) GetSize() int {
	return 100
}

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
