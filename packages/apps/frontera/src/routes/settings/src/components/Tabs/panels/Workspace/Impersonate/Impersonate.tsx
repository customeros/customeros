import { useMemo } from 'react';

import { observer } from 'mobx-react-lite';
import { SwitchWorkspaceUsecase } from '@domain/usecases/settings-impersonate/switch-workspace.usecase';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';

export const Impersonate = observer(() => {
  const store = useStore();

  const swithcWorkspaceUsecase = useMemo(() => {
    return new SwitchWorkspaceUsecase();
  }, []);

  return (
    <div className='px-6 pb-4 pt-2 max-w-[500px] border-r border-gray-200 h-full'>
      <Menu>
        <MenuButton>
          <Button colorScheme='primary'>Switch workspace</Button>
        </MenuButton>
        <MenuList side='bottom' align='start'>
          {store.common?.impersonateAccounts?.map((option) => (
            <MenuItem
              key={option?.tenant}
              onClick={() => {
                swithcWorkspaceUsecase.execute(option?.tenant);
              }}
            >
              {option?.tenant}
            </MenuItem>
          ))}
        </MenuList>
      </Menu>
    </div>
  );
});
