package rangefinder

import "fmt"
import "testing"

func main() {
	fmt.Println("vim-go")
}

func TestGetLaserDistance(t *testing.T) {
	result := GetLaserDistance(3.6, 1.38)
	fmt.Println(result)
}
