package Fibonacci

import "testing"

func Test(t *testing.T)	{
	
	num := Fibo(4)
	
	if num!=2
		t.Error("Expected 2, but got",num)
	}
}