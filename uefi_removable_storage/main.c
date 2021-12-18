#include <Uefi.h>
#include <Library\UefiApplicationEntryPoint.h>
#include <Library\UefiBootServicesTableLib.h>
#include <Library\DebugLib.h>
#include <Library\UefiLib.h>
#include <Library\MemoryAllocationLib.h>
#include <Library\DevicePathLib.h>
#include <Protocol\BlockIo.h>

EFI_STATUS GetDevicePathForRemovableDevice(
	IN EFI_HANDLE *Handle,
	OUT EFI_DEVICE_PATH_PROTOCOL **DevicePath)
{
	EFI_STATUS Status;
	EFI_BLOCK_IO_PROTOCOL *BlockIoProtocol;
	EFI_DEVICE_PATH_PROTOCOL *DevicePathProtocol;
	EFI_DEVICE_PATH_PROTOCOL *Current;
	EFI_DEVICE_PATH_PROTOCOL *Last;

	ASSERT(Handle != NULL);
	ASSERT(DevicePath != NULL);

	// Get block io protocol for handle
	Status = gBS->HandleProtocol(
		Handle,
		&gEfiBlockIoProtocolGuid,
		(VOID**)&BlockIoProtocol);
	if (EFI_ERROR(Status))
		return Status;

	// check removable flag
	if (BlockIoProtocol->Media->RemovableMedia != TRUE)
	{
		*DevicePath = NULL;
		return EFI_SUCCESS;
	}

	// Get device path protocol for handle
	Status = gBS->HandleProtocol(
		Handle,
		&gEfiDevicePathProtocolGuid,
		(VOID**)&DevicePathProtocol);
	
	// check last node of device path
	Current = DevicePathProtocol;
	Last = NULL;
	while (IsDevicePathEnd(Current) != TRUE)
	{
		Last = Current;
		Current = NextDevicePathNode(Current);
	}

	if (Last == NULL ||
		Last->Type != MESSAGING_DEVICE_PATH)
	{
		*DevicePath = NULL;
		return EFI_SUCCESS;
	}

	*DevicePath = DevicePathProtocol;

	return EFI_SUCCESS;
}

EFI_STATUS LocateBlockIoList(
	OUT UINTN *ListLength,
	OUT EFI_HANDLE **List)
{
	EFI_STATUS Status;
	UINTN BufferSize;
	VOID *Buffer;
	
	ASSERT(ListLength != NULL);
	ASSERT(List != NULL);

	Status = EFI_BUFFER_TOO_SMALL;
	BufferSize = 0;

	while (Status == EFI_BUFFER_TOO_SMALL) {
		Buffer = AllocatePool(BufferSize);

		Status = gBS->LocateHandle(
			ByProtocol,
			&gEfiBlockIoProtocolGuid,
			NULL,
			&BufferSize,
			Buffer);
		if (!EFI_ERROR(Status))
			break;

		FreePool(Buffer);
	}

	if (EFI_ERROR(Status))
		return Status;

	*ListLength = BufferSize / sizeof(EFI_HANDLE);
	*List = Buffer;

	return EFI_SUCCESS;
}

VOID FreeBlockIoList(
	IN EFI_HANDLE *List)
{
	ASSERT(List != NULL);

	FreePool(List);
}

EFI_STATUS PrintRemovableMedia()
{
	EFI_STATUS Status;
	UINTN BlockIoListLength;
	EFI_HANDLE *BlockIoList;
	EFI_DEVICE_PATH_PROTOCOL *DevicePathProtocol;
	CHAR16 *DevicePath;
	UINTN Counter;

	Status = LocateBlockIoList(
		&BlockIoListLength,
		&BlockIoList);
	if (EFI_ERROR(Status))
		return Status;

	Counter = 0;

	for (UINTN Index = 0; Index < BlockIoListLength; Index++)
	{
		Status = GetDevicePathForRemovableDevice(
			BlockIoList[Index],
			&DevicePathProtocol);
		if (EFI_ERROR(Status))
		{
			FreeBlockIoList(BlockIoList);
			return Status;
		}

		if (DevicePathProtocol != NULL)
		{
			DevicePath = ConvertDevicePathToText(
				DevicePathProtocol,
				TRUE,
				TRUE);
			if (EFI_ERROR(Status))
			{
				FreeBlockIoList(BlockIoList);
				return Status;
			}

			Counter++;
			Print(L"%d - %s\n", Counter, DevicePath);

			FreePool(DevicePath);
		}
	}

	if (Counter == 0)
	{
		Print(L"Not found removable devices\n");
	}
	else
	{
		Print(L"-----------------------------\n");
		Print(L"%d removable devices\n", Counter);
	}

	FreeBlockIoList(BlockIoList);

	return EFI_SUCCESS;
}

EFI_STATUS EFIAPI UefiMain(
	IN EFI_HANDLE ImageHandle,
	IN EFI_SYSTEM_TABLE *SystemTable)
{
	EFI_STATUS Status;

	Status = PrintRemovableMedia();

	if (EFI_ERROR(Status))
		ErrorPrint(L"Error status 0x%x", Status);
		
	return EFI_SUCCESS;
}