package printer

import (
	"testing"

	corev1 "k8s.io/api/core/v1"

	"github.com/heptio/developer-dash/internal/conversion"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tolerationTable(descriptions ...string) *component.Table {
	table := component.NewTableWithRows(
		"Taints and Tolerations",
		component.NewTableCols("Description"),
		[]component.TableRow{})

	for _, description := range descriptions {
		table.Add(component.TableRow{"Description": component.NewText(description)})
	}

	return table
}

func Test_TolerationDescriber_Create(t *testing.T) {
	cases := []struct {
		name        string
		tolerations []corev1.Toleration
		expected    *component.Table
		isErr       bool
	}{
		{
			name: "key,value",
			tolerations: []corev1.Toleration{
				{
					Key:   "key",
					Value: "value",
				},
			},
			expected: tolerationTable("Schedule on nodes with key:value taint."),
		},
		{
			name: "multiple tolerations",
			tolerations: []corev1.Toleration{
				{
					Key:   "key1",
					Value: "value1",
				},
				{
					Key:   "key2",
					Value: "value2",
				},
			},
			expected: tolerationTable(
				"Schedule on nodes with key1:value1 taint.",
				"Schedule on nodes with key2:value2 taint.",
			),
		},
		{
			name: "key,value",
			tolerations: []corev1.Toleration{
				{
					Key:    "key",
					Value:  "value",
					Effect: corev1.TaintEffectNoSchedule,
				},
			},
			expected: tolerationTable("Schedule on nodes with key:value:NoSchedule taint."),
		},
		{
			name: "effect",
			tolerations: []corev1.Toleration{
				{
					Effect: corev1.TaintEffectNoSchedule,
				},
			},
			expected: tolerationTable("Schedule on nodes with NoSchedule taint."),
		},
		{
			name: "key,value with evict secs",
			tolerations: []corev1.Toleration{
				{
					Key:               "key",
					Value:             "value",
					TolerationSeconds: conversion.PtrInt64(3600),
				},
			},
			expected: tolerationTable("Schedule on nodes with key:value taint. Evict after 3600 seconds."),
		},
		{
			name: "key exists",
			tolerations: []corev1.Toleration{
				{
					Key:      "key",
					Operator: corev1.TolerationOpExists,
				},
			},
			expected: tolerationTable("Schedule on nodes with key taint."),
		},
		{
			name: "exists with no key",
			tolerations: []corev1.Toleration{
				{
					Operator: corev1.TolerationOpExists,
				},
			},
			expected: tolerationTable("Schedule on all nodes."),
		},
		{
			name: "unsupported toleration",
			tolerations: []corev1.Toleration{
				{
					Key: "key",
				},
			},
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			podSpec := corev1.PodSpec{
				Tolerations: tc.tolerations,
			}

			got, err := printTolerations(podSpec)
			if tc.isErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}
