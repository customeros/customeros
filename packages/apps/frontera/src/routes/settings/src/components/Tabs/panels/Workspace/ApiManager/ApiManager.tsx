import { observer } from 'mobx-react-lite';

import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Copy03 } from '@ui/media/icons/Copy03.tsx';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip.tsx';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';

function convertUuid(uuid: string): string {
  return `${uuid.slice(0, 4)}${'•'.repeat(32)}${uuid.slice(-4)}`;
}

export const ApiManager = observer(() => {
  const store = useStore();
  const [_, copyToClipboard] = useCopyToClipboard();

  const apiKey = store.settings.tenantApiKey;

  const formattedKey = apiKey?.length ? convertUuid(apiKey) : '';

  return (
    <article className='px-6 pb-4 max-w-[500px] h-full overflow-y-auto  border-r border-gray-200'>
      <div className='flex flex-col '>
        <div className='flex justify-between items-center pt-2 sticky top-0 bg-gray-25 '>
          <h1 className='text-gray-700 font-semibold text-base'>API</h1>
        </div>
        <p className='mb-4 text-sm'>
          Create an API key to unlock CustomerOS’s API and start building your
          own custom applications
        </p>
      </div>

      <div
        role='button'
        onClick={() =>
          copyToClipboard(apiKey, 'API key copied to your clipboard')
        }
        className='py-1 min-h-[32px] max-h-[32px] mb-1 border text-sm rounded-md border-gray-200 flex justify-between items-center'
      >
        <div className='flex-grow mx-3'>
          <p>{formattedKey}</p>
        </div>

        <Tooltip align='end' side='bottom' label={'Copy API key'}>
          <IconButton
            size='xxs'
            variant='ghost'
            aria-label='Copy API key'
            className='mr-1.5 rounded'
            icon={<Copy03 className='text-xs' />}
            onClick={() =>
              copyToClipboard(apiKey, 'API key copied to your clipboard')
            }
          />
        </Tooltip>
      </div>
    </article>
  );
});
