/*
 * Copyright (C) 2020 Tycho Softworks.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#ifndef BORDEAUX_IPC_
#define BORDEAUX_IPC_

#ifdef  __cplusplus
#include <atomic>
#include <cstdint>
#include <ctime>
#else
#include <time.h>
#include <stdint.h>
#include <stdbool.h>
#include <stdatomic.h>
#endif
#include <unistd.h>

#define IPC_URI_MAX 256     // NOLINT
#define IPC_VERSION 1       // NOLINT
#define MQD_INVALID ((mqd_t)(-1))   // NOLINT

typedef enum { // NOLINT
    IPC_STATUS = 1,
    IPC_LEVEL,      // verbosity
    IPC_RELOAD,
}   ipc_type_t;

struct ipc_event {
    ipc_type_t type;
    uint8_t ver;
    uint8_t res0;
    union {
        uint32_t level;     // logging level
    } body;
};

struct ipc_session {
    time_t created;
    char status[8];
    char sid[10];
    char caller[32];
    char dialed[32];
    uint32_t lines, steps, audio, video;
    enum {ISC_UNKNOWN = 0, ISC_INVITING, ISC_RINGING, ISC_CONNECTED, ISC_CLOSING} state;
    enum {ISC_ANY = 0, ISC_INBOUND, ISC_OUTBOUND} type;
};

struct ipc_system {
#ifdef  __cplusplus
    std::atomic<uint32_t> used;
#else
    _Atomic uint32_t used;
#endif
    time_t started;
    uint32_t series, level, limit;
    char identity[IPC_URI_MAX];
    bool active;    // registered...
};
#endif
