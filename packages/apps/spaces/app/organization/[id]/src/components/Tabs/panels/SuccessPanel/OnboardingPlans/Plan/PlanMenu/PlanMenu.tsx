import { useRouter } from 'next/navigation';

import { Map01 } from '@ui/media/icons/Map01';
import { IconButton } from '@ui/form/IconButton';
import { Archive } from '@ui/media/icons/Archive';
import { PlusSquare } from '@ui/media/icons/PlusSquare';
import { DotsVertical } from '@ui/media/icons/DotsVertical';
import {
  Menu,
  MenuItem,
  MenuList,
  MenuButton,
  MenuDivider,
} from '@ui/overlay/Menu';

interface PlanMenuProps {
  id: string;
  isOpen: boolean;
  onOpen: () => void;
  onClose: () => void;
  onRemovePlan: () => void;
  onAddMilestone: () => void;
}

export const PlanMenu = ({
  id,
  isOpen,
  onOpen,
  onClose,
  onRemovePlan,
  onAddMilestone,
}: PlanMenuProps) => {
  const router = useRouter();

  const handleEditMasterPlan = () => {
    router.push(`/settings?tab=master-plans&show=active&planId=${id}`);
  };

  return (
    <Menu isOpen={isOpen} onClose={onClose} placement='bottom-end'>
      <MenuButton
        size='xs'
        as={IconButton}
        variant='ghost'
        color='gray.500'
        onClick={onOpen}
        fontWeight='normal'
        icon={<DotsVertical color='gray.400' />}
      />
      <MenuList minW='12rem'>
        <MenuItem
          icon={<PlusSquare color='gray.500' />}
          onClick={onAddMilestone}
        >
          Add milestone
        </MenuItem>
        <MenuItem icon={<Archive color='gray.500' />} onClick={onRemovePlan}>
          Archive
        </MenuItem>
        <MenuDivider
          mx='2'
          borderBottom='unset'
          borderTop='1px dashed'
          borderColor='gray.300'
        />
        <MenuItem
          icon={<Map01 color='gray.500' />}
          onClick={handleEditMasterPlan}
        >
          Edit master plan
        </MenuItem>
      </MenuList>
    </Menu>
  );
};
