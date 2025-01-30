import { useNavigate } from 'react-router-dom';

import { useKeyBindings } from 'rooks';
import { observer } from 'mobx-react-lite';
import { AddSearchOrganizationsUsecase } from '@domain/usecases/command-menu/add-search-organizations.usecase';

import { cn } from '@ui/utils/cn';
import { Avatar } from '@ui/media/Avatar';
import { Spinner } from '@ui/feedback/Spinner';
import { useStore } from '@shared/hooks/useStore';
import { User03 } from '@ui/media/icons/User03.tsx';
import { PlusCircle } from '@ui/media/icons/PlusCircle.tsx';
import { ArrowNarrowRight } from '@ui/media/icons/ArrowNarrowRight';
import {
  Command,
  CommandItem,
  CommandInput,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

const usecase = new AddSearchOrganizationsUsecase();

export const AddNewOrganization = observer(() => {
  const store = useStore();
  const navigate = useNavigate();

  const handleClose = () => {
    store.ui.commandMenu.clearContext();
    store.ui.commandMenu.toggle('AddNewOrganization');
    store.ui.commandMenu.clearCallback();
  };

  useKeyBindings({
    Escape: handleClose,
  });

  return (
    <Command shouldFilter={false}>
      <CommandInput
        value={usecase.searchTerm}
        onValueChange={usecase.setSearchTerm}
        placeholder='Search by name or website'
        dataTest={'organizations-create-new-org-org-name'}
        label={<p className='font-medium'>Search 300,000+ companies</p>}
        rightElement={
          usecase.isValidatingDomain && (
            <Spinner
              size='sm'
              label='validating...'
              className='text-gray-300 fill-gray-500 size-3'
            />
          )
        }
        bottomAccessory={
          (usecase.domainValidationError ||
            usecase.domainValidationMessage) && (
            <p
              className={cn(
                'text-xs text-gray-500',
                usecase.domainValidationError.length > 0 && 'text-error-500',
              )}
            >
              {usecase.domainValidationError || usecase.domainValidationMessage}
            </p>
          )
        }
      />

      <CommandCancelIconButton onClose={handleClose} />

      <Command.List>
        {usecase.mixedOptions.length === 0 &&
          !usecase.isLoading &&
          !usecase.isValidatingDomain &&
          !usecase.domainValidationError &&
          usecase.searchTerm.length > 0 && (
            <CommandItem
              data-test={'add-org-modal-add-org'}
              leftAccessory={<PlusCircle className='text-primary-700' />}
              onSelect={() => {
                usecase.addNewOrganization();
              }}
            >
              Add {usecase.searchTerm}
            </CommandItem>
          )}

        {usecase.mixedOptions.length > 0 &&
          usecase.mixedOptions.map((option, index) => {
            const isBeingAdded =
              usecase.isAddingOrganization &&
              usecase.addingOrganizationId === option.id;

            return (
              <CommandItem
                className='group'
                key={option?.id || `option-${index}`}
                onSelect={() => {
                  if (option.source === 'tenant') {
                    navigate('/organization/' + option.id);
                    store.ui.commandMenu.toggle('AddNewOrganization');
                    usecase.reset();

                    return;
                  }
                  usecase.addNewOrganization(option);
                }}
                leftAccessory={
                  <Avatar
                    size='xxs'
                    textSize='xxs'
                    name={option.name}
                    variant='outlineSquare'
                    icon={<User03 className='text-primary-700' />}
                    src={
                      (option?.iconUrl ? option.iconUrl : option?.logoUrl) ||
                      undefined
                    }
                  />
                }
                rightAccessory={
                  option.source === 'tenant' ? (
                    <ArrowNarrowRight className='invisible !text-gray-500 group-hover:visible group-data-[selected="true"]:visible' />
                  ) : isBeingAdded ? (
                    <Spinner
                      size='sm'
                      label='adding...'
                      className='!text-gray-300 !fill-gray-500 size-3'
                    />
                  ) : (
                    <PlusCircle className='invisible !text-gray-500 group-hover:visible group-data-[selected="true"]:visible' />
                  )
                }
              >
                <div className='flex items-center gap-2 truncate'>
                  <span className='truncate'>{option?.name || 'Unnamed'}</span>
                  <span>•</span>
                  {option.website && (
                    <span className='truncate'>
                      {getFormattedLink(option.website)}
                    </span>
                  )}
                </div>
              </CommandItem>
            );
          })}
      </Command.List>
    </Command>
  );
});

const getFormattedLink = (url: string): string => {
  return url.replace(/^(https?:\/\/)?(www\.)?([^/?#]+).*/i, '$3');
};
