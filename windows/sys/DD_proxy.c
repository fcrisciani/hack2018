/*++

Copyright (c) Microsoft Corporation. All rights reserved

Abstract:

   This file implements the classifyFn callout functions for the flow-established callouts.

Environment:

    Kernel mode

--*/

#include <ntddk.h>

#pragma warning(push)
#pragma warning(disable:4201)       // unnamed struct/union

#include <fwpsk.h>

#pragma warning(pop)

#include <fwpmk.h>

#include "DD_proxy.h"


void
DDProxyFlowEstablishedClassify(
   _In_ const FWPS_INCOMING_VALUES* inFixedValues,
   _In_ const FWPS_INCOMING_METADATA_VALUES* inMetaValues,
   _Inout_opt_ void* layerData,
   _In_opt_ const void* classifyContext,
   _In_ const FWPS_FILTER* filter,
   _In_ UINT64 flowContext,
   _Inout_ FWPS_CLASSIFY_OUT* classifyOut
   )
{
   NTSTATUS status = STATUS_SUCCESS;

   BOOLEAN locked = FALSE;

   KLOCK_QUEUE_HANDLE flowListLockHandle;

   DD_PROXY_FLOW_INFO* flow = NULL;

   UNREFERENCED_PARAMETER(layerData);
   UNREFERENCED_PARAMETER(classifyContext);
   UNREFERENCED_PARAMETER(flowContext);
   UNREFERENCED_PARAMETER(filter);

   flow = ExAllocatePoolWithTag(
                        NonPagedPool,
                        sizeof(DD_PROXY_FLOW_INFO),
                        DD_PROXY_FLOW_CONTEXT_POOL_TAG
                        );

   if (flow == NULL)
   {
      status = STATUS_NO_MEMORY;
      goto Exit;
   }

   RtlZeroMemory(flow, sizeof(DD_PROXY_FLOW_INFO));

   flow->ipv4LocalAddr =
	   RtlUlongByteSwap(
		   inFixedValues->incomingValue\
		   [FWPS_FIELD_ALE_FLOW_ESTABLISHED_V4_IP_LOCAL_ADDRESS].value.uint32
	   );

   flow->portLocal =
	       inFixedValues->incomingValue\
		   [FWPS_FIELD_ALE_FLOW_ESTABLISHED_V4_IP_LOCAL_PORT].value.uint16
	   ;

   flow->ipv4RemoteAddr =
	   RtlUlongByteSwap(
		   inFixedValues->incomingValue\
		   [FWPS_FIELD_ALE_FLOW_ESTABLISHED_V4_IP_REMOTE_ADDRESS].value.uint32
	   );

   flow->portRemote =
		   inFixedValues->incomingValue\
		   [FWPS_FIELD_ALE_FLOW_ESTABLISHED_V4_IP_REMOTE_PORT].value.uint16
	   ;

   flow->protocol = inFixedValues->incomingValue\
		   [FWPS_FIELD_ALE_FLOW_ESTABLISHED_V4_IP_PROTOCOL].value.uint8;

   NT_ASSERT(FWPS_IS_METADATA_FIELD_PRESENT(inMetaValues, 
                                         FWPS_METADATA_FIELD_FLOW_HANDLE));
   flow->flowId = inMetaValues->flowHandle;

   KeAcquireInStackQueuedSpinLock(
      &gFlowListLock,
      &flowListLockHandle
      );

   locked = TRUE;

   if (!gDriverUnloading)
   {
      InsertHeadList(&gFlowList, &flow->listEntry);
   }
   flow = NULL;

   classifyOut->actionType = FWP_ACTION_PERMIT;

   if (filter->flags & FWPS_FILTER_FLAG_CLEAR_ACTION_RIGHT)
   {
	   classifyOut->rights &= ~FWPS_RIGHT_ACTION_WRITE;
   }
Exit:

   if(locked)
   {
      KeReleaseInStackQueuedSpinLock(&flowListLockHandle);
   }

   if (flow != NULL)
   {
      ExFreePoolWithTag(flow, DD_PROXY_FLOW_CONTEXT_POOL_TAG);
   }

   return;
}


NTSTATUS
DDProxyFlowEstablishedNotify(
   _In_ FWPS_CALLOUT_NOTIFY_TYPE notifyType,
   _In_ const GUID* filterKey,
   _Inout_ const FWPS_FILTER* filter
   )
{
   UNREFERENCED_PARAMETER(notifyType);
   UNREFERENCED_PARAMETER(filterKey);
   UNREFERENCED_PARAMETER(filter);

   return STATUS_SUCCESS;
}


