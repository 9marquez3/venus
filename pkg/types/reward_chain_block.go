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

type PrivateKey struct {
	Data [32]byte
}

func (p *PrivateKey) Bytes() []byte {
	return p.Data[:]
}

func (p *PrivateKey) GetBlsPrivate() *bls.PrivateKey {
	key := bls.KeyFromBytes(p.Bytes())
	return &key
}

func NewPrivateKeyFromBytes(b []byte) *PrivateKey {
	var d [32]byte
	copy(d[:], b)
	return &PrivateKey{d}
}

type PublicKey struct {
	Data [48]byte
}

func (p *PublicKey) Bytes() []byte {
	return p.Data[:]
}

func (p *PublicKey) GetBlsPublic() *bls.PublicKey {
	key, _ := bls.NewPublicKey(p.Bytes())
	return &key
}

func NewPublicKeyFromBytes(b []byte) *PublicKey {
	var d [48]byte
	copy(d[:], b)
	return &PublicKey{d}
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
	PoolPublicKey *PublicKey
	// pool_contract_puzzle_hash: Optional[bytes32]
	PoolContractPuzzleHash *HashData
	// plot_public_key: G1Element
	PlotPublicKey *PublicKey
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

func (p ProofOfSpace) CalculatePlotIdPk(poolContractPuzzleHash, plotPublicKey *PublicKey) []byte {
	return bls.CalculatePlotIdPk(*poolContractPuzzleHash.GetBlsPublic(), *plotPublicKey.GetBlsPublic())
}

func (p ProofOfSpace) CalculatePlotIdPh(poolContractPuzzleHash *HashData, plotPublicKey *PublicKey) []byte {
	return bls.CalculatePlotIdPh(poolContractPuzzleHash.GetHashDate(), *plotPublicKey.GetBlsPublic())
}

func (p ProofOfSpace) GeneratePlotPublicKey(localPk, farmerPk *PublicKey, includeTaproot bool) *PublicKey {
	publicKey := bls.GeneratePlotPublicKey(*localPk.GetBlsPublic(), *farmerPk.GetBlsPublic(), includeTaproot)
	return NewPublicKeyFromBytes(publicKey.Bytes())
}

func (p ProofOfSpace) GenerateTaprootSk(localPk, farmerPk *PublicKey) *PrivateKey {
	privateKey := bls.GenerateTaprootSk(*localPk.GetBlsPublic(), *farmerPk.GetBlsPublic())
	return NewPrivateKeyFromBytes(privateKey.Bytes())
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
	Data [100]byte
}

func (c ClassGroupElement) FromBytes(d []byte) ClassGroupElement {
	var data ClassGroupElement
	copy(data.Data[:], d)
	return data
}

func (c ClassGroupElement) GetDefaultElement() ClassGroupElement {
	return c.FromBytes([]byte{0x08})
}

func (c ClassGroupElement) GetSize() int {
	return 100
}

// SubEpochSummary
// class SubEpochSummary(Streamable):
type SubEpochSummary struct {
	// prevSubepochSummaryHash: bytes32
	PrevSubepochSummaryHash *HashData
	// rewardChainHash: bytes32  # hash of reward chain at end of last segment
	RewardChainHash *HashData
	// numBlocksOverflow: uint8  # How many more blocks than 384*(N-1)
	NumBlocksOverflow uint8
	// newDifficulty: Optional[uint64]  # Only once per epoch (diff adjustment)
	NewDifficulty uint64
	// newSubSlotIters: Optional[uint64]  # Only once per epoch (diff adjustment)
	NewSubSlotIters uint64
}
