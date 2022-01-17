package utils

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func ParseOIDsToID(oids []string) []int {
	var result []int

	for _, oid := range oids {
		i := strings.LastIndex(oid, "-")
		if i > 0 {
			id, err := strconv.Atoi(oid[i+1:])
			if err != nil {
				log.Print("oid not valid %v", err)
			}
			result = append(result, id)
		}
	}
	return result
}

func ParseIDToOID(ids []int, entityType string) []string {
	var result []string

	for _, id := range ids {
		result = append(result, fmt.Sprintf("%s-%v", entityType, id))
	}
	return result
}

func GetObjectTypeFromOID(oid string) string {
	i := strings.LastIndex(oid, "-")
	return oid[0:i]
}

func GenerateOID(oType string, id interface{}) string {
	return fmt.Sprintf("%s-%v", oType, id)
}

func GetIDFromOID(oid string) int {
	result := strings.Split(oid, "-")
	if len(result) < 2 {
		return 0
	}
	id, err := strconv.Atoi(result[1])
	if err != nil {
		return 0
	}
	return id
}
