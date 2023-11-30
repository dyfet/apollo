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

#ifndef COVENTRY_IPC_
#define COVENTRY_IPC_

#ifdef  __cplusplus
#include <cstdint>
#include <ctime>
#include <atomic>
#else
#include <time.h>
#include <stdint.h>
#include <stdbool.h>
#include <stdatomic.h>
#endif

#include <unistd.h>
#include <sys/socket.h>

// PBX IPC protocol version...
#define PBX_VERSION 1   // NOLINT
#define MQD_INVALID ((mqd_t)(-1))   // NOLINT

/* Known PBX requests */
typedef enum {  // NOLINT
    PBX_MESSAGE = 1,
    PBX_LEVEL,
    PBX_RELOAD,
    PBX_ACTIVE,
    PBX_ASSIGN,
    PBX_CLEARS,
    PBX_PRESENCE,
} pbx_type_t;

/* PBX mqueue messages */
struct pbx_msg {
    pbx_type_t type;
    uint8_t ver;
    uint8_t res0;
    union {
        uint8_t level;
        struct {
            char to[32 + 1];
            char text[160 + 1];
            char subj[64 + 1];
        } message;
        struct {
            uint32_t ext;
            char key[32];
            char value[96];
        } registry;
    } body;
};

/* Server registry for each extension */
struct pbx_reg {
#ifdef  __cplusplus
    std::atomic<uint32_t> count;
#else
    _Atomic uint32_t count;
#endif
    time_t activated, expires;       /* 0 if inactive */
    struct sockaddr_storage address;
    char agent[32];
    char name[32];
    char id[40];
    char token[40];
    uint16_t lines;
    volatile enum {GONE = 0, HERE, AWAY, DND} presence;
    struct {
        bool invite : 1;
        bool message : 1;
    } flags;
};

/* Server system info */
struct pbx_sys {
    time_t started;
    pid_t pid;                      /* server pid */
    uint32_t series;                /* reload series */
    uint8_t version;
    uint8_t level;
    uint16_t first;
    char state[16];
    char realm[64];
    char digest[16];
    struct pbx_reg trunk;
};

/* Each server call; mapped to call capacity */
struct pbx_call {
    uint64_t id;                    /* 0 if inactive */
    time_t created;
    uint16_t active, ringing, segment;
    char caller[64], dialed[64], remote[128];
    enum {INVITING, RINGING, CONNECTED, CLOSING} state;
    enum {LOCAL, INBOUND, OUTGOING, NONE} type;
};
#endif
