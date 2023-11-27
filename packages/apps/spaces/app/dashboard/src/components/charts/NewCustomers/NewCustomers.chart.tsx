'use client';

import { LinearGradient } from '@visx/gradient';
import { XYChart, Tooltip, BarSeries, AnimatedAxis } from '@visx/xychart';

import { useToken } from '@ui/utils';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';

type NewCustomersDatum = {
  x: string;
  value: number;
};

const mockData: NewCustomersDatum[] = [
  {
    x: 'Jan',
    value: 20,
  },
  {
    x: 'Feb',
    value: 50,
  },
  {
    x: 'Mar',
    value: 68,
  },
  {
    x: 'Apr',
    value: 55,
  },
  {
    x: 'May',
    value: 73,
  },
  {
    x: 'Jun',
    value: 69,
  },
  {
    x: 'Jul',
    value: 84,
  },
  {
    x: 'Aug',
    value: 85,
  },
  {
    x: 'Sep',
    value: 80,
  },
  {
    x: 'Oct',
    value: 87,
  },
  {
    x: 'Nov',
    value: 90,
  },
  {
    x: 'Dec',
    value: 95,
  },
];

interface NewCustomersProps {
  width?: number;
  height?: number;
}

const getX = (d: NewCustomersDatum) => d.x;

const NewCustomersChart = ({ width }: NewCustomersProps) => {
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
        <BarSeries
          dataKey='Newly Contracted'
          radius={4}
          data={mockData}
          radiusAll
          xAccessor={(d) => d.x}
          yAccessor={(d) => d.value}
          colorAccessor={({ x }) => primary600}
        />

        <AnimatedAxis
          orientation='bottom'
          hideAxisLine
          hideTicks
          tickLabelProps={{
            fontWeight: 'medium',
            fontFamily: `var(--font-barlow)`,
          }}
        />
        <Tooltip
          snapTooltipToDatumY
          snapTooltipToDatumX
          style={{
            position: 'absolute',
            padding: '8px',
            background: gray700,
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
              <Flex flexDir='column'>
                <Text color='white' fontWeight='normal'>
                  {xLabel}
                  {': '}${value}
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
