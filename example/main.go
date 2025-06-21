package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rafa-mori/selfrestart"
)

func main() {
	fmt.Println("🚀 SelfRestart Example Application")
	fmt.Println("===================================")

	// Criar uma instância do SelfRestart
	sr := selfrestart.New()

	// Verificar se Go está instalado
	if !sr.IsGolangInstalled() {
		fmt.Println("❌ Go não está instalado ou a instalação falhou")
		os.Exit(1)
	}

	fmt.Println("✅ Go está instalado e funcionando")

	// Mostrar informações da plataforma
	platformInfo := sr.GetPlatformInfo()
	fmt.Printf("🖥️  Plataforma: %s/%s\n", platformInfo.OS, platformInfo.Arch)

	// Mostrar PID atual
	pid := sr.GetCurrentPID()
	fmt.Printf("🆔 PID atual: %d\n", pid)

	// Configurar um canal para capturar sinais de sistema
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

	fmt.Println("\n📋 Comandos disponíveis:")
	fmt.Println("  - Ctrl+C: Sair normalmente")
	fmt.Println("  - SIGUSR1 (kill -USR1 <pid>): Reiniciar aplicação")
	fmt.Println("\n⏳ Aplicação rodando... Pressione Ctrl+C para sair ou envie SIGUSR1 para reiniciar")

	// Contador para demonstrar que a aplicação está rodando
	counter := 0
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case sig := <-sigChan:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				fmt.Println("\n👋 Recebido sinal de saída. Finalizando...")
				return
			case syscall.SIGUSR1:
				fmt.Println("\n🔄 Recebido sinal de reinício. Reiniciando aplicação...")
				if err := sr.Restart(); err != nil {
					log.Fatalf("❌ Erro ao reiniciar: %v", err)
				}
				// Após chamar Restart(), o processo atual deve terminar
				fmt.Println("👋 Processo atual finalizando para permitir reinício...")
				os.Exit(0)
			}
		case <-ticker.C:
			counter++
			fmt.Printf("⏰ Aplicação rodando há %d ciclos (PID: %d)\n", counter, pid)
		}
	}
}
