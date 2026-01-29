package logger

import (
	"backend-blog/config"
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
)

// TraceHandler 实现 slog.Handler 接口
type TraceHandler struct {
	slog.Handler
}

// Handle 拦截日志记录，注入 Trace ID
func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	//if ctx != nil {
	//	if traceID, ok := ctx.Value(pkg.TraceKey).(string); ok {
	//		r.AddAttrs(slog.String("trace_id", traceID))
	//	}
	//}
	return h.Handler.Handle(ctx, r)
}

// FanoutHandler 分发日志到多个 Handler
type FanoutHandler struct {
	handlers []slog.Handler
}

func (h *FanoutHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, l) {
			return true
		}
	}
	return false
}

func (h *FanoutHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, r.Level) {
			_ = handler.Handle(ctx, r)
		}
	}
	return nil
}

func (h *FanoutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return &FanoutHandler{handlers: newHandlers}
}

func (h *FanoutHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return &FanoutHandler{handlers: newHandlers}
}

// PrettyHandler 自定义文本格式 Handler
type PrettyHandler struct {
	w        io.Writer
	colorize bool
	opts     *slog.HandlerOptions
	attrs    []slog.Attr
	group    string
}

func NewPrettyHandler(w io.Writer, colorize bool, opts *slog.HandlerOptions) *PrettyHandler {
	return &PrettyHandler{w: w, colorize: colorize, opts: opts}
}

func (h *PrettyHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return l >= h.opts.Level.Level()
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	timeStr := r.Time.Format("2006-01-02 15:04:05")
	level := r.Level.String()
	msg := r.Message

	// Colorize Level
	if h.colorize {
		switch r.Level {
		case slog.LevelDebug:
			level = Cyan + level + Reset
		case slog.LevelInfo:
			level = Green + level + Reset
		case slog.LevelWarn:
			level = Yellow + level + Reset
		case slog.LevelError:
			level = Red + level + Reset
		}
	}

	// Build Attributes string
	var attrsStr strings.Builder
	// Handle attributes from WithAttrs
	for _, a := range h.attrs {
		h.appendAttr(&attrsStr, a)
	}
	// Handle record attributes
	r.Attrs(func(a slog.Attr) bool {
		h.appendAttr(&attrsStr, a)
		return true
	})

	// Add Source if enabled
	sourceStr := ""
	if h.opts.AddSource && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		sourceStr = fmt.Sprintf(" %s:%d", filepath.Base(f.File), f.Line)
	}

	// Format: Time [Level] Message Source Attrs
	logLine := fmt.Sprintf("%s [%s] %s%s %s\n", timeStr, level, msg, sourceStr, attrsStr.String())

	h.w.Write([]byte(logLine))
	return nil
}

func (h *PrettyHandler) appendAttr(sb *strings.Builder, a slog.Attr) {
	if a.Key == "" {
		return
	}
	// Simple key=value format
	sb.WriteString(fmt.Sprintf("%s=%v ", a.Key, a.Value.Any()))
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := append(h.attrs[:len(h.attrs):len(h.attrs)], attrs...)
	return &PrettyHandler{w: h.w, colorize: h.colorize, opts: h.opts, attrs: newAttrs, group: h.group}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	// Simplified group support: just prefix or ignore for now as text format is simple
	return &PrettyHandler{w: h.w, colorize: h.colorize, opts: h.opts, attrs: h.attrs, group: name}
}

var manager *LogManager

func Setup(cfg config.Log, isDev bool) {
	manager = &LogManager{
		dirPath:  cfg.Path,
		filename: cfg.Name,
		maxSize:  cfg.MaxSize * 1024 * 1024,
	}

	if err := manager.rotate(); err != nil {
		log.Fatalf("Logger init failed: %v", err)
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug, // Default to debug for dev
	}

	// Console Handler: Colored
	consoleHandler := NewPrettyHandler(os.Stderr, true, opts)

	// File Handler: Plain Text (No Color)
	fileHandler := NewPrettyHandler(manager, false, opts)

	// Combine
	fanout := &FanoutHandler{handlers: []slog.Handler{consoleHandler, fileHandler}}

	// Trace Middleware
	handler := &TraceHandler{fanout}

	// Set Default
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

// --- 以下保留你原有的 LogManager 逻辑，无需改动 ---

type LogManager struct {
	mu          sync.Mutex
	dirPath     string
	filename    string
	maxSize     int64
	currentFile *os.File
	currentSize int64
	currentDate string
}

func (m *LogManager) Write(p []byte) (n int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.shouldRotate(len(p)) {
		if err := m.rotate(); err != nil {
			return 0, err
		}
	}
	return m.currentFile.Write(p)
}

func (m *LogManager) shouldRotate(writeLen int) bool {
	today := time.Now().Format("2006-01-02")
	return today != m.currentDate || m.currentSize+int64(writeLen) > m.maxSize
}

func (m *LogManager) rotate() error {
	now := time.Now()
	today := now.Format("2006-01-02")
	folderPath := filepath.Join(m.dirPath, today)
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return err
	}
	if m.currentFile != nil {
		_ = m.currentFile.Close()
	}
	nextIndex := m.findNextFileIndex(folderPath)
	filePath := filepath.Join(folderPath, fmt.Sprintf("%s_%d.log", m.filename, nextIndex))
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	m.currentFile = file
	m.currentDate = today
	info, _ := file.Stat()
	m.currentSize = info.Size()
	return nil
}

func (m *LogManager) findNextFileIndex(dir string) int {
	entries, _ := os.ReadDir(dir)
	maxIdx := -1
	baseName := m.filename + "_"
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), baseName) && strings.HasSuffix(entry.Name(), ".log") {
			numStr := strings.TrimSuffix(strings.TrimPrefix(entry.Name(), baseName), ".log")
			if idx, err := strconv.Atoi(numStr); err == nil && idx > maxIdx {
				maxIdx = idx
			}
		}
	}
	if maxIdx != -1 {
		lastFile := filepath.Join(dir, fmt.Sprintf("%s_%d.log", m.filename, maxIdx))
		if info, err := os.Stat(lastFile); err == nil && info.Size() < m.maxSize {
			return maxIdx
		}
	}
	return maxIdx + 1
}
