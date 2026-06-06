// Package coreerr define os erros sentinelas da biblioteca.
//
// Centraliza todas as definições de erros conhecidos (timeout, input inválido,
// cancelamento, erros de exportação) para que consumidores da API
// possam realizar comparações via errors.Is() com segurança.
//
// Principais exportações:
// - ErrInvalidInput: Erros de formato em IP ou configurações.
// - ErrTimeout: Erros indicando que uma operação expirou o tempo.
// - ErrCancelled: Indica que o contexto da varredura foi cancelado.
// - ErrExport: Erros durante a serialização de dados.
package coreerr
