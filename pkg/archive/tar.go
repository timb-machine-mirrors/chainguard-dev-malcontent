package archive

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chainguard-dev/clog"
	"github.com/ulikunitz/xz"
)

// extractTar extracts .apk and .tar* archives.
func ExtractTar(ctx context.Context, d string, f string) error {
	logger := clog.FromContext(ctx).With("dir", d, "file", f)
	logger.Debug("extracting tar")

	// Check if the file is valid
	_, err := os.Stat(f)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	filename := filepath.Base(f)
	tf, err := os.Open(f)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer tf.Close()
	// Set offset to the file origin regardless of type
	_, err = tf.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to seek to start: %w", err)
	}

	var tr *tar.Reader

	switch {
	case strings.Contains(f, ".apk") || strings.Contains(f, ".tar.gz") || strings.Contains(f, ".tgz"):
		gzStream, err := gzip.NewReader(tf)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzStream.Close()
		tr = tar.NewReader(gzStream)
	case strings.Contains(filename, ".tar.xz"):
		xzStream, err := xz.NewReader(tf)
		if err != nil {
			return fmt.Errorf("failed to create xz reader: %w", err)
		}
		tr = tar.NewReader(xzStream)
	case strings.Contains(filename, ".xz"):
		xzStream, err := xz.NewReader(tf)
		if err != nil {
			return fmt.Errorf("failed to create xz reader: %w", err)
		}
		uncompressed := strings.Trim(filepath.Base(f), ".xz")
		target := filepath.Join(d, uncompressed)
		if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
			return fmt.Errorf("failed to create directory for file: %w", err)
		}

		// #nosec G115 // ignore Type conversion which leads to integer overflow
		// header.Mode is int64 and FileMode is uint32
		f, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()
		if _, err = io.Copy(f, xzStream); err != nil {
			return fmt.Errorf("failed to write decompressed xz output: %w", err)
		}
		return nil
	case strings.Contains(filename, ".tar.bz2") || strings.Contains(filename, ".tbz"):
		br := bzip2.NewReader(tf)
		tr = tar.NewReader(br)
	default:
		tr = tar.NewReader(tf)
	}

	for {
		header, err := tr.Next()

		if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		clean := filepath.Clean(header.Name)
		if filepath.IsAbs(clean) || strings.Contains(clean, "../") {
			return fmt.Errorf("path is absolute or contains a relative path traversal: %s", clean)
		}

		target := filepath.Join(d, clean)
		if !IsValidPath(target, d) {
			return fmt.Errorf("invalid file path: %s", target)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			// #nosec G115 // ignore Type conversion which leads to integer overflow
			// header.Mode is int64 and FileMode is uint32
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			// #nosec G115 // ignore Type conversion which leads to integer overflow
			// header.Mode is int64 and FileMode is uint32
			out, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}

			if _, err := io.Copy(out, io.LimitReader(tr, maxBytes)); err != nil {
				out.Close()
				return fmt.Errorf("failed to copy file: %w", err)
			}

			if err := out.Close(); err != nil {
				return fmt.Errorf("failed to close file: %w", err)
			}
		case tar.TypeSymlink:
			// Skip symlinks for targets that do not exist
			_, err = os.Readlink(target)
			if os.IsNotExist(err) {
				continue
			}
			// Ensure that symlinks are not relative path traversals
			// #nosec G305 // L208 handles the check
			linkReal, err := filepath.EvalSymlinks(filepath.Join(d, header.Linkname))
			if err != nil {
				return fmt.Errorf("failed to evaluate symlink: %w", err)
			}
			if !IsValidPath(target, d) {
				return fmt.Errorf("symlink points outside temporary directory: %s", linkReal)
			}
			if err := os.Symlink(linkReal, target); err != nil {
				return fmt.Errorf("failed to create symlink: %w", err)
			}
		}
	}
	return nil
}