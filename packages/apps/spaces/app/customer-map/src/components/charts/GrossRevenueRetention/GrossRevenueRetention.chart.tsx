'use client';

import { useMemo, useCallback } from 'react';

import { set } from 'date-fns';
import { localPoint } from '@visx/event';
import { curveLinear } from '@visx/curve';
import { MarkerCircle } from '@visx/marker';
import { Axis, Orientation } from '@visx/axis';
import { LinearGradient } from '@visx/gradient';
import { scaleUtc, scaleLinear } from '@visx/scale';
import { timeFormat } from '@visx/vendor/d3-time-format';
import { max, extent, bisector } from '@visx/vendor/d3-array';
import { Bar, Line, LinePath, AreaClosed } from '@visx/shape';
import { useTooltip, TooltipWithBounds } from '@visx/tooltip';

import { useToken } from '@ui/utils';
import { Flex } from '@ui/layout/Flex';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';

import { mockData } from './mock';
import { getMonthLabel } from '../util';

export type GrossRevenueRetentionDatum = {
  month: number;
  value: number;
};

interface GrossRevenueRetentionProps {
  width: number;
  height?: number;
  hasContracts?: boolean;
  data: GrossRevenueRetentionDatum[];
}

const margin = {
  top: 10,
  right: 0,
  bottom: 10,
  left: 0,
};

const height = 200;
const axisHeight = 8;

const getDate = (d: GrossRevenueRetentionDatum) =>
  set(new Date(), { month: d.month - 1, date: 1 });
const bisectDate = bisector<GrossRevenueRetentionDatum, Date>((d) =>
  getDate(d),
).left;
const getY = (d: GrossRevenueRetentionDatum) => d.value;

const GrossRevenueRetention = ({
  width,
  hasContracts,
  data: _data = [],
}: GrossRevenueRetentionProps) => {
  const data = hasContracts ? _data : mockData;
  const [primary600, gray300, gray700] = useToken('colors', [
    hasContracts ? 'primary.600' : 'gray.300',
    'gray.300',
    'gray.700',
  ]);
  const {
    tooltipTop,
    tooltipLeft,
    tooltipOpen,
    showTooltip,
    hideTooltip,
    tooltipData,
  } = useTooltip<GrossRevenueRetentionDatum>();

  const innerHeight = height - margin.top - margin.bottom - axisHeight;
  const innerWidth = width - margin.left - margin.right;

  const scaleX = useMemo(
    () =>
      scaleUtc({
        range: [margin.left, innerWidth],
        domain: extent(data, getDate) as [Date, Date],
      }),
    [innerWidth, data],
  );
  const scaleY = useMemo(
    () =>
      scaleLinear({
        range: [innerHeight, margin.top],
        domain: [0, max(data, (d) => d.value) || 0],
        nice: true,
      }),
    [innerHeight, data],
  );

  const handleTooltip = useCallback(
    (
      event:
        | React.TouchEvent<SVGRectElement>
        | React.MouseEvent<SVGRectElement>,
    ) => {
      const { x } = localPoint(event) || { x: 0 };

      const x0 = scaleX.invert(x);
      const index = bisectDate(data, x0, 1);
      const d0 = data[index - 1];
      const d1 = data[index];
      let d = d0;
      if (d1 && getDate(d1)) {
        d =
          x0.valueOf() - getDate(d0).valueOf() >
          getDate(d1).valueOf() - x0.valueOf()
            ? d1
            : d0;
      }

      showTooltip({
        tooltipData: d,
        tooltipLeft: scaleX(getDate(d)),
        tooltipTop: scaleY(getY(d)),
      });
    },
    [showTooltip, scaleY, scaleX, data],
  );

  return (
    <div style={{ position: 'relative' }}>
      <svg width={width || 500} height={height} style={{ overflow: 'visible' }}>
        <LinearGradient
          fromOpacity={0}
          toOpacity={hasContracts ? 0.3 : 0.8}
          to={'white'}
          from={primary600}
          id='revenue-retention-gradient'
        />
        <MarkerCircle
          id='revenue-retention-marker-circle'
          fill={primary600}
          size={2}
          refX={2}
          strokeWidth={1}
          stroke='white'
        />
        <MarkerCircle
          id='revenue-retention-marker-circle-end'
          stroke={primary600}
          size={2}
          refX={2}
          strokeWidth={1}
          fill='white'
        />
        <AreaClosed<GrossRevenueRetentionDatum>
          data={data}
          x={(d) => scaleX(getDate(d))}
          y={(d) => scaleY(d.value) ?? 0}
          yScale={scaleY}
          strokeWidth={0}
          stroke={primary600}
          fill='url(#revenue-retention-gradient)'
          pointerEvents='none'
        />

        <LinePath<GrossRevenueRetentionDatum>
          data={data}
          curve={curveLinear}
          x={(d) => scaleX(getDate(d))}
          y={(d) => scaleY(getY(d)) ?? 0}
          strokeWidth={2}
          stroke={primary600}
          shapeRendering='geometricPrecision'
          markerMid='url(#revenue-retention-marker-circle)'
          markerStart='url(#revenue-retention-marker-circle)'
          markerEnd='url(#revenue-retention-marker-circle-end)'
        />
        <Bar
          x={0}
          y={0}
          rx={14}
          width={width}
          height={height}
          fill='transparent'
          onMouseLeave={hideTooltip}
          onTouchMove={handleTooltip}
          onMouseMove={handleTooltip}
          onTouchStart={handleTooltip}
        />
        <Axis
          hideTicks
          hideAxisLine
          scale={scaleX}
          top={innerHeight + axisHeight}
          tickValues={data.map(getDate)}
          orientation={Orientation.bottom}
          tickFormat={(d) => timeFormat('%b')(d as Date)}
          tickLabelProps={{
            fontSize: 12,
            fill: gray700,
            fontWeight: 'medium',
            fontFamily: `var(--font-barlow)`,
          }}
        />
        {tooltipOpen && tooltipData && (
          <g>
            <Line
              from={{ x: tooltipLeft, y: 0 }}
              to={{ x: tooltipLeft, y: innerHeight }}
              stroke={gray300}
              strokeWidth={1.5}
              pointerEvents='none'
              strokeDasharray='4'
            />
            <circle
              cx={tooltipLeft}
              cy={tooltipTop}
              r={6}
              fill={primary600}
              stroke='white'
              strokeWidth={2}
              pointerEvents='none'
            />
          </g>
        )}
      </svg>
      <Flex w='full' position='relative'>
        {tooltipData && tooltipOpen && (
          <TooltipWithBounds
            key={Math.random()}
            style={{
              top: -axisHeight - 16,
              left: tooltipLeft ?? 0,
              position: 'absolute',
              minWidth: 72,
              fontSize: '14px',
              textAlign: 'center',
              borderRadius: '8px',
              padding: '8px',
              background: gray700,
              color: 'white',
              whiteSpace: 'nowrap',
              transform:
                tooltipData.month === data[0].month
                  ? undefined
                  : tooltipData.month === data[data.length - 1].month
                  ? 'translateX(-100%)'
                  : 'translateX(-50%)',
            }}
          >
            {`${getMonthLabel(tooltipData.month)}: ${
              hasContracts ? formatCurrency(tooltipData.value) : 'No data yet'
            }`}
          </TooltipWithBounds>
        )}
      </Flex>
    </div>
  );
};

export default GrossRevenueRetention;
