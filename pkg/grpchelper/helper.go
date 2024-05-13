package grpchelper

import "strings"

// GetService returns <package>.<service> from fullMethodName
func GetService(fullMethodName string) string {
	fullMethodName = strings.TrimPrefix(fullMethodName, "/")
	parts := strings.Split(fullMethodName, "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}

// GetAddress returns miniresolver:<package>.<service>
func GetAddress(fullMethodName string) string {
	return "miniresolver:" + GetService(fullMethodName)
}
