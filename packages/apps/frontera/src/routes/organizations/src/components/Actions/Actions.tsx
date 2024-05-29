import { useState, useEffect } from 'react';

import { Store } from '@store/store';

import { cn } from '@ui/utils/cn.ts';
import { X } from '@ui/media/icons/X.tsx';
import { Button } from '@ui/form/Button/Button';
import { Copy07 } from '@ui/media/icons/Copy07';
import { useStore } from '@shared/hooks/useStore';
import { Archive } from '@ui/media/icons/Archive';
import { UserX01 } from '@ui/media/icons/UserX01.tsx';
import { TableInstance } from '@ui/presentation/Table';
import { HeartHand } from '@ui/media/icons/HeartHand.tsx';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { Organization, OrganizationStage } from '@graphql/types';
import { CoinsStacked01 } from '@ui/media/icons/CoinsStacked01.tsx';
import { ConfirmDeleteDialog } from '@ui/overlay/AlertDialog/ConfirmDeleteDialog/ConfirmDeleteDialog';

interface TableActionsProps {
  tableName?: string;
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
  tableName,
}: TableActionsProps) => {
  const { open: isOpen, onOpen, onClose } = useDisclosure();
  const [targetId, setTargetId] = useState<string | null>(null);

  const selection = table.getState().rowSelection;
  const selectedIds = Object.keys(selection).map(
    (k) => organizationIds?.[parseInt(k)],
  );
  const selectCount = selectedIds.length;
  const store = useStore();

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

  const handleRelationshipUpdate = (option: OrganizationStage) => {
    selectedIds.map((e) => {
      const organization = store.organizations.value.get(e);

      organization?.update((org) => {
        org.stage = option;

        return org;
      });
    });
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

  useEffect(() => {
    clearSelection();
  }, [tableName]);

  if (!selectCount && !targetId) return null;

  const classname =
    'font-normal bg-gray-700 text-gray-25 hover:bg-gray-800 focus:bg-gray-800 hover:text-gray-25 border-none rounded-none';

  return (
    <>
      <div
        className='flex items-center justify-center left-[50%] absolute bottom-[32px]'
        style={{
          left: `calc(50% - 12.5rem)`, // 12.5 is fixed width of sidebar
          right: `calc(50% - 12.5rem)`, // 12.5 is fixed width of sidebar
        }}
      >
        <div className='inline-flex rounded-md shadow-sm' role='group'>
          <Button
            onClick={clearSelection}
            size='sm'
            colorScheme='gray'
            rightIcon={<X className='text-inherit size-3' />}
            className='font-normal bg-gray-700 text-gray-25 hover:bg-gray-800 focus:bg-gray-800 hover:text-gray-25 border-none rounded-none rounded-l-md'
          >
            {`${selectCount} selected`}
          </Button>
          <Button
            onClick={onOpen}
            size='sm'
            colorScheme='gray'
            leftIcon={<Archive className='text-inherit size-3' />}
            className={cn(classname, {
              'rounded-r-md': tableName !== 'Leads' && selectCount === 1,
            })}
          >
            Archive
          </Button>
          {selectCount > 1 && (
            <Button
              colorScheme='gray'
              size='sm'
              leftIcon={<Copy07 className='text-inherit size-3' />}
              onClick={handleMergeOrganizations}
              className={cn(classname, {
                'rounded-r-md': tableName !== 'Leads',
              })}
            >
              Merge
            </Button>
          )}

          {tableName === 'Leads' && (
            <>
              <Button
                colorScheme='gray'
                size='sm'
                leftIcon={<UserX01 className='text-inherit size-3' />}
                onClick={() =>
                  handleRelationshipUpdate(OrganizationStage.Unqualified)
                }
                className={classname}
              >
                Unqualify
              </Button>
              <Button
                colorScheme='gray'
                size='sm'
                leftIcon={<HeartHand className='text-inherit size-3' />}
                onClick={() =>
                  handleRelationshipUpdate(OrganizationStage.Target)
                }
                className={classname}
              >
                Nurture
              </Button>

              <Button
                colorScheme='gray'
                size='sm'
                leftIcon={<CoinsStacked01 className='text-inherit size-3' />}
                onClick={() =>
                  handleRelationshipUpdate(OrganizationStage.Engaged)
                }
                className='font-normal bg-gray-700 text-gray-25 hover:bg-gray-800 focus:bg-gray-800  hover:text-gray-25 border-none rounded-none rounded-r-md'
              >
                Opportunity
              </Button>
            </>
          )}
        </div>
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
