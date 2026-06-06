// Package ports implementa varredura concorrente de portas TCP.
//
// Utiliza tentativas de conexão com controle de simultaneidade interno
// via semáforos para evitar exaustão de descritores de arquivo,
// retornando de forma determinística portas que aceitam conexões ativas.
//
// Principais exportações:
// - ScanPorts: Varre uma lista de portas de um host de modo concorrente.
package ports
