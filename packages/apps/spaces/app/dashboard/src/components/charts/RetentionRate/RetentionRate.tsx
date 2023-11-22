'use client';

import ParentSize from '@visx/responsive/lib/components/ParentSize';
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

type Datum = {
  x: string;
  values: {
    renewed: number;
    churned: number;
  };
};

const mockData: Datum[] = [
  {
    x: 'Jan',
    values: {
      renewed: 20,
      churned: -10,
    },
  },
  {
    x: 'Feb',
    values: {
      renewed: 50,
      churned: -20,
    },
  },
  {
    x: 'Mar',
    values: {
      renewed: 68,
      churned: -30,
    },
  },
  {
    x: 'Apr',
    values: {
      renewed: 55,
      churned: -40,
    },
  },
  {
    x: 'May',
    values: {
      renewed: 73,
      churned: -50,
    },
  },
  {
    x: 'Jun',
    values: {
      renewed: 80,
      churned: -60,
    },
  },
  {
    x: 'Jul',
    values: {
      renewed: 85,
      churned: -70,
    },
  },
  {
    x: 'Aug',
    values: {
      renewed: 90,
      churned: -80,
    },
  },
  {
    x: 'Sep',
    values: {
      renewed: 95,
      churned: -90,
    },
  },
  {
    x: 'Oct',
    values: {
      renewed: 100,
      churned: -95,
    },
  },
  {
    x: 'Nov',
    values: {
      renewed: 71,
      churned: -65,
    },
  },
  {
    x: 'Dec',
    values: {
      renewed: 95,
      churned: -25,
    },
  },
];

interface RetentionRateProps {}

const getX = (d: Datum) => d.x;

const RetentionRate = (_props: RetentionRateProps) => {
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
    <ParentSize>
      {({ width }) => (
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
                data={mockData}
                radiusTop
                xAccessor={(d) => d.x}
                yAccessor={(d) => d.values.renewed}
                colorAccessor={({ x }) => colorScale.Renewed}
              />
              <AnimatedBarSeries
                dataKey='Churned'
                radius={4}
                data={mockData}
                radiusBottom
                xAccessor={(d) => d.x}
                yAccessor={(d) => d.values.churned}
                colorAccessor={({ x }) => colorScale.Churned}
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
                const xLabel = getX(tooltipData?.nearestDatum?.datum as Datum);
                const values = (tooltipData?.nearestDatum?.datum as Datum)
                  .values;

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
      )}
    </ParentSize>
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
