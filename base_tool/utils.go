package basetool

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	gcp_cloud_storage "github.com/blueturbo-ad/go-utils/gcp_cloud_tool/gcp_cloud_storage"
	"github.com/blueturbo-ad/go-utils/zap_loggerex"
)

func ReadGCPCloudStorageFile(filePath string) ([]byte, error) {
	zap_loggerex.GetSingleton().Debug("bid_stdout_logger", "read file from GCP cloud storage %v", filePath)

	client := gcp_cloud_storage.GetSingleton().GetClient("dsp_bucket")
	if client == nil {
		return nil, fmt.Errorf("cloud storage client is nil")
	}

	reader, err := client.Object(filePath).NewReader(context.Background())
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func StringToInt64(str string) (int64, error) {
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}
func StringToUInt32(str string) (uint32, error) {
	value, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(value), nil
}

func StringToUInt64(str string) (uint64, error) {
	value, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// " ", "\n", "\t", "\r"
func RemoveWhitespace(s string) string {
	specialChars := []string{" ", "\n", "\t", "\r"}
	return RemoveCharacters(s, specialChars)
}

func RemoveCharacters(s string, chars []string) string {
	for _, c := range chars {
		s = strings.ReplaceAll(s, c, "")
	}

	return s
}
