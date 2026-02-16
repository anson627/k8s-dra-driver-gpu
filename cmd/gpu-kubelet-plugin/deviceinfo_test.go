/*
 * Copyright (c) 2024, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	resourceapi "k8s.io/api/resource/v1"
	"k8s.io/dynamic-resource-allocation/deviceattribute"
	"k8s.io/utils/ptr"
)

func TestGpuInfo_GetDevice_NumaNode(t *testing.T) {
	testCases := []struct {
		description      string
		numaNode         int
		expectNumaAttr   bool
		expectedNumaNode int64
	}{
		{
			description:      "valid NUMA node 0",
			numaNode:         0,
			expectNumaAttr:   true,
			expectedNumaNode: 0,
		},
		{
			description:      "valid NUMA node 1",
			numaNode:         1,
			expectNumaAttr:   true,
			expectedNumaNode: 1,
		},
		{
			description:      "valid NUMA node 3",
			numaNode:         3,
			expectNumaAttr:   true,
			expectedNumaNode: 3,
		},
		{
			description:    "invalid NUMA node -1 (not available)",
			numaNode:       -1,
			expectNumaAttr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			gpuInfo := &GpuInfo{
				UUID:                  "GPU-12345678-1234-1234-1234-123456789abc",
				minor:                 0,
				memoryBytes:           40 * 1024 * 1024 * 1024, // 40GB
				productName:           "NVIDIA A100",
				brand:                 "Tesla",
				architecture:          "Ampere",
				cudaComputeCapability: "8.0",
				driverVersion:         "525.85.12",
				cudaDriverVersion:     "12.0",
				pcieBusID:             "0000:00:1e.0",
				numaNode:              tc.numaNode,
				health:                Healthy,
			}

			device := gpuInfo.GetDevice()

			// Check gpu.nvidia.com/numaNode attribute
			numaAttr, hasNumaAttr := device.Attributes["gpu.nvidia.com/numaNode"]
			require.Equal(t, tc.expectNumaAttr, hasNumaAttr, "gpu.nvidia.com/numaNode attribute presence mismatch")

			if tc.expectNumaAttr {
				require.NotNil(t, numaAttr.IntValue, "gpu.nvidia.com/numaNode IntValue should not be nil")
				require.Equal(t, tc.expectedNumaNode, *numaAttr.IntValue, "gpu.nvidia.com/numaNode value mismatch")
			}
		})
	}
}

func TestGpuInfo_GetDevice_PcieRootAttr(t *testing.T) {
	testCases := []struct {
		description      string
		pcieRootAttr     *deviceattribute.DeviceAttribute
		expectPcieRoot   bool
		expectedPcieRoot string
	}{
		{
			description: "standard pcieRoot available",
			pcieRootAttr: &deviceattribute.DeviceAttribute{
				Name:  resourceapi.QualifiedName(deviceattribute.StandardDeviceAttributePrefix + "pcieRoot"),
				Value: resourceapi.DeviceAttribute{StringValue: ptr.To("pci0000:00")},
			},
			expectPcieRoot:   true,
			expectedPcieRoot: "pci0000:00",
		},
		{
			description: "parent PCI device pcieRoot fallback",
			pcieRootAttr: &deviceattribute.DeviceAttribute{
				Name:  resourceapi.QualifiedName(deviceattribute.StandardDeviceAttributePrefix + "pcieRoot"),
				Value: resourceapi.DeviceAttribute{StringValue: ptr.To("0000:00:01.0")},
			},
			expectPcieRoot:   true,
			expectedPcieRoot: "0000:00:01.0",
		},
		{
			description:    "no pcieRoot available",
			pcieRootAttr:   nil,
			expectPcieRoot: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			gpuInfo := &GpuInfo{
				UUID:                  "GPU-12345678-1234-1234-1234-123456789abc",
				minor:                 0,
				memoryBytes:           40 * 1024 * 1024 * 1024,
				productName:           "NVIDIA A100",
				brand:                 "Tesla",
				architecture:          "Ampere",
				cudaComputeCapability: "8.0",
				driverVersion:         "525.85.12",
				cudaDriverVersion:     "12.0",
				pcieBusID:             "0000:00:1e.0",
				pcieRootAttr:          tc.pcieRootAttr,
				numaNode:              -1,
				health:                Healthy,
			}

			device := gpuInfo.GetDevice()

			pcieRootAttrName := resourceapi.QualifiedName(deviceattribute.StandardDeviceAttributePrefix + "pcieRoot")
			pcieRootAttr, hasPcieRoot := device.Attributes[pcieRootAttrName]
			require.Equal(t, tc.expectPcieRoot, hasPcieRoot, "pcieRoot attribute presence mismatch")

			if tc.expectPcieRoot {
				require.NotNil(t, pcieRootAttr.StringValue, "pcieRoot StringValue should not be nil")
				require.Equal(t, tc.expectedPcieRoot, *pcieRootAttr.StringValue, "pcieRoot value mismatch")
			}
		})
	}
}

func TestGpuInfo_PartDevAttributes_NumaNode(t *testing.T) {
	testCases := []struct {
		description      string
		numaNode         int
		expectNumaAttr   bool
		expectedNumaNode int64
	}{
		{
			description:      "valid NUMA node 0",
			numaNode:         0,
			expectNumaAttr:   true,
			expectedNumaNode: 0,
		},
		{
			description:      "valid NUMA node 2",
			numaNode:         2,
			expectNumaAttr:   true,
			expectedNumaNode: 2,
		},
		{
			description:    "invalid NUMA node -1 (not available)",
			numaNode:       -1,
			expectNumaAttr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			gpuInfo := &GpuInfo{
				UUID:                  "GPU-12345678-1234-1234-1234-123456789abc",
				minor:                 0,
				memoryBytes:           40 * 1024 * 1024 * 1024,
				productName:           "NVIDIA A100",
				brand:                 "Tesla",
				architecture:          "Ampere",
				cudaComputeCapability: "8.0",
				driverVersion:         "525.85.12",
				cudaDriverVersion:     "12.0",
				pcieBusID:             "0000:00:1e.0",
				numaNode:              tc.numaNode,
				health:                Healthy,
			}

			attrs := gpuInfo.PartDevAttributes()

			numaAttr, hasNumaAttr := attrs["gpu.nvidia.com/numaNode"]
			require.Equal(t, tc.expectNumaAttr, hasNumaAttr, "gpu.nvidia.com/numaNode attribute presence mismatch")

			if tc.expectNumaAttr {
				require.NotNil(t, numaAttr.IntValue, "gpu.nvidia.com/numaNode IntValue should not be nil")
				require.Equal(t, tc.expectedNumaNode, *numaAttr.IntValue, "gpu.nvidia.com/numaNode value mismatch")
			}
		})
	}
}

func TestIsPCIBusID(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		// Valid PCI bus IDs
		{input: "0000:00:00.0", expected: true},
		{input: "0001:00:00.0", expected: true},
		{input: "ffff:ff:1f.7", expected: true},
		{input: "ABCD:EF:12.3", expected: true},
		{input: "abcd:ef:12.3", expected: true},
		// Invalid formats
		{input: "", expected: false},
		{input: "0000:00:00", expected: false},    // missing function
		{input: "0000:00:00.00", expected: false}, // too long
		{input: "000:00:00.0", expected: false},   // short domain
		{input: "00000:00:00.0", expected: false}, // long domain
		{input: "0000-00-00.0", expected: false},  // wrong separator
		{input: "0000:00:00:0", expected: false},  // wrong function separator
		{input: "gggg:00:00.0", expected: false},  // invalid hex
		{input: "pci0000:00", expected: false},    // pci prefix (not BDF)
		{input: "numa-0", expected: false},        // NUMA format
		{input: "vmbus-guid", expected: false},    // VMBUS format
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := isPCIBusID(tc.input)
			require.Equal(t, tc.expected, result, "isPCIBusID(%q) should return %v", tc.input, tc.expected)
		})
	}
}

func TestGpuInfo_GetDevice_PcieRootAttr_ParentBridge(t *testing.T) {
	// Test that parent bridge format is accepted as a valid pcieRoot value
	testCases := []struct {
		description      string
		pcieRootAttr     *deviceattribute.DeviceAttribute
		expectedPcieRoot string
	}{
		{
			description: "parent bridge pcieRoot (PCI bus ID format)",
			pcieRootAttr: &deviceattribute.DeviceAttribute{
				Name:  resourceapi.QualifiedName(deviceattribute.StandardDeviceAttributePrefix + "pcieRoot"),
				Value: resourceapi.DeviceAttribute{StringValue: ptr.To("0000:00:01.0")},
			},
			expectedPcieRoot: "0000:00:01.0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			gpuInfo := &GpuInfo{
				UUID:                  "GPU-12345678-1234-1234-1234-123456789abc",
				minor:                 0,
				memoryBytes:           40 * 1024 * 1024 * 1024,
				productName:           "NVIDIA A100",
				brand:                 "Tesla",
				architecture:          "Ampere",
				cudaComputeCapability: "8.0",
				driverVersion:         "525.85.12",
				cudaDriverVersion:     "12.0",
				pcieBusID:             "0001:00:00.0",
				pcieRootAttr:          tc.pcieRootAttr,
				numaNode:              0,
				health:                Healthy,
			}

			device := gpuInfo.GetDevice()

			pcieRootAttrName := resourceapi.QualifiedName(deviceattribute.StandardDeviceAttributePrefix + "pcieRoot")
			pcieRootAttr, hasPcieRoot := device.Attributes[pcieRootAttrName]
			require.True(t, hasPcieRoot, "pcieRoot attribute should be present")
			require.NotNil(t, pcieRootAttr.StringValue, "pcieRoot StringValue should not be nil")
			require.Equal(t, tc.expectedPcieRoot, *pcieRootAttr.StringValue, "pcieRoot value mismatch")
		})
	}
}
