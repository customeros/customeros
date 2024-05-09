import { useState, useEffect } from 'react';

import { Organization } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { Copy07 } from '@ui/media/icons/Copy07';
import { Archive } from '@ui/media/icons/Archive';
// import { ButtonGroup } from '@ui/form/ButtonGroup';
import { TableInstance } from '@ui/presentation/Table';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog/ConfirmDeleteDialog';

import { useOrganizationsPageMethods } from '../../hooks/useOrganizationsPageMethods';

interface TableActionsProps {
  allOrganizationsIds: string[];
  table: TableInstance<Organization>;
}

export const TableActions = ({
  table,
  allOrganizationsIds,
}: TableActionsProps) => {
  const { open: isOpen, onOpen, onClose } = useDisclosure();
  const [targetIndex, setTargetIndex] = useState<string | null>(null);
  const { hideOrganizations, mergeOrganizations } =
    useOrganizationsPageMethods();

  const selection = table.getState().rowSelection;
  const selectCount = Object.keys(selection).length;

  const handleMergeOrganizations = () => {
    const primaryId = (allOrganizationsIds as string[])[Number(targetIndex)];
    const selectedIds = Object.keys(selection).map(
      (k) => (allOrganizationsIds as string[])[Number(k)],
    );
    const mergeIds = selectedIds.filter((id) => id !== primaryId);

    if (!primaryId || !mergeIds.length) return;

    mergeOrganizations.mutate(
      {
        primaryOrganizationId: primaryId,
        mergedOrganizationIds: mergeIds,
      },
      {
        onSuccess: () => {
          table.resetRowSelection();
        },
      },
    );
  };

  const handleHideOrganizations = () => {
    const selectedIds = Object.keys(selection)
      .map((k) => (allOrganizationsIds as string[])[Number(k)])
      .filter(Boolean);

    hideOrganizations.mutate(
      {
        ids: selectedIds,
      },
      {
        onSuccess: () => {
          onClose();
          table.resetRowSelection();
        },
      },
    );
  };

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
      <div
        className='flex items-center justify-center left-[50%] absolute bottom-[32px]'
        style={{
          left: `calc(50% - 12.5rem)`, // 12.5 is fixed width of sidebar
        }}
      >
        {/* <ButtonGroup size='md' isAttached left='-50%' position='relative'> */}
        <Button
          onClick={onOpen}
          colorScheme='gray'
          leftIcon={<Archive className='text-inherit' />}
          className='bg-gray-700 text-gray-25 hover:bg-gray-800 hover:text-gray-25'
        >
          {`Archive ${
            selectCount > 1 ? `these ${selectCount}` : ' this organization'
          }`}
        </Button>
        {selectCount > 1 && (
          <Button
            colorScheme='gray'
            leftIcon={<Copy07 className='text-inherit' />}
            onClick={handleMergeOrganizations}
            className='bg-gray-700 text-gray-25 hover:bg-gray-800 hover:text-gray-25'
          >
            {`Merge these ${selectCount}`}
          </Button>
        )}
        {/* </ButtonGroup> */}
      </div>

      <ConfirmDeleteDialog
        isOpen={isOpen}
        icon={<Archive />}
        onClose={onClose}
        confirmButtonLabel={'Archive'}
        onConfirm={handleHideOrganizations}
        isLoading={hideOrganizations.isPending}
        loadingButtonLabel='Archiving'
        label={`Archive selected ${
          selectCount === 1 ? 'organization' : 'organizations'
        }?`}
      />
    </>
  );
};
