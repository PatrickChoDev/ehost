package utils

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

type IPOrHost int

const (
	IP IPOrHost = iota
	HOST
	INVALID
)

// IsIPorHostEntry checks if the input is an IP address or a host entry (hostname or FQDN)
func IsIPorHostEntry(input string) IPOrHost {
	// Check if input is a valid IP address
	if net.ParseIP(input) != nil {
		return IP
	}

	// Regular expression to check if input is a valid host entry (hostname or FQDN)
	hostEntryRegex := `^([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*\.)*[a-zA-Z0-9]+(-[a-zA-Z0-9]+)*$`
	matched, _ := regexp.MatchString(hostEntryRegex, input)
	if matched {
		return HOST
	}

	return INVALID
}

// Contains checks if a slice contains a specific string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// processHostsFile is a helper function to read, find the header, and write the hosts file
func processHostsFile(hostsFilePath string, processFunc func([]string, int) ([]string, error)) error {
	// Ensure the header exists
	if err := EnsureHeader(hostsFilePath); err != nil {
		return err
	}

	// Read the hosts file
	file, err := os.ReadFile(hostsFilePath)
	if err != nil {
		return fmt.Errorf("failed to read hosts file: %w", err)
	}

	// Split the file into lines
	lines := strings.Split(string(file), "\n")

	// Find the header line
	headerIndex := -1
	for i, line := range lines {
		if strings.Contains(line, "####### Managed by EHost #######") {
			headerIndex = i
			break
		}
	}

	if headerIndex == -1 {
		return fmt.Errorf("header line not found in hosts file")
	}

	// Process the lines
	newLines, err := processFunc(lines, headerIndex)
	if err != nil {
		return err
	}

	// Write the updated hosts file
	err = os.WriteFile(hostsFilePath, []byte(strings.Join(newLines, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to hosts file: %w", err)
	}

	return nil
}

// AddEntry adds an entry to the hosts file, skip if the entry pair already existed
func AddEntry(hostsFilePath, ip string, fqdn string) error {
	return processHostsFile(hostsFilePath, func(lines []string, headerIndex int) ([]string, error) {
		// Check if the entry already exists
		for _, line := range lines[headerIndex+1:] {
			if strings.Contains(line, ip) && strings.Contains(line, fqdn) {
				return lines, nil // Entry already exists, skip adding
			}
		}

		// Append the new entry after the header line
		newEntry := fmt.Sprintf("%s %s", ip, fqdn)
		return append(lines[:headerIndex+1], append([]string{newEntry}, lines[headerIndex+1:]...)...), nil
	})
}

// RemoveEntry removes an entry from the hosts file, skip if the entry pair does not exist
func RemoveEntry(hostsFilePath, ip string, fqdn string) error {
	return processHostsFile(hostsFilePath, func(lines []string, headerIndex int) ([]string, error) {
		var newLines []string
		entryFound := false
		for _, line := range lines[headerIndex+1:] {
			if strings.Contains(line, ip) && strings.Contains(line, fqdn) {
				entryFound = true
				continue
			}
			newLines = append(newLines, line)
		}

		if !entryFound {
			return lines, nil // Entry does not exist, skip removing
		}

		return append(lines[:headerIndex+1], newLines...), nil
	})
}

// RemoveAllHostname removes all entries with the specified FQDN from the hosts file
func RemoveAllHostname(hostsFilePath, fqdn string) error {
	return processHostsFile(hostsFilePath, func(lines []string, headerIndex int) ([]string, error) {
		var newLines []string
		entryFound := false
		for _, line := range lines[headerIndex+1:] {
			if strings.Contains(line, fqdn) {
				entryFound = true
				continue
			}
			newLines = append(newLines, line)
		}

		if !entryFound {
			return lines, nil // Entry does not exist, skip removing
		}

		return append(lines[:headerIndex+1], newLines...), nil
	})
}

// RemoveAllIP removes all entries with the specified IP from the hosts file
func RemoveAllIP(hostsFilePath, ip string) error {
	return processHostsFile(hostsFilePath, func(lines []string, headerIndex int) ([]string, error) {
		var newLines []string
		entryFound := false
		for _, line := range lines[headerIndex+1:] {
			if strings.Contains(line, ip) {
				entryFound = true
				continue
			}
			newLines = append(newLines, line)
		}

		if !entryFound {
			return lines, nil // Entry does not exist, skip removing
		}

		return append(lines[:headerIndex+1], newLines...), nil
	})
}

// GetEntriesByIP returns all entries in the hosts file that contains the specified IP address
func GetEntriesByIP(hostsFilePath, ip string) ([]string, error) {
	// Read the hosts file
	file, err := os.ReadFile(hostsFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read hosts file: %w", err)
	}

	// Return all entries with the specified IP address
	lines := strings.Split(string(file), "\n")
	var entries []string
	for _, line := range lines {
		if strings.Contains(line, ip) {
			entries = append(entries, line)
		}
	}

	return entries, nil
}

// GetEntriesByHostname returns all entries in the hosts file that contains the specified FQDN
func GetEntriesByHostname(hostsFilePath, fqdn string) ([]string, error) {
	// Read the hosts file
	file, err := os.ReadFile(hostsFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read hosts file: %w", err)
	}

	// Return all entries with the specified FQDN
	lines := strings.Split(string(file), "\n")
	var entries []string
	for _, line := range lines {
		if strings.Contains(line, fqdn) {
			entries = append(entries, line)
		}
	}

	return entries, nil
}

// EnsureHeader ensures the header line exists in the hosts file
func EnsureHeader(hostsFilePath string) error {
	header := "####### Managed by EHost #######"

	// Read the hosts file
	file, err := os.ReadFile(hostsFilePath)
	if err != nil {
		return fmt.Errorf("failed to read hosts file: %w", err)
	}

	// Check if the header already exists
	if strings.Contains(string(file), header) {
		return nil // Header already exists, skip adding
	}

	// Append the header to the file
	newFile := fmt.Sprintf("%s\n\n\n\n\n%s\n", string(file), header)

	// Write the updated hosts file
	err = os.WriteFile(hostsFilePath, []byte(newFile), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to hosts file: %w", err)
	}

	return nil
}
