import { Search } from './src/components/Search';
import { ProspectsBoard } from './src/components/ProspectsBoard';

export const ProspectsBoardPage = () => {
  return (
    <div className='flex w-full items-start '>
      <div className='w-full '>
        <Search />
        <ProspectsBoard />
      </div>
    </div>
  );
};
