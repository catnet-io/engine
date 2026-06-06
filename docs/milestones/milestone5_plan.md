# CatNet Core - Milestone 5 Plan (Consumers & Asynchronous Architecture)

## Objetivos do Milestone 5

Após o hardening da base do projeto (Milestone 4), o foco deste Milestone é preparar a arquitetura do motor de varredura para cenários de alto desempenho com múltiplos consumidores UI/CLI, onde as garantias de não-bloqueio se tornam essenciais.

### 1. Assincronicidade Genuína de Eventos
**Problema Atual:** O `EventCallback` é acionado de forma síncrona pela própria goroutine de processamento. Consumidores pesados (ex: uma UI em Wails ou uma renderização complexa de TUI) que segurem o callback podem atrasar as goroutines de varredura ativas, gerando starvation nas threads que poderiam estar trabalhando em novos hosts.
**Ações Propostas:**
- Implementar um mecanismo interno em `pkg/engine` utilizando canais cacheados (buffered channels) para despachar eventos.
- Lançar uma goroutine de "event dispatcher" separada que drena o canal e chama os `EventCallback` do consumidor de forma desacoplada da fila de rede.
- Adicionar verificações em testes para garantir que mesmo com `time.Sleep` dentro do callback, o motor atinge sua velocidade nominal de varredura.

### 2. Validação Fina de Context Cancellation
**Problema Atual:** O contexto aborta a varredura e evita novos processamentos, mas o sistema de log e rastreio não comprova estatisticamente que requisições TCP residuais são destruídas rapidamente ao receber `context.Canceled`.
**Ações Propostas:**
- Propagar `context.Context` explícito para chamadas de `discovery` (ICMP/DNS) e `ports` (TCP Dial) para hard-cancellation.
- Escrever testes que utilizam contextos cancelados logo na inicialização ou no meio do processamento, verificando contagem de conexões de socket ativas para certificar zero vazamento.

### 3. Integração e Rollout aos Consumidores Oficiais
**Problema Atual:** A CLI (`catnet`), Desktop GUI (`catnet-scanner`) e TUI (`catnet-tui`) possivelmente utilizam abstrações duplicadas ou estão defasadas em relação às otimizações estruturais que fizemos em structs como `DeviceInfo` (como a remoção do campo `OpenPortsCount`).
**Ações Propostas:**
- Realizar a atualização dos repositórios consumidores apontando o `go.mod` para a nova tag deste core.
- Substituir lógicas redundantes nas camadas clientes para favorecer exclusivamente o motor.
- Adequar todas as referências ao campo `OpenPortsCount` para utilizarem o método `.PortCount()`.

## Critérios de Aceite
- Múltiplos testes de estresse demonstram que gargalos no callback do cliente impactam em no máximo 5% o tempo base de varredura total.
- Vazamento de goroutines e sockets provado falso através da métrica do pacote `net/http/pprof` durante stress-testes com cancelamento prematuro.
- As ferramentas oficiais de consumo `catnet` funcionam com todas as validações end-to-end com sucesso.
