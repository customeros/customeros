import { useLocalStorage } from 'usehooks-ts';

interface PreviewCardProps {
  children: React.ReactNode;
}

export const PreviewCard = ({ children }: PreviewCardProps) => {
  const [previewCard] = useLocalStorage('previewCard', false);

  return (
    <div
      data-state={previewCard ? 'open' : 'closed'}
      className='data-[state=open]:animate-slideLeftAndFade data-[state=closed]:animate-slideRightAndFade flex flex-col max-w-[390px] min-w-[350px] border border-r-0 border-t-0 border-gray-200 bg-white w-[390px]'
    >
      {children}
    </div>
  );
};
