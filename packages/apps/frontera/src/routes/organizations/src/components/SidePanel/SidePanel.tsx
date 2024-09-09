import { useSearchParams } from 'react-router-dom';

import { useStore } from '@shared/hooks/useStore';

import { Icp } from '../Workflows/ICP';

export const SidePanel = () => {
  const store = useStore();

  const [searchParams] = useSearchParams();

  const preset = searchParams.get('preset');

  const tableViewName = store.tableViewDefs.getById(preset || '')?.value.name;

  return (
    <div className='flex'>
      {/* {tableViewName === 'Contacts' && (
        <div className='min-w-[200px] bg-white border-l border-t flex flex-col py-4 px-2'>
          <PersonasFlowProfileMenu />
        </div>
      )} */}

      <div className='min-w-[525px] w-[550px] bg-white  py-4 px-6 flex flex-col h-[100vh] border-t border-l animate-slideLeft'>
        {tableViewName === 'Targets' && <Icp />}
        {/* {tableViewName === 'Contacts' && <PersonasFlowProfile />} */}
      </div>
    </div>
  );
};
