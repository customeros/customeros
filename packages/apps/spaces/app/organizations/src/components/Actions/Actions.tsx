import { memo, useState, useEffect } from 'react';

import { Button } from '@ui/form/Button';
import { useDisclosure } from '@ui/utils';
import { Center } from '@ui/layout/Center';
import { Organization } from '@graphql/types';
import { Copy07 } from '@ui/media/icons/Copy07';
import { Archive } from '@ui/media/icons/Archive';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { TableInstance, RowSelectionState } from '@ui/presentation/Table';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog';

interface TableActionsProps {
  isArchiving: boolean;
  table: TableInstance<Organization>;
  onMergeOrganizations: (
    targetIndex: string,
    selection: RowSelectionState,
  ) => void;
  onArchiveOrganizations: (selection: RowSelectionState) => void;
}

export const TableActions = memo(
  ({
    table,
    isArchiving,
    onMergeOrganizations,
    onArchiveOrganizations,
  }: TableActionsProps) => {
    const { isOpen, onOpen, onClose } = useDisclosure();
    const [targetIndex, setTargetIndex] = useState<string | null>(null);

    const selection = table.getState().rowSelection;
    const selectCount = Object.keys(selection).length;

    useEffect(() => {
      if (selectCount === 1) {
        const [index] = Object.entries(selection)[0];
        setTargetIndex(index);
      }
      if (selectCount < 1) {
        setTargetIndex(null);
      }
    }, [selectCount]);

    if (!selectCount && !targetIndex) return null;

    return (
      <>
        <Center left='50%' position='absolute' bottom='32px'>
          <ButtonGroup size='md' isAttached left='-50%' position='relative'>
            <Button
              bg='gray.700'
              color='white'
              leftIcon={<Archive />}
              onClick={onOpen}
              _hover={{
                bg: 'gray.800',
              }}
            >
              {`Archive ${
                selectCount > 1 ? `these ${selectCount}` : ' this company'
              }`}
            </Button>
            {selectCount > 1 && (
              <Button
                bg='gray.700'
                color='white'
                leftIcon={<Copy07 />}
                _hover={{
                  bg: 'gray.800',
                }}
                onClick={() => {
                  onMergeOrganizations(targetIndex as string, selection);
                  table.resetRowSelection();
                }}
              >
                {`Merge these ${selectCount}`}
              </Button>
            )}
          </ButtonGroup>
        </Center>

        <ConfirmDeleteDialog
          isOpen={isOpen}
          icon={<Archive />}
          onClose={onClose}
          isLoading={isArchiving}
          confirmButtonLabel={'Archive'}
          onConfirm={() => {
            onArchiveOrganizations(selection);
            onClose();
            table.resetRowSelection();
          }}
          label={`Archive selected ${
            selectCount === 1 ? 'organization' : 'organizations'
          }?`}
        />
      </>
    );
  },
  (prev, next) => {
    const prevCount = Object.keys(prev.table.getState().rowSelection).length;
    const nextCount = Object.keys(next.table.getState().rowSelection).length;
    return prevCount !== nextCount;
  },
);
