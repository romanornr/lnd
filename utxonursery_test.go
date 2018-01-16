// +build !rpctest

package main

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/roasbeef/btcd/btcec"
	"github.com/roasbeef/btcd/chaincfg/chainhash"
	"github.com/roasbeef/btcd/txscript"
	"github.com/roasbeef/btcd/wire"
	"github.com/roasbeef/btcutil"
	"github.com/viacoin/lnd/lnwallet"
)

var (
	outPoints = []wire.OutPoint{
		{
			Hash: [chainhash.HashSize]byte{
				0x51, 0xb6, 0x37, 0xd8, 0xfc, 0xd2, 0xc6, 0xda,
				0x48, 0x59, 0xe6, 0x96, 0x31, 0x13, 0xa1, 0x17,
				0x2d, 0xe7, 0x93, 0xe4, 0xb7, 0x25, 0xb8, 0x4d,
				0x1f, 0xb, 0x4c, 0xf9, 0x9e, 0xc5, 0x8c, 0xe9,
			},
			Index: 9,
		},
		{
			Hash: [chainhash.HashSize]byte{
				0xb7, 0x94, 0x38, 0x5f, 0x2d, 0x1e, 0xf7, 0xab,
				0x4d, 0x92, 0x73, 0xd1, 0x90, 0x63, 0x81, 0xb4,
				0x4f, 0x2f, 0x6f, 0x25, 0x88, 0xa3, 0xef, 0xb9,
				0x6a, 0x49, 0x18, 0x83, 0x31, 0x98, 0x47, 0x53,
			},
			Index: 49,
		},
		{
			Hash: [chainhash.HashSize]byte{
				0x81, 0xb6, 0x37, 0xd8, 0xfc, 0xd2, 0xc6, 0xda,
				0x63, 0x59, 0xe6, 0x96, 0x31, 0x13, 0xa1, 0x17,
				0xd, 0xe7, 0x95, 0xe4, 0xb7, 0x25, 0xb8, 0x4d,
				0x1e, 0xb, 0x4c, 0xfd, 0x9e, 0xc5, 0x8c, 0xe9,
			},
			Index: 23,
		},
	}

	keys = [][]byte{
		{0x04, 0x11, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a,
			0x01, 0x6b, 0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc, 0x1e,
			0xb6, 0x8a, 0x38, 0x2e, 0x97, 0xb1, 0x48, 0x2e, 0xca,
			0xd7, 0xb1, 0x48, 0xa6, 0x90, 0x9a, 0x5c, 0xb2, 0xe0,
			0xea, 0xdd, 0xfb, 0x84, 0xcc, 0xf9, 0x74, 0x44, 0x64,
			0xf8, 0x2e, 0x16, 0x0b, 0xfa, 0x9b, 0x8b, 0x64, 0xf9,
			0xd4, 0xc0, 0x3f, 0x99, 0x9b, 0x86, 0x43, 0xf6, 0x56,
			0xb4, 0x12, 0xa3,
		},
		{0x07, 0x11, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a,
			0x01, 0x6b, 0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc, 0x1e,
			0xb6, 0x8a, 0x38, 0x2e, 0x97, 0xb1, 0x48, 0x2e, 0xca,
			0xd7, 0xb1, 0x48, 0xa6, 0x90, 0x9a, 0x5c, 0xb2, 0xe0,
			0xea, 0xdd, 0xfb, 0x84, 0xcc, 0xf9, 0x74, 0x44, 0x64,
			0xf8, 0x2e, 0x16, 0x0b, 0xfa, 0x9b, 0x8b, 0x64, 0xf9,
			0xd4, 0xc0, 0x3f, 0x99, 0x9b, 0x86, 0x43, 0xf6, 0x56,
			0xb4, 0x12, 0xa3,
		},
		{0x02, 0xce, 0x0b, 0x14, 0xfb, 0x84, 0x2b, 0x1b,
			0xa5, 0x49, 0xfd, 0xd6, 0x75, 0xc9, 0x80, 0x75, 0xf1,
			0x2e, 0x9c, 0x51, 0x0f, 0x8e, 0xf5, 0x2b, 0xd0, 0x21,
			0xa9, 0xa1, 0xf4, 0x80, 0x9d, 0x3b, 0x4d,
		},
	}

	signDescriptors = []lnwallet.SignDescriptor{
		{
			SingleTweak: []byte{
				0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02,
				0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02,
				0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02,
				0x02, 0x02, 0x02, 0x02, 0x02,
			},
			WitnessScript: []byte{
				0x00, 0x14, 0xee, 0x91, 0x41, 0x7e, 0x85, 0x6c, 0xde,
				0x10, 0xa2, 0x91, 0x1e, 0xdc, 0xbd, 0xbd, 0x69, 0xe2,
				0xef, 0xb5, 0x71, 0x48,
			},
			Output: &wire.TxOut{
				Value: 5000000000,
				PkScript: []byte{
					0x41, // OP_DATA_65
					0x04, 0xd6, 0x4b, 0xdf, 0xd0, 0x9e, 0xb1, 0xc5,
					0xfe, 0x29, 0x5a, 0xbd, 0xeb, 0x1d, 0xca, 0x42,
					0x81, 0xbe, 0x98, 0x8e, 0x2d, 0xa0, 0xb6, 0xc1,
					0xc6, 0xa5, 0x9d, 0xc2, 0x26, 0xc2, 0x86, 0x24,
					0xe1, 0x81, 0x75, 0xe8, 0x51, 0xc9, 0x6b, 0x97,
					0x3d, 0x81, 0xb0, 0x1c, 0xc3, 0x1f, 0x04, 0x78,
					0x34, 0xbc, 0x06, 0xd6, 0xd6, 0xed, 0xf6, 0x20,
					0xd1, 0x84, 0x24, 0x1a, 0x6a, 0xed, 0x8b, 0x63,
					0xa6, // 65-byte signature
					0xac, // OP_CHECKSIG
				},
			},
			HashType: txscript.SigHashAll,
		},
		{
			SingleTweak: []byte{
				0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02,
				0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02,
				0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02,
				0x02, 0x02, 0x02, 0x02, 0x02,
			},
			WitnessScript: []byte{
				0x00, 0x14, 0xee, 0x91, 0x41, 0x7e, 0x85, 0x6c, 0xde,
				0x10, 0xa2, 0x91, 0x1e, 0xdc, 0xbd, 0xbd, 0x69, 0xe2,
				0xef, 0xb5, 0x71, 0x48,
			},
			Output: &wire.TxOut{
				Value: 5000000000,
				PkScript: []byte{
					0x41, // OP_DATA_65
					0x04, 0xd6, 0x4b, 0xdf, 0xd0, 0x9e, 0xb1, 0xc5,
					0xfe, 0x29, 0x5a, 0xbd, 0xeb, 0x1d, 0xca, 0x42,
					0x81, 0xbe, 0x98, 0x8e, 0x2d, 0xa0, 0xb6, 0xc1,
					0xc6, 0xa5, 0x9d, 0xc2, 0x26, 0xc2, 0x86, 0x24,
					0xe1, 0x81, 0x75, 0xe8, 0x51, 0xc9, 0x6b, 0x97,
					0x3d, 0x81, 0xb0, 0x1c, 0xc3, 0x1f, 0x04, 0x78,
					0x34, 0xbc, 0x06, 0xd6, 0xd6, 0xed, 0xf6, 0x20,
					0xd1, 0x84, 0x24, 0x1a, 0x6a, 0xed, 0x8b, 0x63,
					0xa6, // 65-byte signature
					0xac, // OP_CHECKSIG
				},
			},
			HashType: txscript.SigHashAll,
		},
		{
			SingleTweak: []byte{
				0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02,
				0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02,
				0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02,
				0x02, 0x02, 0x02, 0x02, 0x02,
			},
			WitnessScript: []byte{
				0x00, 0x14, 0xee, 0x91, 0x41, 0x7e, 0x85, 0x6c, 0xde,
				0x10, 0xa2, 0x91, 0x1e, 0xdc, 0xbd, 0xbd, 0x69, 0xe2,
				0xef, 0xb5, 0x71, 0x48,
			},
			Output: &wire.TxOut{
				Value: 5000000000,
				PkScript: []byte{
					0x41, // OP_DATA_65
					0x04, 0xd6, 0x4b, 0xdf, 0xd0, 0x9e, 0xb1, 0xc5,
					0xfe, 0x29, 0x5a, 0xbd, 0xeb, 0x1d, 0xca, 0x42,
					0x81, 0xbe, 0x98, 0x8e, 0x2d, 0xa0, 0xb6, 0xc1,
					0xc6, 0xa5, 0x9d, 0xc2, 0x26, 0xc2, 0x86, 0x24,
					0xe1, 0x81, 0x75, 0xe8, 0x51, 0xc9, 0x6b, 0x97,
					0x3d, 0x81, 0xb0, 0x1c, 0xc3, 0x1f, 0x04, 0x78,
					0x34, 0xbc, 0x06, 0xd6, 0xd6, 0xed, 0xf6, 0x20,
					0xd1, 0x84, 0x24, 0x1a, 0x6a, 0xed, 0x8b, 0x63,
					0xa6, // 65-byte signature
					0xac, // OP_CHECKSIG
				},
			},
			HashType: txscript.SigHashAll,
		},
	}

	kidOutputs = []kidOutput{
		{
			breachedOutput: breachedOutput{
				amt:         btcutil.Amount(13e7),
				outpoint:    outPoints[0],
				witnessType: lnwallet.CommitmentTimeLock,
			},
			originChanPoint:  outPoints[1],
			blocksToMaturity: uint32(100),
			confHeight:       uint32(1770001),
		},

		{
			breachedOutput: breachedOutput{
				amt:         btcutil.Amount(24e7),
				outpoint:    outPoints[1],
				witnessType: lnwallet.CommitmentTimeLock,
			},
			originChanPoint:  outPoints[0],
			blocksToMaturity: uint32(50),
			confHeight:       uint32(22342321),
		},

		{
			breachedOutput: breachedOutput{
				amt:         btcutil.Amount(2e5),
				outpoint:    outPoints[2],
				witnessType: lnwallet.CommitmentTimeLock,
			},
			originChanPoint:  outPoints[2],
			blocksToMaturity: uint32(12),
			confHeight:       uint32(34241),
		},
	}

	babyOutputs = []babyOutput{
		{
			kidOutput: kidOutputs[0],
			expiry:    3829,
			timeoutTx: timeoutTx,
		},
		{
			kidOutput: kidOutputs[1],
			expiry:    85903,
			timeoutTx: timeoutTx,
		},
		{
			kidOutput: kidOutputs[2],
			expiry:    4,
			timeoutTx: timeoutTx,
		},
	}

	// Dummy timeout tx used to test serialization, borrowed from btcd
	// msgtx_test
	timeoutTx = &wire.MsgTx{
		Version: 1,
		TxIn: []*wire.TxIn{
			{
				PreviousOutPoint: wire.OutPoint{
					Hash: chainhash.Hash{
						0xa5, 0x33, 0x52, 0xd5, 0x13, 0x57, 0x66, 0xf0,
						0x30, 0x76, 0x59, 0x74, 0x18, 0x26, 0x3d, 0xa2,
						0xd9, 0xc9, 0x58, 0x31, 0x59, 0x68, 0xfe, 0xa8,
						0x23, 0x52, 0x94, 0x67, 0x48, 0x1f, 0xf9, 0xcd,
					},
					Index: 19,
				},
				SignatureScript: []byte{},
				Witness: [][]byte{
					{ // 70-byte signature
						0x30, 0x43, 0x02, 0x1f, 0x4d, 0x23, 0x81, 0xdc,
						0x97, 0xf1, 0x82, 0xab, 0xd8, 0x18, 0x5f, 0x51,
						0x75, 0x30, 0x18, 0x52, 0x32, 0x12, 0xf5, 0xdd,
						0xc0, 0x7c, 0xc4, 0xe6, 0x3a, 0x8d, 0xc0, 0x36,
						0x58, 0xda, 0x19, 0x02, 0x20, 0x60, 0x8b, 0x5c,
						0x4d, 0x92, 0xb8, 0x6b, 0x6d, 0xe7, 0xd7, 0x8e,
						0xf2, 0x3a, 0x2f, 0xa7, 0x35, 0xbc, 0xb5, 0x9b,
						0x91, 0x4a, 0x48, 0xb0, 0xe1, 0x87, 0xc5, 0xe7,
						0x56, 0x9a, 0x18, 0x19, 0x70, 0x01,
					},
					{ // 33-byte serialize pub key
						0x03, 0x07, 0xea, 0xd0, 0x84, 0x80, 0x7e, 0xb7,
						0x63, 0x46, 0xdf, 0x69, 0x77, 0x00, 0x0c, 0x89,
						0x39, 0x2f, 0x45, 0xc7, 0x64, 0x25, 0xb2, 0x61,
						0x81, 0xf5, 0x21, 0xd7, 0xf3, 0x70, 0x06, 0x6a,
						0x8f,
					},
				},
				Sequence: 0xffffffff,
			},
		},
		TxOut: []*wire.TxOut{
			{
				Value: 395019,
				PkScript: []byte{ // p2wkh output
					0x00, // Version 0 witness program
					0x14, // OP_DATA_20
					0x9d, 0xda, 0xc6, 0xf3, 0x9d, 0x51, 0xe0, 0x39,
					0x8e, 0x53, 0x2a, 0x22, 0xc4, 0x1b, 0xa1, 0x89,
					0x40, 0x6a, 0x85, 0x23, // 20-byte pub key hash
				},
			},
		},
	}
)

func init() {
	// Finish initializing our test vectors by parsing the desired public keys and
	// properly populating the sign descriptors of all baby and kid outputs.
	for i := range signDescriptors {
		pk, err := btcec.ParsePubKey(keys[i], btcec.S256())
		if err != nil {
			panic(fmt.Sprintf("unable to parse pub key during init: %v", err))
		}
		signDescriptors[i].PubKey = pk

		kidOutputs[i].signDesc = signDescriptors[i]
		babyOutputs[i].kidOutput.signDesc = signDescriptors[i]

	}
}

func TestDeserializeKidsList(t *testing.T) {
	var b bytes.Buffer
	for _, kid := range kidOutputs {
		if err := kid.Encode(&b); err != nil {
			t.Fatalf("unable to serialize and add kid output to "+
				"list: %v", err)
		}
	}

	kidList, err := deserializeKidList(&b)
	if err != nil {
		t.Fatalf("unable to deserialize kid output list: %v", err)
	}

	for i := range kidOutputs {
		if !reflect.DeepEqual(&kidOutputs[i], kidList[i]) {
			t.Fatalf("kidOutputs don't match \n%+v\n%+v",
				&kidOutputs[i], kidList[i])
		}
	}
}

func TestKidOutputSerialization(t *testing.T) {
	for i, kid := range kidOutputs {
		var b bytes.Buffer
		if err := kid.Encode(&b); err != nil {
			t.Fatalf("Encode #%d: unable to serialize "+
				"kid output: %v", i, err)
		}

		var deserializedKid kidOutput
		if err := deserializedKid.Decode(&b); err != nil {
			t.Fatalf("Decode #%d: unable to deserialize "+
				"kid output: %v", i, err)
		}

		if !reflect.DeepEqual(kid, deserializedKid) {
			t.Fatalf("DeepEqual #%d: unexpected kidOutput, "+
				"want %+v, got %+v",
				i, kid, deserializedKid)
		}
	}
}

func TestBabyOutputSerialization(t *testing.T) {
	for i, baby := range babyOutputs {
		var b bytes.Buffer
		if err := baby.Encode(&b); err != nil {
			t.Fatalf("Encode #%d: unable to serialize "+
				"baby output: %v", i, err)
		}

		var deserializedBaby babyOutput
		if err := deserializedBaby.Decode(&b); err != nil {
			t.Fatalf("Decode #%d: unable to deserialize "+
				"baby output: %v", i, err)
		}

		if !reflect.DeepEqual(baby, deserializedBaby) {
			t.Fatalf("DeepEqual #%d: unexpected babyOutput, "+
				"want %+v, got %+v",
				i, baby, deserializedBaby)
		}

	}
}
