package lexorank

import (
	"errors"
	"fmt"
	"strings"
)

var (
	DifferentBaseErr = errors.New("expected numbers of same numeral sys")

	zeroMag = []byte{0}
	oneMag  = []byte{1}
)

type LexoInteger struct {
	sys  LexoNumeralSystem
	sign int
	mag  []byte
}

func LexoIntegerParse(strFull string, system LexoNumeralSystem) (*LexoInteger, error) {
	str := strFull
	sign := 1
	if len(strFull) > 0 {
		switch strFull[0] {
		case system.GetPositiveChar():
			str = strFull[1:]
		case system.GetNegativeChar():
			sign = -1
			str = strFull[1:]
		}
	}
	var mag = make([]byte, len(str))
	magIndex := 0
	for strIndex := len(str) - 1; strIndex >= 0; strIndex-- {
		digit, err := system.Digit(str[strIndex])
		if err != nil {
			return nil, fmt.Errorf("to digit error: %w", err)
		}
		mag[magIndex] = digit
		magIndex++
	}
	return makeLexoInteger(system, sign, mag), nil
}

func makeLexoInteger(system LexoNumeralSystem, sign int, mag []byte) *LexoInteger {
	var actualLength int
	for actualLength = len(mag); actualLength > 0 && mag[actualLength-1] == 0; actualLength-- {
		// ignore
	}
	if actualLength == 0 {
		return lexoIntegerZero(system)
	}
	if actualLength == len(mag) {
		return NewLexoInteger(system, sign, mag)
	}
	return NewLexoInteger(system, sign, mag[:actualLength])
}

func NewLexoInteger(sys LexoNumeralSystem, sign int, mag []byte) *LexoInteger {
	return &LexoInteger{
		sys:  sys,
		sign: sign,
		mag:  mag,
	}
}

func lexoIntegerZero(sys LexoNumeralSystem) *LexoInteger {
	return NewLexoInteger(sys, 0, zeroMag)
}

func lexoIntegerOne(sys LexoNumeralSystem) *LexoInteger {
	return NewLexoInteger(sys, 1, oneMag)
}

func (d *LexoInteger) IsZero() bool {
	return d.sign == 0 && len(d.mag) == 1 && d.mag[0] == 0
}

func (d *LexoInteger) String() string {
	if d.IsZero() {
		return string(d.sys.Char(0))
	}
	var sb strings.Builder
	if d.sign == -1 {
		sb.WriteByte(d.sys.GetNegativeChar())
	}
	for idx := len(d.mag) - 1; idx >= 0; idx-- {
		sb.WriteByte(d.sys.Char(d.mag[idx]))
	}
	return sb.String()
}

func (d *LexoInteger) GetMag(idx int) byte {
	return d.mag[idx]
}

func (d *LexoInteger) ShiftLeft(times int) *LexoInteger {
	if times == 0 {
		return d
	}
	if times < 0 {
		return d.ShiftRight(-times)
	}
	newMag := append(make([]byte, times), d.mag...)
	return makeLexoInteger(d.sys, d.sign, newMag)
}

func (d *LexoInteger) ShiftRight(times int) *LexoInteger {
	if times == 0 {
		return d
	}
	if len(d.mag)-times <= 0 {
		return lexoIntegerZero(d.sys)
	}
	if times < 0 {
		return d.ShiftLeft(-times)
	}
	newMag := d.mag[times:len(d.mag)]
	return makeLexoInteger(d.sys, d.sign, newMag)
}

func (d *LexoInteger) GetSystem() LexoNumeralSystem {
	return d.sys
}

func (d *LexoInteger) Compare(other *LexoInteger) int {
	switch {
	case d == other:
		return 0
	case other == nil:
		return 1
	case d.sign == 1 && other.sign == 1:
		return d.cmpMag(other)
	case d.sign == -1 && other.sign == -1:
		return other.cmpMag(d)
	case d.sign == 1:
		return 1
	case other.sign == 1:
		return -1
	}
	return 0
}

func (d *LexoInteger) Add(other *LexoInteger) (*LexoInteger, error) {
	if d.sys.GetBase() != other.sys.GetBase() {
		return nil, DifferentBaseErr
	}
	if d.IsZero() {
		return other, nil
	}
	if other.IsZero() {
		return d, nil
	}
	if d.sign != other.sign {
		if d.sign == -1 {
			val, _ := d.negate().Sub(other)
			return val.negate(), nil
		}
		val, _ := d.Sub(other.negate())
		return val, nil
	}
	return makeLexoInteger(d.sys, d.sign, d.addMag(other)), nil
}

func (d *LexoInteger) Sub(other *LexoInteger) (*LexoInteger, error) {
	if d.sys.GetBase() != other.sys.GetBase() {
		return nil, DifferentBaseErr
	}
	if d.IsZero() {
		return other.negate(), nil
	}
	if other.IsZero() {
		return d, nil
	}
	if d.sign != other.sign {
		if d.sign == -1 {
			val, _ := d.negate().Add(other)
			return val.negate(), nil
		}
		val, _ := d.Add(other.negate())
		return val, nil
	}
	cmp := d.cmpMag(other)
	if cmp == 0 {
		return lexoIntegerZero(d.sys), nil
	}
	if cmp < 0 {
		return makeLexoInteger(d.sys, -d.sign, other.subMag(d)), nil
	}
	return makeLexoInteger(d.sys, d.sign, d.subMag(other)), nil
}

func (d *LexoInteger) negate() *LexoInteger {
	return makeLexoInteger(d.sys, -d.sign, d.mag)
}

func (d *LexoInteger) Multiply(other *LexoInteger) (*LexoInteger, error) {
	if d.sys.GetBase() != other.sys.GetBase() {
		return nil, DifferentBaseErr
	}
	sign := d.sign * other.sign
	switch {
	case d.IsZero():
		return d, nil
	case other.IsZero():
		return other, nil
	case d.isOne():
		return makeLexoInteger(d.sys, sign, other.mag), nil
	case other.isOne():
		return makeLexoInteger(d.sys, sign, d.mag), nil
	}
	newMag := d.multiplyMag(other)
	return makeLexoInteger(d.sys, sign, newMag), nil
}

func (d *LexoInteger) isOne() bool {
	return len(d.mag) == 1 && d.mag[0] == 1
}

func (d *LexoInteger) Equals(other *LexoInteger) bool {
	if d == other {
		return true
	}
	if other == nil {
		return false
	}
	return d.sys.GetBase() == other.sys.GetBase() && d.Compare(other) == 0
}

func (d *LexoInteger) addMag(other *LexoInteger) []byte {
	estimatedSize := len(d.mag)
	if estimatedSize < len(other.mag) {
		estimatedSize = len(other.mag)
	}
	result := make([]byte, estimatedSize)
	carry := byte(0)
	for i := 0; i < estimatedSize; i++ {
		leftNum, rightNum := byte(0), byte(0)
		if i < len(d.mag) {
			leftNum = d.mag[i]
		}
		if i < len(other.mag) {
			rightNum = other.mag[i]
		}
		sum := leftNum + rightNum + carry
		carry = 0
		if sum >= d.sys.GetBase() {
			sum -= d.sys.GetBase()
			carry = 1
		}
		result[i] = sum
	}
	if carry > 0 {
		result = append(result, carry)
	}
	return result
}

func (d *LexoInteger) multiplyMag(other *LexoInteger) []byte {
	result := make([]byte, len(d.mag)+len(other.mag))
	base := int(d.sys.GetBase())
	for li, lb := range d.mag {
		for ri, rb := range other.mag {
			resultIndex := li + ri
			resultInt := int(result[resultIndex]) + int(lb)*int(rb)
			for ; resultInt >= base; resultInt -= base {
				result[resultIndex+1]++
			}
			result[resultIndex] = byte(resultInt)
		}
	}
	return result
}

func (d *LexoInteger) subMag(other *LexoInteger) []byte {
	rComplement := other.complement(len(d.mag))
	rSum := d.addMag(rComplement)
	rSum[len(rSum)-1] = 0
	return NewLexoInteger(d.sys, d.sign, rSum).addMag(lexoIntegerOne(d.sys))
}

func (d *LexoInteger) complement(digits int) *LexoInteger {
	newMag := make([]byte, digits)
	maxDigit := d.sys.GetBase() - 1
	for i := 0; i < len(d.mag); i++ {
		newMag[i] = maxDigit - d.mag[i]
	}
	for i := len(d.mag); i < digits; i++ {
		newMag[i] = maxDigit
	}
	return makeLexoInteger(d.sys, d.sign, newMag)
}

func (d *LexoInteger) cmpMag(other *LexoInteger) int {
	if len(d.mag) < len(other.mag) {
		return -1
	}
	if len(d.mag) > len(other.mag) {
		return 1
	}
	for idx := len(d.mag) - 1; idx >= 0; idx-- {
		if d.mag[idx] < other.mag[idx] {
			return -1
		}
		if d.mag[idx] > other.mag[idx] {
			return 1
		}
	}
	return 0
}
