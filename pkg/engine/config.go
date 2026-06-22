package engine

import "context"

// ScanConfig define os parâmetros de uma varredura.
type ScanConfig struct {
	DefaultPorts  []int `json:"defaultPorts"`
	PortTimeoutMs int   `json:"portTimeoutMs"`
	PingTimeoutMs int   `json:"pingTimeoutMs"`

	// MaxThreads define o nível de paralelismo da varredura.
	// O motor impõe um limite máximo rigoroso de 256 threads para prevenir exaustão
	// de sockets no host (ulimit issues) e um mínimo de 1.
	MaxThreads int `json:"maxThreads"`

	// FingerprintProvider permite injetar lógica customizada de fingerprinting.
	// Se nulo, o motor usará o pacote padrão pkg/fingerprint.
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
		DefaultPorts:  []int{22, 80, 443, 139, 445, 3389},
		PortTimeoutMs: 500,
		PingTimeoutMs: 1000,
		MaxThreads:    64,
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
