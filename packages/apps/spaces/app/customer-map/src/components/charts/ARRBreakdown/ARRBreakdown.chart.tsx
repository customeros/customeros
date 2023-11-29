'use client';

import { PatternLines } from '@visx/pattern';
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
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

import { Legend } from '../../Legend';
import { getMonthLabel } from '../util';

export type ARRBreakdownDatum = {
  month: number;
  upsells: number;
  churned: number;
  renewals: number;
  downgrades: number;
  cancellations: number;
  newlyContracted: number;
};

const _mockData: ARRBreakdownDatum[] = [
  {
    month: 1,
    cancellations: -5,
    churned: -15,
    downgrades: -8,
    newlyContracted: 12,
    renewals: 18,
    upsells: 7,
  },
  {
    month: 2,
    cancellations: -8,
    churned: -10,
    downgrades: -5,
    newlyContracted: 14,
    renewals: 20,
    upsells: 3,
  },
  {
    month: 3,
    cancellations: -12,
    churned: -7,
    downgrades: -10,
    newlyContracted: 8,
    renewals: 15,
    upsells: 5,
  },
  {
    month: 4,
    cancellations: -6,
    churned: -12,
    downgrades: -7,
    newlyContracted: 10,
    renewals: 22,
    upsells: 9,
  },
  {
    month: 5,
    cancellations: -10,
    churned: -8,
    downgrades: -15,
    newlyContracted: 13,
    renewals: 16,
    upsells: 4,
  },
  {
    month: 6,
    cancellations: -14,
    churned: -5,
    downgrades: -9,
    newlyContracted: 11,
    renewals: 14,
    upsells: 6,
  },
  {
    month: 7,
    cancellations: -9,
    churned: -11,
    downgrades: -12,
    newlyContracted: 9,
    renewals: 17,
    upsells: 8,
  },
  {
    month: 8,
    cancellations: -11,
    churned: -14,
    downgrades: -6,
    newlyContracted: 15,
    renewals: 19,
    upsells: 2,
  },
  {
    month: 9,
    cancellations: -7,
    churned: -9,
    downgrades: -11,
    newlyContracted: 12,
    renewals: 21,
    upsells: 10,
  },
  {
    month: 10,
    cancellations: -13,
    churned: -13,
    downgrades: -8,
    newlyContracted: 7,
    renewals: 13,
    upsells: 5,
  },
  {
    month: 11,
    cancellations: -8,
    churned: -6,
    downgrades: -14,
    newlyContracted: 10,
    renewals: 18,
    upsells: 7,
  },
  {
    month: 12,
    cancellations: -10,
    churned: -10,
    downgrades: -10,
    newlyContracted: 10,
    renewals: 10,
    upsells: 10,
  },
];

interface ARRBreakdownProps {
  width: number;
  height?: number;
  data: ARRBreakdownDatum[];
}

const getX = (d: ARRBreakdownDatum) => getMonthLabel(d.month);

const ARRBreakdown = ({ width, data }: ARRBreakdownProps) => {
  const [
    gray700,
    greenLight200,
    greenLight400,
    warning300,
    warning600,
    warning950,
    greenLight700,
    greenLight500,
  ] = useToken('colors', [
    'gray.700',
    'greenLight.200',
    'greenLight.400',
    'warning.300',
    'warning.600',
    'warning.950',
    'greenLight.700',
    'greenLight.500',
  ]);

  const colorScale = {
    NewlyContracted: greenLight700,
    Renewals: greenLight500,
    Upsells: greenLight200,
    Downgrades: warning300,
    Cancellations: warning600,
    Churned: warning950,
  };

  const legendData = [
    {
      label: 'Newly contracted',
      color: colorScale.NewlyContracted,
    },
    {
      label: 'Renewals',
      color: colorScale.Renewals,
    },
    {
      label: 'Upsells',
      color: colorScale.Upsells,
      borderColor: greenLight400,
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

  const getBarColor = (key: keyof typeof colorScale, barIndex: number) =>
    barIndex === data.length - 1 ? `url(#stripes-${key})` : colorScale[key];

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
        {Object.entries(colorScale).map(([key, color]) => (
          <PatternLines
            key={key}
            id={`stripes-${key}`}
            height={8}
            width={8}
            stroke={color}
            strokeWidth={2}
            orientation={['diagonal']}
          />
        ))}
        <BarStack offset='diverging'>
          <BarSeries
            dataKey='Churned'
            data={data}
            xAccessor={(d) => getMonthLabel(d.month)}
            yAccessor={(d) => -d.churned}
            colorAccessor={(_, i) => getBarColor('Churned', i)}
          />
          <BarSeries
            dataKey='Cancelations'
            data={data}
            xAccessor={(d) => getMonthLabel(d.month)}
            yAccessor={(d) => -d.cancellations}
            colorAccessor={(_, i) => getBarColor('Cancellations', i)}
          />
          <BarSeries
            dataKey='Downgrades'
            data={data}
            radiusBottom
            radius={4}
            xAccessor={(d) => getMonthLabel(d.month)}
            yAccessor={(d) => -d.downgrades}
            colorAccessor={(_, i) => getBarColor('Downgrades', i)}
          />

          <BarSeries
            dataKey='Newly Contracted'
            data={data}
            xAccessor={(d) => getMonthLabel(d.month)}
            yAccessor={(d) => d.newlyContracted}
            colorAccessor={(_, i) => getBarColor('NewlyContracted', i)}
          />
          <BarSeries
            dataKey='Renewals'
            data={data}
            xAccessor={(d) => getMonthLabel(d.month)}
            yAccessor={(d) => d.renewals}
            colorAccessor={(_, i) => getBarColor('Renewals', i)}
          />
          <BarSeries
            dataKey='Upsells'
            data={data}
            radius={4}
            radiusTop
            xAccessor={(d) => getMonthLabel(d.month)}
            yAccessor={(d) => d.upsells}
            colorAccessor={(_, i) => getBarColor('Upsells', i)}
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
            fontSize: 12,
            fontWeight: 'medium',
            fontFamily: `var(--font-barlow)`,
          }}
        />
        <Tooltip
          snapTooltipToDatumY
          snapTooltipToDatumX
          style={{
            position: 'absolute',
            padding: '8px 12px',
            background: gray700,
            borderRadius: '8px',
          }}
          renderTooltip={({ tooltipData }) => {
            const xLabel = getX(
              tooltipData?.nearestDatum?.datum as ARRBreakdownDatum,
            );
            const values = tooltipData?.nearestDatum
              ?.datum as ARRBreakdownDatum;

            const sumPositives =
              values.newlyContracted + values.renewals + values.upsells;
            const sumNegatives =
              values.churned + values.cancellations + values.downgrades;

            const totalSum = sumPositives - sumNegatives;

            return (
              <Flex flexDir='column'>
                <Flex justify='space-between' align='center'>
                  <Text color='white' fontWeight='semibold' fontSize='sm'>
                    {xLabel}
                  </Text>
                  <Text color='white' fontWeight='semibold' fontSize='sm'>
                    {formatCurrency(totalSum)}
                  </Text>
                </Flex>

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
                    label='Newly contracted'
                    value={values.newlyContracted}
                    color={colorScale.NewlyContracted}
                  />
                  <TooltipEntry
                    label='Churned'
                    value={values.churned}
                    color={colorScale.Churned}
                  />
                  <TooltipEntry
                    label='Cancellations'
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
      <Flex>
        <Text color='white' fontSize='sm'>
          {formatCurrency(value)}
        </Text>
      </Flex>
    </Flex>
  );
};

export default ARRBreakdown;
