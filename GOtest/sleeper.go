package Sleeper

import(
	"fmt"
	"time"
)

func sleep(duration int) {

	<-time.After(time.Second*time.Duration(duration))
}

func main()	{
	t1 := time.Now().Second()
	fmt.Print(t1)
	fmt.Print("\nUnder sleep...")
	sleep(10)
	t2 := time.Now().Second()
	fmt.Print(\nt2)
	if t2-t1 == 10{
		fmt.Print("\nWait Success!")
	}
}