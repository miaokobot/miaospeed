package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"path"
	"regexp"
)

func FindAndExtract(gzipReader io.Reader, filters ...regexp.Regexp) (map[string][]byte, error) {
	uncompressedStream, err := gzip.NewReader(gzipReader)
	if err != nil {
		return nil, err
	}
	defer uncompressedStream.Close()

	tarReader := tar.NewReader(uncompressedStream)
	result := map[string][]byte{}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			return result, nil
		} else if err != nil {
			return nil, err
		}

		switch header.Typeflag {
		case tar.TypeReg:
			matched := false
			for _, r := range filters {
				if r.MatchString(header.Name) {
					matched = true
					break
				}
			}

			if matched {
				buf := &bytes.Buffer{}
				if _, err := io.Copy(buf, tarReader); err == nil {
					result[path.Base(header.Name)] = buf.Bytes()
				}
			}
		}
	}
}
