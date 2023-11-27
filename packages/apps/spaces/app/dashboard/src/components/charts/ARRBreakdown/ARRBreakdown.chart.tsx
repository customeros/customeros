'use client';

import {
  XYChart,
  Tooltip,
  BarStack,
  BarSeries,
  AnimatedGrid,
  AnimatedAxis,
} from '@visx/xychart';

import { useToken } from '@ui/utils';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';

import { Legend } from '../../Legend';
import { getMonthLabel } from '../util';

export type ARRBreakdownDatum = {
  month: number;
  values: {
    upsells: number;
    churned: number;
    renewals: number;
    downgrades: number;
    cancellations: number;
    newlyContracted: number;
  };
};

const _mockData: ARRBreakdownDatum[] = [
  {
    month: 1,
    values: {
      cancellations: -5,
      churned: -15,
      downgrades: -8,
      newlyContracted: 12,
      renewals: 18,
      upsells: 7,
    },
  },
  {
    month: 2,
    values: {
      cancellations: -8,
      churned: -10,
      downgrades: -5,
      newlyContracted: 14,
      renewals: 20,
      upsells: 3,
    },
  },
  {
    month: 3,
    values: {
      cancellations: -12,
      churned: -7,
      downgrades: -10,
      newlyContracted: 8,
      renewals: 15,
      upsells: 5,
    },
  },
  {
    month: 4,
    values: {
      cancellations: -6,
      churned: -12,
      downgrades: -7,
      newlyContracted: 10,
      renewals: 22,
      upsells: 9,
    },
  },
  {
    month: 5,
    values: {
      cancellations: -10,
      churned: -8,
      downgrades: -15,
      newlyContracted: 13,
      renewals: 16,
      upsells: 4,
    },
  },
  {
    month: 6,
    values: {
      cancellations: -14,
      churned: -5,
      downgrades: -9,
      newlyContracted: 11,
      renewals: 14,
      upsells: 6,
    },
  },
  {
    month: 7,
    values: {
      cancellations: -9,
      churned: -11,
      downgrades: -12,
      newlyContracted: 9,
      renewals: 17,
      upsells: 8,
    },
  },
  {
    month: 8,
    values: {
      cancellations: -11,
      churned: -14,
      downgrades: -6,
      newlyContracted: 15,
      renewals: 19,
      upsells: 2,
    },
  },
  {
    month: 9,
    values: {
      cancellations: -7,
      churned: -9,
      downgrades: -11,
      newlyContracted: 12,
      renewals: 21,
      upsells: 10,
    },
  },
  {
    month: 10,
    values: {
      cancellations: -13,
      churned: -13,
      downgrades: -8,
      newlyContracted: 7,
      renewals: 13,
      upsells: 5,
    },
  },
  {
    month: 11,
    values: {
      cancellations: -8,
      churned: -6,
      downgrades: -14,
      newlyContracted: 10,
      renewals: 18,
      upsells: 7,
    },
  },
  {
    month: 12,
    values: {
      cancellations: -10,
      churned: -10,
      downgrades: -10,
      newlyContracted: 10,
      renewals: 10,
      upsells: 10,
    },
  },
];

interface ARRBreakdownProps {
  width: number;
  height?: number;
  data: ARRBreakdownDatum[];
}

const getX = (d: ARRBreakdownDatum) => getMonthLabel(d.month);

const ARRBreakdown = ({ width, data }: ARRBreakdownProps) => {
  const [gray700, moss300, warning400, warning600, warning950] = useToken(
    'colors',
    ['gray.700', 'moss.300', 'warning.400', 'warning.600', 'warning.950'],
  );

  const colorScale = {
    NewlyContracted: '#3B7C0F',
    Renewals: '#66C61C',
    Upsells: moss300,
    Downgrades: warning400,
    Cancellations: warning600,
    Churned: warning950,
  };

  const legendData = [
    {
      label: 'Newly Contracted',
      color: colorScale.NewlyContracted,
    },
    {
      label: 'Renewals',
      color: colorScale.Renewals,
    },
    {
      label: 'Upsells',
      color: colorScale.Upsells,
    },
    {
      label: 'Downgrades',
      color: colorScale.Downgrades,
    },
    {
      label: 'Cancellations',
      color: colorScale.Cancellations,
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
        <BarStack>
          <BarSeries
            dataKey='Churned'
            data={data}
            xAccessor={(d) => getMonthLabel(d.month)}
            yAccessor={(d) => -d.values.churned}
            colorAccessor={({ month }) => colorScale.Churned}
          />
          <BarSeries
            dataKey='Cancelations'
            data={data}
            xAccessor={(d) => getMonthLabel(d.month)}
            yAccessor={(d) => -d.values.cancellations}
            colorAccessor={({ month }) => colorScale.Cancellations}
          />
          <BarSeries
            dataKey='Downgrades'
            data={data}
            radiusBottom
            radius={4}
            xAccessor={(d) => getMonthLabel(d.month)}
            yAccessor={(d) => -d.values.downgrades}
            colorAccessor={({ month }) => colorScale.Downgrades}
          />

          <BarSeries
            dataKey='Newly Contracted'
            data={data}
            xAccessor={(d) => getMonthLabel(d.month)}
            yAccessor={(d) => d.values.newlyContracted}
            colorAccessor={({ month }) => colorScale.NewlyContracted}
          />
          <BarSeries
            dataKey='Renewals'
            data={data}
            xAccessor={(d) => getMonthLabel(d.month)}
            yAccessor={(d) => d.values.renewals}
            colorAccessor={({ month }) => colorScale.Renewals}
          />
          <BarSeries
            dataKey='Upsells'
            data={data}
            radius={4}
            radiusTop
            xAccessor={(d) => getMonthLabel(d.month)}
            yAccessor={(d) => d.values.upsells}
            colorAccessor={({ month }) => colorScale.Upsells}
          />
        </BarStack>

        <AnimatedGrid
          columns={false}
          numTicks={1}
          lineStyle={{ stroke: 'white', strokeWidth: 2 }}
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
              tooltipData?.nearestDatum?.datum as ARRBreakdownDatum,
            );
            const values = (
              tooltipData?.nearestDatum?.datum as ARRBreakdownDatum
            ).values;

            return (
              <Flex flexDir='column'>
                <Text color='white' fontWeight='normal'>
                  {xLabel}
                </Text>

                <Flex direction='column'>
                  <TooltipEntry
                    label='Upsells'
                    value={values.upsells}
                    color={colorScale.Upsells}
                  />
                  <TooltipEntry
                    label='Renewals'
                    value={values.renewals}
                    color={colorScale.Renewals}
                  />
                  <TooltipEntry
                    label='Newly Contracted'
                    value={values.newlyContracted}
                    color={colorScale.NewlyContracted}
                  />
                  <TooltipEntry
                    label='Churned'
                    value={values.churned}
                    color={colorScale.Churned}
                  />
                  <TooltipEntry
                    label='Cancelations'
                    value={values.cancellations}
                    color={colorScale.Cancellations}
                  />
                  <TooltipEntry
                    label='Downgrades'
                    value={values.downgrades}
                    color={colorScale.Downgrades}
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

export default ARRBreakdown;
