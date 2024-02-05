import { FilePond, registerPlugin } from 'react-filepond';
import React, { useRef, useState, useEffect, ChangeEvent } from 'react';

import axios from 'axios';
import Compressor from 'compressorjs';
import { useWillUnmount } from 'rooks';
import { ExtFile } from '@files-ui/core';
import { renderToString } from 'react-dom/server';
import { useQueryClient } from '@tanstack/react-query';
import FilePondPluginImagePreview from 'filepond-plugin-image-preview';
import { useTenantSettingsQuery } from '@settings/graphql/getTenantSettings.generated';
import { useUpdateTenantSettingsMutation } from '@settings/graphql/updateTenantSettings.generated';
import {
  Dropzone,
  ImagePreview,
  FileInputButton,
  FilesUiProvider,
} from '@files-ui/react';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { Image03 } from '@ui/media/icons/Image03';
import { Upload01 } from '@ui/media/icons/Upload01';
import { Image as ChakraImage } from '@ui/media/Image';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCustomerLogo } from '@shared/state/CustomerLogo.atom';

registerPlugin(FilePondPluginImagePreview);
interface LogoUploaderProps {}

export const LogoUploader: React.FC<LogoUploaderProps> = () => {
  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const pond = useRef<FilePond | null>(null);

  const { data: tenantSettingsData } = useTenantSettingsQuery(client);
  const queryKey = useTenantSettingsQuery.getKey();
  const [{ logoUrl, dimensions }, setLogoUrl] = useCustomerLogo();
  const [progress, setProgress] = useState(0);

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

  const [files, setFiles] = React.useState<any[]>([]);

  function FetchLogo({
    id,
    setLogoUrl,
  }: {
    id: string;
    setLogoUrl: (args: any) => void;
  }) {
    return fetch(`/fs/file/${id}/download`)
      .then(async (response: any) => {
        const blob = await response.blob();
        console.log('🏷️ ----- response: ', response);
        const reader = new FileReader();
        reader.onload = function () {
          const img = new Image();
          img.src = reader.result as string;
          const dataUrl = reader.result as any;
          if (dataUrl) {
            setLogoUrl({
              logoUrl: dataUrl,
              dimensions: {
                width: img.width || '136',
                height: img.height || '36',
              },
            });
          } else {
            console.log('🏷️ ----- : ERROR');
          }
        };
        reader.readAsDataURL(blob);
      })
      .catch((reason: any) => {
        console.log('🏷️ ----- : ERROR CATCH', reason);
      });
  }

  const updateFiles = async (incomingFiles: ExtFile[]) => {
    // Do something with the files
    // Assuming the first file is what you want to compress and upload
    const file = incomingFiles[0];
    setFiles(incomingFiles);
    if (!file?.file) {
      return;
    }
    new Compressor(file.file, {
      quality: 0.92, // Set the quality for compression
      success: async (result) => {
        setFiles([result]);

        const formData = new FormData();
        formData.append('file', result, `${(result as File).name}`); // Append the compressed file

        // Send the compressed image file to server with fetch API

        try {
          const response = await axios.post('fs/file', formData, {
            headers: {
              'Content-Type': 'multipart/form-data',
            },
            onUploadProgress: function (progressEvent) {
              setProgress((progressEvent.loaded / progressEvent.total) * 100);
            },
          });
          if (response.status === 200) {
            // Assuming updateTenantSettingsMutation is a function you've defined to update your application state
            updateTenantSettingsMutation.mutate({
              input: {
                patch: true,
                logoUrl: response.data.previewUrl, // Ensure this matches the structure of your API response
              },
            });
          }
        } catch (error) {
          console.error('Error uploading file:', error);
        }
      },
      error(err) {
        console.log('Compression error:', err.message);
      },
    });
  };
  const removeFile = (id) => {
    setFiles(files.filter((x) => x.id !== id));
  };
  useWillUnmount(() => {
    queryClient.cancelQueries({ queryKey });
  });

  //

  useEffect(() => {
    if (tenantSettingsData?.tenantSettings?.logoUrl) {
      const uuidRegex =
        /[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/;
      const match = `${tenantSettingsData?.tenantSettings?.logoUrl}`.match(
        uuidRegex,
      );
      FetchLogo({ id: match[0], setLogoUrl });
      // setFiles([{
      //   source: `${tenantSettingsData?.tenantSettings?.logoUrl}`,
      // }])
    }
  }, [tenantSettingsData?.tenantSettings?.logoUrl]);

  // useEffect(() => {
  //   if (logoUrl) {
  //     setFiles([
  //       {
  //         source: `${logoUrl}`,
  //         options: {
  //           type: 'local',
  //         },
  //       },
  //     ]);
  //   }
  // }, [logoUrl]);
  function genLabel() {
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
            strokeLineJoin='round'
          />
        </svg>
        <p className='filepond-idle-label-text'>
          <span className='filepond--label-action'>Click to upload</span> or
          drag and drop
          <p className='filepond-sizes'>PNG, JPG or GIF (Max 3MB)</p>
        </p>
      </div>,
    );
  }

  const genLogoLabel = () => {
    return renderToString(
      <Box position='relative'>
        <ChakraImage
          src={`${logoUrl}`}
          alt=''
          width={dimensions.width || 136}
          height={dimensions.height || 36}
          style={{ objectFit: 'contain' }}
        />
      </Box>,
    );
  };

  return (
    <>
      <Flex justifyContent='space-between'>
        <Text color='gray.600' fontSize='sm' fontWeight='semibold' mb={4}>
          Organization logo
        </Text>
        {logoUrl && (
          <Box
            sx={{
              '& .material-button-root.text-6:hover': {
                backgroundColor: 'var(--chakra-colors-primary-50)',
              },
            }}
          >
            <FileInputButton
              variant='text'
              value={files}
              onChange={updateFiles}
              accept={'image/*'}
              maxFileSize={150000}
              maxFiles={1}
              autoClean
              disableRipple
              style={{
                padding: 0,

                maxHeight: '21px',
                color: 'var(--chakra-colors-primary-600)',
                letterSpacing: 'normal',
                textTransform: 'unset',
                fontSize: '12px',
              }}
            >
              Upload new
            </FileInputButton>
          </Box>
        )}
      </Flex>

      <FilePond
        ref={(ref) => {
          pond.current = ref;
        }}
        files={files}
        onupdatefiles={setFiles}
        server={{
          url: '/fs/file',

          timeout: 70000,
          load: (source, load, error, progress, abort, headers) => {
            const myRequest = new Request(source);
            fetch(myRequest).then(function (response) {
              response.blob().then(function (myBlob) {
                load(myBlob);
              });
            });
          },
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
            request.onload = function () {
              if (request.status >= 200 && request.status < 300) {
                // the load method accepts either a string (id) or an object
                load(request.responseText);
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
        labelIdle={logoUrl ? genLogoLabel() : genLabel()}
        stylePanelAspectRatio={'4:1'}
      />
    </>
  );
};
