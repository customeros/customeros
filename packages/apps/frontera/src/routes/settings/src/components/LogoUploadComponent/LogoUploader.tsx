import { useState } from 'react';

import { observer } from 'mobx-react-lite';

import { cn } from '@ui/utils/cn';
import { X } from '@ui/media/icons/X';
import { Image } from '@ui/media/Image/Image';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { toastError } from '@ui/presentation/Toast';
import { Upload01 } from '@ui/media/icons/Upload01';
import { ghostButton } from '@ui/form/Button/Button.variants';
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

export const LogoUploader = observer(() => {
  const store = useStore();

  const [file, setFile] = useState<File | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isDragging, setIsDragging] = useState(false);

  const handelLoad = () => setIsLoading(true);
  const clearLoad = () => setIsLoading(false);

  const handleError = (_refId: number, error: string) => {
    clearLoad();
    setFile(null);
    toastError(error, 'upload-file');
  };

  const handleLoadEnd = () => {
    setFile(null);
    clearLoad();
  };

  const handleTenantLogoUpdate = (_refId: number, res: unknown) => {
    const { id } = res as UploadResponse;

    store.settings.tenant.update((value) => {
      value.logoRepositoryFileId = id;

      return value;
    });
    clearLoad();
  };

  const handleTenantLogoRemove = () => {
    store.settings.tenant.update((value) => {
      value.logoRepositoryFileId = '';

      return value;
    });
    setFile(null);
  };

  return (
    <div className='flex flex-col'>
      <div className='flex justify-between items-center'>
        <p className='text-sm text-gray-900 w-fit whitespace-nowrap font-semibold'>
          Organization logo
        </p>

        <FileUploadTrigger
          apiBaseUrl='/fs'
          onChange={setFile}
          name='logoUploader'
          onError={handleError}
          onLoadStart={handelLoad}
          onLoadEnd={handleLoadEnd}
          onSuccess={handleTenantLogoUpdate}
          endpointOptions={{
            fileKeyName: 'file',
            uploadUrl: '/file',
          }}
          className={cn(
            ghostButton({ colorScheme: 'gray' }),
            'hover:bg-gray-100 p-1 rounded-lg cursor-pointer',
            isLoading && 'opacity-50 pointer-events-none',
          )}
        >
          <Upload01 />
        </FileUploadTrigger>
      </div>

      <FileDropUploader
        apiBaseUrl='/fs'
        onChange={setFile}
        onError={handleError}
        onLoadStart={handelLoad}
        onLoadEnd={handleLoadEnd}
        onDragOverChange={setIsDragging}
        onSuccess={handleTenantLogoUpdate}
        endpointOptions={{
          fileKeyName: 'file',
          uploadUrl: '/file',
        }}
      >
        {isDragging ? (
          <div className='p-4 border border-dashed border-gray-300 rounded-lg text-center'>
            <p className='text-sm text-gray-500'>
              Drag and drop PNG or JPG (Max 150KB)
            </p>
          </div>
        ) : (
          <div className='min-h-5 pt-2'>
            {!store.settings.tenant.value?.logoRepositoryFileId && !file && (
              <label
                htmlFor='logoUploader'
                className='text-sm text-gray-500 underline cursor-pointer'
              >
                Upload a PNG or JPG (Max 150KB)
              </label>
            )}

            {store.settings.tenant.value?.logoRepositoryFileId && !file && (
              <div className='relative max-h-16 w-fit'>
                <Image
                  className='max-h-16'
                  src={store.settings.tenant.value?.logoRepositoryFileId}
                />
                <IconButton
                  size='xxs'
                  icon={<X />}
                  variant='outline'
                  aria-label='Remove Logo'
                  onClick={handleTenantLogoRemove}
                  className='absolute bg-white bg-opacity-50 -top-0.5 -right-5 rounded-full'
                />
              </div>
            )}

            {!store.settings.tenant.value?.logoRepositoryFileId && file && (
              <Image
                src={`${URL.createObjectURL(file)}`}
                className='max-h-16 animate-pulseOpacity'
              />
            )}
          </div>
        )}
      </FileDropUploader>
    </div>
  );
});
