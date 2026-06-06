// Package discovery implementa a detecção de atributos de hosts na rede.
//
// Oferece abstrações multiplataforma para detecção de liveness via ICMP (Ping),
// resolução reversa de DNS e obtenção de endereços MAC através da tabela ARP local.
// Implementações específicas para Windows e POSIX estão disponíveis internamente.
//
// Principais exportações:
// - Ping: Envia uma requisição ICMP Echo para verificar se um host está ativo.
// - ReverseDNS: Obtém o hostname associado a um endereço IP.
// - GetMAC: Obtém o endereço MAC correspondente a um IP.
package discovery
