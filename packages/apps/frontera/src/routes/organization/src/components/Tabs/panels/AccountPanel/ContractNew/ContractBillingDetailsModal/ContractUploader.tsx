import { useState } from 'react';

import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';

import { cn } from '@ui/utils/cn';
import { Plus } from '@ui/media/icons/Plus';
import { Delete } from '@ui/media/icons/Delete';
import { toastError } from '@ui/presentation/Toast';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Spinner } from '@ui/feedback/Spinner/Spinner';
import { Divider } from '@ui/presentation/Divider/Divider';
import { ghostButton } from '@ui/form/Button/Button.variants';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { FileDropUploader, FileUploadTrigger } from '@ui/form/FileUploader';
import { useGetContractQuery } from '@organization/graphql/getContract.generated';
import { useAddContractAttachmentMutation } from '@organization/graphql/addContractAttachment.generated';
import { useRemoveContractAttachmentMutation } from '@organization/graphql/removeContractAttachment.generated';

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

export const ContractUploader = ({ contractId }: ContractUploaderProps) => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const queryKey = useGetContractQuery.getKey({ id: contractId });

  const [files, setFiles] = useState<{ file: File; refId: number }[]>([]);
  const [loadingIds, setIsLoading] = useState<number[]>([]);
  const [isDragging, setIsDragging] = useState(false);
  const { data: attachments } = useGetContractQuery(
    client,
    { id: contractId },
    { select: (data) => data.contract.attachments },
  );

  const addContractAttachment = useAddContractAttachmentMutation(client, {
    onMutate: (variables) => {
      queryClient.cancelQueries({ queryKey });

      const previousEntries = useGetContractQuery.mutateCacheEntry(
        queryClient,
        { id: contractId },
      )((cache) =>
        produce(cache, (draft) => {
          draft.contract.attachments?.push({
            id: variables.attachmentId,
            fileName: 'temp_file',
            basePath: '',
          });
        }),
      );

      return { previousEntries };
    },
    onError: (_, __, context) => {
      if (context?.previousEntries) {
        queryClient.setQueryData(queryKey, context.previousEntries);
      }
      toastError('Failed to add attachment', 'add-contract-attachment');
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey });
    },
  });

  const removeContractAttachment = useRemoveContractAttachmentMutation(client, {
    onMutate: (variables) => {
      queryClient.cancelQueries({ queryKey });

      const previousEntries = useGetContractQuery.mutateCacheEntry(
        queryClient,
        { id: contractId },
      )((cache) =>
        produce(cache, (draft) => {
          draft.contract.attachments = draft.contract.attachments?.filter(
            (attachment) => attachment.id !== variables.attachmentId,
          );
        }),
      );

      return { previousEntries };
    },
    onError: (_, __, context) => {
      if (context?.previousEntries) {
        queryClient.setQueryData(queryKey, context.previousEntries);
      }
      toastError('Failed to remove attachment', 'remove-contract-attachment');
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey });
    },
  });

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

    addContractAttachment.mutate(
      {
        contractId,
        attachmentId: id,
      },
      {
        onSettled: () => {
          clearLoad(refId);
        },
      },
    );
  };
  const handleRemoveAttachment = (id: string) => {
    removeContractAttachment.mutate({ contractId, attachmentId: id });
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
              ghostButton({ colorScheme: 'gray' }),
              'hover:bg-gray-100 p-1 rounded-lg cursor-pointer',
              loadingIds.length && 'opacity-50 pointer-events-none',
            )}
          >
            <Plus className='size-3' style={{ color: 'black' }} tabIndex={-1} />
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
        {isDragging ? (
          <div className='p-4 border border-dashed border-gray-300 rounded-lg text-center'>
            <p className='text-xs text-gray-500'>
              Drag and drop documents here
            </p>
          </div>
        ) : (
          <div className='min-h-5'>
            {!attachments?.length && !files.length && (
              <label
                htmlFor='contractUpload'
                className='text-base text-gray-500 underline cursor-pointer'
              >
                Upload a document
              </label>
            )}

            {attachments?.map(({ id, fileName }) => (
              <AttachmentItem
                id={id}
                key={id}
                fileName={fileName}
                onRemove={handleRemoveAttachment}
                href={`/fs/file/${id}/download?inline=true`}
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
          </div>
        )}
      </FileDropUploader>
    </div>
  );
};

interface AttachmentItemProps {
  id: string;
  href: string;
  fileName: string;
  isLoading?: boolean;
  onRemove?: (id: string) => void;
}

const AttachmentItem = ({
  id,
  href,
  fileName,
  onRemove,
  isLoading,
}: AttachmentItemProps) => {
  return (
    <div className='flex gap-2 items-center group'>
      <a
        href={href}
        target='_blank'
        rel='noopener noreferrer'
        className='text-base text-gray-500 underline group-hover:text-gray-700'
      >
        {fileName}
      </a>
      {isLoading ? (
        <Spinner
          size='sm'
          label='loading'
          className='text-gray-300 fill-gray-700'
        />
      ) : (
        <Delete
          aria-label='Delete attachment'
          onClick={() => onRemove?.(id)}
          className='hidden size-4 text-gray-500 cursor-pointer group-hover:inline-block hover:text-gray-700'
        />
      )}
    </div>
  );
};
