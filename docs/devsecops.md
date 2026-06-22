# DevSecOps em catnet-core

Este documento descreve as práticas de desenvolvimento seguro e qualidade contínua adotadas pelo `catnet-core`.

## Objetivo

Garantir que a evolução do código ocorra com segurança, qualidade e rapidez, integrando verificações automatizadas desde o desenvolvimento até a entrega.

## Princípios

- segurança como parte do fluxo de desenvolvimento, não como etapa separada
- automação de validações de código, dependência e vulnerabilidades
- feedback rápido em pull requests
- responsabilidade compartilhada entre desenvolvimento, operações e segurança
- manter o código leve e sem dependências desnecessárias

## Automação existente

### GitHub Actions

O repositório já tem os seguintes workflows ativos:

- `ci.yml`
  - roda em `push` e `pull_request` para `main` e `develop`
  - `go mod verify`
  - `go vet ./...`
  - `go test -race -v ./...`
  - `go test -fuzz=FuzzParseRange -fuzztime=10s ./pkg/targets`

- `golangci-lint.yml`
  - valida o estilo e possíveis problemas estáticos com `golangci-lint`

- `govulncheck.yml`
  - roda em `push`, `pull_request` e `workflow_dispatch`
  - agenda semanal para verificar vulnerabilidades de dependências

### Dependabot

- atualiza GitHub Actions e dependências Go semanalmente via `.github/dependabot.yml`

## Práticas recomendadas

- abrir PRs pequenos e focados, com descrição clara
- revisar mudanças de segurança, performance e dependências com atenção
- usar `go test -race` e `go vet` como baseline local antes de abrir PR
- não adicionar dependências externas sem justificativa robusta
- proteger secrets via GitHub Actions e evitar hardcoded credentials no código

## Segurança de código

- manter a sanitização de exportação de dados em `pkg/exporter`
- garantir que parsing de alvos não panic com entrada malformada
- evitar exposição de dados sensíveis em logs e relatórios

## Vulnerabilidades e dependências

- manter `go.mod` e `go.sum` atualizados e verificados
- revisar alertas do `govulncheck` e do Dependabot prontamente
- priorizar correções para dependências com severidade alta ou críticas

## Observabilidade e qualidade

- usar os testes existentes para confirmar comportamento em Windows e POSIX
- manter a documentação atualizada em `docs/` e `README.md`
- registrar quaisquer ajustes de segurança em `docs/remediation_202606.md` ou novo documento de remediação

## Próximos passos sugeridos

- adicionar checks de cobertura de teste se útil para o fluxo
- avaliar ferramentas de análise de código adicionais para segurança estática (por exemplo, `staticcheck`)
- padronizar pull requests com template e checklist de segurança
- documentar procedimentos de resposta a vulnerabilidade no `SECURITY.md`
