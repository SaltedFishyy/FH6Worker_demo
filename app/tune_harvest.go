package main

import (
	"context"
	"fmt"

	"fh6worker/internal/harvest"
	"fh6worker/internal/storage"
)

func (a *App) RunTuneHarvest(options storage.TuneHarvestOptions) (*storage.TuneHarvestRunResult, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	baseCtx := a.ctx
	if baseCtx == nil {
		baseCtx = context.Background()
	}
	ctx, cancel := context.WithCancel(baseCtx)
	a.mu.Lock()
	if a.tuneHarvestCancel != nil {
		a.mu.Unlock()
		cancel()
		return nil, fmt.Errorf("tune harvest is already running")
	}
	active := &tuneHarvestCancellation{cancel: cancel}
	a.tuneHarvestCancel = active
	a.mu.Unlock()
	defer func() {
		a.mu.Lock()
		if a.tuneHarvestCancel == active {
			a.tuneHarvestCancel = nil
		}
		a.mu.Unlock()
		cancel()
	}()
	return harvest.Run(ctx, a.store, options)
}

func (a *App) StopTuneHarvest() error {
	a.mu.Lock()
	active := a.tuneHarvestCancel
	a.mu.Unlock()
	if active != nil {
		active.cancel()
	}
	return nil
}

func (a *App) ListTuneHarvestCandidates(status string, limit int) ([]storage.TuneHarvestCandidate, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.ListTuneHarvestCandidates(status, limit)
}

func (a *App) SearchTuneHarvestCandidates(status string, search string, limit int) ([]storage.TuneHarvestCandidate, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.SearchTuneHarvestCandidates(status, search, limit)
}

func (a *App) ClearTuneHarvestCandidates() (int64, error) {
	if err := a.ensureStore(); err != nil {
		return 0, err
	}
	return a.store.ClearTuneHarvestCandidates()
}

func (a *App) UpdateTuneHarvestCandidateStatus(id int64, status string, reason string) (*storage.TuneHarvestCandidate, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.UpdateTuneHarvestCandidateStatus(id, status, reason)
}

func (a *App) ListFH6Cars() ([]storage.FH6Car, error) {
	if err := a.ensureStore(); err != nil {
		return nil, err
	}
	return a.store.ListFH6Cars()
}
