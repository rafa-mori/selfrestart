package install

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	gl "github.com/rafa-mori/selfrestart/logger"
)

type Assistant struct{}

func NewAssistant() *Assistant {
	return &Assistant{}
}

func (a *Assistant) AskUserConfirmation(message string, timeout time.Duration) bool {
	gl.Log("info", message)
	fmt.Print(message)
	responseCh := make(chan string, 1)
	defer close(responseCh)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToUpper(response))
		responseCh <- response
	}()
	select {
	case response := <-responseCh:
		return response == "Y" || response == ""
	case <-time.After(timeout):
		gl.Log("error","Timeout. No confirmation received.")
		return false
	}
}
