'use client';
import { useRouter } from 'next/navigation';
import React, { useMemo, useState, useCallback } from 'react';

import { sample } from 'lodash';
import { Group } from '@visx/group';
import { Circle } from '@visx/shape';
import { useTooltip, TooltipWithBounds } from '@visx/tooltip';
import genRandomNormalPoints, {
  PointsRange,
} from '@visx/mock-data/lib/generators/genRandomNormalPoints';

import { useToken } from '@ui/utils';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { DateTimeUtils } from '@spaces/utils/date';
import { DashboardCustomerMapState } from '@graphql/types';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

import { Legend } from '../../Legend';
import { useDodge, DodgedCircleData } from './useDodge';

const genTimeValue = () =>
  sample([
    new Date('2023-01-01'),
    new Date('2023-01-15'),
    new Date('2023-01-25'),
    new Date('2023-02-01'),
    new Date('2023-02-15'),
    new Date('2023-02-25'),
    new Date('2023-03-01'),
    new Date('2023-03-15'),
    new Date('2023-03-25'),
    new Date('2023-04-01'),
    new Date('2023-04-15'),
    new Date('2023-04-25'),
    new Date('2023-05-01'),
    new Date('2023-05-15'),
    new Date('2023-05-25'),
    new Date('2023-06-01'),
    new Date('2023-06-15'),
    new Date('2023-06-25'),
    new Date('2023-07-02'),
    new Date('2023-07-15'),
    new Date('2023-07-25'),
    new Date('2023-08-01'),
    new Date('2023-08-15'),
    new Date('2023-08-25'),
    new Date('2023-09-01'),
    new Date('2023-09-15'),
    new Date('2023-09-25'),
    new Date('2023-10-01'),
    new Date('2023-10-15'),
    new Date('2023-10-25'),
    new Date('2023-11-01'),
    new Date('2023-11-15'),
    new Date('2023-11-25'),
    new Date('2023-12-01'),
    new Date('2023-12-15'),
    new Date('2023-12-31'),
  ]);

const points: PointsRange[] = genRandomNormalPoints(300, 0.5).filter(
  (_, i) => i < 300,
);

const pointsWithRadii = points.map(([x, y]) => ({
  x: genTimeValue(),
  y: 0,
  r: sample([2, 15, 6, 12, 22, 5, 10, 15, 25, 7, 120]),
  values: {
    id: `${Math.random()}`,
    name: '',
    status: sample(['OK', 'AT_RISK', 'CHURNED']),
  },
}));
const data = pointsWithRadii;

export type CustomerMapDatum = {
  x: Date;
  y: number;
  r: number;
  values: {
    id: string;
    name: string;
    status: string;
  };
};

const margin = {
  top: 40,
  right: 50,
  bottom: 20,
  left: 50,
};

interface CustomerMapChartProps {
  width: number;
  height: number;
  data: CustomerMapDatum[];
}

const CustomerMapChart = ({
  // data,
  width: outerWidth = 800,
  height: outerHeight = 800,
}: CustomerMapChartProps) => {
  const router = useRouter();
  const [crosshairX, setCrosshairX] = useState(0);
  const [hoveredId, setHoveredId] = useState('');
  const [
    greenLight500,
    greenLight400,
    warning300,
    warning200,
    warm200,
    warm300,
    gray700,
    gray300,
    gray500,
  ] = useToken('colors', [
    'greenLight.500',
    'greenLight.400',
    'warning.300',
    'warning.200',
    'warm.200',
    'warm.300',
    'gray.700',
    'gray.300',
    'gray.500',
  ]);

  const width = outerWidth - margin.left - margin.right;
  const height = outerHeight - margin.top - margin.bottom;
  const {
    tooltipLeft,
    tooltipTop,
    tooltipOpen,
    showTooltip,
    hideTooltip,
    tooltipData,
  } = useTooltip<CustomerMapDatum>();

  const handlePointerMove = useCallback(
    (datum: DodgedCircleData<{ id: string; name: string; status: string }>) =>
      (event: React.PointerEvent<SVGCircleElement>) => {
        const containerX = 'clientX' in event ? event.clientX : 0;
        const containerY = 'clientY' in event ? event.clientY : 0;

        showTooltip({
          tooltipLeft: containerX,
          tooltipTop: containerY,
          tooltipData: datum.data,
        });
      },
    [showTooltip],
  );

  const { transformData, minMaxX } = useDodge<{
    id: string;
    name: string;
    status: string;
  }>({
    width,
    height,
    data,
    marginLeft: margin.left,
    marginRight: margin.right,
  });

  const transformedData = useMemo(
    () => transformData(data),
    [transformData, data],
  );

  const legendData = [
    { color: greenLight500, label: 'All good' },
    { color: warning300, label: 'At risk' },
    { color: warm200, label: 'Churned' },
  ];

  const getCircleColor = (
    status: DashboardCustomerMapState,
    options: { isOutline: boolean } = { isOutline: false },
  ) => {
    const { isOutline } = options;
    switch (status) {
      case DashboardCustomerMapState.Ok:
        return isOutline ? greenLight400 : greenLight500;
      case DashboardCustomerMapState.AtRisk:
        return isOutline ? warning200 : warning300;
      case DashboardCustomerMapState.Churned:
        return isOutline ? warm300 : warm200;
      default:
        return isOutline ? warm300 : warm200;
    }
  };

  if (width < 10) return null;

  return (
    <>
      <Legend data={legendData} />
      <svg width={outerWidth} height={outerHeight}>
        <Group>
          <text x={margin.left} y={outerHeight - margin.bottom} fill={gray500}>
            {DateTimeUtils.format(minMaxX[0].toISOString(), 'd MMM')}
          </text>
          <text
            x={outerWidth / 2 - 32}
            y={outerHeight - margin.bottom}
            fill={gray500}
            fontWeight={600}
          >
            Sign date
          </text>
          <text
            x={outerWidth - (margin.left + margin.right + 12)}
            y={outerHeight - margin.bottom}
            fill={gray500}
          >
            {DateTimeUtils.format(minMaxX[1].toISOString(), 'd MMM')}
          </text>
        </Group>
        {tooltipOpen && (
          <>
            <line
              x1={crosshairX}
              x2={crosshairX}
              stroke={gray300}
              strokeWidth={1.5}
              strokeDasharray={'4 4'}
              y2={outerHeight - margin.bottom - 23}
              y1={0}
            />

            <Group>
              <rect
                x={crosshairX - 35}
                width={70}
                y={outerHeight - margin.bottom - 23}
                height={35}
                fill={gray700}
                rx={8}
              />
              <text x={crosshairX - 23} y={outerHeight - 20} fill='white'>
                {tooltipData?.x
                  ? DateTimeUtils.format(tooltipData?.x.toISOString(), 'd MMM')
                  : 'N/A'}
              </text>
            </Group>
          </>
        )}

        <Group width={width} height={height}>
          {transformedData.map((d, i) => (
            <>
              {hoveredId === d.data.values.id && (
                <Circle
                  key={`circle-hovered-${i}`}
                  cx={d.x}
                  cy={height - 6 - d.y}
                  r={d.r + 4}
                  strokeWidth={3}
                  stroke={getCircleColor(
                    d.data.values.status as DashboardCustomerMapState,
                    { isOutline: true },
                  )}
                  fill='white'
                />
              )}
              <Circle
                key={`circle-${i}`}
                cx={d.x}
                cy={height - 6 - d.y}
                r={d.r}
                fill={getCircleColor(
                  d.data.values.status as DashboardCustomerMapState,
                )}
                onMouseLeave={() => {
                  hideTooltip();
                  setHoveredId('');
                }}
                onMouseEnter={() => {
                  setCrosshairX(d.x);
                  setHoveredId(d.data.values.id);
                }}
                onClick={() => router.push(`/organization/${d.data.values.id}`)}
                onPointerMove={handlePointerMove(d)}
                cursor='pointer'
              />
            </>
          ))}
        </Group>
      </svg>
      {tooltipOpen && (
        <TooltipWithBounds
          key={Math.random()}
          left={tooltipLeft}
          top={tooltipTop}
          style={{
            position: 'absolute',
            padding: '8px',
            background: gray700,
            borderRadius: '8px',
          }}
        >
          <Flex flexDir='column'>
            <Text color='white'>{tooltipData?.values?.name}</Text>
            <Text color='white'>{formatCurrency(tooltipData?.r ?? 0)}</Text>
          </Flex>
        </TooltipWithBounds>
      )}
    </>
  );
};

export default CustomerMapChart;
