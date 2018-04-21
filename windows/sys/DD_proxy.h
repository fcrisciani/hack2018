/*++

Copyright (c) Microsoft Corporation. All rights reserved

Abstract:

   This header files declares common data types and function prototypes used
   throughout the driver.

Environment:

    Kernel mode

--*/

#ifndef _DD_PROXY_H_
#define _DD_PROXY_H_


//
// Pooltags used by this callout driver.
//
#define DD_PROXY_FLOW_CONTEXT_POOL_TAG 'olfD'
#define DD_PROXY_CONTROL_DATA_POOL_TAG 'dcdD'

extern HANDLE gInjectionHandle;

extern LIST_ENTRY gFlowList;
extern KSPIN_LOCK gFlowListLock;

extern UINT32 gCalloutIdV4;

extern BOOLEAN gDriverUnloading;

//
// Utility functions
//

typedef struct DD_PROXY_FLOW_INFO_
{
	LIST_ENTRY listEntry;

	UINT32 ipv4LocalAddr;
	UINT16 portLocal;
	UINT32 ipv4RemoteAddr;
	UINT16 portRemote;
	UINT8 protocol;

	UINT64 flowId;
} DD_PROXY_FLOW_INFO;


void
DDProxyFlowEstablishedClassify(
   _In_ const FWPS_INCOMING_VALUES* inFixedValues,
   _In_ const FWPS_INCOMING_METADATA_VALUES* inMetaValues,
   _Inout_opt_ void* layerData,
   _In_opt_ const void* classifyContext,
   _In_ const FWPS_FILTER* filter,
   _In_ UINT64 flowContext,
   _Inout_ FWPS_CLASSIFY_OUT* classifyOut
   );

NTSTATUS
DDProxyFlowEstablishedNotify(
   _In_ FWPS_CALLOUT_NOTIFY_TYPE notifyType,
   _In_ const GUID* filterKey,
   _Inout_ const FWPS_FILTER* filter
   );

#endif // _DD_PROXY_H_
