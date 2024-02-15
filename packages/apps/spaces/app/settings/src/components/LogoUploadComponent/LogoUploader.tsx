import React, { useRef, useState, useEffect } from 'react';
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

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { Upload01 } from '@ui/media/icons/Upload01';
import { Image as ChakraImage } from '@ui/media/Image';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCustomerLogo } from '@shared/state/CustomerLogo.atom';

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
  const queryKey = useTenantSettingsQuery.getKey();
  const [{ logoUrl, dimensions }, setLogoUrl] = useCustomerLogo();
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
    },
  });

  const [files, setFiles] = React.useState<Array<FilePondFile>>([]);

  const fetchLogo = async ({ id }: { id: string }) => {
    try {
      const response = await fetch(`/fs/file/${id}/download`);
      const blob = await response.blob();
      const reader = new FileReader();
      reader.onload = function () {
        const img = new Image();
        img.src = reader.result as string;
        const dataUrl = reader.result as string;
        if (dataUrl) {
          setLogoUrl({
            logoUrl: dataUrl,
            dimensions: {
              width: img.width || 136,
              height: img.height || 36,
            },
          });
        } else {
          setHasError({ error: 'Error loading logo', file: 'logo' });
        }
      };
      reader.readAsDataURL(blob);
    } catch (reason) {
      setHasError({ error: 'Error loading logo', file: 'logo' });
    }
  };

  useWillUnmount(() => {
    queryClient.cancelQueries({ queryKey });
  });
  useEffect(() => {
    if (tenantSettingsData?.tenantSettings?.logoUrl) {
      const uuidRegex =
        /[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/;
      const match = `${tenantSettingsData?.tenantSettings?.logoUrl}`.match(
        uuidRegex,
      );

      if (match) {
        fetchLogo({ id: match[0] });
      }
    }
  }, [tenantSettingsData?.tenantSettings?.logoUrl]);

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

  return (
    <Box position='relative'>
      <Flex justifyContent='space-between' alignItems='center' mb={2}>
        <Text color='gray.600' fontSize='sm' fontWeight='semibold'>
          Organization logo
        </Text>
        {logoUrl && (
          <IconButton
            variant='ghost'
            aria-label='Upload file'
            size='sm'
            color={'gray.500'}
            icon={<Upload01 />}
            onClick={() => pondRef.current?.browse()}
          />
        )}
      </Flex>

      {logoUrl && !isLoading && !hasError && (
        <Box
          position='relative'
          maxHeight={120}
          width='full'
          display='flex'
          justifyContent='center'
          padding={4}
        >
          <ChakraImage
            src={`${logoUrl}`}
            alt='CustomerOS'
            width={dimensions.width || 136}
            height={dimensions.height || 45}
            style={{ objectFit: 'contain', maxHeight: '40px' }}
          />
        </Box>
      )}

      <Box
        onClick={() => hasError && pondRef.current?.browse()}
        className={
          hasError ? 'filepond-error' : logoUrl ? 'filepond-uploaded' : ''
        }
        sx={{
          '&': {
            position:
              logoUrl && !isLoading && !hasError ? 'absolute' : 'static',
            top: logoUrl && !isLoading && !hasError ? '-9999' : 'auto',
          },
          '& .filepond--root .filepond--drop-label': {
            minHeight:
              files.length || hasError ? `${32}px` : `${120}px !important`,
          },
        }}
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
                  const previewUrl = JSON.parse(request.response).id;
                  updateTenantSettingsMutation.mutate({
                    input: {
                      patch: true,
                      logoUrl: previewUrl, // Ensure this matches the structure of your API response
                    },
                  });
                  const reader = new FileReader();
                  reader.readAsDataURL(file);
                  reader.onload = function () {
                    const img = new Image();
                    img.src = reader.result as string;
                    setLogoUrl({
                      logoUrl: reader.result as string,
                      dimensions: {
                        width: img.width,
                        height: img.height,
                      },
                    });
                    fetchLogo({ id: previewUrl });

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
            if (logoUrl) {
              updateTenantSettingsMutation.mutate({
                input: {
                  patch: true,
                  logoUrl: '',
                },
              });
            }
          }}
          onaddfilestart={(file) => {
            setHasError(null);
            if (logoUrl) {
              updateTenantSettingsMutation.mutate({
                input: {
                  patch: true,
                  logoUrl: '',
                },
              });
            }
          }}
        />
      </Box>
    </Box>
  );
};
