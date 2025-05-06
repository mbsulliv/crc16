//-----------------------------------------------------------------------------

package crc16

import "hash"

//-----------------------------------------------------------------------------

// This file contains the CRC16 implementation of the
// go standard library hash.Hash interface

type Hash16 interface {
	hash.Hash
	Sum16() uint16
}

type digest struct {
	sum uint16
	t   *TTable
}

//-----------------------------------------------------------------------------

// Write adds more data to the running digest.
// It never returns an error.
func (aH *digest) Write(data []byte) (int, error) {
	aH.sum = Update(aH.sum, data, aH.t)
	return len(data), nil
}

//--------------------------------------

// Sum appends the current digest (leftmost byte first, big-endian)
// to b and returns the resulting slice.
// It does not change the underlying digest state.
func (aH digest) Sum(b []byte) []byte {
	s := aH.Sum16()
	return append(b, byte(s>>8), byte(s))
}

//--------------------------------------

// Reset resets the Hash to its initial state.
func (aH *digest) Reset() {
	aH.sum = aH.t.algo.Init
}

//--------------------------------------

// Size returns the number of bytes Sum will return.
func (aH digest) Size() int {
	return 2
}

//--------------------------------------

// BlockSize returns the undelying block size.
// See digest.Hash.BlockSize
func (aH digest) BlockSize() int {
	return 1
}

//--------------------------------------

// Sum16 returns the CRC16 checksum.
func (aH digest) Sum16() uint16 {
	return Complete(aH.sum, aH.t)
}

//--------------------------------------

// New creates a new CRC16 digest for the given table.
func New(t *TTable) Hash16 {
	aH := digest{t: t}
	aH.Reset()
	return &aH
}

//-----------------------------------------------------------------------------
