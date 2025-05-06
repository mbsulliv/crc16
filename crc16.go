//-----------------------------------------------------------------------------

// Package crc16 implements the 16-bit cyclic redundancy check, or CRC-16, checksum.
//
// It provides parameters for the majority of well-known CRC-16 algorithms.
package crc16

import "math/bits"

//-----------------------------------------------------------------------------

// TAlgo represents parameters of CRC-16 algorithms.
// More information about algorithms parametrization and parameter descriptions
// can be found here - http://www.zlib.net/crc_v3.txt
type TAlgo struct {
	Poly   uint16
	Init   uint16
	RefIn  bool
	RefOut bool
	XorOut uint16
	Check  uint16
	Name   string
}

// Predefined CRC-16 algorithms.
// List of algorithms with their parameters borrowed from here -  http://reveng.sourceforge.net/crc-catalogue/16.htm
//
// The variables can be used to create TTable for the selected algorithm.
var (
	CRC16_DECT_R      = TAlgo{0x0589, 0x0000, false, false, 0x0001, 0x007E, "CRC-16/DECT-R"}
	CRC16_DECT_X      = TAlgo{0x0589, 0x0000, false, false, 0x0000, 0x007F, "CRC-16/DECT-X"}
	CRC16_KERMIT      = TAlgo{0x1021, 0x0000, true, true, 0x0000, 0x2189, "CRC-16/KERMIT"}
	CRC16_XMODEM      = TAlgo{0x1021, 0x0000, false, false, 0x0000, 0x31C3, "CRC-16/XMODEM"}
	CRC16_AUG_CCITT   = TAlgo{0x1021, 0x1D0F, false, false, 0x0000, 0xE5CC, "CRC-16/AUG-CCITT"}
	CRC16_TMS37157    = TAlgo{0x1021, 0x89EC, true, true, 0x0000, 0x26B1, "CRC-16/TMS37157"}
	CRC16_RIELLO      = TAlgo{0x1021, 0xB2AA, true, true, 0x0000, 0x63D0, "CRC-16/RIELLO"}
	CRC16_CRC_A       = TAlgo{0x1021, 0xC6C6, true, true, 0x0000, 0xBF05, "CRC-16/CRC-A"}
	CRC16_CCITT_FALSE = TAlgo{0x1021, 0xFFFF, false, false, 0x0000, 0x29B1, "CRC-16/CCITT-FALSE"}
	CRC16_GENIBUS     = TAlgo{0x1021, 0xFFFF, false, false, 0xFFFF, 0xD64E, "CRC-16/GENIBUS"}
	CRC16_MCRF4XX     = TAlgo{0x1021, 0xFFFF, true, true, 0x0000, 0x6F91, "CRC-16/MCRF4XX"}
	CRC16_X_25        = TAlgo{0x1021, 0xFFFF, true, true, 0xFFFF, 0x906E, "CRC-16/X-25"}
	CRC16_DNP         = TAlgo{0x3D65, 0x0000, true, true, 0xFFFF, 0xEA82, "CRC-16/DNP"}
	CRC16_EN_13757    = TAlgo{0x3D65, 0x0000, false, false, 0xFFFF, 0xC2B7, "CRC-16/EN-13757"}
	CRC16_ARC         = TAlgo{0x8005, 0x0000, true, true, 0x0000, 0xBB3D, "CRC-16/ARC"}
	CRC16_BUYPASS     = TAlgo{0x8005, 0x0000, false, false, 0x0000, 0xFEE8, "CRC-16/BUYPASS"}
	CRC16_MAXIM       = TAlgo{0x8005, 0x0000, true, true, 0xFFFF, 0x44C2, "CRC-16/MAXIM"}
	CRC16_DDS_110     = TAlgo{0x8005, 0x800D, false, false, 0x0000, 0x9ECF, "CRC-16/DDS-110"}
	CRC16_MODBUS      = TAlgo{0x8005, 0xFFFF, true, true, 0x0000, 0x4B37, "CRC-16/MODBUS"}
	CRC16_USB         = TAlgo{0x8005, 0xFFFF, true, true, 0xFFFF, 0xB4C8, "CRC-16/USB"}
	CRC16_T10_DIF     = TAlgo{0x8BB7, 0x0000, false, false, 0x0000, 0xD0DB, "CRC-16/T10-DIF"}
	CRC16_TELEDISK    = TAlgo{0xA097, 0x0000, false, false, 0x0000, 0x0FB3, "CRC-16/TELEDISK"}
	CRC16_CDMA2000    = TAlgo{0xC867, 0xFFFF, false, false, 0x0000, 0x4C06, "CRC-16/CDMA2000"}
)

// TTable is a 256-word table representing polinomial and algorithm settings for efficient processing.
type TTable struct {
	algo TAlgo
	data [256]uint16
}

//-----------------------------------------------------------------------------

// MakeTable returns the TTable constructed from the specified algorithm.
func MakeTable(aAlgo TAlgo) *TTable {
	vTable := new(TTable)
	vTable.algo = aAlgo
	for n := 0; n < 256; n++ {
		crc := uint16(n) << 8
		for i := 0; i < 8; i++ {
			bit := (crc & 0x8000) != 0
			crc <<= 1
			if bit {
				crc ^= aAlgo.Poly
			}
		}
		vTable.data[n] = crc
	}
	return vTable
}

//--------------------------------------

// Init returns the initial value for CRC register corresponding to the specified algorithm.
func Init(aTable *TTable) uint16 {
	return aTable.algo.Init
}

//--------------------------------------

// Update returns the result of adding the bytes in data to the crc.
func Update(crc uint16, data []byte, aTable *TTable) uint16 {
	for _, d := range data {
		if aTable.algo.RefIn {
			d = bits.Reverse8(d)
		}
		crc = crc<<8 ^ aTable.data[byte(crc>>8)^d]
	}
	return crc
}

//--------------------------------------

// Complete returns the result of CRC calculation and post-calculation processing of the crc.
func Complete(crc uint16, aTable *TTable) uint16 {
	if aTable.algo.RefOut {
		return bits.Reverse16(crc) ^ aTable.algo.XorOut
	}
	return crc ^ aTable.algo.XorOut
}

//--------------------------------------

// Checksum returns CRC checksum of data using scpecified algorithm represented by the TTable.
func Checksum(data []byte, aTable *TTable) uint16 {
	crc := Init(aTable)
	crc = Update(crc, data, aTable)
	return Complete(crc, aTable)
}

//-----------------------------------------------------------------------------
