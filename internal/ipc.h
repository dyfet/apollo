#include <mqueue.h>
#include <sys/mman.h>
#include <fcntl.h>
#include <stdlib.h>
#include <stddef.h>
#include <string.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <netinet/in.h>

#include "coventry.h"
typedef struct pbx_reg pbx_reg_t;
typedef struct pbx_msg pbx_msg_t;
typedef	struct pbx_sys pbx_sys_t;
typedef struct pbx_call pbx_call_t;

#include "bordeaux.h"
typedef struct ipc_session ipc_session_t;
typedef struct ipc_message ipc_message_t;
typedef struct ipc_system ipc_system_t;
typedef struct ipc_event ipc_event_t;

// utility function for registry access to sys
pbx_sys_t *registry_sys(pbx_reg_t *map, size_t count) {
    return (pbx_sys_t*)(&map[count]);
}

char *registry_agent(int id, pbx_reg_t *map) {
    pbx_reg_t *entry = &map[id - 10];
    return entry->agent;
}

char *registry_presence(int id, pbx_reg_t *map) {
    pbx_reg_t *entry = &map[id - 10];
    if(entry->count >= entry->lines)
        return "busy";

    if(entry->count)
        return "call";

    switch(entry->presence) {
    case GONE:
        return "gone";
    case DND:
        return "dnd";
    case AWAY:
        return "away";
    default:
        return "here";
    }
}

int registry_count(int id, pbx_reg_t *map) {
    pbx_reg_t *entry = &map[id - 10];
    return atomic_load(&entry->count);
}

char *registry_host(int id, pbx_reg_t *map) {
    pbx_reg_t *entry = &map[id - 10];
    char host[128];
    memset(host, 0, sizeof(host));
    struct sockaddr *addr = (struct sockaddr *)&entry->address;
    switch(addr->sa_family) {
    case AF_INET:
        getnameinfo(addr, sizeof(struct sockaddr_in), host, sizeof(host), NULL, 0, NI_NUMERICHOST);
        break;
    case AF_INET6:
        getnameinfo(addr, sizeof(struct sockaddr_in6), host, sizeof(host), NULL, 0, NI_NUMERICHOST);
        break;
    default:
        return strdup("unknown");
    }
    return strdup(host);
}

// verify user token in mapped registry, form id:...
int verify_user(pbx_reg_t *map, const char *token) {
    int uid = ((token[0] - '0') * 10) + (token[1] - '0');
    if (uid < 10 || uid > 89 || !map)
        return 0;

    pbx_reg_t *entry = &map[uid - 10];
    if (strcmp(token, entry->token) != 0)
        return 0;

    return uid;
}

// reload coventry
int reload_coventry(const char *path) {
    int status = 0;
    pbx_msg_t msg;

    msg.type = PBX_RELOAD;
    msg.ver = PBX_VERSION;
    mqd_t mq = mq_open(path, O_WRONLY, 0660, NULL);
    if (mq == MQD_INVALID)
        return -1;

    status = mq_send(mq, (const char *)&msg, sizeof(msg), 0);
    mq_close(mq);
    return status;
}

pbx_reg_t *registry_map(size_t size, int shm) {
    return mmap(NULL, size, PROT_READ, MAP_SHARED, shm, 0);
}
