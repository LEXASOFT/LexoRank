package lexorank

import (
	"errors"
	"fmt"
	"strings"
)

type LexoDecimal struct {
	mag   *LexoInteger
	scale int
}

func NewLexoDecimal(mag *LexoInteger, scale int) *LexoDecimal {
	return &LexoDecimal{mag: mag, scale: scale}
}

func LexoDecimalMake(i *LexoInteger, scale int) *LexoDecimal {
	if i.IsZero() {
		return NewLexoDecimal(i, 0)
	}
	zeroCount := 0
	for idx := 0; idx < scale && i.GetMag(idx) == byte(0); idx++ {
		zeroCount++
	}
	newInteger := i.ShiftRight(zeroCount)
	newScale := scale - zeroCount
	return NewLexoDecimal(newInteger, newScale)
}

func LexoDecimalParse(str string, system LexoNumeralSystem) (*LexoDecimal, error) {
	partialIndex := strings.IndexByte(str, system.GetRadixPointChar())
	if strings.LastIndexByte(str, system.GetRadixPointChar()) != partialIndex {
		return nil, fmt.Errorf("more than one %q", system.GetRadixPointChar())
	}
	if partialIndex < 0 {
		i, err := LexoIntegerParse(str, system)
		if err != nil {
			return nil, fmt.Errorf("parse integer: %w", err)
		}
		return LexoDecimalMake(i, 0), nil
	}
	intStr := str[:partialIndex] + str[partialIndex+1:]
	i, err := LexoIntegerParse(intStr, system)
	if err != nil {
		return nil, fmt.Errorf("parse integer: %w", err)
	}
	return LexoDecimalMake(i, len(str)-1-partialIndex), nil
}

func (d *LexoDecimal) String() string {
	intStr := d.mag.String()
	if d.scale == 0 {
		return intStr
	}
	head := intStr[0]
	shift := 0
	if head == d.mag.GetSystem().GetPositiveChar() || head == d.mag.GetSystem().GetNegativeChar() {
		shift = 1
	}
	var sb strings.Builder
	radixPosition := len(intStr) - d.scale + shift
	for idx := 0; idx < radixPosition; idx++ {
		sb.WriteByte(intStr[idx])
	}
	sb.WriteByte(d.mag.sys.GetRadixPointChar())
	for idx := radixPosition; idx < len(intStr); idx++ {
		sb.WriteByte(intStr[idx])
	}
	return sb.String()
}

func (d *LexoDecimal) Sub(other *LexoDecimal) *LexoDecimal {
	thisMag := d.mag
	thisScale := d.scale
	for ; thisScale < other.scale; thisScale++ {
		thisMag = thisMag.ShiftLeft(1)
	}
	otherMag := other.mag
	for otherScale := other.scale; thisScale > otherScale; otherScale++ {
		otherMag = otherMag.ShiftLeft(1)
	}
	sub, _ := thisMag.Sub(otherMag)
	return LexoDecimalMake(sub, thisScale)
}

func (d *LexoDecimal) GetSystem() LexoNumeralSystem {
	return d.mag.GetSystem()
}

func (d *LexoDecimal) GetScale() int {
	return d.scale
}

func (d *LexoDecimal) SetScale(scale int) *LexoDecimal {
	if scale >= d.scale {
		return d
	}
	if scale < 0 {
		scale = 0
	}
	diff := d.scale - scale
	newMag := d.mag.ShiftRight(diff)
	return LexoDecimalMake(newMag, scale)
}

func (d *LexoDecimal) SetScaleWithCeiling(scale int) *LexoDecimal {
	if scale >= d.scale {
		return d
	}
	if scale < 0 {
		scale = 0
	}
	diff := d.scale - scale
	newMag := d.mag.ShiftRight(diff)
	newMag, _ = newMag.Add(lexoIntegerOne(newMag.GetSystem()))
	return LexoDecimalMake(newMag, scale)
}

func (d *LexoDecimal) Compare(other *LexoDecimal) int {
	if d == other {
		return 0
	}
	if other == nil {
		return 1
	}
	tMag, oMag := d.mag, other.mag
	if d.scale > other.scale {
		oMag = oMag.ShiftLeft(d.scale - other.scale)
	} else if d.scale < other.scale {
		tMag = tMag.ShiftLeft(other.scale - d.scale)
	}
	return tMag.Compare(oMag)
}

func (d *LexoDecimal) Add(other *LexoDecimal) *LexoDecimal {
	thisMag := d.mag
	thisScale := d.scale
	otherMag := other.mag
	var otherScale int
	for otherScale = other.scale; thisScale < otherScale; thisScale++ {
		thisMag = thisMag.ShiftLeft(1)
	}
	for ; thisScale > otherScale; otherScale++ {
		otherMag = otherMag.ShiftLeft(1)
	}
	newMag, _ := thisMag.Add(otherMag)
	return LexoDecimalMake(newMag, thisScale)
}

func (d *LexoDecimal) Multiply(other *LexoDecimal) *LexoDecimal {
	newMag, _ := d.mag.Multiply(other.mag)
	return LexoDecimalMake(newMag, d.scale+other.scale)
}

func (d *LexoDecimal) Equals(other *LexoDecimal) bool {
	if d == other {
		return true
	}
	if other == nil {
		return false
	}
	return d.scale == other.scale && d.mag.Equals(other.mag)
}

func (d *LexoDecimal) Ceil() *LexoInteger {
	if d.Exact() {
		return d.mag
	}
	floor := d.Floor()
	result, _ := floor.Add(lexoIntegerOne(floor.GetSystem()))
	return result
}

func (d *LexoDecimal) Exact() bool {
	if d.scale == 0 {
		return true
	}
	for i := 0; i < d.scale; i++ {
		if d.mag.GetMag(i) != 0 {
			return false
		}
	}
	return true
}

func (d *LexoDecimal) Floor() *LexoInteger {
	return d.mag.ShiftRight(d.scale)
}

func (d *LexoDecimal) Between(other *LexoDecimal) (*LexoDecimal, error) {
	if d.GetSystem().GetBase() != other.GetSystem().GetBase() {
		return nil, errors.New("expected same system")
	}
	left, right := d, other
	if d.GetScale() < other.GetScale() {
		right = other.SetScale(d.GetScale())
		if d.Compare(right) >= 0 {
			return d.middle(other), nil
		}
	}
	if d.GetScale() > right.GetScale() {
		left = d.SetScaleWithCeiling(right.GetScale())
		if left.Compare(right) >= 0 {
			return d.middle(other), nil
		}
	}
	var nRight *LexoDecimal
	for scale := left.GetScale(); scale > 0; right = nRight {
		nScale1 := scale - 1
		nLeft1 := left.SetScaleWithCeiling(nScale1)
		nRight = right.SetScale(nScale1)
		cmp := nLeft1.Compare(nRight)
		if cmp == 0 {
			return d.checkMid(other, nLeft1), nil
		}
		if nLeft1.Compare(nRight) > 0 {
			break
		}
		scale = nScale1
		left = nLeft1
	}
	mid := d.checkMid(other, left.middle(right))
	nScale := 0
	for mScale := mid.GetScale(); mScale > 0; mScale = nScale {
		nScale = mScale - 1
		nMid := mid.SetScale(nScale)
		if d.Compare(nMid) >= 0 || nMid.Compare(other) >= 0 {
			break
		}
		mid = nMid
	}
	return mid, nil
}

func (d *LexoDecimal) checkMid(other, mid *LexoDecimal) *LexoDecimal {
	switch {
	case d.Compare(mid) >= 0:
		return d.middle(other)
	case mid.Compare(other) >= 0:
		return d.middle(other)
	}
	return mid
}

func (d *LexoDecimal) middle(other *LexoDecimal) *LexoDecimal {
	sum := d.Add(other)
	mid := sum.Multiply(d.half())
	scale := other.GetScale()
	if d.GetScale() > other.GetScale() {
		scale = d.GetScale()
	}
	if mid.GetScale() > scale {
		roundDown := mid.SetScale(scale)
		if roundDown.Compare(d) > 0 {
			return roundDown
		}
		roundUp := mid.SetScaleWithCeiling(scale)
		if roundUp.Compare(other) < 0 {
			return roundUp
		}
	}
	return mid
}

func (d *LexoDecimal) half() *LexoDecimal {
	system := d.GetSystem()
	mid := system.GetBase() / 2
	mag := NewLexoInteger(system, 1, []byte{mid})
	return LexoDecimalMake(mag, 1)
}
