'use client';
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
        <AnimatedBarStack>
          <AnimatedBarSeries
            dataKey='Renewed'
            radius={4}
            data={data}
            radiusTop
            xAccessor={(d) => getMonthLabel(d.month)}
            yAccessor={(d) => d.values.renewed}
            colorAccessor={({ month }) => colorScale.Renewed}
          />
          <AnimatedBarSeries
            dataKey='Churned'
            radius={4}
            data={data}
            radiusBottom
            xAccessor={(d) => getMonthLabel(d.month)}
            yAccessor={(d) => -d.values.churned}
            colorAccessor={({ month }) => colorScale.Churned}
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
                <Text color='white' fontWeight='normal'>
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
  value: string | number;
}) => {
  return (
    <Flex align='center' gap='4'>
      <Flex align='center' flex='1' gap='2'>
        <Flex
          w='3'
          h='3'
          bg={color}
          borderRadius='full'
          border='1px solid white'
        />
        <Text color='white'>{label}</Text>
      </Flex>
      <Flex minW='10' justify='flex-start'>
        <Text color='white'>${value}</Text>
      </Flex>
    </Flex>
  );
};

export default RetentionRate;
