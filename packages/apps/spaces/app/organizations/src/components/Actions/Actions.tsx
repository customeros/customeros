import { Button } from '@ui/form/Button';
import { Organization } from '@graphql/types';
import { Copy07 } from '@ui/media/icons/Copy07';
import { Archive } from '@ui/media/icons/Archive';
import { ButtonGroup } from '@ui/form/ButtonGroup';
import { TableInstance, RowSelectionState } from '@ui/presentation/Table';

interface TableActionsProps {
  table: TableInstance<Organization>;
  selection: RowSelectionState;
  onMergeOrganizations: () => void;
  onArchiveOrganizations: () => void;
}

export const TableActions = ({
  table,
  selection,
  onMergeOrganizations,
  onArchiveOrganizations,
}: TableActionsProps) => {
  const selectCount = Object.keys(selection).length;

  if (!selectCount) return null;

  return (
    <ButtonGroup size='md' isAttached left='-50%' position='relative'>
      <Button
        bg='gray.700'
        color='white'
        leftIcon={<Archive />}
        onClick={onArchiveOrganizations}
        // borderRight='1px solid'
        // borderRightColor='gray.500'
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
            onMergeOrganizations();
            table.resetRowSelection();
          }}
        >
          {`Merge these ${selectCount}`}
        </Button>
      )}
    </ButtonGroup>
  );
};
