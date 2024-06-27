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

var MAX_CPU_SOCKET = 4
var MAX_CPU_THREADS = 768
var MAX_GPU = 24

type AMDParams struct {
	CoreEnergy [MAX_CPU_THREADS]float64
	SocketEnergy [MAX_CPU_SOCKET]float64
	CoreBoost [MAX_CPU_THREADS]float64
	SocketPower [MAX_CPU_SOCKET]float64
	PowerLimit [MAX_CPU_SOCKET]float64
	ProchotStatus [MAX_CPU_SOCKET]float64
	Sockets uint
	Threads uint
	ThreadsPerCore uint
	NumGPUs uint
	GPUDevId [MAX_GPU]float64
	GPUPowerCap [MAX_GPU]float64
	GPUPowerAvg [MAX_GPU]float64
	GPUTemperature [MAX_GPU]float64
	GPUSCLK [MAX_GPU]float64
	GPUMCLK [MAX_GPU]float64
	GPUUsage [MAX_GPU]float64
	GPUMemoryUsage [MAX_GPU]float64
}

func (amdParams *AMDParams) Init() {

	amdParams.Sockets = -1
	amdParams.Threads = -1
	amdParams.ThreadsPerCore = -1
	amdParams.NumGPUs = -1

	for socketloopCounter := 0; socketloopCounter < MAX_CPU_SOCKET; socketloopCounter++	{ 
		amdParams.SocketEnergy = -1
		amdParams.SocketPower = -1
		amdParams.PowerLimit = -1
		amdParams.ProchotStatus = -1
	}
	
	for logicalCoreloopCounter := 0; logicalCoreloopCounter < MAX_CPU_THREADS; logicalCoreloopCounter++	{ 
		amdParams.CoreEnergy = -1
		amdParams.CoreBoost = -1
	}
	
	for gpuloopCounter := 0; gpuloopCounter < MAX_GPU; gpuloopCounter++	{ 
		amdParams.GPUDevId = -1
		amdParams.GPUPowerCap = -1
		amdParams.GPUPowerAvg = -1
		amdParams.GPUTemperature = -1
		amdParams.GPUSCLK = -1
		amdParams.GPUMCLK = -1
		amdParams.GPUUsage = -1
		amdParams.GPUMemoryUsage = -1
	}
}

func Scan() (AMDParams) {

	var stat AMDParams
	stat.Init()
	
	value64 := int(0)
	value32 := int(0)
	
	value16 := uint16(0)

	if 1 == goamdsmi.GO_cpu_init() {
		num_sockets := int(goamdsmi.GO_cpu_number_of_sockets_get())
		num_threads := int(goamdsmi.GO_cpu_number_of_threads_get())
		num_threads_per_core := int(goamdsmi.GO_cpu_threads_per_core_get())

		stat.Sockets = uint(num_sockets)
		stat.Threads = uint(num_threads)
		stat.ThreadsPerCore = uint(num_threads_per_core)

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
		stat.NumGPUs = uint(num_gpus)

		for i := 0; i < num_gpus ; i++ {
			value16 = int(goamdsmi.GO_gpu_dev_id_get(i))
			stat.GPUDevId[i] = float64(value16)
			value16 = 0

			value64 = int(goamdsmi.GO_gpu_dev_power_cap_get(i))
			stat.GPUPowerCap[i] = float64(value64)
			value64 = 0

			value64 = int(goamdsmi.GO_gpu_dev_power_ave_get(i))
			stat.GPUPowerAvg[i] = float64(value64)
			value64 = 0

			//Get the value for GPU current temperature. Sensor = 0(GPU), Metric = 0(current)
			value64 = int(goamdsmi.GO_gpu_dev_temp_metric_get(i, 0, 0))
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

