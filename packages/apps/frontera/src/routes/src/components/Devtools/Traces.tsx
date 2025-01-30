import { observer } from 'mobx-react-lite';
import { Tracer, TraceRecord } from '@infra/tracer';

import { Icon } from '@ui/media/Icon';
import { IconButton } from '@ui/form/IconButton';

const Trace = observer(({ trace }: { trace: TraceRecord }) => {
  return (
    <div className='ml-4 border-l pl-0'>
      <div
        className='flex flex-col cursor-pointer'
        onClick={() => Tracer.toggleExpand(trace.id)}
      >
        <div className='flex items-center justify-between hover:bg-grayModern-50'>
          <div className='flex items-center gap-1'>
            <IconButton
              size='xxs'
              variant='ghost'
              aria-label='toggle expand'
              onClick={() => Tracer.toggleExpand(trace.id)}
              icon={
                <Icon
                  name={
                    Tracer.expandedTraces.has(trace.id)
                      ? 'chevron-down'
                      : 'chevron-right'
                  }
                />
              }
            />
            <span className='text-xs font-medium'>{trace.name}</span>
          </div>
          <span className='text-xs text-grayModern-500'>
            {trace.duration ?? '~'}ms
          </span>
        </div>
        {Tracer.expandedTraces.has(trace.id) && (
          <div>
            {trace?.attributes && (
              <div className='pl-5 bg-primary-50'>
                <p className='text-xs font-medium text-primary-700'>
                  Attributes:
                </p>
                <pre className='text-xs text-primary-700'>
                  {JSON.stringify(trace.attributes, null, 2)}
                </pre>
              </div>
            )}
            {!!trace?.result && (
              <div className='pl-5 bg-greenLight-50'>
                <p className='text-xs font-medium text-greenLight-700'>
                  Result:
                </p>
                <pre className='text-xs text-greenLight-700'>
                  {JSON.stringify(trace?.result, null, 2)}
                </pre>
              </div>
            )}
          </div>
        )}
      </div>
      {Tracer.expandedTraces.has(trace.id) && trace.children.length > 0 && (
        <div className='ml-4'>
          {trace.children.map((t) => (
            <Trace trace={t} key={t.id} />
          ))}
        </div>
      )}
    </div>
  );
});

export const Traces = observer(() => {
  return (
    <div className='w-full'>
      <p className='font-medium underline mb-2'>Traces</p>
      {Tracer.traces
        .filter((trace) => trace.parentId === null)
        .map((trace) => (
          <Trace trace={trace} key={trace.id} />
        ))}
    </div>
  );
});
