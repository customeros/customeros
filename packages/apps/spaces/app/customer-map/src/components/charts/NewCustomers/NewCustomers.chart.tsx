'use client';

import { PatternLines } from '@visx/pattern';
import { LinearGradient } from '@visx/gradient';
import { XYChart, Tooltip, BarSeries, AnimatedAxis } from '@visx/xychart';

import { useToken } from '@ui/utils';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';

import { getMonthLabel } from '../util';

export type NewCustomersDatum = {
  month: number;
  value: number;
};

const _mockData: NewCustomersDatum[] = [
  {
    month: 1,
    value: 20,
  },
  {
    month: 2,
    value: 50,
  },
  {
    month: 3,
    value: 68,
  },
  {
    month: 4,
    value: 55,
  },
  {
    month: 5,
    value: 73,
  },
  {
    month: 6,
    value: 69,
  },
  {
    month: 7,
    value: 84,
  },
  {
    month: 8,
    value: 85,
  },
  {
    month: 9,
    value: 80,
  },
  {
    month: 10,
    value: 87,
  },
  {
    month: 11,
    value: 90,
  },
  {
    month: 12,
    value: 95,
  },
];

interface NewCustomersProps {
  width?: number;
  height?: number;
  data: NewCustomersDatum[];
}

const getX = (d: NewCustomersDatum) => getMonthLabel(d.month);

const NewCustomersChart = ({ width, data }: NewCustomersProps) => {
  const [primary600, gray700] = useToken('colors', ['primary.600', 'gray.700']);

  return (
    <>
      <Flex h='24px' />
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
          from={primary600}
          id='visx-area-gradient'
        />
        <PatternLines
          id='stripes'
          height={8}
          width={8}
          stroke={primary600}
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
            i === data.length - 1 ? 'url(#stripes)' : primary600
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
            background: gray700,
            borderRadius: '8px',
          }}
          offsetTop={-50}
          offsetLeft={-30}
          renderTooltip={({ tooltipData }) => {
            const xLabel = getX(
              tooltipData?.nearestDatum?.datum as NewCustomersDatum,
            );
            const value = (
              tooltipData?.nearestDatum?.datum as NewCustomersDatum
            ).value;

            return (
              <Flex flexDir='column'>
                <Text fontSize='sm' color='white' fontWeight='normal'>
                  {xLabel}
                  {': '}
                  {value}
                </Text>
              </Flex>
            );
          }}
        />
      </XYChart>
    </>
  );
};

export default NewCustomersChart;
