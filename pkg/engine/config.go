package engine

import "context"

// ScanConfig define os parâmetros de uma varredura.
type ScanConfig struct {
	// DefaultPorts é a lista de portas TCP escaneadas em cada host vivo.
	// Padrão: [22 (SSH), 80 (HTTP), 443 (HTTPS), 139 (NetBIOS), 445 (SMB), 3389 (RDP)].
	DefaultPorts []int `json:"defaultPorts"`

	// PortTimeoutMs é o timeout em milissegundos para cada tentativa de conexão TCP.
	// Padrão: 500ms. Limite: 1–10000ms (Sanitize clampeia valores fora do range).
	PortTimeoutMs int `json:"portTimeoutMs"`

	// PingTimeoutMs é o timeout em milissegundos para o ping ICMP de liveness.
	// Padrão: 1000ms. Limite: 1–10000ms.
	PingTimeoutMs int `json:"pingTimeoutMs"`

	// MaxThreads define o nível de paralelismo da varredura (goroutines simultâneas).
	// Padrão: 64. Limite máximo absoluto: 256 (imposto pelo engine em StartScan).
	// Valores acima de 256 ou abaixo de 1 são silenciosamente clampeados para 16.
	MaxThreads int `json:"maxThreads"`

	// FingerprintProvider permite injetar lógica customizada de fingerprinting.
	// Se nil, o motor usa pkg/fingerprint com heurísticas padrão de TTL, banner e OUI.
	// Útil para testes (mock) ou extensão de capacidades de detecção.
	FingerprintProvider FingerprintProvider `json:"-"`
}

// FingerprintData contém os resultados da detecção.
type FingerprintData struct {
	OS         string
	OSFamily   string
	DeviceType string
	Vendor     string
}

// FingerprintProvider define o contrato para heurísticas de detecção de SO e dispositivos.
type FingerprintProvider interface {
	Fingerprint(ctx context.Context, ip, mac string, ttl int, ports []int, timeoutMs int) FingerprintData
}

// DefaultConfig retorna uma ScanConfig com valores padrão conservadores.
func DefaultConfig() ScanConfig {
	return ScanConfig{
		DefaultPorts:        []int{22, 80, 443, 139, 445, 3389},
		PortTimeoutMs:       500,
		PingTimeoutMs:       1000,
		MaxThreads:          64,
		FingerprintProvider: nil, // Usará default em StartScan se nulo
	}
}

// Sanitize corrige valores fora de limites seguros.
// O próprio motor executa essa sanitização defensivamente no StartScan,
// mas pode ser invocada manualmente para refletir os limites na interface do cliente.
func (c *ScanConfig) Sanitize() {
	if c.MaxThreads <= 0 || c.MaxThreads > 256 {
		c.MaxThreads = 16
	}
	if c.PortTimeoutMs <= 0 || c.PortTimeoutMs > 10000 {
		c.PortTimeoutMs = 500
	}
	if c.PingTimeoutMs <= 0 || c.PingTimeoutMs > 10000 {
		c.PingTimeoutMs = 1000
	}
}
