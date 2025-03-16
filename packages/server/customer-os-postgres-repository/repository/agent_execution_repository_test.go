package postgres_repository

import (
	"context"
	"testing"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/pkg/errors"

	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

func TestAgentExecutionRepository_CreateAndGet(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase()(t)

	// Create a new AgentExecution record to save
	exec := &postgresentity.AgentExecution{
		AgentID:      utils.StringPtr("agent-001"),
		Tenant:       "tenant1",
		TriggerEvent: "Test Event",
		Status:       enum.AgentExecutionPending, // default status
		MaxRetries:   3,
	}

	// Save the AgentExecution record
	savedExec, err := repositories.AgentExecutionRepository.Create(ctx, *exec)
	if err != nil {
		t.Fatalf("failed to create agent execution: %v", err)
	}
	if savedExec.ID == "" {
		t.Errorf("expected agent execution to have an ID, got empty")
	}

	// Retrieve the AgentExecution record
	fetched, err := repositories.AgentExecutionRepository.GetById(ctx, savedExec.ID)
	if err != nil {
		t.Fatalf("failed to get agent execution: %v", err)
	}
	if fetched == nil {
		t.Fatalf("expected agent execution to be fetched, got nil")
	}
	if fetched.Tenant != "tenant1" {
		t.Errorf("expected tenant 'tenant1', got '%s'", fetched.Tenant)
	}
	if fetched.TriggerEvent != "Test Event" {
		t.Errorf("expected TriggerEvent 'Test Event', got '%s'", fetched.TriggerEvent)
	}
}

func TestAgentExecutionRepository_Find(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase()(t)

	// Create a record to later find
	exec := &postgresentity.AgentExecution{
		AgentID:      utils.StringPtr("agent-002"),
		Tenant:       "tenant1",
		TriggerEvent: "FindTest",
		Status:       enum.AgentExecutionPending,
		MaxRetries:   3,
	}
	_, err := repositories.AgentExecutionRepository.Create(ctx, *exec)
	if err != nil {
		t.Fatalf("failed to create record: %v", err)
	}

	// Use Find to retrieve the record based on criteria
	criteria := postgresentity.AgentExecution{
		AgentID:      utils.StringPtr("agent-002"),
		TriggerEvent: "FindTest",
	}
	found, err := repositories.AgentExecutionRepository.Find(ctx, criteria)
	if err != nil {
		t.Fatalf("failed to find record: %v", err)
	}
	if found == nil {
		t.Fatal("expected record to be found, got nil")
	}
	if found.TriggerEvent != "FindTest" {
		t.Errorf("expected TriggerEvent 'FindTest', got %s", found.TriggerEvent)
	}
}

func TestAgentExecutionRepository_Fail(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase()(t)

	exec := &postgresentity.AgentExecution{
		AgentID:      utils.StringPtr("agent-003"),
		Tenant:       "tenant1",
		TriggerEvent: "FailTest",
		Status:       enum.AgentExecutionPending,
		MaxRetries:   3,
	}
	saved, err := repositories.AgentExecutionRepository.Create(ctx, *exec)
	if err != nil {
		t.Fatalf("failed to create record: %v", err)
	}

	errMsg := "failure occurred"
	if err := repositories.AgentExecutionRepository.Fail(ctx, saved.ID, errMsg); err != nil {
		t.Fatalf("Fail method failed: %v", err)
	}

	fetched, err := repositories.AgentExecutionRepository.GetById(ctx, saved.ID)
	if err != nil {
		t.Fatalf("GetById failed: %v", err)
	}
	if fetched.Status != enum.AgentExecutionError {
		t.Errorf("expected status %s, got %s", enum.AgentExecutionError, fetched.Status)
	}
	if fetched.ErrorMessage == nil || *fetched.ErrorMessage != errMsg {
		t.Errorf("expected error message '%s', got %v", errMsg, fetched.ErrorMessage)
	}
}

func TestAgentExecutionRepository_Pending(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase()(t)

	// Create a record with an error status initially
	exec := &postgresentity.AgentExecution{
		AgentID:      utils.StringPtr("agent-004"),
		Tenant:       "tenant1",
		TriggerEvent: "PendingTest",
		Status:       enum.AgentExecutionError,
		MaxRetries:   3,
	}
	saved, err := repositories.AgentExecutionRepository.Create(ctx, *exec)
	if err != nil {
		t.Fatalf("failed to create record: %v", err)
	}

	// Set the execution status back to pending
	if err := repositories.AgentExecutionRepository.Pending(ctx, saved.ID); err != nil {
		t.Fatalf("Pending method failed: %v", err)
	}

	fetched, err := repositories.AgentExecutionRepository.GetById(ctx, saved.ID)
	if err != nil {
		t.Fatalf("GetById failed: %v", err)
	}
	if fetched.Status != enum.AgentExecutionPending {
		t.Errorf("expected status %s, got %s", enum.AgentExecutionPending, fetched.Status)
	}
}

func TestAgentExecutionRepository_Finish(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase()(t)

	// Create a record with status Running (only running executions can be finished)
	exec := &postgresentity.AgentExecution{
		AgentID:      utils.StringPtr("agent-005"),
		Tenant:       "tenant1",
		TriggerEvent: "FinishTest",
		Status:       enum.AgentExecutionRunning,
		MaxRetries:   3,
	}
	saved, err := repositories.AgentExecutionRepository.Create(ctx, *exec)
	if err != nil {
		t.Fatalf("failed to create record: %v", err)
	}

	if err := repositories.AgentExecutionRepository.Finish(ctx, saved.ID); err != nil {
		t.Fatalf("Finish method failed: %v", err)
	}

	fetched, err := repositories.AgentExecutionRepository.GetById(ctx, saved.ID)
	if err != nil {
		t.Fatalf("GetById failed: %v", err)
	}
	if fetched.Status != enum.AgentExecutionCompleted {
		t.Errorf("expected status %s, got %s", enum.AgentExecutionCompleted, fetched.Status)
	}
}

func TestAgentExecutionRepository_Completed(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase()(t)

	exec := &postgresentity.AgentExecution{
		AgentID:      utils.StringPtr("agent-006"),
		Tenant:       "tenant1",
		TriggerEvent: "CompletedTest",
		Status:       enum.AgentExecutionPending,
		MaxRetries:   3,
	}
	saved, err := repositories.AgentExecutionRepository.Create(ctx, *exec)
	if err != nil {
		t.Fatalf("failed to create record: %v", err)
	}

	goalAchieved := true
	updated, err := repositories.AgentExecutionRepository.Completed(ctx, saved.ID, &goalAchieved)
	if err != nil {
		t.Fatalf("Completed method failed: %v", err)
	}
	if updated.Status != enum.AgentExecutionCompleted {
		t.Errorf("expected status %s, got %s", enum.AgentExecutionCompleted, updated.Status)
	}
	if updated.GoalAchieved == nil || *updated.GoalAchieved != goalAchieved {
		t.Errorf("expected GoalAchieved %v, got %v", goalAchieved, updated.GoalAchieved)
	}
	if updated.CompletedAt == nil {
		t.Error("expected CompletedAt to be set")
	}
}

func TestAgentExecutionRepository_ScheduleRetry_FirstTime(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase()(t)

	exec := &postgresentity.AgentExecution{
		AgentID:      utils.StringPtr("agent-007"),
		Tenant:       "tenant1",
		TriggerEvent: "ScheduleRetryTest",
		Status:       enum.AgentExecutionPending,
		MaxRetries:   3,
		NextRetryAt:  nil,
		RetryCount:   0,
	}
	saved, err := repositories.AgentExecutionRepository.Create(ctx, *exec)
	if err != nil {
		t.Fatalf("failed to create record: %v", err)
	}

	retryErr := errors.New("temporary failure")
	stateData := map[string]any{"attempt": 1}
	_, err = repositories.AgentExecutionRepository.ScheduleRetry(ctx, saved.ID, retryErr, stateData)
	if err != nil {
		t.Fatalf("ScheduleRetry failed: %v", err)
	}

	fetched, err := repositories.AgentExecutionRepository.GetById(ctx, saved.ID)
	if err != nil {
		t.Fatalf("GetById failed: %v", err)
	}
	if fetched.Status != enum.AgentExecutionRetrying {
		t.Errorf("expected status %s, got %s", enum.AgentExecutionRetrying, fetched.Status)
	}
	if fetched.RetryCount != 0 {
		t.Error("expected RetryCount should remain 0")
	}
	if fetched.NextRetryAt == nil {
		t.Error("expected NextRetryAt to be set")
	}
}

func TestAgentExecutionRepository_ScheduleRetry(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase()(t)

	exec := &postgresentity.AgentExecution{
		AgentID:      utils.StringPtr("agent-007"),
		Tenant:       "tenant1",
		TriggerEvent: "ScheduleRetryTest",
		Status:       enum.AgentExecutionPending,
		MaxRetries:   3,
		NextRetryAt:  utils.NowPtr(),
		RetryCount:   0,
	}
	saved, err := repositories.AgentExecutionRepository.Create(ctx, *exec)
	if err != nil {
		t.Fatalf("failed to create record: %v", err)
	}

	retryErr := errors.New("temporary failure")
	stateData := map[string]any{"attempt": 1}
	_, err = repositories.AgentExecutionRepository.ScheduleRetry(ctx, saved.ID, retryErr, stateData)
	if err != nil {
		t.Fatalf("ScheduleRetry failed: %v", err)
	}

	fetched, err := repositories.AgentExecutionRepository.GetById(ctx, saved.ID)
	if err != nil {
		t.Fatalf("GetById failed: %v", err)
	}
	if fetched.Status != enum.AgentExecutionRetrying {
		t.Errorf("expected status %s, got %s", enum.AgentExecutionRetrying, fetched.Status)
	}
	if fetched.RetryCount < 1 {
		t.Error("expected RetryCount to be incremented")
	}
	if fetched.NextRetryAt == nil {
		t.Error("expected NextRetryAt to be set")
	}
}

func TestAgentExecutionRepository_SaveAsyncState(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase()(t)

	exec := &postgresentity.AgentExecution{
		AgentID:      utils.StringPtr("agent-008"),
		Tenant:       "tenant1",
		TriggerEvent: "AsyncStateTest",
		Status:       enum.AgentExecutionPending,
		MaxRetries:   3,
	}
	saved, err := repositories.AgentExecutionRepository.Create(ctx, *exec)
	if err != nil {
		t.Fatalf("failed to create record: %v", err)
	}

	newStep := "processing"
	newStateData := map[string]any{"progress": 75}
	if err := repositories.AgentExecutionRepository.SaveAsyncState(ctx, saved.ID, newStep, newStateData); err != nil {
		t.Fatalf("SaveAsyncState failed: %v", err)
	}

	fetched, err := repositories.AgentExecutionRepository.GetById(ctx, saved.ID)
	if err != nil {
		t.Fatalf("GetById failed: %v", err)
	}
	if fetched.CurrentStep != newStep {
		t.Errorf("expected CurrentStep %q, got %q", newStep, fetched.CurrentStep)
	}
	if fetched.Status != enum.AgentExecutionPending {
		t.Errorf("expected status %s, got %s", enum.AgentExecutionPending, fetched.Status)
	}
}

func TestAgentExecutionRepository_CompleteStep(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase()(t)

	exec := &postgresentity.AgentExecution{
		AgentID:      utils.StringPtr("agent-009"),
		Tenant:       "tenant1",
		TriggerEvent: "CompleteStepTest",
		Status:       enum.AgentExecutionPending,
		MaxRetries:   3,
		CurrentStep:  "step1",
		StateData:    map[string]any{"data": "value"},
		Checkpoints:  make(map[string]any),
	}
	saved, err := repositories.AgentExecutionRepository.Create(ctx, *exec)
	if err != nil {
		t.Fatalf("failed to create record: %v", err)
	}

	stepResult := map[string]any{"result": "success"}
	if err := repositories.AgentExecutionRepository.CompleteStep(ctx, saved.ID, "step1", stepResult); err != nil {
		t.Fatalf("CompleteStep failed: %v", err)
	}

	fetched, err := repositories.AgentExecutionRepository.GetById(ctx, saved.ID)
	if err != nil {
		t.Fatalf("GetById failed: %v", err)
	}
	if fetched.Checkpoints == nil {
		t.Fatal("expected checkpoints to be initialized")
	}
	cp, exists := fetched.Checkpoints["step1"]
	if !exists {
		t.Error("expected a checkpoint for step1")
	} else {
		resultMap, ok := cp.(map[string]any)
		if !ok || resultMap["result"] != "success" {
			t.Errorf("expected step result 'success', got %v", cp)
		}
	}
	if fetched.CurrentStep != "" {
		t.Errorf("expected CurrentStep to be cleared, got %q", fetched.CurrentStep)
	}
	if len(fetched.StateData) > 0 {
		t.Errorf("expected StateData to be cleared, got %v", fetched.StateData)
	}
}

func TestAgentExecutionRepository_GoalAchieved(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase()(t)

	exec := &postgresentity.AgentExecution{
		AgentID:      utils.StringPtr("agent-010"),
		Tenant:       "tenant1",
		TriggerEvent: "GoalAchievedTest",
		Status:       enum.AgentExecutionPending,
		MaxRetries:   3,
	}
	saved, err := repositories.AgentExecutionRepository.Create(ctx, *exec)
	if err != nil {
		t.Fatalf("failed to create record: %v", err)
	}

	if err := repositories.AgentExecutionRepository.GoalAchieved(ctx, saved.ID, true, utils.StringPtr("123")); err != nil {
		t.Fatalf("GoalAchieved failed: %v", err)
	}

	fetched, err := repositories.AgentExecutionRepository.GetById(ctx, saved.ID)
	if err != nil {
		t.Fatalf("GetById failed: %v", err)
	}
	if fetched.GoalAchieved == nil || !*fetched.GoalAchieved {
		t.Error("expected GoalAchieved to be true")
	}
	if fetched.CompletedAt == nil {
		t.Error("expected CompletedAt to be set")
	}
	if fetched.ImpactedId == nil || *fetched.ImpactedId != "123" {
		t.Errorf("expected ImpactedId '123', got %v", fetched.ImpactedId)
	}
}

func TestAgentExecutionRepository_GetGoalAchievedCountLast30Days(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase()(t)

	agentID := "agent-011"
	now := time.Now()

	execTrue := &postgresentity.AgentExecution{
		AgentID:      utils.StringPtr(agentID),
		Tenant:       "tenant1",
		TriggerEvent: "GoalCountTest",
		Status:       enum.AgentExecutionCompleted,
		MaxRetries:   3,
		GoalAchieved: utils.BoolPtr(true),
		UpdatedAt:    &now,
	}
	execFalse := &postgresentity.AgentExecution{
		AgentID:      utils.StringPtr(agentID),
		Tenant:       "tenant1",
		TriggerEvent: "GoalCountTest",
		Status:       enum.AgentExecutionCompleted,
		MaxRetries:   3,
		GoalAchieved: utils.BoolPtr(false),
		UpdatedAt:    &now,
	}

	// Create two records with goal achieved true and one with false.
	if _, err := repositories.AgentExecutionRepository.Create(ctx, *execTrue); err != nil {
		t.Fatalf("failed to create execTrue: %v", err)
	}
	if _, err := repositories.AgentExecutionRepository.Create(ctx, *execTrue); err != nil {
		t.Fatalf("failed to create second execTrue: %v", err)
	}
	if _, err := repositories.AgentExecutionRepository.Create(ctx, *execFalse); err != nil {
		t.Fatalf("failed to create execFalse: %v", err)
	}

	count, err := repositories.AgentExecutionRepository.GetGoalAchievedCountLast30Days(ctx, agentID)
	if err != nil {
		t.Fatalf("GetGoalAchievedCountLast30Days failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected goal achieved count 2, got %d", count)
	}
}

func TestAgentExecutionRepository_GetExecutionsForRetry(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase()(t)

	agentID := "agent-012"
	pastTime := time.Now().Add(-time.Hour)
	exec := &postgresentity.AgentExecution{
		AgentID:      utils.StringPtr(agentID),
		Tenant:       "tenant1",
		TriggerEvent: "RetryQueryTest",
		Status:       enum.AgentExecutionRetrying,
		MaxRetries:   3,
		RetryCount:   1,
		NextRetryAt:  &pastTime,
	}
	if _, err := repositories.AgentExecutionRepository.Create(ctx, *exec); err != nil {
		t.Fatalf("failed to create record: %v", err)
	}

	results, err := repositories.AgentExecutionRepository.GetExecutionsForRetry(ctx, 10)
	if err != nil {
		t.Fatalf("GetExecutionsForRetry failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 execution for retry, got %d", len(results))
	}
}
