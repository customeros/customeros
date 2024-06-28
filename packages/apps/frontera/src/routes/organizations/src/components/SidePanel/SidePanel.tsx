import { Icp } from '../ICP';

interface SidePanelProps {
  children?: React.ReactNode;
}

export const SidePanel = ({ children }: SidePanelProps) => {
  return (
    <div className='min-w-[600px] w-[600px] bg-white  py-4 px-6 flex flex-col h-[100vh] border-t border-l animate-slideLeftandFade '>
      <Icp />
    </div>
  );
};
