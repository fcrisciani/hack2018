/*++

Copyright (c) Microsoft Corporation. All rights reserved

Abstract:

    Monitor executable

Environment:

    User mode
    
--*/

#include "windows.h"
#include "winioctl.h"
#include "strsafe.h"

#ifndef _CTYPE_DISABLE_MACROS
#define _CTYPE_DISABLE_MACROS
#endif

#include "ws2def.h"

#include <stdio.h>
#include <stdlib.h>

#include "ioctl.h"

VOID
PrintIP(UINT32 ip)
{
	UCHAR bytes[4];
	bytes[3] = ip & 0xFF;
	bytes[2] = (ip >> 8) & 0xFF;
	bytes[1] = (ip >> 16) & 0xFF;
	bytes[0] = (ip >> 24) & 0xFF;
	printf("%u.%u.%u.%u", bytes[3], bytes[2], bytes[1], bytes[0]);
}

int
PrintFlow(DDPROXY_FLOW* Flow)
{
	SYSTEMTIME st;

	if (Flow->protocol == 0) {
		return 0;
	}

	GetSystemTime(&st);

	printf("{ \"@timestamp\" : \"%d-%02d-%02dT%02d:%02d:%02d\", ", st.wYear, st.wMonth, st.wDay, st.wHour, st.wMinute, st.wSecond);
	printf("\"ct.event\" : 1 ,");
	printf("\"src_ip\" : \"");
	PrintIP(Flow->ipv4Local);
	printf("\", \"dest_ip\" : \"");
	PrintIP(Flow->ipv4Remote);
	printf("\", \"orig.l4.sport\": %u, \"orig.l4.dport\": %u, \"orig.ip.protocol\": %u }", 
		Flow->portLocal, Flow->portRemote, Flow->protocol);
	return 1;
}

DWORD
StartMonitoring()
{
	UINT32 flowCount = 20;
	DWORD bytesReturned;
	int BufferSz = FIELD_OFFSET(DDPROXY_FLOWS, Flow[flowCount]);
	DDPROXY_FLOWS* flows = (DDPROXY_FLOWS*) malloc(BufferSz);
	DDPROXY_FLOWS* inflows = (DDPROXY_FLOWS*) malloc(BufferSz);
	UINT32 flowIndex = 0;
	flows->NumberOfFlows = flowCount;
	inflows->NumberOfFlows = flowCount;
	DDPROXY_FLOW* flow = NULL;
	int p = 0;


    HANDLE h = CreateFileW(DDPROXY_DOS_NAME,
                                 GENERIC_READ, 
                                 FILE_SHARE_READ, 
                                 NULL, 
                                 OPEN_EXISTING, 
                                 0, 
                                 NULL);

    if (h == INVALID_HANDLE_VALUE)
    {
		printf("could not get handle\n");
        return GetLastError();
    }

	for (flowIndex = 0; flowIndex < flowCount; flowIndex++) {
		flow = &flows->Flow[flowIndex];
		flow->ipv4Local = 0;
		flow->ipv4Remote = 0;
		flow->portLocal = 0;
		flow->portRemote = 0;
		flow->protocol = 0;
	}

	if (!DeviceIoControl(h, DDPROXY_IOCTL_GET_CONNECTIONS,
		inflows, BufferSz, flows, BufferSz, &bytesReturned, NULL))
	{
		printf("ioctl failed\n");
		return GetLastError();
	}
	// printf("flow count: %d\n", flows->NumberOfFlows);
	// printf("[\n");
	for (flowIndex = 0; flowIndex < flowCount - 1; flowIndex++) {
		flow = &flows->Flow[flowIndex];
		p = PrintFlow(flow);
		if (p == 1) {
			printf("\n");
		}
	}
	flow = &flows->Flow[flowIndex];
	PrintFlow(flow);
	//printf("\n]\n");

	CloseHandle(h);
    return NO_ERROR;
}

BOOL MonitorAppCloseMonitorDevice(
   _In_ HANDLE monitorDevice)
{
    return CloseHandle(monitorDevice);
}


int __cdecl wmain(_In_ int argc, _In_reads_(argc) PCWSTR argv[])
{
   UNREFERENCED_PARAMETER(argc);
   UNREFERENCED_PARAMETER(argv);
   DWORD err = NO_ERROR;
   //printf("Starting to Monitor ..\n");

   err = StartMonitoring();
   //printf("err: %d \n", err);
   if (err != NO_ERROR) {
	   printf("Failed to monitor. Error: %d \n", err);
	   goto Exit;
   }

   //printf("Monitoring started\n");
Exit:
   return (int) err;
}
