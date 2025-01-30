import { useMemo } from 'react';

import { observer } from 'mobx-react-lite';
import { EditOwnerOrganizationUsecase } from '@domain/usecases/organization-details/edit-owner.usecase';

import { cn } from '@ui/utils/cn.ts';
import { Icon } from '@ui/media/Icon';
import { Combobox } from '@ui/form/Combobox';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Popover, PopoverContent, PopoverTrigger } from '@ui/overlay/Popover';

interface OwnerProps {
  id: string;
  dataTest?: string;
}

export const OwnerInput = observer(({ id, dataTest }: OwnerProps) => {
  const ownerUseCase = useMemo(
    () => new EditOwnerOrganizationUsecase(id),
    [id],
  );

  return (
    <>
      <Popover
        open={ownerUseCase.isMenuOpen}
        onOpenChange={(open) => ownerUseCase.toggleMenu(open)}
      >
        <Tooltip label='Owner' align='start'>
          <PopoverTrigger className={cn('flex items-center')}>
            <Icon name='key-01' className='mr-3 text-gray-500' />
            <div
              data-test={dataTest}
              className='flex flex-wrap  w-fit items-center'
            >
              {ownerUseCase.selectedUser ? (
                <div className='text-sm'>{ownerUseCase.selectedUser.label}</div>
              ) : (
                <span className='text-gray-400 text-sm'>Owner</span>
              )}
            </div>
          </PopoverTrigger>
        </Tooltip>
        <PopoverContent align='start' className='min-w-[264px] max-w-[320px]'>
          <Combobox
            isMulti
            options={ownerUseCase.userOptions}
            inputValue={ownerUseCase.searchTerm}
            value={ownerUseCase.selectedUser?.value}
            onInputChange={ownerUseCase.setSearchTerm}
            onChange={(newValue) => {
              ownerUseCase.execute(newValue[0]);
              ownerUseCase.toggleMenu(false);
            }}
            noOptionsMessage={({ inputValue }) => (
              <div className='text-gray-700 px-3 py-1 mt-0.5 rounded-md bg-grayModern-100 gap-1 flex items-center'>
                <span>{`No results matching "${inputValue}"`}</span>
              </div>
            )}
          />
        </PopoverContent>
      </Popover>
    </>
  );
});
