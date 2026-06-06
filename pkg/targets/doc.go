// Package targets fornece funcionalidades para parsing de alvos de rede.
//
// O pacote é responsável por interpretar formatos como CIDR, intervalos de
// IPs (dash ranges) e endereços individuais, expandindo-os em listas de IPs
// válidas prontas para varredura pelo engine.
//
// Principais exportações:
// - ParseRange: Analisa e converte uma string de alvo em uma lista de IPs.
package targets
