// Package engine orquestra a execução concorrente de varreduras na rede.
//
// O pacote é o ponto de entrada principal do catnet-core, fornecendo a
// função StartScan para gerenciar pools de goroutines, timeouts de
// execução e a emissão de eventos assíncronos durante a varredura.
//
// Principais exportações:
// - StartScan: Inicia uma varredura de rede.
// - ScanConfig: Configurações como limites de concorrência e timeouts.
// - EventCallback: Tipo para recebimento de eventos de progresso e resultados.
package engine
