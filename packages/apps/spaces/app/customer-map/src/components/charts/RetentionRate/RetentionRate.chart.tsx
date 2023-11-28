'use client';
import { PatternLines } from '@visx/pattern';
import {
  XYChart,
  Tooltip,
  AnimatedAxis,
  AnimatedGrid,
  AnimatedBarStack,
  AnimatedBarSeries,
} from '@visx/xychart';

import { useToken } from '@ui/utils';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

import { Legend } from '../../Legend';
import { getMonthLabel } from '../util';

export type RetentionRateDatum = {
  month: number;
  values: {
    renewed: number;
    churned: number;
  };
};

const _mockData: RetentionRateDatum[] = [
  {
    month: 1,
    values: {
      renewed: 20,
      churned: -10,
    },
  },
  {
    month: 2,
    values: {
      renewed: 50,
      churned: -20,
    },
  },
  {
    month: 3,
    values: {
      renewed: 68,
      churned: -30,
    },
  },
  {
    month: 4,
    values: {
      renewed: 55,
      churned: -40,
    },
  },
  {
    month: 5,
    values: {
      renewed: 73,
      churned: -50,
    },
  },
  {
    month: 6,
    values: {
      renewed: 80,
      churned: -60,
    },
  },
  {
    month: 7,
    values: {
      renewed: 85,
      churned: -70,
    },
  },
  {
    month: 8,
    values: {
      renewed: 90,
      churned: -80,
    },
  },
  {
    month: 9,
    values: {
      renewed: 95,
      churned: -90,
    },
  },
  {
    month: 10,
    values: {
      renewed: 100,
      churned: -95,
    },
  },
  {
    month: 11,
    values: {
      renewed: 71,
      churned: -65,
    },
  },
  {
    month: 12,
    values: {
      renewed: 95,
      churned: -25,
    },
  },
];

interface RetentionRateProps {
  width?: number;
  height?: number;
  data: RetentionRateDatum[];
}

const getX = (d: RetentionRateDatum) => getMonthLabel(d.month);

const RetentionRate = ({ data, width }: RetentionRateProps) => {
  const [gray700, warning950] = useToken('colors', ['gray.700', 'warning.950']);

  const colorScale = {
    Renewed: '#66C61C',
    Churned: warning950,
  };

  const legendData = [
    {
      label: 'Renewed',
      color: colorScale.Renewed,
    },
    {
      label: 'Churned',
      color: colorScale.Churned,
    },
  ];

  return (
    <>
      <Legend data={legendData} />
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
        <PatternLines
          id='stripes-renewed'
          height={8}
          width={8}
          stroke={colorScale.Renewed}
          strokeWidth={2}
          orientation={['diagonal']}
        />
        <PatternLines
          id='stripes-churned'
          height={8}
          width={8}
          stroke={colorScale.Churned}
          strokeWidth={2}
          orientation={['diagonal']}
        />
        <AnimatedBarStack>
          <AnimatedBarSeries
            dataKey='Renewed'
            radius={4}
            data={data}
            radiusTop
            xAccessor={(d) => getMonthLabel(d.month)}
            yAccessor={(d) => d.values.renewed}
            colorAccessor={(_, i) =>
              i === data.length - 1
                ? 'url(#stripes-renewed)'
                : colorScale.Renewed
            }
          />
          <AnimatedBarSeries
            dataKey='Churned'
            radius={4}
            data={data}
            radiusBottom
            xAccessor={(d) => getMonthLabel(d.month)}
            yAccessor={(d) => -d.values.churned}
            colorAccessor={(_, i) =>
              i === data.length - 1
                ? 'url(#stripes-churned)'
                : colorScale.Churned
            }
          />
        </AnimatedBarStack>

        <AnimatedAxis
          orientation='bottom'
          hideAxisLine
          hideTicks
          tickLabelProps={{
            fontWeight: 'medium',
            fontFamily: `var(--font-barlow)`,
          }}
        />
        <AnimatedGrid
          columns={false}
          numTicks={1}
          lineStyle={{ stroke: 'white', strokeWidth: 2 }}
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
              tooltipData?.nearestDatum?.datum as RetentionRateDatum,
            );
            const values = (
              tooltipData?.nearestDatum?.datum as RetentionRateDatum
            ).values;

            return (
              <Flex flexDir='column'>
                <Text color='white' fontWeight='normal' fontSize='sm'>
                  {xLabel}
                </Text>

                <Flex direction='column'>
                  <TooltipEntry
                    label='Renewed'
                    value={values.renewed}
                    color={colorScale.Renewed}
                  />
                  <TooltipEntry
                    label='Churned'
                    value={values.churned}
                    color={colorScale.Churned}
                  />
                </Flex>
              </Flex>
            );
          }}
        />
      </XYChart>
    </>
  );
};

const TooltipEntry = ({
  color,
  label,
  value,
}: {
  color: string;
  label: string;
  value: number;
}) => {
  return (
    <Flex align='center' gap='4'>
      <Flex align='center' flex='1' gap='2'>
        <Flex
          w='2'
          h='2'
          bg={color}
          borderRadius='full'
          border='1px solid white'
        />
        <Text color='white' fontSize='sm'>
          {label}
        </Text>
      </Flex>
      <Flex justify='flex-start'>
        <Text color='white' fontSize='sm'>
          {formatCurrency(value)}
        </Text>
      </Flex>
    </Flex>
  );
};

export default RetentionRate;
