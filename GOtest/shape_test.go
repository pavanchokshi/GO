package Shape
import "testing"

func TestRectPeri(t *testing.T)	{
	var s Shaper
	r := Rect{3,5}
	s = r
	var rectPeri, wantRectPeri float64
	rectPeri = s.Peri()
	wantRectPeri = 16
	if rectPeri != wantRectPeri	{
		t.Errorf("Expected %.2f, got %.2f",wantRectPeri,rectPeri)
	}	
}

func TestCircPeri(t *testing.T)	{
	var s Shaper
	c := Circ{5}
	s = c
	var CircPeri, wantCircPeri float64
	CircPeri = s.Peri()
	wantCircPeri = 31.4
	if CircPeri > wantCircPeri	{
		t.Errorf("Expected %.2f, got %.2f",wantCircPeri,CircPeri)
	}	
}