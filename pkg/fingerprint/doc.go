// Package fingerprint provides heuristic and banner-based operating system
// and device type detection mechanisms.
//
// Principais exportações:
// - Fingerprint: Orquestra detecção de SO, tipo de dispositivo e vendor.
// - GrabBanners: Coleta banners de portas abertas via conexão TCP.
// - GuessOSFromTTL: Detecta família de SO a partir do valor TTL.
// - VendorFromMAC: Identifica fabricante a partir do prefixo OUI do MAC.
// - OsFromBanners: Infere SO e tipo de dispositivo a partir de banners coletados.
package fingerprint
