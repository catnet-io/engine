# Arquitetura: catnet-core

O `catnet-core` é o motor compartilhado de varredura e descoberta de rede para o ecossistema CatNet. Ele foi projetado para ser leve, não possuir dependências externas de terceiros (apenas a biblioteca padrão do Go) e funcionar tanto no Windows quanto em sistemas POSIX.

A ausência de abstrações de interface desnecessárias e dependências garante que ele sirva como um núcleo de alta performance para os componentes upstream (`catnet-scanner`, `catnet-tui`, `catnet-cli`).

## Estrutura de Pacotes

A base do código é dividida em dois pacotes principais:

### 1. `pkg/scanner`
Este pacote gerencia toda a lógica de varredura, orquestração concorrente e chamadas de rede baixo-nível.
- **Orquestração Concorrente (`scan.go`)**: Implementa a função `StartScan`, utilizando um pool de goroutines parametrizado pelo `ScanConfig`. O controle de fluxo é garantido por `context.Context`, `sync.Mutex` e `atomic.Bool` para evitar scans paralelos conflitantes.
- **Primitivas de Rede (`net.go`)**: Funções unificadas e públicas para varredura (como `Ping`, `ScanPorts`, `GetMAC` e `ReverseDNS`).
- **Resolução Específica de Sistema Operacional (`os_windows.go` e `os_posix.go`)**: 
  - *Windows*: Utiliza `syscall` para acessar a API do Windows (`iphlpapi.dll` via `SendARP`) evitando o spawn de subprocessos lentos sempre que possível.
  - *POSIX*: Faz fallback para comandos do sistema operacional (`ping` e `arp` via `os/exec`) assegurando compatibilidade geral.
- **Tratamento de Ranges (`utils.go`)**: Lógica de parsing de IPs capaz de entender e traduzir sub-redes CIDR (`/24`), intervalos com hífen (`192.168.1.10-20`) e IPs unitários em blocos compatíveis para varredura.

### 2. `pkg/exporter`
Separa estritamente a varredura da serialização. Recebe as structs completas (`scanner.DeviceInfo`) e as formata em padrões de mercado, garantindo a integridade dos dados gerados.
- Formatos suportados: `JSON`, `XML` e `CSV`.
- Segurança: A função de exportação para CSV traz sanitização embutida para mitigar vulnerabilidades de Injeção de CSV (CSV Injection), filtrando prefixos de execução maliciosa originados nas resoluções de Hostname e vendor.

## Diagrama de Execução

```mermaid
flowchart TD
    App[Frontend / CLI / TUI] -->|1. Passa Lista de IPs e ScanConfig| ScannerCore(pkg/scanner)
    
    subgraph ScannerCore
      StartScan --> |Goroutine Pool| Workers
      Workers --> net.Ping
      Workers --> net.ScanPorts
      Workers --> net.GetMAC
      Workers --> net.ReverseDNS
    end

    ScannerCore -->|2. Emite onResult e onProgress callbacks| App

    App -->|3. Passa []DeviceInfo para exportação| ExporterCore(pkg/exporter)
    
    subgraph ExporterCore
      ExportJSON
      ExportXML
      ExportCSV --> |Sanitização contra CSV Injection| CSV
    end
```
