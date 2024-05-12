package grpchelper

import "strings"

func GetService(fullMethodName string) string {
	fullMethodName = strings.TrimPrefix(fullMethodName, "/")
	parts := strings.Split(fullMethodName, "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}

func GetAddress(fullMethodName string) string {
	return "miniresolver:" + GetService(fullMethodName)
}
