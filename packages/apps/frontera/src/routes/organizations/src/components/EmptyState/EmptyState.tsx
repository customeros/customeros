import { useNavigate, useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { EmptyTable } from '@ui/media/logos/EmptyTable';

import HalfCirclePattern from '../../../../src/assets/HalfCirclePattern';

export const EmptyState = observer(() => {
  const store = useStore();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const preset = searchParams?.get('preset');
  const isAllOrgsTable =
    store.tableViewDefs?.toArray().find((e) => e.value.name === 'All orgs')
      ?.value?.id === preset;
  const handleCreateOrganization = () => {
    store.organizations.create();
  };

  const options =
    !preset || isAllOrgsTable
      ? {
          title: "Let's get started",
          description:
            'Start seeing your customer conversations all in one place by adding an organization',
          buttonLabel: 'Add Organization',
          onClick: handleCreateOrganization,
        }
      : preset === 'portfolio'
      ? {
          title: 'No organizations assigned to you yet',
          description:
            'Currently, you have not been assigned to any organizations.\n' +
            '\n' +
            'Head to your list of organizations and assign yourself as an owner to one of them.',
          buttonLabel: 'Go to Organizations',
          onClick: () => {
            navigate(`/finder`);
          },
        }
      : {
          title: 'No organizations created yet',
          description:
            'Currently, there are no organizations created yet.\n' +
            '\n' +
            'Head to your list of organizations and create one.',
          buttonLabel: 'Go to Organizations',
          onClick: () => {
            navigate(`/finder`);
          },
        };

  return (
    <div className='flex items-center justify-center h-full bg-white'>
      <div className='flex flex-col h-[500px] w-[500px]'>
        <div className='flex relative'>
          <EmptyTable className='w-[152px] h-[120px] absolute top-[25%] right-[35%]' />
          <HalfCirclePattern height={500} width={500} />
        </div>
        <div className='flex flex-col text-center items-center top-[5vh] transform translate-y-[-230px]'>
          <p className='text-gray-900 text-md font-semibold'>{options.title}</p>
          <p className='max-w-[400px] text-sm text-gray-600 my-1'>
            {options.description}
          </p>

          <Button
            onClick={options.onClick}
            className='mt-2 min-w-min text-sm'
            variant='outline'
          >
            {options.buttonLabel}
          </Button>
        </div>
      </div>
    </div>
  );
});
