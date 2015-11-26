package Sleeper
import (
	"time"
	"testing"
)
func Test(t *testing.T)	{

	t1 := time.Now().Second()
	sleep(5)
	t2 := time.Now().Second()
	
	if (t2-t1)!=5	{
		t.Error("Expected 10, but got",t2-t1)
	}
}