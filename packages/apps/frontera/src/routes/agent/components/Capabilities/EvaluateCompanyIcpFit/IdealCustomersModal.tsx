import { useRef } from 'react';

import { useKey } from 'rooks';
import { observer } from 'mobx-react-lite';
import { EditIcpDomainsUsecase } from '@domain/usecases/agents/capabilities/edit-icp-domains.usecase';

import { cn } from '@ui/utils/cn';
import { Avatar } from '@ui/media/Avatar';
import { Check } from '@ui/media/icons/Check';
import { Spinner } from '@ui/feedback/Spinner';
import { User03 } from '@ui/media/icons/User03';
import { useModKey } from '@shared/hooks/useModKey';
import { PlusCircle } from '@ui/media/icons/PlusCircle';
import { getFormattedLink } from '@utils/getExternalLink';
import { useOutsideClick } from '@ui/utils/hooks/useOutsideClick';
import {
  Modal,
  ModalBody,
  ModalPortal,
  ModalContent,
  ModalOverlay,
} from '@ui/overlay/Modal';
import {
  Command,
  CommandItem,
  CommandInput,
  CommandCancelIconButton,
} from '@ui/overlay/CommandMenu';

export const IdealCustomersModal = observer(
  ({
    isOpen,
    onClose,
    usecase,
  }: {
    isOpen: boolean;
    onClose: () => void;
    usecase: EditIcpDomainsUsecase;
  }) => {
    const commandRef = useRef(null);

    useOutsideClick({
      ref: commandRef,
      handler: () => {
        onClose();
      },
    });

    useKey('Escape', (e) => {
      e.stopPropagation();
      usecase.reset();
      onClose();
    });

    useModKey('Enter', (e) => {
      e.stopPropagation();
    });

    return (
      <Modal open={isOpen}>
        <ModalPortal>
          <ModalOverlay
            className='z-[5001]'
            onKeyDown={(e) => e.key !== 'Escape' && e.stopPropagation()}
          >
            <ModalBody>
              <ModalContent ref={commandRef}>
                <Command shouldFilter={false}>
                  <CommandInput
                    value={usecase.searchTerm}
                    onValueChange={usecase.setSearchTerm}
                    placeholder='Search by name or website'
                    dataTest={'organizations-create-new-org-org-name'}
                    label={
                      <p className='font-medium'>
                        Search 300,000+ companies ICP
                      </p>
                    }
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
                            usecase.domainValidationError.length > 0 &&
                              'text-error-500',
                          )}
                        >
                          {usecase.domainValidationError ||
                            usecase.domainValidationMessage}
                        </p>
                      )
                    }
                  />

                  <CommandCancelIconButton
                    onClose={() => {
                      usecase.reset();
                      onClose();
                    }}
                  />

                  <Command.List>
                    {usecase.mixedOptions.length === 0 &&
                      !usecase.isLoading &&
                      !usecase.isValidatingDomain &&
                      !usecase.domainValidationError &&
                      usecase.searchTerm.length > 0 && (
                        <CommandItem
                          data-test={'add-org-modal-add-org'}
                          leftAccessory={<Check className='text-primary-700' />}
                          onSelect={() => {
                            usecase.selectCustom(usecase.searchTerm);
                          }}
                          rightAccessory={
                            <PlusCircle className='invisible !text-gray-500 group-hover:visible group-data-[selected="true"]:visible' />
                          }
                        >
                          Add {usecase.searchTerm}
                        </CommandItem>
                      )}

                    {usecase.mixedOptions.length > 0 &&
                      usecase.mixedOptions.map((option, index) => {
                        if (!option.website) return null;

                        return (
                          <CommandItem
                            className='group'
                            key={option?.id || `option-${index}`}
                            onSelect={() => {
                              usecase.select(option?.website ?? '');
                            }}
                            rightAccessory={
                              usecase.icpCompanyExamples.has(option.website) ? (
                                <Check className='text-primary-700' />
                              ) : (
                                <PlusCircle className='invisible !text-gray-500 group-hover:visible group-data-[selected="true"]:visible' />
                              )
                            }
                            leftAccessory={
                              <Avatar
                                size='xxs'
                                textSize='xxs'
                                name={option.name}
                                variant='outlineSquare'
                                icon={<User03 className='text-primary-700' />}
                                src={
                                  (option?.iconUrl
                                    ? option.iconUrl
                                    : option?.logoUrl) || undefined
                                }
                              />
                            }
                          >
                            <div className='flex items-center gap-2 truncate'>
                              <span className='truncate'>
                                {option?.name || 'Unnamed'}
                              </span>
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
              </ModalContent>
            </ModalBody>
          </ModalOverlay>
        </ModalPortal>
      </Modal>
    );
  },
);
