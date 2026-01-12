package generator

import (
	"context"
	"testing"

	"lesmotsdatche/internal/generator/fill"
	"lesmotsdatche/internal/generator/languagepack"
	"lesmotsdatche/internal/generator/llm"
)

func TestOrchestrator_CreateDefaultTemplate(t *testing.T) {
	langPack := languagepack.NewFrenchPack()
	config := DefaultConfig()
	config.GridSize = [2]int{7, 7}

	mock := llm.NewMockClient()
	validatingClient := llm.NewValidatingClient(mock, llm.DefaultConfig())

	orch := NewOrchestrator(validatingClient, langPack, nil, config)

	template := orch.createDefaultTemplate()

	if len(template) != 7 {
		t.Errorf("expected 7 rows, got %d", len(template))
	}
	if len(template[0]) != 7 {
		t.Errorf("expected 7 cols, got %d", len(template[0]))
	}

	// Check center is not a block
	if template[3][3].IsBlock() {
		t.Error("center cell should not be a block")
	}

	// Count blocks for reasonable density
	blockCount := 0
	for _, row := range template {
		for _, cell := range row {
			if cell.IsBlock() {
				blockCount++
			}
		}
	}

	density := float64(blockCount) / 49.0
	if density > 0.3 {
		t.Errorf("block density too high: %.2f", density)
	}
}

func TestOrchestrator_ShouldBeBlock(t *testing.T) {
	config := DefaultConfig()
	config.GridSize = [2]int{11, 11}

	mock := llm.NewMockClient()
	validatingClient := llm.NewValidatingClient(mock, llm.DefaultConfig())

	orch := NewOrchestrator(validatingClient, languagepack.NewFrenchPack(), nil, config)

	// Center should never be a block
	if orch.shouldBeBlock(5, 5, 11, 11) {
		t.Error("center should not be a block")
	}
}

func TestOrchestrator_BuildSlotInfos(t *testing.T) {
	config := DefaultConfig()
	mock := llm.NewMockClient()
	validatingClient := llm.NewValidatingClient(mock, llm.DefaultConfig())

	orch := NewOrchestrator(validatingClient, languagepack.NewFrenchPack(), nil, config)

	slots := []fill.Slot{
		{ID: 0, Length: 4},
		{ID: 1, Length: 5},
	}

	fillResult := &fill.Result{
		Words: map[int]string{
			0: "TEST",
			1: "WORDS",
		},
	}

	infos := orch.buildSlotInfos(slots, fillResult)

	if len(infos) != 2 {
		t.Errorf("expected 2 slot infos, got %d", len(infos))
	}

	if infos[0].Answer != "TEST" {
		t.Errorf("expected answer 'TEST', got %q", infos[0].Answer)
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.MaxAttempts <= 0 {
		t.Error("MaxAttempts should be positive")
	}
	if config.Timeout <= 0 {
		t.Error("Timeout should be positive")
	}
	if config.TargetDifficulty < 1 || config.TargetDifficulty > 5 {
		t.Error("TargetDifficulty should be 1-5")
	}
	if config.MinQAScore <= 0 {
		t.Error("MinQAScore should be positive")
	}
	if config.GridSize[0] <= 0 || config.GridSize[1] <= 0 {
		t.Error("GridSize should be positive")
	}
}

func TestSortClues(t *testing.T) {
	// Test is internal but we can test the sorting behavior through the result
	// This is a placeholder for more comprehensive tests
}

// Integration test (skipped by default, requires API key)
func TestOrchestrator_Generate_Integration(t *testing.T) {
	t.Skip("Integration test requires API key")

	// This would be an integration test that actually calls the LLM
	// For now, we skip it in automated tests
	ctx := context.Background()
	_ = ctx
}
