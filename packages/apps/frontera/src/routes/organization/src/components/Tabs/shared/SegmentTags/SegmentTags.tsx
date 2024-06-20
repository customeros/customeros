import { observer } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { PieChart02 } from '@ui/media/icons/PieChart02';
import { MultiCreatableSelect } from '@ui/form/MultiCreatableSelect';

export const SegmentTags = observer(() => {
  const store = useStore();

  const options = store.tags?.toArray().map((tag) => ({
    value: tag.value.id,
    label: tag.value.name,
  }));

  return (
    <div className='flex items-center gap-2'>
      <PieChart02 className='text-gray-500 mt-1' />
      <MultiCreatableSelect
        className=''
        defaultOptions={options || []}
        placeholder='Segment tags'
        classNames={{
          multiValueLabel: () =>
            'bg-grayModern-100 rounded-s-md ps-1 pe-1 cursor-pointer',
          multiValueRemove: () =>
            'bg-grayModern-100 hover:bg-grayModern-200 ps-0.5 rounded-e-md pe-0.5 text-grayModern-400 hover:text-warning-700',
        }}
        backspaceRemovesValue
        isMulti
        onCreateOption={(value) => {
          store.tags?.create(undefined, {
            onSucces: (serverId) => {
              store.tags?.value.get(serverId)?.update((tag) => {
                tag.name = value;

                return tag;
              });
            },
          });
        }}
        // value={store.organization.value.segments.map((segment) => ({
        //   value: segment.id,
        //   label: segment.name,
        // }))}
        // onChange={(e) => {
        //   store.organization.value.update((org) => {
        //     org.segments = e.map((tag) => ({
        //       id: tag.value,
        //       name: tag.label,
        //     }));

        //     return org;
        //   });
        // }}
      />
    </div>
  );
});
