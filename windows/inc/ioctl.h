/*++

Copyright (c) Microsoft Corporation. All rights reserved

Abstract:

    ddproxy IOCTL header

Environment:

    Kernel mode
    
--*/

#pragma once

#define DDPROXY_DEVICE_NAME     L"\\Device\\DDProxy"
#define DDPROXY_SYMBOLIC_NAME   L"\\DosDevices\\Global\\DDProxy"
#define DDPROXY_DOS_NAME   L"\\\\.\\DDProxy"

typedef struct _DDPROXY_FLOW
{
	UINT32  ipv4Remote;
	UINT16  portRemote;
	UINT32  ipv4Local;
	UINT16  portLocal;
	UINT32  protocol;
} DDPROXY_FLOW;

typedef struct _DDPROXY_FLOWS
{
	UINT32                  NumberOfFlows;
	DDPROXY_FLOW            Flow[1];
} DDPROXY_FLOWS;

#define	DDPROXY_IOCTL_GET_CONNECTIONS  CTL_CODE(FILE_DEVICE_NETWORK, 0x1, METHOD_BUFFERED, FILE_ANY_ACCESS)

