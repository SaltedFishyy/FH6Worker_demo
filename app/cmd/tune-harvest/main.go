package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"fh6worker/internal/harvest"
	"fh6worker/internal/storage"
)

func main() {
	sourcesFlag := flag.String("sources", "jsr,codmunity", "comma-separated sources: jsr,codmunity,forzafire")
	dryRun := flag.Bool("dry-run", true, "collect without writing candidates to the local database")
	limit := flag.Int("limit", 80, "maximum ForzaFire detail pages to inspect")
	jsonOut := flag.Bool("json", false, "print full JSON result")
	flag.Parse()

	store, err := storage.OpenDefault()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer store.Close()

	result, err := harvest.Run(context.Background(), store, storage.TuneHarvestOptions{
		Sources:        splitSources(*sourcesFlag),
		DryRun:         *dryRun,
		LimitPerSource: *limit,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if *jsonOut {
		encoded, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(encoded))
		return
	}
	fmt.Printf("found=%d saved=%d pending=%d rejected=%d warnings=%d\n", result.Found, result.Saved, result.Pending, result.Rejected, len(result.Warnings))
	for _, candidate := range result.Candidates {
		fmt.Printf("%s %s %s %s score=%.3f %s\n",
			candidate.Source, storage.FormatTuneShareCode(candidate.ShareCode), candidate.CarName, candidate.UseCase, candidate.MatchScore, candidate.MatchReason)
	}
}

func splitSources(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
