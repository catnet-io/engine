# Fase 3: Desktop GUI Refactoring (React + Wails)

## Visão Geral
Com o pacote do `catnet-core` completamente desacoplado, os *bindings* em Go do Wails (`app.go`) sofreram mudanças nas estruturas de dados. O front-end em React precisa ser ajustado para consumir os novos tipos (como `results.HostResult`) e se comunicar corretamente com a engine através dos novos canais de eventos.

## Proposed Changes

### 1. Regeneração dos Bindings
- **[COMANDO]**: Rodaremos `wails generate module` na pasta do `catnet-scanner` para forçar a criação dos novos modelos TypeScript (que substituirão `scanner.DeviceInfo` por `results.HostResult` e `profile.ScanProfile`).

### 2. Atualização do `App.tsx`
- **Importações:** Atualizar os `imports` para buscar de `../wailsjs/go/models.ts` sob o novo namespace `results` e `profile`.
- **Estado (State):** O array de dispositivos (`devices`) passará a ser do tipo `results.HostResult[]` ao invés de `DeviceInfo[]`.
- **Mapeamento de Dados (JSX):**
  - Trocar `dev.isAlive` por `dev.alive`.
  - Trocar `dev.openPortsCount` por uma contagem real `dev.open_ports?.length || 0`.
  - Trocar `dev.openPorts` por `dev.open_ports`.
- **Configuração de Scan:** Atualizar a função `handleScan` para instanciar o novo objeto de configuração esperado (`profile.ScanProfile`) com `concurrency` em vez de `maxThreads`, e `timeout_ms` em vez de `portTimeoutMs`.

### 3. Melhorias UX Residuais
- Assegurar que os botões "Auto Detect", "Start", "Stop" e "Export" operem perfeitamente sobre a nova tipagem, preservando os *toasts* e a *Debug Log* visual que já estava pronta.
- Manter o CSS atual (`cyber-table`, `glass-panel`) que já dá a identidade *hacker/synthwave* da aplicação, garantindo apenas que a tabela reflita fielmente as colunas de `Vendor` e `RTTMs` caso você queira exibi-los no futuro.

## Verification Plan
1. Após alterar o TSX, vamos compilar e rodar `wails build` (ou simular a transpilação localmente rodando um `npm run build` na pasta frontend) para confirmar que nenhum erro de TypeScript sobreviveu.
2. Inspecionar o *binding* para ter certeza de que o canal de exportação `ExportResults(devices)` funciona com o novo array Go de `HostResult`.
