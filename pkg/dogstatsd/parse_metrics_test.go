// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package dogstatsd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseMetricSample(rawSample []byte) (dogstatsdMetricSample, error) {
	parser := newParser(newFloat64ListPool())
	return parser.parseMetricSample(rawSample)
}

const epsilon = 0.00001

func TestParseGauge(t *testing.T) {
	sample, err := parseMetricSample([]byte("daemon:666|g"))

	assert.NoError(t, err)

	assert.Equal(t, "daemon", sample.name)
	assert.Equal(t, 666.0, sample.value)
	assert.InEpsilon(t, 666.0, sample.value, epsilon)
	require.Nil(t, sample.values)
	assert.Equal(t, gaugeType, sample.metricType)
	assert.Equal(t, []string{}, sample.tags.Sorted())
	assert.InEpsilon(t, 1.0, sample.sampleRate, epsilon)
}

func TestParseGaugeMultiple(t *testing.T) {
	sample, err := parseMetricSample([]byte("daemon:666:777|g"))

	assert.NoError(t, err)

	assert.Equal(t, "daemon", sample.name)
	assert.Len(t, sample.values, 2)
	assert.InEpsilon(t, 666.0, sample.values[0], epsilon)
	assert.InEpsilon(t, 777.0, sample.values[1], epsilon)
	assert.Equal(t, gaugeType, sample.metricType)
	assert.Equal(t, []string{}, sample.tags.Sorted())
	assert.InEpsilon(t, 1.0, sample.sampleRate, epsilon)
}

func TestParseCounter(t *testing.T) {
	sample, err := parseMetricSample([]byte("daemon:21|c"))

	assert.NoError(t, err)

	assert.Equal(t, "daemon", sample.name)
	assert.InEpsilon(t, 21.0, sample.value, epsilon)
	require.Nil(t, sample.values)
	assert.Equal(t, countType, sample.metricType)
	assert.Equal(t, []string{}, sample.tags.Sorted())
	assert.InEpsilon(t, 1.0, sample.sampleRate, epsilon)
}

func TestParseCounterMultiple(t *testing.T) {
	sample, err := parseMetricSample([]byte("daemon:666:777|c"))

	assert.NoError(t, err)

	assert.Equal(t, "daemon", sample.name)
	assert.Len(t, sample.values, 2)
	assert.InEpsilon(t, 666.0, sample.values[0], epsilon)
	assert.InEpsilon(t, 777.0, sample.values[1], epsilon)
	assert.Equal(t, countType, sample.metricType)
	assert.Equal(t, []string{}, sample.tags.Sorted())
	assert.InEpsilon(t, 1.0, sample.sampleRate, epsilon)
}

func TestParseCounterWithTags(t *testing.T) {
	sample, err := parseMetricSample([]byte("custom_counter:1|c|#protocol:http,bench"))

	assert.NoError(t, err)

	assert.Equal(t, "custom_counter", sample.name)
	assert.InEpsilon(t, 1.0, sample.value, epsilon)
	require.Nil(t, sample.values)
	assert.Equal(t, countType, sample.metricType)
	assert.Equal(t, []string{"bench", "protocol:http"}, sample.tags.Sorted())
	assert.InEpsilon(t, 1.0, sample.sampleRate, epsilon)
}

func TestParseHistogram(t *testing.T) {
	sample, err := parseMetricSample([]byte("daemon:21|h"))

	assert.NoError(t, err)

	assert.Equal(t, "daemon", sample.name)
	assert.InEpsilon(t, 21.0, sample.value, epsilon)
	require.Nil(t, sample.values)
	assert.Equal(t, histogramType, sample.metricType)
	assert.Equal(t, []string{}, sample.tags.Sorted())
	assert.InEpsilon(t, 1.0, sample.sampleRate, epsilon)
}

func TestParseHistogramrMultiple(t *testing.T) {
	sample, err := parseMetricSample([]byte("daemon:21:22|h"))

	assert.NoError(t, err)

	assert.Equal(t, "daemon", sample.name)
	assert.Len(t, sample.values, 2)
	assert.InEpsilon(t, 21.0, sample.values[0], epsilon)
	assert.InEpsilon(t, 22.0, sample.values[1], epsilon)
	assert.Equal(t, histogramType, sample.metricType)
	assert.Equal(t, []string{}, sample.tags.Sorted())
	assert.InEpsilon(t, 1.0, sample.sampleRate, epsilon)
}

func TestParseTimer(t *testing.T) {
	sample, err := parseMetricSample([]byte("daemon:21|ms"))

	assert.NoError(t, err)

	assert.Equal(t, "daemon", sample.name)
	assert.InEpsilon(t, 21.0, sample.value, epsilon)
	require.Nil(t, sample.values)
	assert.Equal(t, timingType, sample.metricType)
	assert.Equal(t, []string{}, sample.tags.Sorted())
	assert.InEpsilon(t, 1.0, sample.sampleRate, epsilon)
}

func TestParseTimerMultiple(t *testing.T) {
	sample, err := parseMetricSample([]byte("daemon:21:22|ms"))

	assert.NoError(t, err)

	assert.Equal(t, "daemon", sample.name)
	assert.Len(t, sample.values, 2)
	assert.InEpsilon(t, 21.0, sample.values[0], epsilon)
	assert.InEpsilon(t, 22.0, sample.values[1], epsilon)
	assert.Equal(t, timingType, sample.metricType)
	assert.Equal(t, []string{}, sample.tags.Sorted())
	assert.InEpsilon(t, 1.0, sample.sampleRate, epsilon)
}

func TestParseSet(t *testing.T) {
	sample, err := parseMetricSample([]byte("daemon:abc|s"))

	assert.NoError(t, err)

	assert.Equal(t, "daemon", sample.name)
	assert.Equal(t, "abc", sample.setValue)
	assert.Equal(t, setType, sample.metricType)
	assert.Equal(t, []string{}, sample.tags.Sorted())
	assert.InEpsilon(t, 1.0, sample.sampleRate, epsilon)
}

func TestParseSetMultiple(t *testing.T) {
	// multiple values are not supported for set. ':' can be part of the
	// set value for backward compatibility
	sample, err := parseMetricSample([]byte("daemon:abc:def|s"))

	assert.NoError(t, err)

	assert.Equal(t, "daemon", sample.name)
	assert.Equal(t, "abc:def", sample.setValue)
	assert.Equal(t, setType, sample.metricType)
	assert.Equal(t, []string{}, sample.tags.Sorted())
	assert.InEpsilon(t, 1.0, sample.sampleRate, epsilon)
}

func TestSampleDistribution(t *testing.T) {
	sample, err := parseMetricSample([]byte("daemon:3.5|d"))

	assert.NoError(t, err)

	assert.Equal(t, "daemon", sample.name)
	assert.InEpsilon(t, 3.5, sample.value, epsilon)
	require.Nil(t, sample.values)
	assert.Equal(t, distributionType, sample.metricType)
	assert.Equal(t, []string{}, sample.tags.Sorted())
}

func TestParseDistributionMultiple(t *testing.T) {
	sample, err := parseMetricSample([]byte("daemon:3.5:4.5|d"))

	assert.NoError(t, err)

	assert.Equal(t, "daemon", sample.name)
	assert.Len(t, sample.values, 2)
	assert.InEpsilon(t, 3.5, sample.values[0], epsilon)
	assert.InEpsilon(t, 4.5, sample.values[1], epsilon)
	assert.Equal(t, distributionType, sample.metricType)
	assert.Equal(t, []string{}, sample.tags.Sorted())
}

func TestParseSetUnicode(t *testing.T) {
	sample, err := parseMetricSample([]byte("daemon:♬†øU†øU¥ºuT0♪|s"))

	assert.NoError(t, err)

	assert.Equal(t, "daemon", sample.name)
	assert.Equal(t, "♬†øU†øU¥ºuT0♪", sample.setValue)
	assert.Equal(t, setType, sample.metricType)
	assert.Equal(t, []string{}, sample.tags.Sorted())
	assert.InEpsilon(t, 1.0, sample.sampleRate, epsilon)
}

func TestParseGaugeWithTags(t *testing.T) {
	sample, err := parseMetricSample([]byte("daemon:666|g|#sometag1:somevalue1,sometag2:somevalue2"))

	assert.NoError(t, err)

	assert.Equal(t, "daemon", sample.name)
	assert.InEpsilon(t, 666.0, sample.value, epsilon)
	require.Nil(t, sample.values)
	assert.Equal(t, gaugeType, sample.metricType)
	assert.Equal(t, []string{"sometag1:somevalue1", "sometag2:somevalue2"}, sample.tags.Sorted())
	assert.InEpsilon(t, 1.0, sample.sampleRate, epsilon)
}

func TestParseGaugeWithNoTags(t *testing.T) {
	sample, err := parseMetricSample([]byte("daemon:666|g"))
	assert.NoError(t, err)

	assert.Equal(t, "daemon", sample.name)
	assert.InEpsilon(t, 666.0, sample.value, epsilon)
	require.Nil(t, sample.values)
	assert.Equal(t, gaugeType, sample.metricType)
	assert.Equal(t, []string{}, sample.tags.Sorted())
	assert.InEpsilon(t, 1.0, sample.sampleRate, epsilon)
}

func TestParseGaugeWithSampleRate(t *testing.T) {
	sample, err := parseMetricSample([]byte("daemon:666|g|@0.21"))

	assert.NoError(t, err)

	assert.Equal(t, "daemon", sample.name)
	assert.InEpsilon(t, 666.0, sample.value, epsilon)
	require.Nil(t, sample.values)
	assert.Equal(t, gaugeType, sample.metricType)
	assert.Equal(t, []string{}, sample.tags.Sorted())
	assert.InEpsilon(t, 0.21, sample.sampleRate, epsilon)
}

func TestParseGaugeWithPoundOnly(t *testing.T) {
	sample, err := parseMetricSample([]byte("daemon:666|g|#"))

	assert.NoError(t, err)

	assert.Equal(t, "daemon", sample.name)
	assert.InEpsilon(t, 666.0, sample.value, epsilon)
	require.Nil(t, sample.values)
	assert.Equal(t, gaugeType, sample.metricType)
	assert.Equal(t, []string{}, sample.tags.Sorted())
	assert.InEpsilon(t, 1.0, sample.sampleRate, epsilon)
}

func TestParseGaugeWithUnicode(t *testing.T) {
	sample, err := parseMetricSample([]byte("♬†øU†øU¥ºuT0♪:666|g|#intitulé:T0µ"))

	assert.NoError(t, err)

	assert.Equal(t, "♬†øU†øU¥ºuT0♪", sample.name)
	assert.InEpsilon(t, 666.0, sample.value, epsilon)
	require.Nil(t, sample.values)
	assert.Equal(t, gaugeType, sample.metricType)
	assert.Equal(t, []string{"intitulé:T0µ"}, sample.tags.Sorted())
	assert.InEpsilon(t, 1.0, sample.sampleRate, epsilon)
}

func TestParseMetricError(t *testing.T) {
	// not enough information
	_, err := parseMetricSample([]byte("daemon:666"))
	assert.Error(t, err)

	_, err = parseMetricSample([]byte("daemon:666|"))
	assert.Error(t, err)

	_, err = parseMetricSample([]byte("daemon:|g"))
	assert.Error(t, err)

	_, err = parseMetricSample([]byte(":666|g"))
	assert.Error(t, err)

	_, err = parseMetricSample([]byte("abc666|g"))
	assert.Error(t, err)

	// unknown metadata prefix
	_, err = parseMetricSample([]byte("daemon:666|g|m:test"))
	assert.NoError(t, err)

	// invalid value
	_, err = parseMetricSample([]byte("daemon:abc|g"))
	assert.Error(t, err)

	// invalid metric type
	_, err = parseMetricSample([]byte("daemon:666|unknown"))
	assert.Error(t, err)

	// invalid sample rate
	_, err = parseMetricSample([]byte("daemon:666|g|@abc"))
	assert.Error(t, err)
}
