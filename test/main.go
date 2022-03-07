package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// func main() {
// 	loger := logrus.New()
// 	loger.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
// 	loger.SetOutput(os.Stdout)
// 	// rescueStdout := os.Stdout
// 	// r, w, _ := os.Pipe()
// 	// os.Stdout = w
// 	// w.Close()
// 	// out, _ := ioutil.ReadAll(r)
// 	// os.Stdout = rescueStdout
// 	// fmt.Printf("Captured: %s", out) // prints: Captured: Hello, playground

// 	test(loger.Info, "To capture")
// }

// func test(f func(...interface{}), message string) {
// 	rescueStdout := os.Stdout
// 	r, w, _ := os.Pipe()
// 	os.Stdout = w
// 	w.Close()
// 	out, _ := ioutil.ReadAll(r)
// 	os.Stdout = rescueStdout
// 	fmt.Printf("Captured: %s", out) // prints: Captured: Hello, playground
// 	f(message)
// 	f(message)

// }

func Print(loger *logrus.Logger, w io.Writer, f func(...interface{}), message string) {
	loger.SetOutput(w)
	f(message)
}
func PrintToDiscord(loger *logrus.Logger, f func(...interface{}), message string) {
	var b bytes.Buffer
	Print(loger, &b, f, message)
	fmt.Printf("string is : %v", b.String())

}

func PrintToStdotu(loger *logrus.Logger, f func(...interface{}), message string) {
	Print(loger, os.Stdout, f, message)
}

func PrintToDiscordAndStdOut(loger *logrus.Logger, f func(...interface{}), message string) {
	PrintToStdotu(loger, f, message)
	PrintToDiscord(loger, f, message)
}

func main() {
	loger := logrus.New()
	loger.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	PrintToDiscordAndStdOut(loger, loger.Info, "testing")

	// fmt.Println("print with os.Stdout:")
	// PrintToDiscordAndStdOut(loger, os.Stdout, loger.Info, "testing")
}
