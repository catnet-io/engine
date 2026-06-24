package profile

// ScanProfile define os parâmetros de execução de uma varredura.
// É o tipo de configuração passado pelo caller para Engine.ScanStream.
type ScanProfile struct {
	DefaultPorts []int `json:"default_ports"`
	Concurrency  int   `json:"concurrency"`
	TimeoutMs    int   `json:"timeout_ms"`
}

// DefaultProfile retorna um ScanProfile com valores conservadores e seguros.
func DefaultProfile() ScanProfile {
	return ScanProfile{
		DefaultPorts: []int{22, 80, 443, 139, 445, 3389},
		Concurrency:  64,
		TimeoutMs:    1000,
	}
}

// Sanitize corrige valores fora dos limites aceitáveis.
// Deve ser chamado pelo caller antes de passar o profile para ScanStream.
func (p *ScanProfile) Sanitize() {
	if p.Concurrency <= 0 || p.Concurrency > 256 {
		p.Concurrency = 16
	}
	if p.TimeoutMs <= 0 || p.TimeoutMs > 10000 {
		p.TimeoutMs = 1000
	}
}
