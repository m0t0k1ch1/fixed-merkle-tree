package merkle

import (
	"bytes"
	"crypto/sha256"
	"testing"
)

var testConfig = &Config{
	hasher:       sha256.New(),
	depth:        2,
	hashSize:     32,
	allLeavesNum: 4,
	allNodesNum:  7,
}

func TestNewTree(t *testing.T) {
	type input struct {
		config *Config
		leaves [][]byte
	}
	type output struct {
		root *Node
		err  error
	}
	testCases := []struct {
		name string
		in   input
		out  output
	}{
		{
			"success",
			input{
				testConfig,
				[][]byte{
					[]byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
					[]byte{0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
					[]byte{0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03},
				},
			},
			output{
				&Node{
					b: []byte{
						0x52, 0x36, 0xe8, 0xd7, 0xc1, 0x38, 0x4c, 0x10,
						0xd7, 0x67, 0xb9, 0x2b, 0x8c, 0xd3, 0x3e, 0x56,
						0x9f, 0xf3, 0x3c, 0x17, 0xf2, 0x64, 0x73, 0xc8,
						0xf9, 0x5d, 0xf8, 0x99, 0xf3, 0x7c, 0x47, 0xfc,
					},
					left: &Node{
						b: []byte{
							0x1e, 0x2a, 0xbc, 0x6e, 0x47, 0x7b, 0x5a, 0xc3,
							0xb1, 0x5d, 0x7f, 0x15, 0x39, 0x89, 0xf4, 0x9d,
							0xb2, 0x19, 0xc0, 0x24, 0x4a, 0xc9, 0x4b, 0x9a,
							0x1b, 0x77, 0x8c, 0x9d, 0xbd, 0xd5, 0xb7, 0xe4,
						},
					},
					right: &Node{
						b: []byte{
							0x90, 0x4f, 0x88, 0xd8, 0xf6, 0x4f, 0xf1, 0xaa,
							0x68, 0xed, 0xa4, 0x46, 0xf8, 0xec, 0xf0, 0xeb,
							0x0f, 0xc7, 0x92, 0x52, 0xef, 0x8d, 0x55, 0xca,
							0x3c, 0x17, 0x61, 0x1c, 0x00, 0xb7, 0x5f, 0x6f,
						},
					},
				},
				nil,
			},
		},
		{
			"failure: too many leaves",
			input{
				testConfig,
				[][]byte{nil, nil, nil, nil, nil},
			},
			output{
				nil,
				ErrTooManyLeaves,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in, out := tc.in, tc.out

			tree, err := NewTree(in.config, in.leaves)
			if err != out.err {
				t.Errorf("expected: %v, actual: %v", out.err, err)
			}

			if err == nil {
				rootActual := tree.Root()
				rootExpected := out.root

				testNodesEquality(t, rootExpected, rootActual)
				testNodesEquality(t, rootExpected.Left(), rootActual.Left())
				testNodesEquality(t, rootExpected.Right(), rootActual.Right())
			}
		})
	}
}

func TestMembershipProof(t *testing.T) {
	tree, err := NewTree(
		testConfig,
		[][]byte{
			[]byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
			[]byte{0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02},
			[]byte{0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03, 0x03},
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	type input struct {
		index int
	}
	type output struct {
		proof []byte
		err   error
	}
	testCases := []struct {
		name string
		in   input
		out  output
	}{
		{
			"success",
			input{
				0,
			},
			output{
				[]byte{
					0x10, 0xae, 0x0f, 0xdb, 0xf8, 0xc4, 0xf1, 0xf2,
					0xb5, 0xe7, 0x08, 0xfd, 0x74, 0x78, 0xab, 0xd2,
					0xbf, 0x03, 0xb1, 0x90, 0xed, 0xc8, 0x78, 0xdc,
					0x62, 0xad, 0xa6, 0x45, 0xaa, 0x7e, 0x03, 0x10,
					0x90, 0x4f, 0x88, 0xd8, 0xf6, 0x4f, 0xf1, 0xaa,
					0x68, 0xed, 0xa4, 0x46, 0xf8, 0xec, 0xf0, 0xeb,
					0x0f, 0xc7, 0x92, 0x52, 0xef, 0x8d, 0x55, 0xca,
					0x3c, 0x17, 0x61, 0x1c, 0x00, 0xb7, 0x5f, 0x6f,
				},
				nil,
			},
		},
		{
			"failure: leaf index out of range",
			input{
				4,
			},
			output{
				nil,
				ErrLeafIndexOutOfRange,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in, out := tc.in, tc.out

			proof, err := tree.CreateMembershipProof(in.index)
			if err != out.err {
				t.Errorf("expected: %v, actual: %v", out.err, err)
			}
			if !bytes.Equal(proof, out.proof) {
				t.Errorf("expected: %x, actual: %x", out.proof, proof)
			}

			if len(proof) > 0 {
				for j := 0; j <= tree.config.allLeavesNum; j++ {
					ok, err := tree.VerifyMembershipProof(j, proof)
					if err != nil {
						if j < tree.config.allLeavesNum {
							t.Fatal(err)
						} else if err != ErrLeafIndexOutOfRange {
							t.Fatal(err)
						}
					}
					if j == in.index && !ok {
						t.Errorf("expected: %t, actual: %t", true, ok)
					} else if j != in.index && ok {
						t.Errorf("expected: %t, actual: %t", false, ok)
					}
				}
			}
		})
	}
}
