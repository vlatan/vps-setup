package setup

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vlatan/vps-setup/internal/utils"
)

type SwapDevice struct {
	path        string
	sizeMB      int
	isPartition bool
}

// ChangeSwappiness changes the system swappiness,
// which means changing the free memory treshold (%) at which
// the swap is going to begin to be utilized.
func (s *Setup) ChangeSwappiness() error {

	// Check if the swappiness is valid
	if s.Swappiness < 0 || s.Swappiness > 100 {
		return fmt.Errorf("invalid swappiness: %d", s.Swappiness)
	}

	fmt.Println("Setting up the swappiness...")

	// Make parent directories
	name := "sysctl.d/99-my-swappiness.conf"
	if err := s.Etc.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}

	// Write to the file
	data := fmt.Appendf([]byte{}, "vm.swappiness = %d\n", s.Swappiness)
	if err := s.Etc.WriteFile(name, data, 0644); err != nil {
		return err
	}

	// Load our config
	confPath := filepath.Join(s.Etc.Name(), name)
	return utils.Command("sysctl", "-p", confPath).Run()
}

// CreateSwap creates a swap file and enables the swap only
// if the machine doesn't have swap partition(s) already set up.
func (s *Setup) CreateSwap() error {

	fmt.Println("Creating swap...")

	devices, err := getSwapDevices()
	if err != nil {
		return err
	}

	var hasCorrectSwap bool
	for _, dev := range devices {

		// Exit if partition-based swap exists
		if dev.isPartition {
			return nil
		}

		// Check if swap size is already correct
		if !hasCorrectSwap && dev.sizeMB == s.SwapSizeMB {
			hasCorrectSwap = true
			continue
		}

		// Disable and remove the swap file
		if err := utils.Command("swapoff", dev.path).Run(); err != nil {
			return err
		}

		if err := os.Remove(dev.path); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}

	}

	// If correct swap already exists, we're done
	if hasCorrectSwap {
		return nil
	}

	return s.enableSwap("/swapfile")
}

// enableSwap allocates swap to file and turns on the swap
func (s *Setup) enableSwap(path string) error {
	cmd := utils.Command("fallocate", "-l", fmt.Sprintf("%dM", s.SwapSizeMB), path)
	if err := cmd.Run(); err != nil {
		return err
	}

	if err := os.Chmod(path, 0600); err != nil {
		return err
	}

	if err := utils.Command("mkswap", path).Run(); err != nil {
		return err
	}

	if err := utils.Command("swapon", path).Run(); err != nil {
		return err
	}

	// Make swap permanent by adding to /etc/fstab
	f, err := s.Etc.OpenFile("fstab", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	fstabEntry := fmt.Sprintf("%s none swap sw 0 0\n", path)
	if _, err := f.WriteString(fstabEntry); err != nil {
		return err
	}

	return nil
}

// getSwapDevices gets information about all the swap devices on the machine
func getSwapDevices() ([]SwapDevice, error) {
	file, err := os.Open("/proc/swaps")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var devices []SwapDevice
	scanner := bufio.NewScanner(file)
	scanner.Scan() // Skip header

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 3 {
			sizeKB, _ := strconv.Atoi(fields[2])
			device := SwapDevice{
				path:        fields[0],
				sizeMB:      sizeKB / 1024,
				isPartition: strings.HasPrefix(fields[0], "/dev/"),
			}
			devices = append(devices, device)
		}
	}
	return devices, scanner.Err()
}
