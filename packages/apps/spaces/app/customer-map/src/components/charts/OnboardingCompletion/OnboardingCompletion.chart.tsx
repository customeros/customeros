'use client';

import { PatternLines } from '@visx/pattern';
import { LinearGradient } from '@visx/gradient';
import { XYChart, Tooltip, BarSeries, AnimatedAxis } from '@visx/xychart';

import { mockData } from './mock';
import { getMonthLabel } from '../util';

export type OnboardingCompletionDatum = {
  month: number;
  value: number;
};

interface OnboardingCompletionProps {
  width?: number;
  height?: number;
  hasContracts?: boolean;
  data: OnboardingCompletionDatum[];
}

const getX = (d: OnboardingCompletionDatum) => getMonthLabel(d.month);

const OnboardingCompletionChart = ({
  width,
  data: _data,
  hasContracts,
}: OnboardingCompletionProps) => {
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
        margin={{ top: 12, right: 0, bottom: 20, left: 0 }}
        xScale={{
          type: 'band',
          paddingInner: 0.4,
          paddingOuter: 0.4,
        }}
        yScale={{ type: 'linear' }}
      >
        <LinearGradient
          fromOpacity={0}
          toOpacity={0.3}
          to={'white'}
          from={colors.primary600}
          id='visx-area-gradient'
        />
        <PatternLines
          id='stripes'
          height={8}
          width={8}
          stroke={colors.primary600}
          strokeWidth={2}
          orientation={['diagonal']}
        />
        <BarSeries
          dataKey='Newly contracted'
          radius={4}
          data={data}
          radiusAll
          xAccessor={(d) => getX(d)}
          yAccessor={(d) => d.value}
          colorAccessor={(_, i) =>
            i === data.length - 1 ? 'url(#stripes)' : colors.primary600
          }
        />

        <AnimatedAxis
          orientation='bottom'
          hideAxisLine
          hideTicks
          tickLabelProps={{
            fontSize: 12,
            fontWeight: 'medium',
            fontFamily: `var(--font-barlow)`,
          }}
        />
        <Tooltip
          key={Math.random()}
          snapTooltipToDatumY
          snapTooltipToDatumX
          style={{
            position: 'absolute',
            padding: '8px',
            background: colors.gray700,
            borderRadius: '8px',
          }}
          offsetTop={-50}
          offsetLeft={-30}
          renderTooltip={({ tooltipData }) => {
            const xLabel = getX(
              tooltipData?.nearestDatum?.datum as OnboardingCompletionDatum,
            );
            const value = (
              tooltipData?.nearestDatum?.datum as OnboardingCompletionDatum
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

export default OnboardingCompletionChart;
