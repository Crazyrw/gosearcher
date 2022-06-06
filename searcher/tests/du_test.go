package tests

import (
	"fmt"
	"github.com/ricochet2200/go-disk-usage/du"
	"testing"
)

func TestDu(t *testing.T) {
	usage := du.NewDiskUsage("E:/")
	fmt.Println(usage.Used())
	fmt.Println(usage.Available())
	fmt.Println(usage.Usage())
	fmt.Println(usage.Size())
}
