import NextImage from 'next/image';
import React, { useRef, useState } from 'react';
import { FilePond, FileStatus, registerPlugin } from 'react-filepond';

import { useWillUnmount } from 'rooks';
import { FilePondFile } from 'filepond';
import { renderToString } from 'react-dom/server';
import { useQueryClient } from '@tanstack/react-query';
import FilePondPluginImageResize from 'filepond-plugin-image-resize';
import FilePondPluginImagePreview from 'filepond-plugin-image-preview';
import FilePondPluginValidateSize from 'filepond-plugin-file-validate-size';
import FilePondPluginFileValidateType from 'filepond-plugin-file-validate-type';
import { useTenantSettingsQuery } from '@settings/graphql/getTenantSettings.generated';
import { useUpdateTenantSettingsMutation } from '@settings/graphql/updateTenantSettings.generated';

import { cn } from '@ui/utils/cn';
import { Upload01 } from '@ui/media/icons/Upload01';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

registerPlugin(FilePondPluginImagePreview);
registerPlugin(FilePondPluginValidateSize);
registerPlugin(FilePondPluginImageResize);
registerPlugin(FilePondPluginFileValidateType);

interface LogoUploaderProps {}

export const LogoUploader: React.FC<LogoUploaderProps> = () => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const pondRef = useRef<FilePond | null>(null);

  const { data: tenantSettingsData } = useTenantSettingsQuery(client);
  const { data: globalCacheData } = useGlobalCacheQuery(client);
  const queryKey = useTenantSettingsQuery.getKey();
  const globalCacheQueryKey = useGlobalCacheQuery.getKey();
  const [hasError, setHasError] = useState<null | {
    file: string;
    error: string;
  }>(null);

  const updateTenantSettingsMutation = useUpdateTenantSettingsMutation(client, {
    onMutate: ({ input: { patch, ...newSettings } }) => {
      queryClient.cancelQueries({ queryKey });
      const previousSettings = tenantSettingsData?.tenantSettings;
      queryClient.setQueryData(queryKey, {
        tenantSettings: {
          ...previousSettings,
          ...newSettings,
        },
      });

      return { previousSettings };
    },
    onError: (err, newSettings, context) => {
      queryClient.setQueryData(queryKey, context?.previousSettings);
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey });
      queryClient.invalidateQueries({ queryKey: globalCacheQueryKey });
    },
  });

  const [files, setFiles] = React.useState<Array<FilePondFile>>([]);

  useWillUnmount(() => {
    queryClient.cancelQueries({ queryKey });
  });

  function getDefaultLabel() {
    return renderToString(
      <div className='filepond-label-container'>
        <svg
          width='24'
          height='24'
          viewBox='0 0 24 24'
          fill='none'
          xmlns='http://www.w3.org/2000/svg'
        >
          <path
            d='M21 15V16.2C21 17.8802 21 18.7202 20.673 19.362C20.3854 19.9265 19.9265 20.3854 19.362 20.673C18.7202 21 17.8802 21 16.2 21H7.8C6.11984 21 5.27976 21 4.63803 20.673C4.07354 20.3854 3.6146 19.9265 3.32698 19.362C3 18.7202 3 17.8802 3 16.2V15M17 8L12 3M12 3L7 8M12 3V15'
            stroke='#475467'
            strokeWidth='2'
            strokeLinecap='round'
            strokeLinejoin='round'
          />
        </svg>
        <p className='filepond-idle-label-text'>
          <span className='filepond--label-action'>Click to upload</span> or
          drag and drop
          <p className='filepond-sizes'>PNG or JPG (Max 150KB)</p>
        </p>
      </div>,
    );
  }

  const isLoading =
    pondRef.current?.getFile()?.status === FileStatus.PROCESSING_QUEUED ||
    pondRef.current?.getFile()?.status === FileStatus.PROCESSING ||
    pondRef.current?.getFile()?.status === FileStatus.LOADING;

  const statusClassName = hasError
    ? 'filepond-error'
    : globalCacheData?.global_Cache?.cdnLogoUrl
    ? 'filepond-uploaded'
    : '';

  const position =
    globalCacheData?.global_Cache?.cdnLogoUrl && !isLoading && !hasError
      ? 'absolute'
      : 'static';

  const minHeightClass =
    files.length || hasError ? `${32}px` : `${120}px !important`;

  const topPosition =
    globalCacheData?.global_Cache?.cdnLogoUrl && !isLoading && !hasError
      ? 'top-[-9999px]'
      : 'auto';

  return (
    <div>
      <div className='flex justify-between items-center mb-2'>
        <span className='text-gray-600 text-sm font-semibold'>
          Organization logo
        </span>
        {globalCacheData?.global_Cache?.cdnLogoUrl && (
          <IconButton
            className='pr-0'
            variant='ghost'
            aria-label='Upload file'
            size='sm'
            colorScheme='gray'
            icon={<Upload01 />}
            onClick={() => pondRef.current?.browse()}
          />
        )}
      </div>

      {globalCacheData?.global_Cache?.cdnLogoUrl && !isLoading && !hasError && (
        <div className='relative max-h-[120px] w-full flex p-4 '>
          <NextImage
            src={`${globalCacheData?.global_Cache?.cdnLogoUrl}`}
            alt='CustomerOS'
            width={136}
            height={40}
            style={{
              objectFit: 'contain',
              maxHeight: '40px',
              maxWidth: 'fit-content',
            }}
          />
        </div>
      )}

      <div
        onClick={() => hasError && pondRef.current?.browse()}
        className={cn(
          statusClassName,
          position,
          topPosition,
          '[&_.filepond--root .filepond--drop-label]:',
          minHeightClass,
          '[&_.filepond--image-clip]:m-0',
        )}
      >
        <FilePond
          ref={pondRef}
          // @ts-expect-error ignore for now
          files={files}
          onupdatefiles={setFiles}
          dropOnPage={true}
          dropOnElement={false}
          server={{
            url: '/fs/file',

            timeout: 5000,
            // load: (source, load, error, progress, abort, headers) => {
            //   const myRequest = new Request(source);
            //   fetch(myRequest).then(function (response) {
            //     response.blob().then(function (myBlob) {
            //       load(myBlob);
            //     });
            //   });
            // },
            // fetch: (source, load, error, progress, abort, headers) => {
            //   const myRequest = new Request(source);
            //   fetch(myRequest).then(function (response) {
            //     response.blob().then(function (myBlob) {
            //       load(myBlob);
            //     });
            //   });
            // },

            process: (
              fieldName,
              file,
              metadata,
              load,
              error,
              progress,
              abort,
              transfer,
              options,
            ) => {
              // fieldName is the name of the input field
              // file is the actual file object to send
              const formData = new FormData();
              formData.append('file', file, file.name);
              formData.append('cdnUpload', 'true');

              const request = new XMLHttpRequest();
              request.open('POST', '/fs/file');
              // Should call the progress method to update the progress to 100% before calling load
              // Setting computable to false switches the loading indicator to infinite mode
              request.upload.onprogress = (e) => {
                progress(e.lengthComputable, e.loaded, e.total);
              };
              // Should call the load method when done and pass the returned server file id
              // this server file id is then used later on when reverting or restoring a file
              // so your server knows which file to return without exposing that info to the client
              request.onload = function (ev) {
                if (request.status >= 200 && request.status < 300) {
                  // the load method accepts either a string (id) or an object
                  load(request.responseText);
                  const parsedResponse = JSON.parse(request.response);
                  updateTenantSettingsMutation.mutate({
                    input: {
                      patch: true,
                      logoUrl: parsedResponse?.previewUrl,
                      logoRepositoryFileId: parsedResponse.id,
                    },
                  });
                  const reader = new FileReader();
                  reader.readAsDataURL(file);
                  reader.onload = function () {
                    const img = new Image();
                    img.src = reader.result as string;

                    return reader.result;
                  };
                } else {
                  // Can call the error method if something is wrong, should exit after
                  error('oh no');
                }
              };

              request.send(formData);

              // Should expose an abort method so the request can be cancelled
              return {
                abort: () => {
                  // This function is entered if the user has tapped the cancel button
                  request.abort();

                  // Let FilePond know the request has been cancelled
                  abort();
                },
              };
            },
          }}
          maxFiles={1}
          allowMultiple={false}
          allowReplace={true}
          name='files'
          allowDrop={!files.length}
          acceptedFileTypes={['image/jpg', 'image/png', 'image/jpeg']}
          panelHeight={120}
          imagePreviewMaxFileSize='150KB'
          maxFileSize='150KB'
          imagePreviewMaxInstantPreviewFileSize={150000}
          imagePreviewMaxHeight={32}
          imageResizeTargetWidth={40}
          imageResizeMode='contain'
          imageResizeTargetHeight={32}
          labelIdle={getDefaultLabel()}
          labelFileProcessing={'Uploading'}
          labelMaxFileSizeExceeded={'Your logo needs to be less than 150KB'}
          labelFileWaitingForSize={'Waiting for size'}
          labelFileLoadError={'Upload failed, please try again'}
          labelFileProcessingError={'Upload failed, please try again'}
          labelFileProcessingComplete={'Logo uploaded successfully'}
          labelFileTypeNotAllowed={'Your logo must be a PNG or JPG'}
          credits={false}
          onerror={(error, file) => {
            // @ts-expect-error error file
            setHasError({ error: error?.main, file: file?.file?.name });
          }}
          onremovefile={() => {
            setHasError(null);
            if (globalCacheData?.global_Cache?.cdnLogoUrl) {
              updateTenantSettingsMutation.mutate({
                input: {
                  patch: true,
                  logoUrl: '',
                  logoRepositoryFileId: '',
                },
              });
            }
          }}
          onaddfilestart={(file) => {
            setHasError(null);
            if (globalCacheData?.global_Cache?.cdnLogoUrl) {
              updateTenantSettingsMutation.mutate({
                input: {
                  patch: true,
                  logoUrl: '',
                  logoRepositoryFileId: '',
                },
              });
            }
          }}
        />
      </div>
    </div>
  );
};
