package types

import (
	"fmt"

	bls "github.com/cnc-project/cnc-bls"
	"github.com/filecoin-project/go-state-types/big"
)

// IsOverflowBlock
// def isOverflowBlock(constants: ConsensusConstants, signage_point_index: uint8) -> bool:
//     if signage_point_index >= constants.NUM_SPS_SUB_SLOT:
//         raise ValueError("SP index too high")
//     return signage_point_index >= constants.NUM_SPS_SUB_SLOT - constants.NUM_SP_INTERVALS_EXTRA
func IsOverflowBlock(constants *ConsensusConstants, signagePointIndex uint64) bool {
	if signagePointIndex >= constants.NumSpsSubSlot {
		panic("SP index too high")
	}
	return signagePointIndex >= constants.NumSpsSubSlot-uint64(constants.NumSpIntervalsExtra)
}

// CalculateSpIntervalIters
// def calculateSpIntervalIters(constants: ConsensusConstants, sub_slot_iters: uint64) -> uint64:
//     assert sub_slot_iters % constants.NUM_SPS_SUB_SLOT == 0
//     return uint64(sub_slot_iters // constants.NUM_SPS_SUB_SLOT)
func CalculateSpIntervalIters(constants *ConsensusConstants, subSlotIters uint64) uint64 {
	if subSlotIters%constants.NumSpsSubSlot != 0 {
		panic("subSlotIters % constants.NumSpsSubSlot != 0")
	}
	return subSlotIters / constants.NumSpsSubSlot
}

// CalculateSpIters
// def CalculateSpIters(constants: ConsensusConstants, sub_slot_iters: uint64, signage_point_index: uint8) -> uint64:
//     if signage_point_index >= constants.NUM_SPS_SUB_SLOT:
//         raise ValueError("SP index too high")
//     return uint64(calculate_sp_interval_iters(constants, sub_slot_iters) * signage_point_index)
func CalculateSpIters(constants *ConsensusConstants, subSlotIters uint64, signagePointIndex uint64) uint64 {
	if signagePointIndex >= constants.NumSpsSubSlot {
		panic("SP index too high")
	}
	return CalculateSpIntervalIters(constants, subSlotIters) * signagePointIndex
}

// CalculateIpIters
// def calculateIpIters(
//     constants: ConsensusConstants,
//     sub_slot_iters: uint64,
//     signage_point_index: uint8,
//     required_iters: uint64,
// ) -> uint64:
//     # Note that the SSI is for the block passed in, which might be in the previous epoch
//     sp_iters = CalculateSpIters(constants, sub_slot_iters, signage_point_index)
//     sp_interval_iters: uint64 = calculate_sp_interval_iters(constants, sub_slot_iters)
//     if sp_iters % sp_interval_iters != 0 or sp_iters >= sub_slot_iters:
//         raise ValueError(f"Invalid sp iters {sp_iters} for this ssi {sub_slot_iters}")
//
//     if required_iters >= sp_interval_iters or required_iters == 0:
//         raise ValueError(
//             f"Required iters {required_iters} is not below the sp interval iters {sp_interval_iters} "
//             f"{sub_slot_iters} or not >0."
//         )
//
//     return uint64((sp_iters + constants.NUM_SP_INTERVALS_EXTRA * sp_interval_iters + required_iters) % sub_slot_iters)
func CalculateIpIters(constants *ConsensusConstants, subSlotIters uint64, signagePointIndex uint8, requiredIters uint64) uint64 {
	spIters := CalculateSpIters(constants, subSlotIters, uint64(signagePointIndex))
	spIntervalIters := CalculateSpIntervalIters(constants, subSlotIters)
	if spIters%spIntervalIters != 0 || spIters >= subSlotIters {
		panic(fmt.Sprintf("invalid spiters %d for this ssi %d", spIters, spIntervalIters))
	}

	if requiredIters >= spIntervalIters || requiredIters == 0 {
		panic(fmt.Sprintf("Required iters %d is not below the sp interval iters %d；subSlotIters (%d) or not >0. ",
			requiredIters, spIntervalIters, subSlotIters))
	}
	return (spIters + uint64(constants.NumSpIntervalsExtra)*spIntervalIters + requiredIters) % subSlotIters
}

// CalculateIterationsQuality
//     Calculates the number of iterations from the quality. This is derives as the difficulty times the constant factor
//     times a random number between 0 and 1 (based on quality string), divided by plot size.
func CalculateIterationsQuality(difficultyConstantFactor BigInt, qualityString []byte, size int64, difficulty uint64, ccSpOutputHash *HashData) uint64 {
	//     sp_quality_string: bytes32 = std_hash(quality_string + cc_sp_output_hash)
	spQualityString := bls.Hash256(append(qualityString, ccSpOutputHash.Bytes()...))
	// pow(2, 256)
	pow := big.Exp(NewInt(2), NewInt(256))
	//     iters = uint64(
	//         int(difficulty) * int(difficulty_constant_factor) * int.from_bytes(sp_quality_string, "big", signed=False) / (int(pow(2, 256)) * int(_expected_plot_size(size)))
	//     )
	iters := BigMul(BigDiv(BigMul(BigMul(NewInt(difficulty), difficultyConstantFactor), BigFromBytes(spQualityString)), pow), ExpectedPlotSize(size))
	//     return max(iters, uint64(1))
	return big.Max(iters, NewInt(1)).Uint64()
}

// ExpectedPlotSize
// Given the plot size parameter k (which is between 32 and 59), computes the
// expected size of the plot in bytes (times a constant factor). This is based on efficient encoding
// of the plot, and aims to be scale agnostic, so larger plots don't
// necessarily get more rewards per byte. The +1 is added to give half a bit more space per entry, which
// is necessary to store the entries in the plot.
func ExpectedPlotSize(k int64) BigInt {
	// return ((2 * k) + 1) * (2 ** (k - 1))
	return BigMul(NewInt(uint64(2*k+1)), big.Exp(NewInt(2), NewInt(uint64(k-1))))
}
