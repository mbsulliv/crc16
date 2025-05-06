//-----------------------------------------------------------------------------

package crc16

import (
	"fmt"
	"path"
	"runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

//-----------------------------------------------------------------------------

// Returns the function name of the calling function.
func funcName() string {
	vRet := "?"
	vPc, _, _, vOk := runtime.Caller(1)
	if vOk {
		vRet = path.Base(runtime.FuncForPC(vPc).Name())
	}
	return vRet
}

//-----------------------------------------------------------------------------

func TestMain(aT *testing.T) {
	vCases := []struct {
		Algo *TAlgo
	}{
		{&CRC16_DECT_R},
		{&CRC16_DECT_X},
		{&CRC16_NRSC_5},
		{&CRC16_GSM},
		{&CRC16_KERMIT},
		{&CRC16_XMODEM},
		{&CRC16_SPI_FUJITSU},
		{&CRC16_TMS37157},
		{&CRC16_RIELLO},
		{&CRC16_CRC_A},
		{&CRC16_CCITT_FALSE},
		{&CRC16_GENIBUS},
		{&CRC16_IBM_3740},
		{&CRC16_IBM_SDLC},
		{&CRC16_MCRF4XX},
		{&CRC16_X_25},
		{&CRC16_PROFIBUS},
		{&CRC16_DNP},
		{&CRC16_EN_13757},
		{&CRC16_OPENSAFETY_A},
		{&CRC16_M17},
		{&CRC16_LJ1200},
		{&CRC16_OPENSAFETY_B},
		{&CRC16_ARC},
		{&CRC16_BUYPASS},
		{&CRC16_MAXIM},
		{&CRC16_UMTS},
		{&CRC16_DDS_110},
		{&CRC16_CMS},
		{&CRC16_MODBUS},
		{&CRC16_USB},
		{&CRC16_T10_DIF},
		{&CRC16_TELEDISK},
		{&CRC16_CDMA2000},
	}

	vTestData := []byte("123456789")

	for _, vCase := range vCases {
		Convey(fmt.Sprintf("%s: %s", funcName(), vCase.Algo.Name), aT, func() {
			vTable := MakeTable(*vCase.Algo)
			So(vTable, ShouldNotBeNil)

			vGotCrc := Checksum(vTestData, vTable)
			So(fmt.Sprintf("0x%04X", vGotCrc), ShouldEqual, fmt.Sprintf("0x%04X", vTable.algo.Check))
		})
	}
}

//--------------------------------------

func TestHash(aT *testing.T) {
	Convey(funcName(), aT, func() {
		vTable := MakeTable(CRC16_XMODEM)
		vH := New(vTable)

		fmt.Fprint(vH, "standard")
		fmt.Fprint(vH, " library hash interface")
		vSum1 := vH.Sum16()
		vH.Reset()
		fmt.Fprint(vH, "standard library hash interface")
		vSum2 := vH.Sum16()
		So(vSum1, ShouldEqual, vSum2)

		So(vSum1, ShouldEqual, 0xe698)
		So(vH.Size(), ShouldEqual, 2)

		vBuf := make([]byte, 0, 10)
		vBuf = vH.Sum(vBuf)
		vExpected := []byte{0xe6, 0x98}
		So(len(vBuf), ShouldEqual, 2)
		So(vBuf[0], ShouldEqual, vExpected[0])
		So(vBuf[1], ShouldEqual, vExpected[1])

		So(vH.BlockSize(), ShouldEqual, 1)
	})
}

//-----------------------------------------------------------------------------
