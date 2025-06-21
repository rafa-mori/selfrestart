package selfrestart

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rafa-mori/selfrestart/internal/install"
	"github.com/rafa-mori/selfrestart/internal/platform"
	"github.com/rafa-mori/selfrestart/internal/process"
	"github.com/rafa-mori/selfrestart/internal/restart"
	gl "github.com/rafa-mori/selfrestart/logger"
)

// SelfRestart provides functionality for automatic process restart
type SelfRestart struct {
	installer *install.Installer
	manager   *process.ProcessManager
	restarter *restart.Restarter
}

// New creates a new SelfRestart instance
func New() *SelfRestart {
	return &SelfRestart{
		installer: install.NewInstaller(),
		manager:   process.NewProcessManager(),
		restarter: restart.NewRestarter(),
	}
}

// Module provides access to module information
var Module = &ModuleInfo{
	commandEntry: "go",
	repoName:     "rafa-mori/selfrestart",
}

// ModuleInfo contains module-specific information
type ModuleInfo struct {
	commandEntry string
	repoName     string
}

// GetCommandEntry returns the command entry point
func (m *ModuleInfo) GetCommandEntry() string {
	return m.commandEntry
}

// GetRepoName returns the repository name
func (m *ModuleInfo) GetRepoName() string {
	return m.repoName
}

// IsGolangInstalled checks if Go is installed and prompts for automatic installation if needed
func (sr *SelfRestart) IsGolangInstalled() bool {
	isInstalled, err := sr.installer.IsInPath(Module.GetCommandEntry())
	if err != nil {
		gl.Log("error", fmt.Sprintf("Erro ao verificar instalação do Go: %v", err))
		return false
	}
	
	if isInstalled {
		return true
	}

	return sr.promptForGoInstallation()
}

// promptForGoInstallation displays installation prompt and handles user response
func (sr *SelfRestart) promptForGoInstallation() bool {

	// Exibe uma mensagem de aviso ao usuário
	fmt.Printf("\033[1;36m%s\033[0m", `
┌──────────────────────────────────────────────────────────────┐
│                                                              │
│  O Go não está instalado ou não foi encontrado no PATH.      │
│                                                              │
│  Para continuar, precisamos instalá-lo.                      │
│  Podemos instalar o Go automaticamente, mas você precisa     │
│  confirmar se deseja prosseguir.                             │
│                                                              │
│  Se você deseja instalar o Go automaticamente,               │
│  pressione 'Y' e 'Enter'.                                    │
│  Caso contrário, pressione 'N' e 'Enter' para sair.          │
│                                                              │
└──────────────────────────────────────────────────────────────┘`)
	fmt.Printf("\033[1;36m%s\033[0m", `

  Se você não tem certeza se deseja instalar o Go,            
  pode instalar a versão já compilada diretamente do repo:    
                                                              
  https://github.com/`+Module.GetRepoName()+`/releases/latest   

  `)

	fmt.Printf("\033[1;33m%s\033[0m", `
┌──────────────────────────────────────────────────────────────┐
│                                                              │
│  Pressione 'Y' para sim, instalar o Go automaticamente,      │
│  ou 'N' para sair do programa. (15 segundos para resposta)   │
│                                                              │
└──────────────────────────────────────────────────────────────┘
`)

	fmt.Print("Resposta: ")

	// Canal para receber a resposta do usuário
	// O canal é usado para evitar o bloqueio do terminal enquanto espera a entrada do usuário
	// e permite que o programa continue executando.
	responseCh := make(chan string, 1)
	defer close(responseCh)

	// Lê a resposta do usuário em uma goroutine
	// Isso permite que o programa continue executando enquanto espera a entrada do usuário
	// e evita o bloqueio do terminal. O usuário pode pressionar 'Y' ou 'N', caso contrário
	// o programa irá aguardar 15 segundos antes de encerrar.
	go func() {
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToUpper(response))
		responseCh <- response
	}()

	// Aguarda a resposta do usuário ou timeout de 15 segundos
	select {
	case response := <-responseCh:
		return response == "Y"
	case <-time.After(15 * time.Second):
		fmt.Println("\nTempo esgotado. Docker não será ativado.")
		return false
	}
}

// Restart restarts the current process
func (sr *SelfRestart) Restart() error {
	binPath, err := sr.getCurrentBinaryPath()
	if err != nil {
		return fmt.Errorf("erro ao obter caminho do binário atual: %v", err)
	}

	pid := sr.manager.GetCurrentPID()
	if pid <= 0 {
		return fmt.Errorf("PID inválido: %d", pid)
	}

	gl.Log("info", fmt.Sprintf("Reiniciando processo %d com binário %s", pid, binPath))

	// Cria e executa o script de reinício
	if err := sr.restarter.CreateAndExecRestartScript(pid, binPath); err != nil {
		return fmt.Errorf("erro ao criar e executar script de reinício: %v", err)
	}

	return nil
}

// GetCurrentPID returns the current process ID
func (sr *SelfRestart) GetCurrentPID() int {
	return sr.manager.GetCurrentPID()
}

// KillCurrentProcess kills the current process
func (sr *SelfRestart) KillCurrentProcess() error {
	return sr.manager.KillCurrentProcess()
}

// IsProcessRunning checks if a process with the given PID is running
func (sr *SelfRestart) IsProcessRunning(pid int) (bool, error) {
	return sr.manager.IsProcessRunning(pid)
}

// InstallGo installs Go on Unix systems
func (sr *SelfRestart) InstallGo() (bool, error) {
	return sr.installer.InstallGoUnix()
}

// GetPlatformInfo returns information about the current platform
func (sr *SelfRestart) GetPlatformInfo() platform.PlatformInfo {
	return platform.GetHostPlatform()
}

// getCurrentBinaryPath returns the path to the current binary
func (sr *SelfRestart) getCurrentBinaryPath() (string, error) {
	binPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("erro ao obter caminho do binário atual: %v", err)
	}
	return binPath, nil
}

func init() {
	gl.Log("notice", "Inicializando módulo de reinício automático")
	gl.SetDebug(true) // Ativa o modo de depuração para logs
	gl.Log("notice", "Módulo de reinício automático inicializado com sucesso")
	gl.Log("info", "selfrestart Versão: v1.0.0")
}
