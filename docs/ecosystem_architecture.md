# CatNet Ecosystem Architecture

## Visão Geral

O projeto **CatNet** evoluiu de um monólito amarrado a uma interface Desktop (Wails + React) para um ecossistema distribuído, limpo e modular. O foco da reestruturação foi isolar as regras de domínio de rede em um núcleo agnóstico, permitindo o desenvolvimento de múltiplas interfaces de consumo.

A arquitetura atual é composta por 4 repositórios principais:

| Repositório | Descrição | Papel na Arquitetura |
|-------------|------------|----------------------|
| `catnet-core` | Núcleo de domínio em Go. | Contém o motor de scanner, parsers, resoluções DNS/MAC e modelos. |
| `catnet` | Command-Line Interface (CLI) | Interface primária de terminal, foca em velocidade e scripts (DevOps). |
| `catnet-scanner`| Desktop GUI (Wails/React)| Consome o core para prover usabilidade e visualização acessível. |
| `catnet-tui` | Terminal UI (Bubble Tea)| Visão interativa para uso contínuo no terminal (Dashboard). |

---

## 1. O Núcleo: `catnet-core`

Foi criado para não conhecer ou depender de nenhuma biblioteca visual. Todo o fluxo ocorre através de pacotes tipados:

- **`pkg/scan`:** O `Engine` de varredura. Roda sobre `goroutines` e resolve ping, ports, DNS e MAC (distinguindo `os_windows` e `os_posix`).
- **`pkg/events`:** Eventos disparados pelo engine (`ScanStarted`, `HostDiscovered`, `ScanProgress`, `ScanCompleted`) entregues via _Go Channels_.
- **`pkg/results`:** Entidades canônicas de domínio como `HostResult` e `ScanResult`. Substituem as structs acopladas do Wails (ex: `DeviceInfo`).
- **`pkg/targets`:** Responsável pelo parser inteligente das entradas (`ParseRange`), mapeando desde IPs únicos até faixas de CIDR ou Dash Ranges.
- **`pkg/export`:** Exportadores em JSON e CSV, com sanitização inclusa.

### Padrão de Integração
Qualquer consumidor (CLI, GUI) deve invocar o scanner via `engine.ScanStream()` fornecendo um channel local `chan<- events.Event`. O consumidor implementa sua própria `goroutine` para ler este canal e retransmitir/desenhar na sua interface respectiva.

---

## 2. A Ponte Desktop: `catnet-scanner`

Originalmente continha os pacotes `pkg/scanner` e `pkg/exporter`.
**Pós-Refatoração:**
- Os pacotes físicos internos foram deletados.
- Passou a importar o módulo `github.com/mendsec/catnet-core`.
- O `app.go` atua puramente como *bindings*. Ele traduz a chamada de `StartScan` da UI React, invoca o core remotamente, consome o *channel* de eventos e chama o `runtime.EventsEmit(ctx)` nativo do Wails.
- O resultado é uma interface fluida, onde o Go envia os dados para o Javascript sob demanda.

---

## 3. O Cliente CLI: `catnet`

Implementado com `spf13/cobra`.
**Comportamento:**
- Usa o paradigma de subcomandos (`catnet scan`, `catnet version`).
- Trata os argumentos primários, submete ao `catnet-core`, traduz os *events* como saídas textuais amigáveis de terminal (`fmt.Printf`) com progressão.
- Útil para CI/CD ou uso via SSH.

---

## Estratégia de Desenvolvimento Local (Workspace)

Para permitir agilidade sem gerar atritos com versionamento prematuro (Git Tags), os diretórios `catnet` e `catnet-scanner` estão configurados localmente com a diretiva `replace` no arquivo `go.mod`:
```go
require github.com/mendsec/catnet-core v0.0.0
replace github.com/mendsec/catnet-core => ../catnet-core
```
Esta configuração obriga o compilador Go a usar o diretório local relativo ao invés de baixar a dependência da nuvem, sendo ideal para desenvolvimento iterativo do Ecossistema.

> **Data da Refatoração:** Junho de 2026.
> **Autor:** CatNet AI / Fábio Mendes.
