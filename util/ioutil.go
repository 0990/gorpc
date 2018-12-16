package util

import "io"

func WriteFull(writer io.Writer, buf []byte) error {
	total := len(buf)

	for pos := 0; pos < total; {
		n, err := writer.Write(buf[pos:])
		if err != nil {
			return err
		}
		pos += n
	}
	return nil
}
