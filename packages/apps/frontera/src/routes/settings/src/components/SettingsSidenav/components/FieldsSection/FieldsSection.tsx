import { CursorBox } from '@ui/media/icons/CursorBox';
import { SidenavItem } from '@shared/components/RootSidenav/components/SidenavItem';

interface WorkspaceSectionProps {
  handleItemClick: (item: string) => () => void;
  checkIsActive: (
    path: string,
    options?: { preset: string | Array<string> },
  ) => boolean;
}

export const FieldsSection = ({
  handleItemClick,
  checkIsActive,
}: WorkspaceSectionProps) => {
  return (
    <div className='flex flex-col gap-2'>
      <div className='flex items-center gap-2 px-3'>
        <CursorBox className='w-4 h-4 text-gray-500' />
        <span className='text-sm  text-gray-700'>Fields</span>
      </div>
      <div className='ml-[23px]'>
        <SidenavItem
          label='Organization'
          isActive={checkIsActive('organizations')}
          onClick={handleItemClick('organizations')}
        />
        {/* <SidenavItem
          label='Contacts'
          isActive={checkIsActive('contacts')}
          onClick={handleItemClick('contacts')}
        /> */}
      </div>
    </div>
  );
};
