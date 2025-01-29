import { useMemo } from 'react';

import { observer } from 'mobx-react-lite';
import { EditOrganizationOwnerUsecase } from '@domain/usecases/organization-about-panel/edit-organization-owner.usecase';

import { cn } from '@ui/utils/cn';
import { Combobox } from '@ui/form/Combobox';
import { Key01 } from '@ui/media/icons/Key01';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { Popover, PopoverContent, PopoverTrigger } from '@ui/overlay/Popover';

interface OwnerProps {
  id: string;
  dataTest?: string;
}

export const OwnerInput = observer(({ id, dataTest }: OwnerProps) => {
  const usecase = useMemo(() => new EditOrganizationOwnerUsecase(id), [id]);

  return (
    <>
      <Popover
        open={usecase.isMenuOpen}
        onOpenChange={(open) => usecase.toggleMenu(open)}
      >
        <Tooltip label='Owner' align='start'>
          <PopoverTrigger className={cn('flex items-center')}>
            <Key01 className='mr-3 text-gray-500' />
            <div
              data-test={dataTest}
              className='flex flex-wrap  w-fit items-center'
            >
              {usecase.selectedUser ? (
                <div className='text-sm'>{usecase.selectedUser.label}</div>
              ) : (
                <span className='text-gray-400 text-sm'>Owner</span>
              )}
            </div>
          </PopoverTrigger>
        </Tooltip>
        <PopoverContent align='start' className='min-w-[264px] max-w-[320px]'>
          <Combobox
            placeholder='Owner'
            value={usecase.selectedUser}
            options={usecase.userOptions}
            inputValue={usecase.searchTerm}
            onInputChange={usecase.setSearchTerm}
            onChange={(newValue) => {
              usecase.select(newValue);
              usecase.toggleMenu(false);
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
