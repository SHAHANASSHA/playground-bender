package cleanup_manager

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func CleanupManager(ctx context.Context, interval, maxAge time.Duration) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("unable to create Docker client: %w", err)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			containers, err := cli.ContainerList(ctx, container.ListOptions{All: false})
			if err != nil {
				fmt.Println("Error listing containers:", err)
				continue
			}

			now := time.Now()
			log.Printf("Starting cleanup manager...%d", len(containers))
			for _, c := range containers {
				created := time.Unix(c.Created, 0)
				if runtime := now.Sub(created); runtime > maxAge {
					name := "(no-name)"
					if len(c.Names) > 0 {
						name = c.Names[0]
					}
					fmt.Printf("⚠️ Container %s (ID %s) has been running for %v\n", name, c.ID[:12], runtime)
				}
			}
		}
	}
}
