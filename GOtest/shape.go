package main
import	(
	"fmt"
	"math"
)

type Shaper interface	{
	
	Area() float64
	Perimeter() float64
}

type Rect struct	{
	length, breadth float64
}

// Area function for Rectangle 
func (r Rect) Area() float64	{
	return r.length*r.breadth
}
// Perimeter function for Rectagle 
func (r Rect) Perimeter() float64	{
	return 2*(r.length+r.breadth)
}

type Circle struct	{
	radius float64
}

// Area function for Circle 
func (c Circle) Area() float64	{
	return math.Pi*c.radius*c.radius
}
// Perimeter function for Circle 
func (c Circle) Perimeter() float64	{
	return 2*math.Pi*c.radius
}

func main()	{
	var s Shaper
	r := Rect{3,5}
	c := Circle{4}
	s = r
	fmt.Printf("Rectangle with length: %.2f and breadth: %.2f has area:%.2f\n",r.length,r.breadth,s.Area())	
	fmt.Printf("Rectangle with Length=%.2f and Breadth=%.2f has perimieter:%.2f\n",r.length,r.breadth,s.Perimeter())	
	
	//assigning the circle refernce type to Shaper variable 
	s = c
	fmt.Printf("Circle with Radius=%.2f has area:%.2f\n",c.radius,s.Area())	
	fmt.Printf("Circle with Radius=%.2f is:%.2f has perimeter\n",c.radius,s.Perimeter())
}