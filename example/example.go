package main

import (
	"fmt"
	"path/filepath"

	"github.com/getlantern/golog"
	"github.com/getlantern/sysproxy"
)

var log = golog.LoggerFor("example")

func main() {
	helperFullPath := "sysproxy-cmd"
	iconFullPath, _ := filepath.Abs("./icon.png")
	log.Debugf("Using icon at %v", iconFullPath)
	err := sysproxy.EnsureHelperToolPresent(helperFullPath, "Input your password and save the world!", iconFullPath)
	if err != nil {
		fmt.Printf("Error EnsureHelperToolPresent: %s\n", err)
		return
	}
	off, err := sysproxy.On("localhost:12345")
	if err != nil {
		fmt.Printf("Error set proxy: %s\n", err)
		return
	}
	fmt.Println("proxy set, hit enter to continue (or kill the parent process)...")
	var i int
	fmt.Scanf("%d\n", &i)
	off()
}
