package database

import (
	"database/sql"
	"sync"
	"sync/atomic"
	"time"
)

type PoolStats struct {
	OpenConnections   int
	InUse             int
	Idle              int
	WaitCount         int64
	WaitDuration      string
	MaxIdleClosed     int64
	MaxLifetimeClosed int64
	// Performance Metrics
	AvgQueryTime     string
	QueriesPerMinute float64
	MaxQueryDuration string
}

type metrics struct {
	totalDuration int64
	totalQueries  int64
	maxDuration   int64
	startTime     time.Time
}

var (
	registry sync.Map // Map[*sql.DB]*metrics
)

// getTracker internal helper to ensure a tracker exists for the DB
func getTracker(db *sql.DB) *metrics {
	m, _ := registry.LoadOrStore(db, &metrics{startTime: time.Now()})
	return m.(*metrics)
}

// RecordQuery must be called manually after your database operations.
// It maps the duration to the specific DB instance.
func RecordQuery(db *sql.DB, duration time.Duration) {
	t := getTracker(db)
	dur := int64(duration)

	atomic.AddInt64(&t.totalDuration, dur)
	atomic.AddInt64(&t.totalQueries, 1)

	for {
		currentMax := atomic.LoadInt64(&t.maxDuration)
		if dur <= currentMax {
			break
		}
		if atomic.CompareAndSwapInt64(&t.maxDuration, currentMax, dur) {
			break
		}
	}
}

func Stats(db *sql.DB) PoolStats {
	stats := db.Stats()
	t := getTracker(db)

	totalQueries := atomic.LoadInt64(&t.totalQueries)
	totalDuration := atomic.LoadInt64(&t.totalDuration)
	maxDuration := atomic.LoadInt64(&t.maxDuration)

	// Calculate Average Query Time: $$Avg = \frac{\sum Duration}{N}$$
	avgTime := time.Duration(0)
	if totalQueries > 0 {
		avgTime = time.Duration(totalDuration / totalQueries)
	}

	// Calculate Throughput: $$QPM = \frac{N}{Minutes_{uptime}}$$
	uptimeMinutes := time.Since(t.startTime).Minutes()
	qpm := 0.0
	if uptimeMinutes > 0 {
		qpm = float64(totalQueries) / uptimeMinutes
	}

	return PoolStats{
		OpenConnections:   stats.OpenConnections,
		InUse:             stats.InUse,
		Idle:              stats.Idle,
		WaitCount:         stats.WaitCount,
		WaitDuration:      stats.WaitDuration.String(),
		MaxIdleClosed:     stats.MaxIdleClosed,
		MaxLifetimeClosed: stats.MaxLifetimeClosed,
		// Performance Mapping
		AvgQueryTime:     avgTime.String(),
		QueriesPerMinute: qpm,
		MaxQueryDuration: time.Duration(maxDuration).String(),
	}
}
