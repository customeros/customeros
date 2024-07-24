import { useState } from 'react';

import { observer } from 'mobx-react-lite';
import { ContractStore } from '@store/Contracts/Contract.store.ts';

import { cn } from '@ui/utils/cn.ts';
import { Plus } from '@ui/media/icons/Plus.tsx';
import { useStore } from '@shared/hooks/useStore';
import { Delete } from '@ui/media/icons/Delete.tsx';
import { Button } from '@ui/form/Button/Button.tsx';
import { toastError } from '@ui/presentation/Toast';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { Spinner } from '@ui/feedback/Spinner/Spinner.tsx';
import { Divider } from '@ui/presentation/Divider/Divider.tsx';
import { outlineButton } from '@ui/form/Button/Button.variants.ts';
import { FileDropUploader, FileUploadTrigger } from '@ui/form/FileUploader';

type UploadResponse = {
  id: string;
  size: number;
  cdnUrl: string;
  fileName: string;
  mimeType: string;
  previewUrl: string;
  downloadUrl: string;
};

interface ContractUploaderProps {
  contractId: string;
}

export const ContractUploader = observer(
  ({ contractId }: ContractUploaderProps) => {
    const { contracts } = useStore();
    const contractStore = contracts.value.get(contractId) as ContractStore;
    const [files, setFiles] = useState<{ file: File; refId: number }[]>([]);
    const [loadingIds, setIsLoading] = useState<number[]>([]);
    const [_isDragging, setIsDragging] = useState(false);

    const attachments = contractStore?.tempValue?.attachments;
    const handelLoad = (refId: number) =>
      setIsLoading((prev) => [...prev, refId]);
    const clearLoad = (refId: number) =>
      setIsLoading((prev) => prev.filter((id) => id !== refId));
    const handleError = (refId: number, error: string) => {
      clearLoad(refId);
      setFiles((prev) => prev.filter((file) => file.refId !== refId));
      toastError(error, 'upload-file');
    };
    const handleLoadEnd = (refId: number) => {
      clearLoad(refId);
      setFiles((prev) => prev.filter((file) => file.refId !== refId));
    };

    const handleAddAttachment = (refId: number, res: unknown) => {
      const { id } = res as UploadResponse;

      contractStore.addAttachment(id).then(() => {
        clearLoad(refId);
      });
    };

    const handleRemoveAttachment = (id: string) => {
      contractStore.removeAttachment(id);
    };

    return (
      <div className='flex flex-col'>
        <div className='flex relative items-center h-8 '>
          <p className='text-sm text-gray-500 after:border-t-2 w-fit whitespace-nowrap mr-2'>
            Contracts & documents
          </p>
          <Divider />
          <Tooltip
            hasArrow
            side='bottom'
            align='center'
            label='Upload a document'
          >
            <FileUploadTrigger
              name='contractUpload'
              apiBaseUrl='/fs'
              endpointOptions={{
                fileKeyName: 'file',
                uploadUrl: '/file',
              }}
              onChange={(file, refId) => {
                setFiles((prev) => [...prev, { file, refId }]);
              }}
              onError={handleError}
              onLoadStart={handelLoad}
              onLoadEnd={handleLoadEnd}
              onSuccess={handleAddAttachment}
              className={cn(
                'p-1 rounded-md cursor-pointer ml-[5px] outline-none focus:outline-none',
                loadingIds.length && 'opacity-50 pointer-events-none ',
                outlineButton({ colorScheme: 'gray' }),
              )}
            >
              <Plus className='size-3 outline-none' tabIndex={-1} />
            </FileUploadTrigger>
          </Tooltip>
        </div>

        <FileDropUploader
          apiBaseUrl='/fs'
          endpointOptions={{
            fileKeyName: 'file',
            uploadUrl: '/file',
          }}
          onChange={(file, refId) => {
            setFiles((prev) => [...prev, { file, refId }]);
          }}
          onError={handleError}
          onLoadStart={handelLoad}
          onLoadEnd={handleLoadEnd}
          onSuccess={handleAddAttachment}
          onDragOverChange={setIsDragging}
        >
          <div className='min-h-5 gap-2'>
            {!attachments?.length && !files.length && (
              <label
                htmlFor='contractUpload'
                className='text-base text-gray-500 underline cursor-pointer'
              ></label>
            )}

            {attachments?.map(({ id, fileName }) => (
              <AttachmentItem
                id={id}
                key={id}
                fileName={fileName}
                onRemove={handleRemoveAttachment}
                href={`/fs/file/${id}/download`}
              />
            ))}

            {files.map(({ file, refId }) => (
              <AttachmentItem
                href='#'
                key={refId}
                fileName={file.name}
                id={refId.toString()}
                isLoading={loadingIds.includes(refId)}
              />
            ))}
            <div className='p-4 border border-dashed border-gray-300 rounded-lg text-center mt-2'>
              <p className='text-sm text-gray-500'>
                <label
                  htmlFor='contractUpload'
                  className='text-sm text-gray-500 underline cursor-pointer'
                >
                  Click to upload{' '}
                </label>
                or Drag and drop
              </p>
            </div>
          </div>
        </FileDropUploader>
      </div>
    );
  },
);

interface AttachmentItemProps {
  id: string;
  href: string;
  fileName: string;
  isLoading?: boolean;
  onRemove?: (id: string) => void;
}

const AttachmentItem = observer(
  ({ id, fileName, onRemove, isLoading }: AttachmentItemProps) => {
    const { files } = useStore();

    const handleDownload = () => {
      files.downloadAttachment(id, fileName);
      files.clear(id);
    };

    return (
      <div className='flex  items-center group mt-2 mb-3'>
        <Button
          variant='ghost'
          size='xs'
          className={
            'text-base pt-0 pb-0 leading-none font-normal text-gray-500 underline hover:bg-transparent focus:bg-transparent group-hover:text-gray-700'
          }
          onClick={handleDownload}
        >
          {fileName}
        </Button>

        {isLoading ? (
          <Spinner
            size='sm'
            label='loading'
            className='text-gray-300 fill-gray-700'
          />
        ) : (
          <Delete
            role='button'
            aria-label='Delete attachment'
            onClick={() => onRemove?.(id)}
            className='hidden size-4 text-gray-500 cursor-pointer group-hover:inline-block hover:text-gray-700'
          />
        )}
      </div>
    );
  },
);
