/*
 * MIT-X11 Open Source License
 *
 * Copyright (c) 2022, Advanced Micro Devices, Inc.
 * All rights reserved.
 *
 * Developed by:
 *
 *                 AMD Research and AMD Software Development
 *
 *                 Advanced Micro Devices, Inc.
 *
 *                 www.amd.com
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or
 * sellcopies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 *  - The above copyright notice and this permission notice shall be included in
 *    all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 * Except as contained in this notice, the name of the Advanced Micro Devices,
 * Inc. shall not be used in advertising or otherwise to promote the sale, use
 * or other dealings in this Software without prior written authorization from
 * the Advanced Micro Devices, Inc.
 *
 */

// Package cpustat provides an example parser for Linux CPU utilization statistics.
package collect
import "github.com/MuthusamyRamalingam/go_amd_smi"
//import "github.com/amd/go_amd_smi"

var MAX_CPU_SOCKETS = 4
var MAX_CPU_THREADS = 768
var MAX_GPUS = 24

type AMDParams struct {
	CoreEnergy [768]float64
	SocketEnergy [4]float64
	CoreBoost [768]float64
	SocketPower [4]float64
	PowerLimit [4]float64
	ProchotStatus [4]float64
	Sockets int
	Threads int
	ThreadsPerCore int
	NumGPUs int
	GPUDevId [24]float64
	GPUPowerCap [24]float64
	GPUPower [24]float64
	GPUTemperature [24]float64
	GPUSCLK [24]float64
	GPUMCLK [24]float64
	GPUUsage [24]float64
	GPUMemoryUsage [24]float64
}

func (amdParams *AMDParams) Init() {

	amdParams.Sockets = 0
	amdParams.Threads = 0
	amdParams.ThreadsPerCore = 0
	amdParams.NumGPUs = 0

	for socketLoopCounter := 0; socketLoopCounter < MAX_CPU_SOCKETS; socketLoopCounter++	{ 
		amdParams.SocketEnergy[socketLoopCounter] = -1
		amdParams.SocketPower[socketLoopCounter] = -1
		amdParams.PowerLimit[socketLoopCounter] = -1
		amdParams.ProchotStatus[socketLoopCounter] = -1
	}
	
	for logicalCoreLoopCounter := 0; logicalCoreLoopCounter < MAX_CPU_THREADS; logicalCoreLoopCounter++	{ 
		amdParams.CoreEnergy[logicalCoreLoopCounter] = -1
		amdParams.CoreBoost[logicalCoreLoopCounter] = -1
	}
	
	for gpuLoopCounter := 0; gpuLoopCounter < MAX_GPUS; gpuLoopCounter++	{ 
		amdParams.GPUDevId[gpuLoopCounter] = -1
		amdParams.GPUPowerCap[gpuLoopCounter] = -1
		amdParams.GPUPower[gpuLoopCounter] = -1
		amdParams.GPUTemperature[gpuLoopCounter] = -1
		amdParams.GPUSCLK[gpuLoopCounter] = -1
		amdParams.GPUMCLK[gpuLoopCounter] = -1
		amdParams.GPUUsage[gpuLoopCounter] = -1
		amdParams.GPUMemoryUsage[gpuLoopCounter] = -1
	}
}

func Scan() (AMDParams) {

	var stat AMDParams
	stat.Init()
	
	value64 := int(0)
	value32 := int(0)
	
	value16 := int(0)

	if 1 == goamdsmi.GO_cpu_init() {
		num_sockets := int(goamdsmi.GO_cpu_number_of_sockets_get())
		num_threads := int(goamdsmi.GO_cpu_number_of_threads_get())
		num_threads_per_core := int(goamdsmi.GO_cpu_threads_per_core_get())

		stat.Sockets = int(num_sockets)
		stat.Threads = int(num_threads)
		stat.ThreadsPerCore = int(num_threads_per_core)

		for i := 0; i < num_threads ; i++ {
			value64 = int(goamdsmi.GO_cpu_core_energy_get(i))
			stat.CoreEnergy[i] = float64(value64)
			value64 = 0

			value32 = int(goamdsmi.GO_cpu_core_boostlimit_get(i))
			stat.CoreBoost[i] = float64(value32)
			value32 = 0
		}

		for i := 0; i < num_sockets ; i++ {
			value64 = int(goamdsmi.GO_cpu_socket_energy_get(i))
			stat.SocketEnergy[i] = float64(value64)
			value64 = 0

			value32 = int(goamdsmi.GO_cpu_socket_power_get(i))
			stat.SocketPower[i] = float64(value32)
			value32 = 0

			value32 = int(goamdsmi.GO_cpu_socket_power_cap_get(i))
			stat.PowerLimit[i] = float64(value32)
			value32 = 0

			value32 = int(goamdsmi.GO_cpu_prochot_status_get(i))
			stat.ProchotStatus[i] = float64(value32)
			value32 = 0
		}
	}

	
	if 1 == goamdsmi.GO_gpu_init() {	
		num_gpus := int(goamdsmi.GO_gpu_num_monitor_devices())
		stat.NumGPUs = int(num_gpus)

		for i := 0; i < num_gpus ; i++ {
			value16 = int(goamdsmi.GO_gpu_dev_id_get(i))
			stat.GPUDevId[i] = float64(value16)
			value16 = 0

			value64 = int(goamdsmi.GO_gpu_dev_power_cap_get(i))
			stat.GPUPowerCap[i] = float64(value64)
			value64 = 0

			value64 = int(goamdsmi.GO_gpu_dev_power_get(i))
			stat.GPUPower[i] = float64(value64)
			value64 = 0

			//Get the value for GPU current temperature. Sensor = 0(GPU), Metric = 0(current)
			value64 = int(goamdsmi.GO_gpu_dev_temp_metric_get(i, 0, 0))
			if -1 == value64 {
				//Sensor = 1 (GPU Junction Temp)
				value64 = int(goamdsmi.GO_gpu_dev_temp_metric_get(i, 1, 0))
			}
			stat.GPUTemperature[i] = float64(value64)
			value64 = 0

			value64 = int(goamdsmi.GO_gpu_dev_gpu_clk_freq_get_sclk(i))
			stat.GPUSCLK[i] = float64(value64)
			value64 = 0

			value64 = int(goamdsmi.GO_gpu_dev_gpu_clk_freq_get_mclk(i))
			stat.GPUMCLK[i] = float64(value64)
			value64 = 0

			value64 = int(goamdsmi.GO_gpu_dev_gpu_busy_percent_get(i))
			stat.GPUUsage[i] = float64(value64)
			value64 = 0

			value64 = int(goamdsmi.GO_gpu_dev_gpu_memory_busy_percent_get(i))
			stat.GPUMemoryUsage[i] = float64(value64)
			value64 = 0
		}
	}

	return stat
}

