package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func RunScheduler(cfg *Config) {
	fmt.Printf("Go Directory Backup version: %s\n", Version)
	BackupNow(cfg) // Immediate backup on start
	for {
		// Reload config before each backup
		newCfg, err := LoadOrInitConfig()
		if err == nil {
			cfg = newCfg
		}
		next := nextRun(cfg.SyncTime)
		dur := time.Until(next)
		fmt.Printf("Next backup at %s\n", next.Format(time.RFC1123))
		time.Sleep(dur)
		BackupNow(cfg)
	}
}

func nextRun(syncTime string) time.Time {
	if d, err := time.ParseDuration(syncTime); err == nil {
		return time.Now().Add(d)
	}
	t := time.Now()
	parts := strings.Split(syncTime, ":")
	if len(parts) == 2 {
		hour, _ := strconv.Atoi(parts[0])
		min, _ := strconv.Atoi(parts[1])
		next := time.Date(t.Year(), t.Month(), t.Day(), hour, min, 0, 0, t.Location())
		if next.Before(t) {
			next = next.Add(24 * time.Hour)
		}
		return next
	}
	return t.Add(24 * time.Hour)
}

func BackupNow(cfg *Config) {
	dateFolder := time.Now().Format("2006-01-02")
	target := filepath.Join(cfg.TargetDir, dateFolder)
	os.MkdirAll(target, 0755)
	for _, src := range cfg.SourceDirs {
		base := filepath.Base(src)
		dst := filepath.Join(target, base)
		fmt.Printf("Backing up %s to %s\n", src, dst)
		CopyDir(src, dst)
	}
	fmt.Println("Backup complete.")
}

func CopyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}
		return copyFile(path, targetPath, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()
	dstF, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstF.Close()
	_, err = io.Copy(dstF, srcF)
	if err != nil {
		return err
	}
	return os.Chmod(dst, mode)
}
