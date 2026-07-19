package evidence

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

type BundleComponent struct {
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Path        string    `json:"path"`
	Checksum    string    `json:"checksum"`
	GeneratedAt time.Time `json:"generated_at"`
}

type BundleManifest struct {
	Version         string            `json:"version"`
	GeneratedAt     time.Time         `json:"generated_at"`
	GeneratedBy     string            `json:"generated_by"`
	Environment     string            `json:"environment"`
	CommitSHA       string            `json:"commit_sha"`
	Components      []BundleComponent `json:"components"`
	PolicyCount     int               `json:"policy_count"`
	RedTeamPassRate float64           `json:"redteam_pass_rate"`
	SLOCompliant    map[string]bool   `json:"slo_compliant"`
	Checksum        string            `json:"checksum"`
}

type BundleConfig struct {
	OutputDir          string
	AEGISVersion       string
	Environment        string
	Actor              string
	PolicyStorePath    string
	RedTeamResultsPath string
	ADRsPath           string
}

type Bundler struct {
	config BundleConfig
	logger *zap.Logger
}

func NewBundler(config BundleConfig, logger *zap.Logger) *Bundler {
	return &Bundler{
		config: config,
		logger: logger,
	}
}

func (b *Bundler) Generate(ctx context.Context) (*BundleManifest, string, error) {
	manifest := &BundleManifest{
		Version:         b.config.AEGISVersion,
		GeneratedAt:     time.Now(),
		GeneratedBy:     b.config.Actor,
		Environment:     b.config.Environment,
		CommitSHA:       os.Getenv("GIT_COMMIT"),
		Components:      []BundleComponent{},
		SLOCompliant:    map[string]bool{"gateway_latency": true, "scanner_latency": true, "availability": true},
		RedTeamPassRate: 100.0,
	}

	tempFile, err := os.CreateTemp("", "aegis-bundle-*.zip")
	if err != nil {
		return nil, "", fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	zw := zip.NewWriter(tempFile)

	// Here we would walk b.config.PolicyStorePath, b.config.ADRsPath, etc.
	// For now, we mock it.

	// Write mock manifest to calculate checksum later
	manifestBytes, _ := json.Marshal(manifest)
	manF, _ := zw.Create("manifest.json")
	manF.Write(manifestBytes)

	if err := zw.Close(); err != nil {
		return nil, "", fmt.Errorf("closing zip: %w", err)
	}
	tempFile.Close()

	hash := sha256.New()
	f, _ := os.Open(tempFile.Name())
	io.Copy(hash, f)
	f.Close()
	manifest.Checksum = hex.EncodeToString(hash.Sum(nil))

	outName := fmt.Sprintf("aegis-evidence-%d-%s.zip", time.Now().Unix(), "mock")
	outPath := filepath.Join(b.config.OutputDir, outName)

	if err := os.MkdirAll(b.config.OutputDir, 0755); err != nil {
		return nil, "", fmt.Errorf("creating output dir: %w", err)
	}
	
	// Copy tempfile to output
	in, err := os.Open(tempFile.Name())
	if err != nil {
		return nil, "", fmt.Errorf("opening temp file: %w", err)
	}
	defer in.Close()
	out, err := os.Create(outPath)
	if err != nil {
		return nil, "", fmt.Errorf("creating output file: %w", err)
	}
	defer out.Close()
	io.Copy(out, in)

	return manifest, outPath, nil
}

func (b *Bundler) addFileToZip(zw *zip.Writer, srcPath, destName string) (string, error) {
	fileToZip, err := os.Open(srcPath)
	if err != nil {
		return "", fmt.Errorf("opening file: %w", err)
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return "", fmt.Errorf("stat file: %w", err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return "", fmt.Errorf("creating zip header: %w", err)
	}
	header.Name = destName
	header.Method = zip.Deflate

	writer, err := zw.CreateHeader(header)
	if err != nil {
		return "", fmt.Errorf("creating zip entry: %w", err)
	}

	hash := sha256.New()
	multiWriter := io.MultiWriter(writer, hash)
	if _, err := io.Copy(multiWriter, fileToZip); err != nil {
		return "", fmt.Errorf("copying file: %w", err)
	}

	b.logger.Info("Added file to bundle", zap.String("file", destName), zap.Int64("size", info.Size()))
	return hex.EncodeToString(hash.Sum(nil)), nil
}
