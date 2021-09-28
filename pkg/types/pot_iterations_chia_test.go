package types

import (
	"encoding/hex"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func hexString32(s string) [32]byte {
	if len(s) != 64 {
		panic("len s != 64")
	}
	bytes, _ := hex.DecodeString(s)
	a := [32]byte{}
	copy(a[:], bytes)
	return a
}

func hexString(s string) []byte {
	if len(s) != 64 {
		panic("len s != 64")
	}
	bytes, _ := hex.DecodeString(s)
	a := make([]byte, 32)
	copy(a, bytes)
	return a
}

var DefConsensusConstants = ConsensusConstants{
	SlotBlocksTarget:               32,
	MinBlocksPerChallengeBlock:     16,
	MaxSubSlotBlocks:               128,
	NumSpsSubSlot:                  64,
	SubSlotItersStarting:           big.Exp(NewInt(2), NewInt(27)).Uint64(),
	DifficultyConstantFactor:       big.Exp(NewInt(2), NewInt(67)),
	DifficultyStarting:             7,
	DifficultyChangeMaxFactor:      3,
	SubEpochBlocks:                 384,
	EpochBlocks:                    4608,
	SignificantBits:                8,
	DiscriminantSizeBits:           1024,
	NumberZeroBitsPlotFilter:       9,
	MinPlotSize:                    32,
	MaxPlotSize:                    50,
	SubSlotTimeTarget:              600,
	NumSpIntervalsExtra:            3,
	MaxFutureTime:                  5 * 60,
	NumberOfTimestamps:             11,
	GenesisChallenge:               hexString32("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"),
	AggSigMeAdditionalData:         hexString("ccd5bb71183532bff220ba46c268991a3ff07eb358e8255a65c30a2dce0e5fbb"),
	GenesisPreFarmPoolPuzzleHash:   hexString32("d23da14695a188ae5708dd152263c4db883eb27edeb936178d4d988b8f3ce5fc"),
	GenesisPreFarmFarmerPuzzleHash: hexString32("3d8765d3a597ec1d99663f6c9816d915b9f68613ac94009884c4addaefcce6af"),
	MaxVdfWitnessSize:              64,
	MempoolBlockBuffer:             50,
	MaxCoinAmount:                  math.MaxUint64,
	MaxBlockCostClvm:               11000000000,
	CostPerByte:                    12000,
	WeightProofThreshold:           2,
	WeightProofRecentBlocks:        1000,
	MaxBlockCountPerRequests:       32,
	RustConditionChecker:           730000 + 138000,
	BlocksCacheSize:                4608 + (128 * 4),
	NetworkType:                    0,
	MaxGeneratorSize:               1000000,
	MaxGeneratorRefListSize:        512,
	PoolSubSlotIters:               37600000000,
}

var testCon = DefConsensusConstants

func init() {
	testCon.NumSpsSubSlot = 32
	testCon.SubSlotTimeTarget = 300
}

// def test_is_overflow_block(self):
//        assert not is_overflow_block(test_constants, uint8(27))
//        assert not is_overflow_block(test_constants, uint8(28))
//        assert is_overflow_block(test_constants, uint8(29))
//        assert is_overflow_block(test_constants, uint8(30))
//        assert is_overflow_block(test_constants, uint8(31))
//        with raises(ValueError):
//            assert is_overflow_block(test_constants, uint8(32))
//
func TestIsOverflowBlock(t *testing.T) {
	assert.False(t, IsOverflowBlock(&testCon, uint64(27)))
	assert.False(t, IsOverflowBlock(&testCon, uint64(28)))
	assert.True(t, IsOverflowBlock(&testCon, uint64(29)))
	assert.True(t, IsOverflowBlock(&testCon, uint64(30)))
	assert.True(t, IsOverflowBlock(&testCon, uint64(31)))
	defer func() {
		err := recover()
		assert.NotNil(t, err)
	}()
	assert.True(t, IsOverflowBlock(&testCon, uint64(32)))

}

//    def test_calculate_sp_iters(self):
//        ssi: uint64 = uint64(100001 * 64 * 4)
//        with raises(ValueError):
//            calculate_sp_iters(test_constants, ssi, uint8(32))
//        calculate_sp_iters(test_constants, ssi, uint8(31))
func TestCalculateSpIters(t *testing.T) {
	ssi := uint64(100001 * 64 * 4)
	CalculateSpIters(&testCon, ssi, uint64(31))
	defer func() {
		err := recover()
		assert.NotNil(t, err)
	}()
	CalculateSpIters(&testCon, ssi, uint64(32))
}

//    def test_calculate_ip_iters(self):
//        ssi: uint64 = uint64(100001 * 64 * 4)
//        sp_interval_iters = ssi // test_constants.NUM_SPS_SUB_SLOT
//
//        with raises(ValueError):
//            // Invalid signage point index
//            calculate_ip_iters(test_constants, ssi, uint8(123), uint64(100000))
//
//        sp_iters = sp_interval_iters * 13
//
//        with raises(ValueError):
//            // required_iters too high
//            calculate_ip_iters(test_constants, ssi, sp_interval_iters, sp_interval_iters)
//
//        with raises(ValueError):
//            // required_iters too high
//            calculate_ip_iters(test_constants, ssi, sp_interval_iters, sp_interval_iters * 12)
//
//        with raises(ValueError):
//            // required_iters too low (0)
//            calculate_ip_iters(test_constants, ssi, sp_interval_iters, uint64(0))
//
//        required_iters = sp_interval_iters - 1
//        ip_iters = calculate_ip_iters(test_constants, ssi, uint8(13), required_iters)
//        assert ip_iters == sp_iters + test_constants.NUM_SP_INTERVALS_EXTRA * sp_interval_iters + required_iters
//
//        required_iters = uint64(1)
//        ip_iters = calculate_ip_iters(test_constants, ssi, uint8(13), required_iters)
//        assert ip_iters == sp_iters + test_constants.NUM_SP_INTERVALS_EXTRA * sp_interval_iters + required_iters
//
//        required_iters = uint64(int(ssi * 4 / 300))
//        ip_iters = calculate_ip_iters(test_constants, ssi, uint8(13), required_iters)
//        assert ip_iters == sp_iters + test_constants.NUM_SP_INTERVALS_EXTRA * sp_interval_iters + required_iters
//        assert sp_iters < ip_iters
//
//        // Overflow
//        sp_iters = sp_interval_iters * (test_constants.NUM_SPS_SUB_SLOT - 1)
//        ip_iters = calculate_ip_iters(
//            test_constants,
//            ssi,
//            uint8(test_constants.NUM_SPS_SUB_SLOT - 1),
//            required_iters,
//        )
//        assert ip_iters == (sp_iters + test_constants.NUM_SP_INTERVALS_EXTRA * sp_interval_iters + required_iters) % ssi
//        assert sp_iters > ip_iters
func TestCalculateIpIters(t *testing.T) {
	ssi := uint64(100001 * 64 * 4)
	spIntervalIters := ssi / testCon.NumSpsSubSlot

	type args struct {
		constants         *ConsensusConstants
		subSlotIters      uint64
		signagePointIndex uint8
		requiredIters     uint64
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "001",
			args: args{
				constants:         &testCon,
				subSlotIters:      ssi,
				signagePointIndex: 123,
				requiredIters:     100_000,
			},
			want: true,
		},
		{
			name: "002",
			args: args{
				constants:         &testCon,
				subSlotIters:      ssi,
				signagePointIndex: uint8(spIntervalIters),
				requiredIters:     spIntervalIters,
			},
			want: true,
		},
		{
			name: "003",
			args: args{
				constants:         &testCon,
				subSlotIters:      ssi,
				signagePointIndex: uint8(spIntervalIters),
				requiredIters:     spIntervalIters * 12,
			},
			want: true,
		},
		{
			name: "004",
			args: args{
				constants:         &testCon,
				subSlotIters:      ssi,
				signagePointIndex: uint8(spIntervalIters),
				requiredIters:     0,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				err := recover()
				if tt.want {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err)
				}
			}()
			CalculateIpIters(tt.args.constants, tt.args.subSlotIters, tt.args.signagePointIndex, tt.args.requiredIters)
		})
	}

	tests2 := []struct {
		name string
		args args
	}{
		{
			name: "011",
			args: args{
				constants:         &testCon,
				subSlotIters:      ssi,
				signagePointIndex: 13,
				requiredIters:     spIntervalIters - 1,
			},
		},
		{
			name: "012",
			args: args{
				constants:         &testCon,
				subSlotIters:      ssi,
				signagePointIndex: 13,
				requiredIters:     1,
			},
		},
		{
			name: "013",
			args: args{
				constants:         &testCon,
				subSlotIters:      ssi,
				signagePointIndex: 13,
				requiredIters:     ssi / 300 * 4,
			},
		},
	}

	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			ipIters := CalculateIpIters(tt.args.constants, tt.args.subSlotIters, tt.args.signagePointIndex, tt.args.requiredIters)
			assert.True(t, ipIters == spIntervalIters*13+uint64(testCon.NumSpIntervalsExtra)*spIntervalIters+tt.args.requiredIters)
			if tt.name == "013" {
				assert.True(t, spIntervalIters*13 < ipIters)
			}
		})
	}

	spIters := spIntervalIters * (testCon.NumSpsSubSlot - 1)
	requiredIters := ssi / 300 * 4
	ipIters := CalculateIpIters(&testCon, ssi, uint8(testCon.NumSpsSubSlot-1), requiredIters)
	assert.True(t, ipIters == (spIters+uint64(testCon.NumSpIntervalsExtra)*spIntervalIters+requiredIters)%ssi)
	assert.True(t, spIters > ipIters)
}
