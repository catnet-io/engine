package diff

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mendsec/catnet-core/pkg/results"
)

type HostStatus string

const (
	StatusNew       HostStatus = "NEW"
	StatusLost      HostStatus = "LOST"
	StatusChanged   HostStatus = "CHANGED"
	StatusUnchanged HostStatus = "UNCHANGED"
)

type HostDiff struct {
	IP       string     `json:"ip"`
	Hostname string     `json:"hostname"`
	Status   HostStatus `json:"status"`
	Details  string     `json:"details"`
}

// Compare analyzes two scan reports and returns a list of differences.
func Compare(oldReport, newReport *results.ScanReport) []HostDiff {
	var diffs []HostDiff

	oldMap := make(map[string]results.DeviceInfo)
	if oldReport != nil {
		for _, d := range oldReport.Devices {
			oldMap[d.IP] = d
		}
	}

	newMap := make(map[string]results.DeviceInfo)
	if newReport != nil {
		for _, d := range newReport.Devices {
			newMap[d.IP] = d
		}
	}

	// 1. Check for NEW and CHANGED/UNCHANGED devices
	for ip, newDev := range newMap {
		oldDev, exists := oldMap[ip]
		if !exists {
			diffs = append(diffs, HostDiff{
				IP:       ip,
				Hostname: newDev.Hostname,
				Status:   StatusNew,
				Details:  "Host came online",
			})
			continue
		}

		if !oldDev.IsAlive && newDev.IsAlive {
			diffs = append(diffs, HostDiff{
				IP:       ip,
				Hostname: newDev.Hostname,
				Status:   StatusNew,
				Details:  "Host came online (was dead)",
			})
			continue
		}

		if oldDev.IsAlive && !newDev.IsAlive {
			// Handled in the LOST loop
			continue
		}

		// Both existed and were alive (or both dead, but usually scans only store alive/attempted)
		changes := comparePorts(oldDev.OpenPorts, newDev.OpenPorts)
		if len(changes) > 0 {
			diffs = append(diffs, HostDiff{
				IP:       ip,
				Hostname: newDev.Hostname,
				Status:   StatusChanged,
				Details:  strings.Join(changes, "; "),
			})
		} else {
			diffs = append(diffs, HostDiff{
				IP:       ip,
				Hostname: newDev.Hostname,
				Status:   StatusUnchanged,
				Details:  "No changes",
			})
		}
	}

	// 2. Check for LOST devices
	for ip, oldDev := range oldMap {
		newDev, exists := newMap[ip]
		if !exists || (oldDev.IsAlive && !newDev.IsAlive) {
			diffs = append(diffs, HostDiff{
				IP:       ip,
				Hostname: oldDev.Hostname,
				Status:   StatusLost,
				Details:  "Host went offline",
			})
		}
	}

	// Sort by IP for consistent output
	sort.Slice(diffs, func(i, j int) bool {
		return diffs[i].IP < diffs[j].IP
	})

	return diffs
}

func comparePorts(oldPorts, newPorts []int) []string {
	oldSet := make(map[int]bool)
	for _, p := range oldPorts {
		oldSet[p] = true
	}

	newSet := make(map[int]bool)
	for _, p := range newPorts {
		newSet[p] = true
	}

	var changes []string
	var opened []string
	for p := range newSet {
		if !oldSet[p] {
			opened = append(opened, fmt.Sprintf("%d", p))
		}
	}
	if len(opened) > 0 {
		changes = append(changes, "Opened ports: "+strings.Join(opened, ", "))
	}

	var closed []string
	for p := range oldSet {
		if !newSet[p] {
			closed = append(closed, fmt.Sprintf("%d", p))
		}
	}
	if len(closed) > 0 {
		changes = append(changes, "Closed ports: "+strings.Join(closed, ", "))
	}

	return changes
}
