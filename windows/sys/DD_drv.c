/*++

Copyright (c) Microsoft Corporation. All rights reserved

Abstract:

   Connection tracking filter.

Environment:

    Kernel mode

--*/

#include <ntddk.h>
#include <wdf.h>

#pragma warning(push)
#pragma warning(disable:4201)

#include <fwpsk.h>

#pragma warning(pop)

#include <fwpmk.h>

#include <ws2ipdef.h>
#include <in6addr.h>
#include <ip2string.h>

#include "DD_proxy.h"
#include "ioctl.h"

#define INITGUID
#include <guiddef.h>

// Callout and sublayer GUIDs

// ee93719d-ad5d-48c9-ae46-7270367d205d
DEFINE_GUID(
    DD_PROXY_FLOW_ESTABLISHED_CALLOUT_V4,
    0xee93719d,
    0xad5d,
    0x48c9,
    0xae, 0x46, 0x72, 0x70, 0x36, 0x7d, 0x20, 0x5d
);


// 0104fd7e-c825-414e-94c9-f0d525bbc169
DEFINE_GUID(
    DD_PROXY_SUBLAYER,
    0x0104fd7e,
    0xc825,
    0x414e,
    0x94, 0xc9, 0xf0, 0xd5, 0x25, 0xbb, 0xc1, 0x69
);


// Callout driver global variables

DEVICE_OBJECT* gWdmDevice;

HANDLE gEngineHandle;
UINT32 gFlowEstablishedCalloutIdV4, gCalloutIdV4;
UINT32 gFlowEstablishedCalloutIdV6, gCalloutIdV6;

HANDLE gInjectionHandle;

LIST_ENTRY gFlowList;
KSPIN_LOCK gFlowListLock;

LIST_ENTRY gPacketQueue;
KSPIN_LOCK gPacketQueueLock;
KEVENT gPacketQueueEvent;

BOOLEAN gDriverUnloading = FALSE;
void* gThreadObj;

DRIVER_INITIALIZE DriverEntry;
EVT_WDF_DRIVER_UNLOAD EvtDriverUnload;


// Callout driver implementation

NTSTATUS
DDProxyRegisterFlowEstablishedCallouts(
   _In_ const GUID* layerKey,
   _In_ const GUID* calloutKey,
   _Inout_ void* deviceObject,
   _Out_ UINT32* calloutId
   )
{
   NTSTATUS status = STATUS_SUCCESS;

   FWPS_CALLOUT sCallout = {0};
   FWPM_CALLOUT mCallout = {0};
   FWPM_FILTER filter =    {0};

   FWPM_DISPLAY_DATA displayData = {0};
   BOOLEAN calloutRegistered = FALSE;

   sCallout.calloutKey = *calloutKey;
   sCallout.classifyFn = DDProxyFlowEstablishedClassify;
   sCallout.notifyFn = DDProxyFlowEstablishedNotify;

   status = FwpsCalloutRegister(
               deviceObject,
               &sCallout,
               calloutId
               );
   if (!NT_SUCCESS(status))
   {
      goto Exit;
   }
   calloutRegistered = TRUE;

   displayData.name = L"Datagram-Data Proxy Flow-Established Callout";
   displayData.description = L"Intercepts flow creations";

   mCallout.calloutKey = *calloutKey;
   mCallout.displayData = displayData;
   mCallout.applicableLayer = *layerKey;

   status = FwpmCalloutAdd(
	   gEngineHandle,
	   &mCallout,
	   NULL,
	   NULL
   );

   filter.layerKey = *layerKey;
   filter.displayData.name = L"Flow-Established Filter";
   filter.displayData.description = L"Filter to record flows";

   filter.action.type = FWP_ACTION_CALLOUT_TERMINATING;
   filter.action.calloutKey = *calloutKey;
   filter.numFilterConditions = 0;
   filter.subLayerKey = DD_PROXY_SUBLAYER;
   filter.weight.type = FWP_EMPTY; // auto-weight.

   status = FwpmFilterAdd(
	   gEngineHandle,
	   &filter,
	   NULL,
	   NULL);

   if (!NT_SUCCESS(status))
   {
	   goto Exit;
   }


Exit:

   if (!NT_SUCCESS(status))
   {
      if (calloutRegistered)
      {
         FwpsCalloutUnregisterById(*calloutId);
         *calloutId = 0;
      }
   }

   return status;
}


NTSTATUS
DDProxyRegisterCallouts(
   _Inout_ void* deviceObject
   )
{
   NTSTATUS status = STATUS_SUCCESS;
   FWPM_SUBLAYER DDProxySubLayer;

   BOOLEAN engineOpened = FALSE;
   BOOLEAN inTransaction = FALSE;

   FWPM_SESSION session = {0};

   session.flags = FWPM_SESSION_FLAG_DYNAMIC;

   status = FwpmEngineOpen(
                NULL,
                RPC_C_AUTHN_WINNT,
                NULL,
                &session,
                &gEngineHandle
                );
   if (!NT_SUCCESS(status))
   {
      goto Exit;
   }
   engineOpened = TRUE;

   status = FwpmTransactionBegin(gEngineHandle, 0);
   if (!NT_SUCCESS(status))
   {
      goto Exit;
   }
   inTransaction = TRUE;

   RtlZeroMemory(&DDProxySubLayer, sizeof(FWPM_SUBLAYER)); 

   DDProxySubLayer.subLayerKey = DD_PROXY_SUBLAYER;
   DDProxySubLayer.displayData.name = L"Datagram-Data Proxy Sub-Layer";
   DDProxySubLayer.displayData.description = 
      L"Sub-Layer for use by Datagram-Data Proxy callouts";
   DDProxySubLayer.flags = 0;
   DDProxySubLayer.weight = FWP_EMPTY; // auto-weight.;

   status = FwpmSubLayerAdd(gEngineHandle, &DDProxySubLayer, NULL);
   if (!NT_SUCCESS(status))
   {
      goto Exit;
   }

   status = DDProxyRegisterFlowEstablishedCallouts(
               &FWPM_LAYER_ALE_FLOW_ESTABLISHED_V4,
               &DD_PROXY_FLOW_ESTABLISHED_CALLOUT_V4,
               deviceObject,
               &gFlowEstablishedCalloutIdV4
               );
   if (!NT_SUCCESS(status))
   {
      goto Exit;
   }

   status = FwpmTransactionCommit(gEngineHandle);
   if (!NT_SUCCESS(status))
   {
      goto Exit;
   }
   inTransaction = FALSE;

Exit:

   if (!NT_SUCCESS(status))
   {
      if (inTransaction)
      {
         FwpmTransactionAbort(gEngineHandle);
         _Analysis_assume_lock_not_held_(gEngineHandle); // Potential leak if "FwpmTransactionAbort" fails
      }
      if (engineOpened)
      {
         FwpmEngineClose(gEngineHandle);
         gEngineHandle = NULL;
      }
   }

   return status;
}

void
DDProxyUnregisterCallouts(void)
{
   FwpmEngineClose(gEngineHandle);
   gEngineHandle = NULL;

   FwpsCalloutUnregisterById(gFlowEstablishedCalloutIdV4);
}

_Function_class_(EVT_WDF_DRIVER_UNLOAD)
_IRQL_requires_same_
_IRQL_requires_max_(PASSIVE_LEVEL)
void
EvtDriverUnload(
   _In_ WDFDRIVER driverObject
   )
{
   KLOCK_QUEUE_HANDLE flowListLockHandle;

   UNREFERENCED_PARAMETER(driverObject);

   KeAcquireInStackQueuedSpinLock(
      &gFlowListLock,
      &flowListLockHandle
      );

   gDriverUnloading = TRUE;

   while (!IsListEmpty(&gFlowList))
   {
	   PLIST_ENTRY Entry = RemoveHeadList(&gFlowList);
	   DD_PROXY_FLOW_INFO* flow = CONTAINING_RECORD(Entry, DD_PROXY_FLOW_INFO, listEntry);
	   ExFreePoolWithTag(flow, DD_PROXY_FLOW_CONTEXT_POOL_TAG);
   }

   KeReleaseInStackQueuedSpinLock(&flowListLockHandle);

   DDProxyUnregisterCallouts();

   FwpsInjectionHandleDestroy(gInjectionHandle);
}


VOID
DDProxyDeviceControl(
	_In_ WDFQUEUE Queue,
	_In_ WDFREQUEST Request,
	_In_ size_t OutputBufferLength,
	_In_ size_t InputBufferLength,
	_In_ ULONG IoControlCode
)
{
	KLOCK_QUEUE_HANDLE flowListLockHandle;
	NTSTATUS status = STATUS_SUCCESS;
	BOOLEAN locked = FALSE;
	UINT32 count = 0;
	DDPROXY_FLOWS* flows = NULL;
	UINT32 sz = 0;

	UNREFERENCED_PARAMETER(Queue);
	UNREFERENCED_PARAMETER(OutputBufferLength);

	switch (IoControlCode)
	{
	case DDPROXY_IOCTL_GET_CONNECTIONS:
	{
		void* pBuffer;

		if (InputBufferLength < sizeof(DDPROXY_FLOWS))
		{
			status = STATUS_INVALID_PARAMETER;
		}
		else
		{
			KeAcquireInStackQueuedSpinLock(
				&gFlowListLock,
				&flowListLockHandle
			);
			locked = TRUE;
			status = WdfRequestRetrieveOutputBuffer(Request, sizeof(DDPROXY_FLOWS), &pBuffer, NULL);
			if (NT_SUCCESS(status))
			{
				flows = (DDPROXY_FLOWS*) pBuffer;
				while (!IsListEmpty(&gFlowList) && (count < 10))
				{
					PLIST_ENTRY Entry = RemoveHeadList(&gFlowList);
					DD_PROXY_FLOW_INFO* flow = CONTAINING_RECORD(Entry, DD_PROXY_FLOW_INFO, listEntry);
					flows->Flow[count].ipv4Remote = flow->ipv4RemoteAddr;
					flows->Flow[count].ipv4Local  = flow->ipv4LocalAddr;
					flows->Flow[count].portLocal  = flow->portLocal;
					flows->Flow[count].portRemote = flow->portRemote;
					flows->Flow[count].protocol   = flow->protocol;
					ExFreePoolWithTag(flow, DD_PROXY_FLOW_CONTEXT_POOL_TAG);
					count++;
				}
				flows->NumberOfFlows = count;
				status = STATUS_SUCCESS;
				sz = FIELD_OFFSET(DDPROXY_FLOWS, Flow[flows->NumberOfFlows]);
			}

			
		}
		break;
	}

	default:
	{
		status = STATUS_INVALID_PARAMETER;
	}
	}

	WdfRequestCompleteWithInformation(Request, status, sz);

	if (locked)
	{
		KeReleaseInStackQueuedSpinLock(&flowListLockHandle);
	}
}

NTSTATUS
DDProxyCtlDriverInit(
	_In_ WDFDEVICE* pDevice
)
{
	NTSTATUS status;
	WDF_IO_QUEUE_CONFIG queueConfig;

	WDF_IO_QUEUE_CONFIG_INIT_DEFAULT_QUEUE(&queueConfig, WdfIoQueueDispatchSequential);
	queueConfig.EvtIoDeviceControl = DDProxyDeviceControl;
	status = WdfIoQueueCreate(*pDevice, &queueConfig, WDF_NO_OBJECT_ATTRIBUTES, NULL);

	return status;
}


//
// Create the minimal WDF Driver and Device objects required for a WFP callout
// driver.
//
NTSTATUS
DDProxyInitDriverObjects(
   _Inout_ DRIVER_OBJECT* driverObject,
   _In_ const UNICODE_STRING* registryPath,
   _Out_ WDFDRIVER* pDriver,
   _Out_ WDFDEVICE* pDevice
   )
{
   NTSTATUS status;
   DECLARE_CONST_UNICODE_STRING(ntDeviceName, DDPROXY_DEVICE_NAME);
   DECLARE_CONST_UNICODE_STRING(symbolicName, DDPROXY_SYMBOLIC_NAME);
   WDF_DRIVER_CONFIG config;
   PWDFDEVICE_INIT pInit = NULL;

   WDF_DRIVER_CONFIG_INIT(&config, WDF_NO_EVENT_CALLBACK);

   config.DriverInitFlags |= WdfDriverInitNonPnpDriver;
   config.EvtDriverUnload = EvtDriverUnload;

   status = WdfDriverCreate(
               driverObject,
               registryPath,
               WDF_NO_OBJECT_ATTRIBUTES,
               &config,
               pDriver
               );

   if (!NT_SUCCESS(status))
   {
      goto Exit;
   }

   pInit = WdfControlDeviceInitAllocate(*pDriver, &SDDL_DEVOBJ_SYS_ALL_ADM_ALL);
   if (!pInit)
   {
      status = STATUS_INSUFFICIENT_RESOURCES;
      goto Exit;
   }

   WdfDeviceInitSetDeviceType(pInit, FILE_DEVICE_NETWORK);
   WdfDeviceInitSetCharacteristics(pInit, FILE_DEVICE_SECURE_OPEN, FALSE);

   status = WdfDeviceInitAssignName(pInit, &ntDeviceName);
   if (!NT_SUCCESS(status))
   {
	   goto Exit;
   }

   status = WdfDeviceCreate(&pInit, WDF_NO_OBJECT_ATTRIBUTES, pDevice);
   if (!NT_SUCCESS(status))
   {
      WdfDeviceInitFree(pInit);
      goto Exit;
   }

   status = WdfDeviceCreateSymbolicLink(*pDevice, &symbolicName);
   if (!NT_SUCCESS(status))
   {
	   WdfDeviceInitFree(pInit);
	   goto Exit;
   }

   status = DDProxyCtlDriverInit(pDevice);
   if (!NT_SUCCESS(status))
   {
	   WdfDeviceInitFree(pInit);
	   goto Exit;
   }

   WdfControlFinishInitializing(*pDevice);

Exit:
   return status;
}


NTSTATUS
DriverEntry(
   DRIVER_OBJECT* driverObject,
   UNICODE_STRING* registryPath
   )
{
   NTSTATUS status;
   WDFDRIVER driver;
   WDFDEVICE device;

   // Request NX Non-Paged Pool when available
   ExInitializeDriverRuntime(DrvRtPoolNxOptIn);

   status = DDProxyInitDriverObjects(
               driverObject,
               registryPath,
               &driver,
               &device
               );

   if (!NT_SUCCESS(status))
   {
      goto Exit;
   }

   InitializeListHead(&gFlowList);
   KeInitializeSpinLock(&gFlowListLock);

   gWdmDevice = WdfDeviceWdmGetDeviceObject(device);
   
   status = DDProxyRegisterCallouts(gWdmDevice);

   if (!NT_SUCCESS(status))
   {
      goto Exit;
   }

Exit:
   
   if (!NT_SUCCESS(status))
   {
      if (gEngineHandle != NULL)
      {
         DDProxyUnregisterCallouts();
      }
   }

   return status;
}
