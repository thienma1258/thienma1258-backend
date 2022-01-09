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

func GetObjectTypeFromOID(oid string) string {
	i := strings.LastIndex(oid, "-")
	return oid[0 : i]
}

func GenerateOID(oType string, id interface{}) string {
	return fmt.Sprintf("%s-%v", oType, id)
}