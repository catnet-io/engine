# Hardening & API Readiness for catnet-core

Este documento descreve as mudanças necessárias para preparar a engine `catnet-core` para consumo no ecossistema CatNet.

## Resumo do Plano

O plano está dividido em 10 tarefas conforme solicitado, focando em segurança, estabilidade da API e documentação.

### 1. Correção de Bug de Ponteiro no EventCallback
- **Problema:** A goroutine em `StartScan` passa um ponteiro da variável do loop (`&di`) para os callbacks síncronos, o que causa memory aliasing se esses eventos forem manipulados assincronamente no futuro.
- **Solução:**
  - Em `pkg/engine/scan.go`, criar uma cópia explícita da variável (`diCopy := di`) antes de despachar o `EventResult`.
  - Passar `nil` para `EventProgress` conforme padrão de uso atual.
  - Adicionar um caso de teste em `pkg/engine/scan_test.go` para validar que ponteiros recebidos no callback não sofrem aliasing se salvos para verificação posterior.

### 2. Remoção de `OpenPortsCount` de `DeviceInfo`
- **Problema:** `OpenPortsCount` e `len(OpenPorts)` carregam a mesma informação, o que gera risco de dessincronização e confusão na desserialização.
- **Solução:**
  - Remover `OpenPortsCount` da struct `DeviceInfo` em `pkg/results/device.go`.
  - Adicionar o método `PortCount() int` em `DeviceInfo`.
  - Remover a lógica de atribuição em `pkg/engine/scan.go`.
  - Limpar as fixtures (`testdata/scan_report_fixture.json`) e testes (`pkg/exporter/exporter_test.go`).
  > [!WARNING]
  > Isso introduzirá uma *Breaking Change* na API pública e será documentado apropriadamente no Changelog.

### 3. Atualização da Versão do Go no CI
- **Problema:** Referência a uma versão inexistente de Go (`1.25.x`) no CI.
- **Solução:**
  - Atualizar `.github/workflows/ci.yml`, `.github/workflows/govulncheck.yml` e `go.mod` para a versão estável atual do Go: `1.23.x` (1.23.0 no go.mod).

### 4. Documentação de Pacotes (`doc.go`)
- **Problema:** Ausência de `doc.go` nos pacotes do projeto prejudica a visibilidade no `pkg.go.dev`.
- **Solução:**
  - Adicionar 8 arquivos `doc.go` detalhando o propósito de cada pacote (`engine`, `results`, `targets`, `discovery`, `ports`, `exporter`, `coreerr`, `internal/netutil`).

### 5. Testes Unitários e Validação em `pkg/discovery`
- **Problema:** As funções em `pkg/discovery/net.go` não possuem testes determinísticos para validação de entrada sem overhead de I/O de rede.
- **Solução:**
  - Criar `pkg/discovery/net_test.go` cobrindo cenários de invalid input (`""`, IPs mal formados, IPv6 não suportados para MAC).

### 6. Timeout Defensivo em `StartScan`
- **Problema:** Cálculo de timeout superestima o tempo real pois não considera o semáforo que limita as goroutines de ScanPort.
- **Solução:**
  - Extrair constante `ScanConcurrency = 10` em `pkg/ports/scanner.go`.
  - Ajustar a fórmula de timeout em `pkg/engine/scan.go` para `portBatches := math.Ceil(len(ports)/10)`.

### 7. Proteção de Parâmetros de Timeout (`osPing`)
- **Problema:** `osPing` assume `timeoutMs > 0`, causando erros em chamadas diretas com limite de tempo incorreto.
- **Solução:**
  - Em `pkg/discovery/os_windows.go`, aplicar um valor padrão de 1000ms se `timeoutMs <= 0`.
  - Em `pkg/discovery/os_posix.go`, converter corretamente milissegundos para os segundos inteiros que a flag `-W` necessita, ignorando a constante fixa de `1` (que era um bug semântico de longo prazo).

### 8. Padronização do Changelog
- **Problema:** O formato atual não segue "Keep a Changelog".
- **Solução:**
  - Reestruturar `CHANGELOG.md` para suportar tags padrão `Added`, `Changed`, `Fixed`, e documentar a breaking change do passo 2 sob `[Unreleased]`.

### 9. Badges e Indicadores no README
- **Problema:** Ausência de health indicators sobre o projeto no `README.md`.
- **Solução:**
  - Incorporar badges do GitHub Actions (CI Status), Go Version, Go Report Card, License, e pkg.go.dev, todos no estilo "flat-square".

### 10. Arquivo de Licença
- **Problema:** O código é open-source na teoria, mas faltam as salvaguardas explícitas de licenciamento.
- **Solução:**
  - Criar o arquivo `LICENSE` com a licença MIT no nome de "Fábio Mendes", 2026.
  - Referenciar a licença no `README.md` e no `CONTRIBUTING.md`.

## Verification Plan

- `go mod tidy` e `go build ./...`
- `go test -race -v ./...`
- Validação visual do `CHANGELOG.md` e badges do `README.md`.
- Gofmt aplicado sobre todos os arquivos Go.

Por favor, revise o plano e me informe se eu posso prosseguir com as alterações e criações dos arquivos detalhadas acima.
