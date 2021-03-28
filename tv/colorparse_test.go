package tv

import (
	"testing"

	"github.com/prospero78/goTV/tv/cons"
)

func TestParserEmpty(t *testing.T) {
	prs := NewColorParser("", cons.ColorBlack, cons.ColorWhite)

	elem := prs.NextElement()

	if elem.Type != ElemEndOfText {
		t.Errorf("Empty string must return end of text (%v vs %v)",
			ElemEndOfText, elem.Type)
	}
}

func TestParserColors(t *testing.T) {
	prs := NewColorParser("a<b:green>c<t:red>d<b:>e<t:>fg\nf",
		cons.ColorBlack, cons.ColorWhite)
	elems := []TextElement{
		{ElemPrintable, 'a', cons.ColorBlack, cons.ColorWhite},
		{ElemBackColor, ' ', cons.ColorBlack, cons.ColorGreen},
		{ElemPrintable, 'c', cons.ColorBlack, cons.ColorGreen},
		{ElemTextColor, 'c', cons.ColorRed, cons.ColorGreen},
		{ElemPrintable, 'd', cons.ColorRed, cons.ColorGreen},
		{ElemBackColor, 'd', cons.ColorRed, cons.ColorWhite},
		{ElemPrintable, 'e', cons.ColorRed, cons.ColorWhite},
		{ElemTextColor, 'e', cons.ColorBlack, cons.ColorWhite},
		{ElemPrintable, 'f', cons.ColorBlack, cons.ColorWhite},
		{ElemPrintable, 'g', cons.ColorBlack, cons.ColorWhite},
		{ElemLineBreak, 'g', cons.ColorBlack, cons.ColorWhite},
		{ElemPrintable, 'f', cons.ColorBlack, cons.ColorWhite},
	}

	idx := 0
	el := prs.NextElement()

	for el.Type != ElemEndOfText {
		if idx >= len(elems) {
			t.Errorf("Size mismatch: string must have only %v items", len(elems))
		}

		if el.Type != elems[idx].Type ||
			(el.Type == ElemPrintable && (el.Ch != elems[idx].Ch || el.Fg != elems[idx].Fg || el.Bg != elems[idx].Bg)) ||
			(el.Type == ElemTextColor && el.Fg != elems[idx].Fg) ||
			(el.Type == ElemBackColor && el.Bg != elems[idx].Bg) {
			t.Errorf("Elements mismatch at %v: {%v, %v, %v, %v} = {%v, %v, %v, %v}",
				idx, el.Type, elems[idx].Type,
				el.Ch, elems[idx].Ch,
				el.Fg, elems[idx].Fg,
				el.Bg, elems[idx].Bg)
		}

		el = prs.NextElement()
		idx += 1
	}
}
