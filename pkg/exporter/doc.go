// Package exporter lida com a formatação e serialização de relatórios.
//
// Suporta a conversão de ScanReport em formatos padronizados (JSON, XML, CSV)
// garantindo que as exportações evitem problemas de segurança como CSV Injection.
//
// Principais exportações:
// - ExportJSON: Exporta relatórios como JSON indentado.
// - ExportXML: Exporta relatórios como XML válido.
// - ExportCSV: Exporta relatórios para CSV, tratando injeção de fórmulas.
package exporter
