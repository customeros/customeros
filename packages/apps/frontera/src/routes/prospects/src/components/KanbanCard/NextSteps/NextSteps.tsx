import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { Textarea } from '@ui/form/Textarea/Textarea';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { ArrowsRight } from '@ui/media/icons/ArrowsRight';

interface NextStepsProps {
  opportunityId: string;
  onToggle: (value: boolean) => void;
  textareaRef: React.RefObject<HTMLTextAreaElement>;
}

export const NextSteps = observer(
  ({ textareaRef, onToggle, opportunityId }: NextStepsProps) => {
    const store = useStore();
    const opportunity = store.opportunities.value.get(opportunityId);
    const value = opportunity?.value.nextSteps;

    const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
      const nextValue = e.target.value;

      opportunity?.update(
        (value) => {
          value.nextSteps = nextValue;

          return value;
        },
        { mutate: !nextValue },
      );
    };

    const handleBlur = () => {
      !value && onToggle(false);
      opportunity?.saveProperty('nextSteps');
    };

    return (
      <Tooltip side='top' align='start' label='Next steps'>
        <div className='flex gap-2 w-full items-start justify-start'>
          <ArrowsRight className='size-4 ml-1 mr-0.5 mt-0.5 text-gray-500' />
          <Textarea
            size='xs'
            value={value}
            ref={textareaRef}
            variant='unstyled'
            onBlur={handleBlur}
            className='leading-5'
            onChange={handleChange}
            placeholder="What's your next step?"
          />
        </div>
      </Tooltip>
    );
  },
);
