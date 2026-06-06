# Relatório de Análise Profunda: catnet-core

Este documento apresenta uma análise técnica e estratégica do projeto `catnet-core`, o motor de escaneamento de rede compartilhado para o ecossistema CatNet.

---

## 1. ANÁLISE DE DESENVOLVIMENTO

### Arquitetura e Design
O projeto adota uma arquitetura modular de biblioteca padrão em Go, desenhada para ser o motor (engine) agnóstico de interface para múltiplos clientes. O padrão principal no fluxo de escaneamento é o de um **Pipeline Concorrente com Worker Pool**, coordenado via `context.Context` para timeout e controle de vida (cancellation), utilizando canais (channels) e `sync.WaitGroup`.

**Análise dos Pacotes:**
- **`pkg/engine`**: Orquestrador principal. O `StartScan` consome alvos de um canal e distribui para uma pool de goroutines, emitindo eventos de progresso via callback.
- **`pkg/results`**: Define a tipagem canônica do sistema (`ScanReport`, `DeviceInfo`), atuando como contrato principal para o ecossistema.
- **`pkg/targets`**: Lida com o *parsing* de IPs, com suporte a notações CIDR e ranges (e.g. dashboards) focado em isolamento da lógica de entrada.
- **`pkg/discovery`**: Encapsula as primitivas de resolução e "liveness" (Ping, GetMAC, ReverseDNS), suportando separação por SO com `build tags` (como `os_windows.go` e `os_posix.go`).
- **`pkg/ports`**: Responsável pelo port scanning.
- **`pkg/exporter`**: Centraliza os formatos de saída (JSON, CSV, XML). Nota positiva para o fato de o JSON ser considerado o "reference schema format".

**Separação internal/ vs pkg/:**
Existe uma divisão estrita. O diretório `pkg/` abriga os contratos de API suportados, enquanto o `internal/` contém apenas `netutil`, com utilitários de validação de IP sem garantia de estabilidade pública. O uso anterior do `pkg/scanner` foi completamente refatorado para pacotes de domínio granulares; o pacote permanece apenas como shim de retrocompatibilidade **deprecated**.

### Qualidade de Código
- **Legibilidade**: Excelente. As responsabilidades estão bem distribuídas. 
- **Segurança**: Existe uma prevenção muito clara e documentada contra **CSV Injection** na função `sanitizeCSVField` do pacote `exporter`. 
- **Error Handling**: A base se apoia em wrap errors (`fmt.Errorf("%w", ...)`), utilizando um pacote dedicado `coreerr` para criar uma taxonomia de erros estruturada (`coreerr.ErrTimeout`, `coreerr.ErrCancelled`, `coreerr.ErrExport`), facilitando a leitura de erros via `errors.Is`.
- **Validações Limite**: No `pkg/engine/scan.go`, o `StartScan` impõe um teto rígido (`maxAllowedThreads = 256`) e recalcula timeouts defensivos para garantir estabilidade do hardware.
- **Comentários**: Códigos essenciais estão documentados. Presença correta do arquivo `tests/doc.go`.

### Infraestrutura de Desenvolvimento
- **Dependências (`go.mod`)**: O projeto requer `go 1.26.4` e, de forma impressionante, não depende de bibliotecas externas (zero-dependency core), o que é um fator gigante de segurança.
- **CI/CD (`.github/workflows`)**: Existem pipelines bem configurados: `ci.yml` rodando `go vet` e `go test -race -v`, garantindo checagem de concorrência. Adicionalmente, há o workflow semanal de `govulncheck` auditando falhas.
- **Testes**: Há forte cultura de testes visível na separação de suítes de testes nos pacotes locais, além de testes E2E/Integração robustos em `tests/integration_test.go` valendo-se do diretório de fixtures `testdata/`.
- **Padrões de Ferramentas**: Conta com configuração estrita de formatação (tab em Go) via `.editorconfig`.

### Documentação
- **`README.md`** & **`ROADMAP.md`**: Claríssimos. Deixam claro os *Non-Goals* (sem integração nativa Cloud, sem exploração, sem bindings de GUI).
- **Contratos (`docs/contracts/`)**: Um verdadeiro diferencial para projetos open-source em versão `0.x.x`. O `compatibility.md` exige que clientes da lib limitem-se a *minor versions* (`v0.1.x`) e tratem JSON fields desconhecidos preventivamente, e o `api-stability.md` define muito bem quais pacotes são intocáveis e quais mudarão.
- **Políticas (`SECURITY.md`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`)**: Diretrizes claras, exigindo abertura de *issues* antes de grandes PRs, além de proibir reporte de vulnerabilidades públicas.

### Performance e Escalabilidade
- **Pre-allocation**: Capacidade pre-alocada de arrays no engine (`report.Devices = make([]results.DeviceInfo, 0, total)`) e tamanho fixo do canal previnem overhead do garbage collector durante escalonamento.
- **Gargalos Atuais**: Conforme os commits e o Roadmap mostram, o escaneamento de portas por host (`ports.ScanPorts` roda serialmente bloqueando a goroutine do worker) gera *worker starvation*. Além disso, os limits de sockets de SO afetam o scanning.
- **Ações de Performance Planejadas**: Implementar Raw Sockets (`M5` no ROADMAP) para ignorar os limites de rede do OS e o port scanner paralelo assíncrono.

---

## 2. ANÁLISE DE NEGÓCIO

### Propósito e Valor do Produto
O `catnet-core` é a espinha dorsal de processamento e descoberta de rede dentro de um ecossistema projetado com um objetivo conciso: prover capacidade de identificação profunda na rede *sem misturar lógica de apresentação*. Seu valor único reside em ser **embedável**, padronizado e determinístico na sua execução e falha (sem vazamentos de goroutines).

### Ecossistema e Integrações
A fragmentação é o trunfo de design (Micro-repositórios baseados em módulos Go).
| Repositório | Papel | Stack Tecnológica |
| --- | --- | --- |
| **`catnet-core`** | Backend de processamento | Go puro |
| **`catnet-scanner`**| Frontend gráfico Desktop | Go + Wails + React |
| **`catnet`** | CLI rápida para SysAdmins | Go p/ Scripts |
| **`catnet-tui`** | UI interativa no terminal | Go + Bubble Tea |

### Modelo de Mercado
- **Mercado Alvo**: Profissionais de cibersegurança, SysAdmins, Blue/Red Teams operando auditoria/asset discovery.
- **Posicionamento**: Como uma alternativa "developer-friendly" ao `nmap` ou ao `masscan`, com foco em facilitar a construção de subprodutos e integrações com Go moderno em pipelines DevOps. Ele compete no mercado de enumeração de open-source com alta extensibilidade.

### Métricas de Adoção & Sustentabilidade
- **Adoção**: Estágio de "lançamento puro". 0 Stars, 0 Forks, lançamento da `v0.1.0` há apenas 3 dias.
- **Risco de Comunidade**: Alto risco *Bus Factor* com apenas **1 Contribuidor** (@mendsec). É essencial fomentar governança. 

### Estratégia de Open Source
A estratégia é blindar a base e evitar "bloatware" (ferramentas pesadas e monolíticas). A restrição de integrações AWS/GCP e exploração de vulnerabilidade ao core mostra maturidade para manter uma biblioteca de rede universal sem complexidade extra, permitindo adoção de ecossistema com plugins adjacentes no futuro.

---

## 3. RECOMENDAÇÕES

### 🛠️ Para Desenvolvimento
- **Prioridade 1: Port Scanner Assíncrono:** Executar a tarefa M5 do Roadmap. O scanner de portas *por-host* reduz criticamente o throughput no worker-pool. Transforme-o em uma pool focada por portas.
- **Prioridade 2: Fuzzing no `ParseRange`:** Introduzir testes de Fuzz nativos em Go (`go test -fuzz`) no `pkg/targets` para garantir que ranges CIDR malformados nunca causem panics indesejados no core.
- **Benchmarking Extensivo:** Adicionar no `tests/` suítes com arquivos tipo `targets_huge.txt` e testes tipo `BenchmarkStartScan` verificando tempos e taxa de alocação `allocs/op`.
- **Limpeza de Refatoração:** Remover as pastas internas sem utilidade que estão vazias ou completar a documentação do `internal/netutil`. 

### 💼 Para Negócio
- **Geração de Valor / Monetização:** Criar um tier **Enterprise** oferecendo suporte a integrações diretas (SOCs) ou então vender a aplicação Desktop (`catnet-scanner`) baseada neste core aberto, semelhante ao que a Tailscale faz em cima do Wireguard.
- **Prospecção:** Divulgar `catnet-core` no Hacker News ("Show HN") e nos subreddits focados em Go e netsec como ferramenta de base de arquitetura limpa. Focar nas qualidades de "zero dependency" e na proteção por "sandboxed concurrency timeouts".
- **Casos de Uso Corporativos:** Criar guias em `docs/examples` integrando a saída JSON (`SchemaVersion 1.0.0`) do core no ElasticSearch, Splunk ou Fluentd para Security Observability.

### 👥 Para Comunidade
- **Facilitar o Onboarding:** Criar *issues* com a tag `good first issue` (Ex: adicionar documentação faltante num pacote, melhorar uma mensagem de erro na flag).
- **Contratos Visualizáveis:** Montar um site de documentação (e.g. Docusaurus ou MkDocs) baseado nos markdowns criados em `docs/contracts/` provando a solidez da API, gerando credibilidade.
- **Divisão de Responsabilidades:** Designar *CODEOWNERS* para os subprojetos para estruturar uma base segura à medida que a comunidade chega.

---

> [!TIP]
> **Resumo Geral:** O `catnet-core` apresenta uma das mais rigorosas documentações de conformidade e arquitetura defensiva (timeouts forçados, thread limits, injeção de erros controlados) vista em repositórios early-stage `0.x.x`. O sucesso dependerá apenas da superação do gargalo do SO (Raw Sockets) e da conversão do uso CLI em estrelas no Github.
