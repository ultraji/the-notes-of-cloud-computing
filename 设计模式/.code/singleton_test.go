package singleton

import "testing"

func TestSingleton1(t *testing.T) {
	sgt1 := GetInstance1()
	sgt2 := GetInstance1()
	if sgt1 == sgt2 {
		t.Logf("sgt1 := GetInstance1() // sgt1(%T) = %p", sgt1, sgt1)
		t.Logf("sgt2 := GetInstance1() // sgt2(%T) = %p", sgt2, sgt2)
	}
}

func TestSingleton2(t *testing.T) {
	sgt1 := GetInstance2()
	sgt2 := GetInstance2()
	if sgt1 == sgt2 {
		t.Logf("sgt1 := GetInstance2() // sgt1(%T) = %p", sgt1, sgt1)
		t.Logf("sgt2 := GetInstance2() // sgt2(%T) = %p", sgt2, sgt2)
	}
}

func TestSingleton3(t *testing.T) {
	sgt1 := GetInstance3()
	sgt2 := GetInstance3()
	if sgt1 == sgt2 {
		t.Logf("sgt1 := GetInstance3() // sgt1(%T) = %p", sgt1, sgt1)
		t.Logf("sgt2 := GetInstance3() // sgt2(%T) = %p", sgt2, sgt2)
	}
}

func TestSingleton4(t *testing.T) {
	sgt1 := GetInstance4()
	sgt2 := GetInstance4()
	if sgt1 == sgt2 {
		t.Logf("sgt1 := GetInstance4() // sgt1(%T) = %p", sgt1, sgt1)
		t.Logf("sgt2 := GetInstance4() // sgt2(%T) = %p", sgt2, sgt2)
	}
}
