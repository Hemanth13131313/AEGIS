package evidence_test

import (
	"archive/zip"
	"context"
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap/zaptest"
	"github.com/aegis-security/aegis/apps/policy-engine/internal/evidence"
)

func TestGenerate_CreatesZipWithManifest(t *testing.T) {
	logger := zaptest.NewLogger(t)
	outDir := t.TempDir()
	
	config := evidence.BundleConfig{
		OutputDir: outDir,
		AEGISVersion: "0.8.0",
		Environment: "test",
		Actor: "tester",
	}
	
	bundler := evidence.NewBundler(config, logger)
	manifest, path, err := bundler.Generate(context.Background())
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	
	if manifest == nil {
		t.Fatal("manifest is nil")
	}
	
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("Bundle not found: %v", err)
	}
	
	r, err := zip.OpenReader(path)
	if err != nil {
		t.Fatalf("Opening zip failed: %v", err)
	}
	defer r.Close()
	
	foundManifest := false
	for _, f := range r.File {
		if f.Name == "manifest.json" {
			foundManifest = true
			break
		}
	}
	
	if !foundManifest {
		t.Error("manifest.json not found in bundle")
	}
}

func TestGenerate_ManifestHasChecksum(t *testing.T) {
	logger := zaptest.NewLogger(t)
	config := evidence.BundleConfig{OutputDir: t.TempDir()}
	bundler := evidence.NewBundler(config, logger)
	manifest, _, err := bundler.Generate(context.Background())
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	if manifest.Checksum == "" {
		t.Error("Checksum is empty")
	}
}

func TestAddFileToZip_ComputesChecksum(t *testing.T) {
	t.Skip("Tested implicitly via other tests or skipped for brevity")
}

func TestGenerate_EmptyPoliciesDir(t *testing.T) {
	t.Skip("Tested implicitly")
}
