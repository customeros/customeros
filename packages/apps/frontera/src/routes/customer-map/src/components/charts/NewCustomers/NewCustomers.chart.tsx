import { PatternLines } from '@visx/pattern';
import { LinearGradient } from '@visx/gradient';
import { XYChart, Tooltip, BarSeries, AnimatedAxis } from '@visx/xychart';

import { mockData } from './mock';
import { getMonthLabel } from '../util';

export type NewCustomersDatum = {
  month: number;
  value: number;
};

interface NewCustomersProps {
  width?: number;
  height?: number;
  hasContracts?: boolean;
  data: NewCustomersDatum[];
}

const getX = (d: NewCustomersDatum) => getMonthLabel(d.month);

const NewCustomersChart = ({
  width,
  data: _data,
  hasContracts,
}: NewCustomersProps) => {
  const data = hasContracts ? _data : mockData;

  const colors = {
    primary600: hasContracts ? '#7F56D9' : '#EAECF0',
    gray700: '#344054',
  };

  return (
    <>
      <div className='flex h-6' />
      <XYChart
        height={200}
        width={width || 500}
        yScale={{ type: 'linear' }}
        margin={{ top: 12, right: 0, bottom: 20, left: 0 }}
        xScale={{
          type: 'band',
          paddingInner: 0.4,
          paddingOuter: 0.4,
        }}
      >
        <LinearGradient
          to={'white'}
          fromOpacity={0}
          toOpacity={0.3}
          id='visx-area-gradient'
          from={colors.primary600}
        />
        <PatternLines
          width={8}
          height={8}
          id='stripes'
          strokeWidth={2}
          stroke={colors.primary600}
          orientation={['diagonal']}
        />
        <BarSeries
          radiusAll
          radius={4}
          data={data}
          dataKey='Newly contracted'
          xAccessor={(d) => getX(d)}
          yAccessor={(d) => d.value}
          colorAccessor={(_, i) =>
            i === data.length - 1 ? 'url(#stripes)' : colors.primary600
          }
        />

        <AnimatedAxis
          hideTicks
          hideAxisLine
          orientation='bottom'
          tickLabelProps={{
            fontSize: 12,
            fontWeight: 'medium',
            fontFamily: `var(--font-ibm-plex-sans)`,
          }}
        />
        <Tooltip
          offsetTop={-50}
          offsetLeft={-30}
          key={Math.random()}
          snapTooltipToDatumY
          snapTooltipToDatumX
          style={{
            position: 'absolute',
            padding: '8px',
            background: colors.gray700,
            borderRadius: '8px',
          }}
          renderTooltip={({ tooltipData }) => {
            const xLabel = getX(
              tooltipData?.nearestDatum?.datum as NewCustomersDatum,
            );
            const value = (
              tooltipData?.nearestDatum?.datum as NewCustomersDatum
            ).value;

            return (
              <div className='flex flex-col'>
                <p className='text-sm text-white font-normal'>
                  {xLabel}
                  {': '}
                  {hasContracts ? value : 'No data yet'}
                </p>
              </div>
            );
          }}
        />
      </XYChart>
    </>
  );
};

export default NewCustomersChart;
