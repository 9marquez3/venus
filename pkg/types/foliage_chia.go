package types

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
)

type TransactionsInfo struct {
	// # Information that goes along with each transaction block
	// generatorRoot: bytes32  # sha256 of the block generator in this block
	GeneratorRoot [32]byte
	// generatorRefsRoot: bytes32  # sha256 of the concatenation of the generator ref list entries
	GeneratorRefsRoot [32]byte
	// aggregatedSignature: G2Element
	AggregatedSignature [SignatureBytes]byte
	// fees: uint64  # This only includes user fees, not block rewards
	Fees uint64
	// cost: uint64  # This is the total cost of this block, including CLVM cost, cost of program size and conditions
	Cost uint64
	// rewardClaimsIncorporated: List[Coin]  # These can be in any order
	RewardClaimsInCorporated []*Coin
}

// Coin this structure is used in the body for the reward and fees genesis coins.
type Coin struct {
	// parentCoinInfo: bytes32  # down with this sort of thing.
	ParentCoinInfo [32]byte
	// puzzleHash: bytes32
	PuzzleHash [32]byte
	// amount: uint64
	Amount uint64
}

// GetHash
// This does not use streamable format for hashing, the amount is
// serialized using CLVM integer format.
//
// Note that int_to_bytes() will prepend a 0 to integers where the most
// significant bit is set, to encode it as a positive number. This
// despite "amount" being unsigned. This way, a CLVM program can generate
// these hashes easily.
func (c *Coin) GetHash() [32]byte {
	bytes := make([]byte, 32*2)
	copy(bytes,c.ParentCoinInfo[:])
	copy(bytes[32:],c.PuzzleHash[:])
	return sha256.Sum256(append(bytes,intToBytes(c.Amount)...))
}

func (c *Coin) Name() [32]byte {
	return c.GetHash()
}


// AsList return List[self.parent_coin_info, self.puzzle_hash, self.amount -> ([]byte)]
func (c *Coin) AsList() [][]byte {
	bytes := make([][]byte, 0)
	bytes = append(bytes, c.ParentCoinInfo[:])
	bytes = append(bytes, c.PuzzleHash[:])
	d := make([]byte,8)
	binary.BigEndian.PutUint64(d,c.Amount)
	bytes = append(bytes, d)
	return bytes
}

func (c *Coin) NameStr() string {
	name := c.Name()
	return hex.EncodeToString(name[:])
}

// 返回最小的 byte 数组
func intToBytes(v uint64) []byte {
	if v == 0 {
		return nil
	}
	d := make([]byte,8)
	binary.BigEndian.PutUint64(d,v)

	for d[0] == 0 {
		d = d[1:]
	}

	data := make([]byte,len(d))
	copy(data,d)
	return data
}

// Foliage
// The entire foliage block, containing signature and the unsigned back pointer
//
// The hash of this is the "header hash". Note that for unfinished blocks, the prev_block_hash
//
// Is the prev from the signage point, and can be replaced with a more recent block
type Foliage struct {
	// prevBlockHash: bytes32
	PrevBlockHash [32]byte
	// rewardBlockHash: bytes32
	RewardBlockHash [32]byte
	// foliageBlockData: FoliageBlockData
	FoliageBlockData *FoliageBlockData
	// foliageBlockDataSignature: G2Element
	FoliageBlockDataSignature [SignatureBytes]byte
	// foliageTransactionBlockHash: Optional[bytes32]
	FoliageTransactionBlockHash [32]byte
	// foliageTransactionBlockSignature: Optional[G2Element] (optional use golang pointer)
	FoliageTransactionBlockSignature []byte
}

// FoliageBlockData
// Part of the block that is signed by the plot key
type FoliageBlockData struct {
	// unfinishedRewardBlockHash: bytes32
	UnfinishedRewardBlockHash [32]byte
	// poolTarget: PoolTarget
	PoolTarget *PoolTarget
	// poolSignature: Optional[G2Element]  # Iff ProofOfSpace has a pool pk (optional use golang pointer)
	PoolSignature []byte
	// farmerRewardPuzzleHash: bytes32
	FarmerRewardPuzzleHash [32]byte
	// extensionData: bytes32  # Used for future updates. Can be any 32 byte value initially
	ExtensionData [32]byte
}

type PoolTarget struct {
	// puzzleHash: bytes32
	PuzzleHash [32]byte
	// maxHeight: uint32  # A max height of 0 means it is valid forever
	MaxHeight uint32
}