# Correções da Análise Técnica - Junho/2026 (catnet-core)

Este documento sumariza a execução do plano de ação para sanar as vulnerabilidades e melhorias técnicas apontadas especificamente para a biblioteca `catnet-core`.

## Resumo das Modificações no `catnet-core`

1. **Atualização da matriz CI (`C1`/`M1`)**
   - Inserimos na matriz de testes (`matrix.os`) do `.github/workflows/ci.yml` a versão mais recente do macOS (`macos-latest`).

2. **Compatibilidade POSIX no Ping (`C1`/`M1`)**
   - Corrigimos a temporização do comando `ping -W` no macOS. A verificação do OS é feita em tempo de execução via `runtime.GOOS == "darwin"`.
   - Como o macOS espera o valor do timeout nativamente em *milissegundos*, os limites não são mais convertidos a *segundos* como no Linux. Acompanhado da criação de um teste `os_posix_test.go`.

3. **Sanitização Reforçada contra CSV Injection (`C3`/`M3`)**
   - Em `pkg/exporter/exporter.go`, atualizamos o sanitizador CSV (`sanitizeCSVField`).
   - Agora ele bloqueia explicitamente a inserção de quebras de linha prefixadas `\n`, protegendo contra tentativas maliciosas de quebrar uma célula e invadir linhas seguintes. Adicionado `strings.ContainsAny(field, "\n\r")`.

4. **Tratamento de Edge Cases via Fuzzing (`C5`/`M5`)**
   - Introduzimos o método nativo Fuzz (`go test -fuzz`) em `pkg/targets/parse_test.go`. O `FuzzParseRange` tenta intencionalmente causar um *Panic* na engine com endereços e CIDRs corrompidos, atestando a robustez da validação de string.

5. **Cancelamento Estrito por Context (`C2`/`M6`)**
   - A biblioteca `context.Context` agora é propagada diretamente até as camadas de bloqueio I/O de rede (`discovery.Ping` e `ports.ScanPorts`).
   - O tempo de reposta de concorrência frente ao acionamento de um cancelamento (CTRL+C) foi reduzido, eliminando o bloqueio interno residual até o fim do TCP handshake.
