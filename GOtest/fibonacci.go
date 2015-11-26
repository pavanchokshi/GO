package Fibonacci

import "fmt"

func Fibo(n int) int {
	if n==1	{
		return 0;
	} else if n==2{
		return 1;
	} else	{
		return (Fibo(n-1)+Fibo(n-2));
	}
}

func main()	{
	var n int
	fmt.Println("Enter number to get required Fibonacci number")
	fmt.Scanf("%d", &n)
	fmt.Printf("The fibonacci number is: ")
	fmt.Println(Fibo(n))
}