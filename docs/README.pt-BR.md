# 🔄 SelfRestart

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/rafa-mori/selfrestart)](https://goreportcard.com/report/github.com/rafa-mori/selfrestart)

[🇺🇸 Read this documentation in English](../README.md)

**SelfRestart** é uma biblioteca Go que permite que aplicações se reiniciem automaticamente de forma segura e elegante. Diferente de outras soluções como `go-selfupdate`, esta biblioteca foca especificamente no **reinício automático** de processos, oferecendo controle fino sobre o ciclo de vida da aplicação.

## 🎯 Características

- ✅ **Reinício Automático**: Reinicia a aplicação preservando argumentos e ambiente
- ✅ **Detecção de Plataforma**: Suporte para Linux, macOS e Windows
- ✅ **Instalação Automática do Go**: Instala o Go automaticamente se necessário
- ✅ **Gestão de Processos**: Controle completo sobre PIDs e sinais
- ✅ **Logging Integrado**: Sistema de logs com diferentes níveis
- ✅ **Modular**: Arquitetura limpa e bem organizada
- ✅ **Thread-Safe**: Seguro para uso em aplicações concorrentes

## 📦 Instalação

```bash
go get github.com/rafa-mori/selfrestart
```

## 🚀 Uso Básico

```go
package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/rafa-mori/selfrestart"
)

func main() {
    // Criar uma instância do SelfRestart
    sr := selfrestart.New()
    
    // Verificar se Go está instalado
    if !sr.IsGolangInstalled() {
        fmt.Println("Go não está instalado")
        os.Exit(1)
    }
    
    // Configurar captura de sinais
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGUSR1)
    
    for {
        select {
        case <-sigChan:
            fmt.Println("Reiniciando aplicação...")
            if err := sr.Restart(); err != nil {
                fmt.Printf("Erro ao reiniciar: %v\\n", err)
                os.Exit(1)
            }
            os.Exit(0) // Termina o processo atual
        }
    }
}
```

## 📚 Documentação da API

### Estrutura Principal

```go
type SelfRestart struct {
    // Campos internos (não expostos)
}
```

### Métodos Principais

#### `New() *SelfRestart`
Cria uma nova instância do SelfRestart.

#### `IsGolangInstalled() bool`
Verifica se o Go está instalado no sistema e oferece instalação automática se necessário.

#### `Restart() error`
Reinicia o processo atual de forma segura.

#### `GetCurrentPID() int`
Retorna o PID do processo atual.

#### `KillCurrentProcess() error`
Finaliza o processo atual de forma graceful.

#### `IsProcessRunning(pid int) (bool, error)`
Verifica se um processo com o PID especificado está rodando.

#### `InstallGo() (bool, error)`
Instala o Go automaticamente em sistemas Unix.

#### `GetPlatformInfo() platform.PlatformInfo`
Retorna informações sobre a plataforma atual (OS/Arquitetura).

## 🏗️ Arquitetura

O projeto é organizado de forma modular:

```
selfrestart/
├── selfrestart.go          # API pública principal
├── example/
│   └── main.go            # Exemplo de uso
├── internal/
│   ├── install/           # Gestão de instalação do Go
│   ├── platform/          # Detecção de plataforma
│   ├── process/           # Gestão de processos
│   └── restart/           # Lógica de reinício
└── logger/                # Sistema de logging
```

### Módulos Internos

- **install**: Responsável pela detecção e instalação automática do Go
- **platform**: Gerencia informações de plataforma e arquitetura
- **process**: Controla processos, PIDs e sinais do sistema
- **restart**: Implementa a lógica de reinício usando scripts temporários
- **logger**: Sistema de logging integrado

## 🔧 Exemplo Avançado

Veja o arquivo `example/main.go` para um exemplo completo que demonstra:

- Captura de sinais do sistema
- Reinício automático via SIGUSR1
- Monitoramento contínuo da aplicação
- Gestão graceful de shutdown

Para executar o exemplo:

```bash
cd example
go run main.go
```

Em outro terminal, teste o reinício:

```bash
# Obter o PID da aplicação
ps aux | grep main

# Enviar sinal de reinício
kill -USR1 <PID>
```

## 🎛️ Configuração

### Variáveis de Ambiente

- `PATH`: Usado para detectar instalação do Go

### Argumentos de Linha de Comando

- `--wait`: Aguarda o script de reinício ser executado completamente

## 🧪 Testando

```bash
# Executar todos os testes
go test ./...

# Executar testes com cobertura
go test -v -race -coverprofile=coverage.out ./...

# Ver relatório de cobertura
go tool cover -html=coverage.out
```

## 🔍 Troubleshooting

### Go não encontrado
```
ERRO: Go não está instalado ou não foi encontrado no PATH
```
**Solução**: A biblioteca oferecerá instalação automática do Go.

### Permissões insuficientes
```
ERRO: Erro ao criar script de reinício: permission denied
```
**Solução**: Execute com permissões adequadas ou verifique se `/tmp` é gravável.

### Processo não reinicia
```
ERRO: Processo não finalizou após sinal de interrupção
```
**Solução**: Verifique se não há bloqueios na aplicação que impeçam o shutdown graceful.

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📋 Roadmap

- [ ] Suporte para Windows Service
- [ ] Integração com systemd (Linux)
- [ ] Backup automático antes do reinício
- [ ] Webhooks para notificações
- [ ] Interface web para monitoramento
- [ ] Métricas e monitoramento
- [ ] Suporte a containers Docker

## 📄 Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🙏 Agradecimentos

- Inspirado no projeto [go-selfupdate](https://github.com/sanbornm/go-selfupdate)
- Construído com as melhores práticas da comunidade Go
- Agradecimentos especiais aos contribuidores

## 📞 Suporte

- 🐛 **Issues**: [GitHub Issues](https://github.com/rafa-mori/selfrestart/issues)
- 💬 **Discussões**: [GitHub Discussions](https://github.com/rafa-mori/selfrestart/discussions)
- 📧 **Email**: [seu-email@exemplo.com]

---

⭐ **Se este projeto foi útil para você, considere dar uma estrela no GitHub!**
