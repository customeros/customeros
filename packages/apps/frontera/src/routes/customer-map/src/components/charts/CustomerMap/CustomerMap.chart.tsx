import { useNavigate } from 'react-router-dom';
import React, { useMemo, useState, useCallback } from 'react';

import { Group } from '@visx/group';
import { Circle } from '@visx/shape';
import { useTooltip, TooltipWithBounds } from '@visx/tooltip';

import { DateTimeUtils } from '@utils/date';
import { DashboardCustomerMapState } from '@graphql/types';
import { formatCurrency } from '@utils/getFormattedCurrencyNumber';

import { mockData } from './mock';
import { Legend } from '../../Legend';
import { useDodge, DodgedCircleData } from './useDodge';

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
  hasContracts?: boolean;
  data: CustomerMapDatum[];
}

const CustomerMapChart = ({
  data: _data,
  hasContracts,
  width: outerWidth = 800,
  height: outerHeight = 800,
}: CustomerMapChartProps) => {
  const data = hasContracts ? _data : mockData;
  const navigate = useNavigate();
  const [crosshairX, setCrosshairX] = useState(0);
  const [hoveredId, setHoveredId] = useState('');

  const colors = {
    greenLight500: hasContracts ? '#66C61C' : '#98A2B3',
    greenLight400: hasContracts ? '#85E13A' : '#D0D5DD',
    yellow400: hasContracts ? '#FAC515' : '#D0D5DD',
    yellow300: hasContracts ? '#FDE272' : '#EAECF0',
    orangeDark800: hasContracts ? '#97180C' : '#98A2B3',
    orangeDark700: hasContracts ? '#BC1B06' : '#D0D5DD',
    warm200: hasContracts ? '#E7E5E4' : '#EAECF0',
    warm300: hasContracts ? '#D7D3D0' : '#D0D5DD',
    gray700: '#344054',
    gray300: '#D0D5DD',
    gray500: '#667085',
  };

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
    { color: colors.greenLight500, label: 'High' },
    { color: colors.yellow400, label: 'Medium' },
    { color: colors.orangeDark800, label: 'Low' },
    { color: colors.warm200, label: 'Churned', borderColor: colors.warm300 },
  ];

  const getCircleColor = (
    status: DashboardCustomerMapState,
    options: { isOutline: boolean } = { isOutline: false },
  ) => {
    const { isOutline } = options;

    switch (status) {
      case DashboardCustomerMapState.Ok:
        return isOutline ? colors.greenLight400 : colors.greenLight500;
      case DashboardCustomerMapState.MediumRisk:
        return isOutline ? colors.yellow300 : colors.yellow400;
      case DashboardCustomerMapState.HighRisk:
        return isOutline ? colors.orangeDark700 : colors.orangeDark800;
      case DashboardCustomerMapState.Churned:
        return isOutline ? colors.warm300 : colors.warm200;
      default:
        return isOutline ? colors.warm300 : colors.warm200;
    }
  };

  if (width < 10) return null;

  return (
    <>
      <Legend
        data={legendData}
        leftElement={
          hasContracts ? (
            <div className='flex items-center text-base'>
              <p className='font-normal text-gray-500  mr-1'>
                Renewal likelihood
              </p>
              <span className='text-gray-500'>· {data.length} customers</span>
            </div>
          ) : undefined
        }
      />
      <svg width={outerWidth} height={outerHeight}>
        <Group>
          <text
            fontSize={14}
            x={margin.left}
            fill={colors.gray500}
            y={outerHeight - margin.bottom}
          >
            {DateTimeUtils.format(minMaxX[0]?.toISOString(), 'd MMM')}
          </text>
          <text
            fontSize={14}
            fontWeight={600}
            fill={colors.gray500}
            x={outerWidth / 2 - 28}
            y={outerHeight - margin.bottom}
          >
            Sign date
          </text>
          <text
            fontSize={14}
            fill={colors.gray500}
            y={outerHeight - margin.bottom}
            x={outerWidth - margin.right - 38}
          >
            {DateTimeUtils.format(minMaxX[1]?.toISOString(), 'd MMM')}
          </text>
        </Group>
        {tooltipOpen && (
          <>
            <line
              y1={0}
              x1={crosshairX}
              x2={crosshairX}
              strokeWidth={1.5}
              stroke={colors.gray300}
              strokeDasharray={'4 4'}
              y2={outerHeight - margin.bottom - 23}
            />

            <Group>
              <rect
                rx={8}
                width={104}
                height={35}
                x={crosshairX - 51}
                fill={colors.gray700}
                y={outerHeight - margin.bottom - 23}
              />
              <text fill='white' x={crosshairX - 39} y={outerHeight - 20}>
                {tooltipData?.x
                  ? DateTimeUtils.format(
                      tooltipData?.x?.toISOString(),
                      'd MMM y',
                    )
                  : 'N/A'}
              </text>
            </Group>
          </>
        )}

        <Group width={width} height={height}>
          {transformedData.map((d, i) => (
            <React.Fragment key={i}>
              {hoveredId === d.data.values.id && (
                <Circle
                  cx={d.x}
                  r={d.r + 4}
                  fill='white'
                  strokeWidth={3}
                  cy={height - 6 - d.y}
                  key={`circle-hovered-${i}`}
                  stroke={getCircleColor(
                    d.data.values.status as DashboardCustomerMapState,
                    { isOutline: true },
                  )}
                />
              )}
              <Circle
                r={d.r}
                cx={d.x}
                cursor='pointer'
                key={`circle-${i}`}
                cy={height - 6 - d.y}
                onPointerMove={handlePointerMove(d)}
                fill={getCircleColor(
                  d.data.values.status as DashboardCustomerMapState,
                )}
                onMouseLeave={() => {
                  hideTooltip();
                  setHoveredId('');
                }}
                onClick={() =>
                  hasContracts && navigate(`/organization/${d.data.values.id}`)
                }
                onMouseEnter={() => {
                  setCrosshairX(d.x);
                  setHoveredId(d.data.values.id);
                }}
              />
            </React.Fragment>
          ))}
        </Group>
      </svg>
      {tooltipOpen && (
        <TooltipWithBounds
          top={tooltipTop}
          left={tooltipLeft}
          key={Math.random()}
          style={{
            position: 'absolute',
            padding: '8px',
            background: colors.gray700,
            borderRadius: '8px',
          }}
        >
          <div className='flex flex-col'>
            <p className='text-white'>
              {hasContracts ? tooltipData?.values?.name : 'No data available'}
            </p>
            <p className='text-white'>
              {formatCurrency(hasContracts ? tooltipData?.r ?? 0 : 0)}
            </p>
          </div>
        </TooltipWithBounds>
      )}
    </>
  );
};

export default CustomerMapChart;
