import { useState, useEffect } from 'react';

import { Store } from '@store/store';

import { Organization } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { Copy07 } from '@ui/media/icons/Copy07';
import { Archive } from '@ui/media/icons/Archive';
// import { ButtonGroup } from '@ui/form/ButtonGroup';
import { TableInstance } from '@ui/presentation/Table';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog/ConfirmDeleteDialog';

interface TableActionsProps {
  organizationIds: string[];
  onHide: (ids: string[]) => void;
  table: TableInstance<Store<Organization>>;
  onMerge: (primaryId: string, mergeIds: string[]) => void;
}

export const TableActions = ({
  table,
  organizationIds,
  onHide,
  onMerge,
}: TableActionsProps) => {
  const { open: isOpen, onOpen, onClose } = useDisclosure();
  const [targetId, setTargetId] = useState<string | null>(null);

  const selection = table.getState().rowSelection;
  const selectedIds = Object.keys(selection).map(
    (k) => organizationIds?.[parseInt(k)],
  );
  const selectCount = selectedIds.length;

  const clearSelection = () => table.resetRowSelection();

  const handleMergeOrganizations = () => {
    const mergeIds = selectedIds.filter((id) => id !== targetId);

    if (!targetId || !mergeIds.length) return;

    onMerge(targetId, mergeIds);
    clearSelection();
  };

  const handleHideOrganizations = () => {
    onHide(selectedIds);
    onClose();
    clearSelection();
  };

  useEffect(() => {
    if (selectCount === 1) {
      setTargetId(selectedIds[0]);
    }
    if (selectCount < 1) {
      setTargetId(null);
    }
  }, [selectCount]);

  if (!selectCount && !targetId) return null;

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
        loadingButtonLabel='Archiving'
        label={`Archive selected ${
          selectCount === 1 ? 'organization' : 'organizations'
        }?`}
      />
    </>
  );
};
