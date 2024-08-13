import { useState } from 'react';

import { cn } from '@ui/utils/cn';
import { X } from '@ui/media/icons/X';
import { Input } from '@ui/form/Input';
import { Image } from '@ui/media/Image/Image';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { toastError } from '@ui/presentation/Toast';
import { ImagePlus } from '@ui/media/icons/ImagePlus';
import { outlineButton } from '@ui/form/Button/Button.variants';
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

export const Workspace = () => {
  const store = useStore();
  const [name, setName] = useState(
    () => store.settings.tenant.value?.workspaceName,
  );

  const [file, setFile] = useState<File | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isDragging, setIsDragging] = useState(false);
  const [_triggerRender, setTriggerRender] = useState(false);

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
      value.workspaceLogo = id;

      return value;
    });
    clearLoad();
  };

  const handleTenantLogoRemove = () => {
    store.settings.tenant.update((value) => {
      value.workspaceLogo = '';

      return value;
    });
    setFile(null);
    setTriggerRender((prev) => !prev);
  };

  const handleNameChange = (value: string) => {
    setName(value);
    store.settings.tenant.update((tenant) => {
      tenant.workspaceName = value;

      return tenant;
    });
  };

  return (
    <div className='px-6 pb-4 pt-2 max-w-[415px] border-r border-gray-200 h-full'>
      <div className='flex flex-col gap-4'>
        <p className='text-gray-700  font-semibold'>Workspace</p>
        <div className='flex flex-col'>
          <div className='flex justify-between items-center'>
            <p className='text-sm text-gray-900 w-fit whitespace-nowrap font-semibold'>
              Workspace logo & name
            </p>
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
              <div className='flex  flex-1 items-center justify-between min-h-5 pt-2'>
                {!store.settings.tenant.value?.workspaceLogo && !file && (
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
                      outlineButton({ colorScheme: 'gray' }),
                      'hover:bg-gray-100 p-1 rounded-md cursor-pointer text-gray-500',
                      isLoading && 'opacity-50 pointer-events-none',
                    )}
                  >
                    <ImagePlus />
                  </FileUploadTrigger>
                )}

                {store.settings.tenant.value?.workspaceLogo && !file && (
                  <div className='relative max-h-16 w-fit group'>
                    <Image
                      className='h-10'
                      src={store.settings.tenant.value?.workspaceLogo}
                    />
                    <IconButton
                      size='xxs'
                      icon={<X />}
                      variant='outline'
                      aria-label='Remove Logo'
                      onClick={handleTenantLogoRemove}
                      className='absolute bg-white bg-opacity-50 -top-[9px] -right-[10px] rounded-full size-5 opacity-0 group-hover:opacity-100'
                    />
                  </div>
                )}

                {!store.settings.tenant.value?.workspaceLogo && file && (
                  <Image
                    src={`${URL.createObjectURL(file)}`}
                    className='max-h-16 animate-pulseOpacity'
                  />
                )}

                <Input
                  variant='unstyled'
                  value={name || ''}
                  className='ml-2.5 flex-2'
                  placeholder='Workspace name'
                  onChange={(e) => handleNameChange(e.target.value)}
                />
              </div>
            )}
          </FileDropUploader>
        </div>
      </div>
    </div>
  );
};
