# Fase 4: CatNet Terminal UI (Bubble Tea)

## Visão Geral
Construir uma interface rica e interativa de terminal para o repositório `catnet-tui` utilizando as consagradas bibliotecas `charmbracelet/bubbletea` e `charmbracelet/lipgloss`. A interface consumirá os canais de eventos assíncronos do núcleo (`catnet-core`) garantindo uma usabilidade de ponta com zero bloqueios (non-blocking UI).

## Arquitetura Bubble Tea
O Bubble Tea segue estritamente o padrão arquitetural Elm (Model, View, Update).
Para interagir com o motor Go do CatNet (que usa `channels` para enviar descobertas de rede), construiremos **Commands** (`tea.Cmd`) que envelopam nosso motor.

### Estrutura da Aplicação
1. **Model (`tui/model.go`)**
   - Armazenará o estado: _Target Input_, lista de `results.HostResult` descobertos, percentual da barra de progresso e *flags* de estado (isScanning, err).
2. **Update (`tui/update.go`)**
   - Lida com eventos de teclado (ex: `<Enter>` para escanear, `<Esc>` para sair) e eventos customizados do sistema (ex: `hostDiscoveredMsg`, `scanProgressMsg`).
3. **View (`tui/view.go`)**
   - Renderização baseada em blocos do `lipgloss`:
     - **Header**: Área de Input estilo *Cyber*.
     - **Body**: Tabela rica listando IP, Status, Hostname e MAC (se similar ao GUI Desktop).
     - **Footer**: Barra de status / Logs de processamento.

### Integração com o `catnet-core`
1. O usuário digita o alvo (ex: `127.0.0.1`) e aperta `<Enter>`.
2. O Bubble Tea invoca a instrução `StartScanCmd`, que instancia o `scan.Engine` do core.
3. Uma `goroutine` monitora o `eventChan` do core. Cada vez que o núcleo acha um host ou avança um percentual, ela chama `program.Send(msg)` enviando um "sinal" de volta ao *Event Loop* principal da TUI.
4. O `Update` atualiza a TUI, e o `View` redesenha o frame instantaneamente.

## Proposed Changes

- **[NEW] `main.go`**: *Bootstrap* do programa via `tea.NewProgram`.
- **[NEW] `tui/model.go`**: Struct central.
- **[NEW] `tui/update.go`**: Lógica de roteamento dos eventos e controle do `catnet-core`.
- **[NEW] `tui/view.go`**: Desenho dos painéis em cores vibrantes usando `lipgloss`.

> [!NOTE]
> De acordo com as diretrizes do projeto, todas essas mudanças serão realizadas exclusivamente na branch local **`develop`** do repositório `catnet-tui`, subidas para a `develop` no remote, e somente incorporadas à `main` se validadas via *Pull Request*.

## Verification Plan
1. Rodar `go run .` na pasta `catnet-tui`.
2. Validar que a renderização inicial não quebra a tela.
3. Realizar um escaneamento em uma faixa de IPs real e garantir que a tabela seja populada dinamicamente, enquanto a UI permanece fluida e interativa.
