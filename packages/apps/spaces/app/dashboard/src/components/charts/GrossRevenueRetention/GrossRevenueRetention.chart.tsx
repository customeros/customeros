'use client';

import { useMemo, useCallback } from 'react';

import { set } from 'date-fns';
import { localPoint } from '@visx/event';
import { curveLinear } from '@visx/curve';
import { MarkerCircle } from '@visx/marker';
import { LinearGradient } from '@visx/gradient';
import { scaleTime, scaleLinear } from '@visx/scale';
import { max, extent, bisector } from '@visx/vendor/d3-array';
import { Bar, Line, LinePath, AreaClosed } from '@visx/shape';
import { useTooltip, TooltipWithBounds } from '@visx/tooltip';

import { useToken } from '@ui/utils';
import { Flex } from '@ui/layout/Flex';
// import { Text } from '@ui/typography/Text';

export type GrossRevenueRetentionDatum = {
  month: number;
  value: number;
};

interface GrossRevenueRetentionProps {
  width: number;
  height?: number;
  data: GrossRevenueRetentionDatum[];
}

const _margin = {
  top: 40,
  right: 50,
  bottom: 20,
  left: 50,
};

const getDate = (d: GrossRevenueRetentionDatum) =>
  set(new Date(), { month: d.month });
const bisectDate = bisector<GrossRevenueRetentionDatum, Date>((d) =>
  getDate(d),
).left;
// const getX = (d: GrossRevenueRetentionDatum) => getMonthLabel(d.month);
const getY = (d: GrossRevenueRetentionDatum) => d.value;

const GrossRevenueRetention = ({ data, width }: GrossRevenueRetentionProps) => {
  const [primary600, gray300, gray700] = useToken('colors', [
    'primary.600',
    'gray.300',
    'gray.700',
  ]);

  const height = 200;

  const {
    tooltipLeft,
    tooltipTop,
    tooltipOpen,
    showTooltip,
    hideTooltip,
    tooltipData,
  } = useTooltip<GrossRevenueRetentionDatum>();

  const scaleX = useMemo(
    () =>
      scaleTime({
        range: [0, width],
        domain: extent(data, getDate) as [Date, Date],
      }),
    [width, data],
  );
  const scaleY = useMemo(
    () =>
      scaleLinear({
        range: [height, 0],
        domain: [0, max(data, (d) => d.value) || 0],
        nice: true,
      }),
    [height, data],
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
        tooltipLeft: x,
        tooltipTop: scaleY(getY(d)),
      });
    },
    [showTooltip, scaleY, scaleX, data],
  );

  return (
    <div style={{ position: 'relative' }}>
      <svg width={width || 500} height={height}>
        <LinearGradient
          fromOpacity={0}
          toOpacity={0.3}
          to={'white'}
          from={primary600}
          id='visx-area-gradient'
        />
        <MarkerCircle
          id='marker-circle'
          fill={primary600}
          size={2}
          refX={2}
          strokeWidth={1}
          stroke='white'
        />
        <MarkerCircle
          id='marker-circle-end'
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
          fill='url(#visx-area-gradient)'
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
          markerMid='url(#marker-circle)'
          markerStart='url(#marker-circle)'
          markerEnd='url(#marker-circle-end)'
        />
        <Bar
          x={0}
          y={0}
          width={width}
          height={height}
          fill='transparent'
          rx={14}
          onMouseLeave={hideTooltip}
          onTouchMove={handleTooltip}
          onMouseMove={handleTooltip}
          onTouchStart={handleTooltip}
        />
        {tooltipOpen && tooltipData && (
          <g>
            <Line
              from={{ x: tooltipLeft, y: 0 }}
              to={{ x: tooltipLeft, y: height }}
              stroke={gray300}
              strokeWidth={2}
              pointerEvents='none'
              strokeDasharray='5,2'
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
      <Flex w='full' position='relative' h='20px'>
        {tooltipData && (
          <TooltipWithBounds
            key={Math.random()}
            style={{
              top: 0,
              left: tooltipLeft,
              position: 'absolute',
              minWidth: 72,
              textAlign: 'center',
              borderRadius: '8px',
              padding: '8px',
              background: gray700,
              color: 'white',
              transform: 'translateX(-50%)',
            }}
          >
            {`$${tooltipData.value}`}
          </TooltipWithBounds>
        )}
      </Flex>
    </div>
  );
};

export default GrossRevenueRetention;
