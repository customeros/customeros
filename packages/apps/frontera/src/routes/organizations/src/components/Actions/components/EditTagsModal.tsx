import React, { useRef, useState, useEffect, MouseEvent } from 'react';

import { useKeyBindings } from 'rooks';
import { observer } from 'mobx-react-lite';

import { DataSource } from '@graphql/types';
import { FeaturedIcon } from '@ui/media/Icon';
import { Tag01 } from '@ui/media/icons/Tag01.tsx';
import { useStore } from '@shared/hooks/useStore';
import { Button } from '@ui/form/Button/Button.tsx';
import { Tags } from '@organization/components/Tabs';
import { Spinner } from '@ui/feedback/Spinner/Spinner.tsx';
import {
  AlertDialog,
  AlertDialogBody,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogPortal,
  AlertDialogContent,
  AlertDialogOverlay,
  AlertDialogCloseButton,
  AlertDialogConfirmButton,
} from '@ui/overlay/AlertDialog/AlertDialog.tsx';

interface ConfirmDeleteDialogProps {
  isOpen: boolean;
  onClose: () => void;
  selectedIds: string[];
  clearSelection: () => void;
}

export const EditTagsModal = observer(
  ({
    isOpen,
    onClose,
    selectedIds,
    clearSelection,
  }: ConfirmDeleteDialogProps) => {
    const store = useStore();
    const confirmRef = useRef<HTMLButtonElement>(null);
    const [url, setUrl] = useState('');
    const [selectedTags, setSelectedTags] = useState<
      Array<{ label: string; value: string }>
    >([]);

    useEffect(() => {
      store.ui.setIsEditingTableCell(isOpen);

      if (isOpen && !url.includes('linkedin.com')) {
        setUrl('');
      }
    }, [isOpen]);

    const handleClose = () => {
      setSelectedTags([]);
      onClose();
    };

    const handleConfirm = (
      e: MouseEvent<HTMLButtonElement> | KeyboardEvent,
    ) => {
      e.preventDefault();
      e.stopPropagation();

      const tags = selectedTags.map((e) => ({
        name: e.label,
        id: e.value,
        appSource: 'organization',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        source: DataSource.Openline,
        metadata: {
          id: e.value,
          source: DataSource.Openline,
          sourceOfTruth: DataSource.Openline,
          appSource: 'organization',
          created: new Date().toISOString(),
          lastUpdated: new Date().toISOString(),
        },
      }));

      store.organizations.updateTags(selectedIds, tags);
      handleClose();
      clearSelection();
    };

    useKeyBindings(
      {
        Escape: handleClose,
        Enter: handleConfirm,
      },
      { when: isOpen },
    );

    const selectCount = selectedIds.length;

    return (
      <AlertDialog isOpen={isOpen} onClose={handleClose} className='z-[99999] '>
        <AlertDialogPortal>
          <AlertDialogOverlay>
            <AlertDialogContent className='rounded-xl bg-no-repeat bg-[url(/backgrounds/organization/circular-bg-pattern.png)]'>
              <FeaturedIcon
                size='lg'
                colorScheme={'primary'}
                className='mt-[13px] ml-[11px]'
              >
                <Tag01 />
              </FeaturedIcon>
              <AlertDialogHeader className='text-lg font-bold mt-4'>
                <p className='pb-0 font-semibold'>
                  Add tags to {selectCount}
                  {selectCount === 1 ? ' organization' : ' organizations'}?
                </p>
                <p className='mt-1 mb-2 text-sm text-gray-700 font-normal'>
                  What tags would you like to add to your selected
                  organizations?{' '}
                </p>
              </AlertDialogHeader>
              <AlertDialogBody>
                <Tags
                  autofocus
                  icon={null}
                  placeholder='Tags'
                  value={selectedTags}
                  closeMenuOnSelect={true}
                  onChange={(e) => setSelectedTags(e)}
                />
              </AlertDialogBody>
              <AlertDialogFooter>
                <AlertDialogCloseButton>
                  <Button
                    size='md'
                    variant='outline'
                    colorScheme={'gray'}
                    className='bg-white w-full'
                  >
                    Cancel
                  </Button>
                </AlertDialogCloseButton>
                <AlertDialogConfirmButton>
                  <Button
                    size='md'
                    ref={confirmRef}
                    variant='outline'
                    className='w-full'
                    colorScheme={'primary'}
                    onClick={handleConfirm}
                    loadingText='Adding tags'
                    isLoading={store.organizations.isLoading}
                    spinner={
                      <Spinner
                        size={'sm'}
                        label='deleting'
                        className='text-error-300 fill-error-700'
                      />
                    }
                  >
                    Add tags
                  </Button>
                </AlertDialogConfirmButton>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialogOverlay>
        </AlertDialogPortal>
      </AlertDialog>
    );
  },
);
