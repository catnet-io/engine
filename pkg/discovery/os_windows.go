//go:build windows

package discovery

import (
	"fmt"
	"net"
	"syscall"
	"unsafe"
)

// osPing faz ping no Windows
func osPing(ip string, timeoutMs int) bool {
	iphlpapi := syscall.NewLazyDLL("iphlpapi.dll")
	icmpCreateFile := iphlpapi.NewProc("IcmpCreateFile")
	icmpSendEcho := iphlpapi.NewProc("IcmpSendEcho")
	icmpCloseHandle := iphlpapi.NewProc("IcmpCloseHandle")

	handle, _, _ := icmpCreateFile.Call()
	if handle == ^uintptr(0) {
		return false
	}
	defer icmpCloseHandle.Call(handle)

	destIP := net.ParseIP(ip).To4()
	if destIP == nil {
		return false
	}

	var destIPUint32 uint32 = uint32(destIP[0]) | uint32(destIP[1])<<8 | uint32(destIP[2])<<16 | uint32(destIP[3])<<24

	var replyData [128]byte
	replySize := uint32(len(replyData))

	ret, _, _ := icmpSendEcho.Call(
		handle,
		uintptr(destIPUint32),
		0,
		0,
		0,
		uintptr(unsafe.Pointer(&replyData[0])),
		uintptr(replySize),
		uintptr(timeoutMs),
	)

	if ret == 0 {
		return false
	}

	status := *(*uint32)(unsafe.Pointer(&replyData[4]))
	return status == 0
}

// osGetMAC obtém o MAC usando SendARP no Windows
func osGetMAC(ip string) string {
	iphlpapi := syscall.NewLazyDLL("iphlpapi.dll")
	sendARP := iphlpapi.NewProc("SendARP")
	destIP := net.ParseIP(ip).To4()
	if destIP == nil {
		return ""
	}
	var destIPUint32 uint32
	destIPUint32 = uint32(destIP[0]) | uint32(destIP[1])<<8 | uint32(destIP[2])<<16 | uint32(destIP[3])<<24
	var mac [6]byte
	macLen := uint32(len(mac))
	// Segurança: mac é um array de tamanho fixo [6]byte alocado na stack.
	// macLen é inicializado com len(mac) == 6 antes da chamada.
	// O acesso via unsafe.Pointer é seguro porque o array não escapa do
	// escopo e seu tamanho é conhecido em tempo de compilação.
	// A validação `macLen == 6` após o retorno garante dados não corrompidos.
	ret, _, _ := sendARP.Call(
		uintptr(destIPUint32),
		0,
		uintptr(unsafe.Pointer(&mac[0])),
		uintptr(unsafe.Pointer(&macLen)),
	)
	if ret == 0 && macLen == 6 {
		return fmt.Sprintf("%02X-%02X-%02X-%02X-%02X-%02X",
			mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
	}
	return ""
}
