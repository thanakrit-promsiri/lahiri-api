package calculator

import (
	"testing"
)

func TestGetNakshatra(t *testing.T) {
	// ทดสอบที่ 0 องศา (ควรเป็นนักษัตรแรก Ashwini, บาทที่ 1, เจ้าเรือน Ketu)
	info, pada, offset := GetNakshatra(0.0)

	if info.Index != 1 {
		t.Errorf("Expected index 1, got %d", info.Index)
	}
	if info.Name != "Ashwini" {
		t.Errorf("Expected nakshatra 'Ashwini', got '%s'", info.Name)
	}
	if info.Lord != "Ketu" {
		t.Errorf("Expected lord 'Ketu', got '%s'", info.Lord)
	}
	if pada != 1 {
		t.Errorf("Expected pada 1, got %d", pada)
	}
	if offset != 0.0 {
		t.Errorf("Expected offset 0.0, got %f", offset)
	}

	// ทดสอบกรณีค่าลองจิจูดเกิน 360 หรือขอบเขตสูงสุด (Revati)
	infoMax, _, _ := GetNakshatra(360.0)
	if infoMax.Name != "Revati" {
		t.Errorf("Expected max nakshatra 'Revati', got '%s'", infoMax.Name)
	}
	if infoMax.Lord != "Mercury" {
		t.Errorf("Expected lord 'Mercury' for Revati, got '%s'", infoMax.Lord)
	}
}

func TestFormatNakshatraString(t *testing.T) {
	// ทดสอบการจัดรูปแบบข้อความที่ 0 องศา
	result := FormatNakshatraString(0.0)
	expected := "Ashwini (Pad 1) [Lord: Ketu]"

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}
