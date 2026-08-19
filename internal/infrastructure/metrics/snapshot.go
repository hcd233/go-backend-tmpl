package metrics

import (
	"strconv"
	"time"

	"github.com/hcd233/aris-api-tmpl/internal/common/constant"
	"github.com/prometheus/client_golang/prometheus"
	metricpb "github.com/prometheus/client_model/go"
)

// Snapshot 单个 instance 在某一时刻的运行时指标快照（写入 Redis 的最小单位）。
//
// 仅存"可直接相加的原值"：gauge 原值 + counter 累计值 + histogram 桶计数；
// 速率与分位的计算全部留给聚合层。
type Snapshot struct {
	TS         int64              `json:"ts"`                   // unix 秒
	Goroutines float64            `json:"goroutines"`           // gauge
	HeapBytes  float64            `json:"heapBytes"`            // gauge
	CPUSeconds float64            `json:"cpuSeconds"`           // counter 累计值 → 聚合层求 CPU%
	LatBuckets map[string]float64 `json:"latBuckets,omitempty"` // le -> 累计计数 → 聚合层求 P95
	LatCount   float64            `json:"latCount"`             // histogram 累计样本数 → 聚合层求 QPS
	ReqTotal   float64            `json:"reqTotal"`             // counter 累计业务请求数 → 聚合层求成功率
	ReqSuccess float64            `json:"reqSuccess"`           // counter 累计 200 请求数 → 聚合层求成功率
}

// SnapshotStore flusher 写入快照所需的存储能力（由 cache.RuntimeMetricsCache 实现）。
type SnapshotStore interface {
	WriteSnapshot(instanceID string, score int64, payload []byte, retentionCutoff int64) error
}

// BuildSnapshot 从 Gatherer 采集当前所有运行时指标，组装成一份快照。
func BuildSnapshot(gatherer prometheus.Gatherer, now time.Time) (*Snapshot, error) {
	families, err := gatherer.Gather()
	if err != nil {
		return nil, err
	}

	byName := make(map[string]*metricpb.MetricFamily, len(families))
	for _, f := range families {
		byName[f.GetName()] = f
	}

	requests := byName[constant.MetricFullHTTPRequests]
	reqSuccess := labeledCounterValue(requests, constant.MetricLabelResult, constant.HTTPResultSuccess)
	snap := &Snapshot{
		TS:         now.Unix(),
		Goroutines: firstGaugeValue(byName[constant.MetricFullGoGoroutines]),
		HeapBytes:  firstGaugeValue(byName[constant.MetricFullGoHeapAlloc]),
		CPUSeconds: firstCounterValue(byName[constant.MetricFullProcessCPU]),
		ReqTotal:   reqSuccess + labeledCounterValue(requests, constant.MetricLabelResult, constant.HTTPResultFailure),
		ReqSuccess: reqSuccess,
	}
	snap.LatBuckets, snap.LatCount = histogramBuckets(byName[constant.MetricFullRequestDuration])
	return snap, nil
}

func firstGaugeValue(f *metricpb.MetricFamily) float64 {
	if f == nil || len(f.GetMetric()) == 0 {
		return 0
	}
	return f.GetMetric()[0].GetGauge().GetValue()
}

func firstCounterValue(f *metricpb.MetricFamily) float64 {
	if f == nil || len(f.GetMetric()) == 0 {
		return 0
	}
	return f.GetMetric()[0].GetCounter().GetValue()
}

// labeledCounterValue 取指定 label 值对应的 counter 累计值；缺失时返回 0。
func labeledCounterValue(f *metricpb.MetricFamily, label, want string) float64 {
	if f == nil || len(f.GetMetric()) == 0 {
		return 0
	}
	for _, m := range f.GetMetric() {
		key := ""
		for _, l := range m.GetLabel() {
			if l.GetName() == label {
				key = l.GetValue()
				break
			}
		}
		if key == want {
			return m.GetCounter().GetValue()
		}
	}
	return 0
}

// histogramBuckets 抽取 histogram 各桶累计计数与样本总数。
func histogramBuckets(f *metricpb.MetricFamily) (buckets map[string]float64, count float64) {
	if f == nil || len(f.GetMetric()) == 0 {
		return nil, 0
	}
	h := f.GetMetric()[0].GetHistogram()
	if h == nil {
		return nil, 0
	}
	buckets = make(map[string]float64, len(h.GetBucket()))
	for _, b := range h.GetBucket() {
		buckets[strconv.FormatFloat(b.GetUpperBound(), 'f', -1, 64)] = float64(b.GetCumulativeCount())
	}
	return buckets, float64(h.GetSampleCount())
}
