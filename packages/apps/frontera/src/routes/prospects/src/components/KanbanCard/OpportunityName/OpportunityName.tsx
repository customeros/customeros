import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { ResizableInput } from '@ui/form/Input/ResizableInput';

interface OpportunityNameProps {
  opportunityId: string;
}

export const OpportunityName = observer(
  ({ opportunityId }: OpportunityNameProps) => {
    const store = useStore();
    const opportunity = store.opportunities.value.get(opportunityId);
    const value = opportunity?.value.name;

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const nextValue = e.target.value;

      opportunity?.update(
        (value) => {
          value.name = nextValue;

          return value;
        },
        { mutate: !nextValue },
      );
    };

    const handleBlur = () => {
      opportunity?.saveProperty('name');
    };

    return (
      <ResizableInput
        size='xs'
        value={value}
        variant='unstyled'
        onBlur={handleBlur}
        onChange={handleChange}
        placeholder='Unnamed opportunity'
        onClick={(e) => (e.target as HTMLInputElement).select()}
        className='font-medium line-clamp-1 max-w-[178px] text-ellipsis'
      />
    );
  },
);
