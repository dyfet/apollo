// Copyright (C) 2023 Tycho Softworks.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package internal

/*
#cgo LDFLAGS: -L/usr/pkg/lib -L/usr/local/lib -lrt
#include "ipc.h"
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"os"
	"unsafe"

	"gitlab.com/tychosoft/service"
)

type IpcInfo struct {
	// common data
	IPCPath string `json:"ipc_path"`

	// coventry ipc sizes
	MsgSize   uintptr `json:"msg_size,omitempty"`
	SysSize   uintptr `json:"sys_size,omitempty"`
	RegSize   uintptr `json:"reg_size,omitempty"`
	CallSize  uintptr `json:"call_size,omitempty"`
	RegCount  uintptr `json:"reg_count,omitempty"`
	CallCount uintptr `json:"call_count,omitempty"`

	// bordeaux ipc sizes
	EventSize    uintptr `json:"event_size,omitempty"`
	SystemSize   uintptr `json:"system_size,omitempty"`
	SessionSize  uintptr `json:"session_size,omitempty"`
	SessionCount uintptr `json:"session_count,omitempty"`
}

var (
	ipcCoventry string
	ipcRegistry uintptr
	registryMap *C.pbx_reg_t = nil
)

func VerifyToken(token string) int {
	if registryMap == nil {
		return 0
	}

	cs_token := C.CString(token)
	defer C.free(unsafe.Pointer(cs_token))
	return int(C.verify_user(registryMap, cs_token))
}

func ReloadCoventry() error {
	cs_path := C.CString(ipcCoventry)
	defer C.free(unsafe.Pointer(cs_path))
	result := int(C.reload_coventry(cs_path))
	if result < 0 {
		return fmt.Errorf("mqueue error %d", result)
	}
	return nil
}

func ipcInit(coventry string) {
	var ipc IpcInfo
	data, err := os.ReadFile(coventry)
	if err != nil {
		service.Fail(90, err)
	}
	err = json.Unmarshal(data, &ipc)
	if err != nil {
		service.Fail(90, err)
	}

	if ipc.MsgSize != unsafe.Sizeof(C.pbx_msg_t{}) ||
		ipc.SysSize != unsafe.Sizeof(C.pbx_sys_t{}) ||
		ipc.RegSize != unsafe.Sizeof(C.pbx_reg_t{}) ||
		ipc.CallSize != unsafe.Sizeof(C.pbx_call_t{}) {
		service.Fail(91, fmt.Errorf("IPC size mismatch"))
	}

	ipcCoventry = ipc.IPCPath
	ipcRegistry = (ipc.RegSize * ipc.RegCount) + ipc.SysSize

	if registryMap != nil {
		C.munmap(unsafe.Pointer(registryMap), C.size_t(ipcRegistry))
		registryMap = nil
	}

	reg_path := C.CString(ipcCoventry + ".registry")
	shm := C.shm_open(reg_path, C.O_RDONLY, 0660)
	defer C.free(unsafe.Pointer(reg_path))
	if shm < C.int(0) {
		service.Warn("shared registry missing")
		return
	}

	registryMap = C.registry_map(C.size_t(ipcRegistry), shm)
	C.close(shm)
	if unsafe.Pointer(registryMap) == C.MAP_FAILED || registryMap == nil {
		registryMap = nil
		service.Warn("shared registry broken")
	}
}

func getRegistry(id int, line *Line) {
	if registryMap == nil || id < 10 || id > 89 {
		return
	}

	line.Agent = C.GoString(C.registry_agent(C.int(id), registryMap))
	line.Count = uint16(C.registry_count(C.int(id), registryMap))
	line.Presence = C.GoString(C.registry_presence(C.int(id), registryMap))
	if len(line.Agent) < 1 {
		line.Agent = "offline"
	}

	cs_host := C.registry_host(C.int(id), registryMap)
	defer C.free(unsafe.Pointer(cs_host))
	line.Host = C.GoString(cs_host)
	if line.Host == "unknown" {
		line.Agent = "offline"
	}

	agentInfo(line)
}
