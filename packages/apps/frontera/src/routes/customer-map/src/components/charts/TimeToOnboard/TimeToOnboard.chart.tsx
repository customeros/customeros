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

import { mockData } from './mock';
import { getMonthLabel } from '../util';

const margin = {
  top: 10,
  right: 0,
  bottom: 10,
  left: 0,
};

const height = 200;
const axisHeight = 8;

export type TimeToOnboardDatum = {
  month: number;
  value: number;
  index: number;
};

interface MrrPerCustomerProps {
  width: number;
  height?: number;
  hasContracts?: boolean;
  data?: TimeToOnboardDatum[];
}

const getDate = (d: TimeToOnboardDatum) => {
  return set(new Date(), { month: d.month - 1, year: d.index });
};
const bisectDate = bisector<TimeToOnboardDatum, Date>((d) => getDate(d)).left;
const getY = (d: TimeToOnboardDatum) => d.value;

const TimeToOnboardChart = ({
  width,
  hasContracts,
  data: _data = [],
}: MrrPerCustomerProps) => {
  const data = hasContracts ? _data : mockData;

  const colors = {
    primary600: hasContracts ? '#7F56D9' : '#D0D5DD',
    gray300: '#D0D5DD',
    gray700: '#344054',
  };
  const {
    tooltipTop,
    tooltipLeft,
    tooltipOpen,
    showTooltip,
    hideTooltip,
    tooltipData,
  } = useTooltip<TimeToOnboardDatum>();

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
      <svg height={height} width={width || 500} style={{ overflow: 'visible' }}>
        <LinearGradient
          to={'white'}
          fromOpacity={0}
          from={colors.primary600}
          id='mrr-per-customer-gradient'
          toOpacity={hasContracts ? 0.3 : 0.8}
        />
        <MarkerCircle
          size={2}
          refX={2}
          stroke='white'
          strokeWidth={1}
          fill={colors.primary600}
          id='mrr-per-customer-marker-circle'
        />
        <MarkerCircle
          size={2}
          refX={2}
          fill='white'
          strokeWidth={1}
          stroke={colors.primary600}
          id='mrr-per-customer-marker-circle-end'
        />
        <AreaClosed<TimeToOnboardDatum>
          data={data}
          yScale={scaleY}
          strokeWidth={0}
          pointerEvents='none'
          stroke={colors.primary600}
          x={(d) => scaleX(getDate(d))}
          y={(d) => scaleY(d.value) ?? 0}
          fill='url(#mrr-per-customer-gradient)'
        />

        <LinePath<TimeToOnboardDatum>
          data={data}
          strokeWidth={2}
          curve={curveLinear}
          stroke={colors.primary600}
          x={(d) => scaleX(getDate(d))}
          y={(d) => scaleY(getY(d)) ?? 0}
          shapeRendering='geometricPrecision'
          markerMid='url(#mrr-per-customer-marker-circle)'
          markerStart='url(#mrr-per-customer-marker-circle)'
          markerEnd='url(#mrr-per-customer-marker-circle-end)'
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
            fill: colors.gray700,
            fontWeight: 'medium',
            fontFamily: `var(--font-ibm-plex-sans)`,
          }}
        />
        {tooltipOpen && tooltipData && (
          <g>
            <Line
              strokeWidth={1.5}
              strokeDasharray='4'
              pointerEvents='none'
              stroke={colors.gray300}
              from={{ x: tooltipLeft, y: 0 }}
              to={{ x: tooltipLeft, y: innerHeight }}
            />
            <circle
              r={6}
              stroke='white'
              cy={tooltipTop}
              strokeWidth={2}
              cx={tooltipLeft}
              pointerEvents='none'
              fill={colors.primary600}
            />
          </g>
        )}
      </svg>
      <div className='flex w-full relative'>
        {tooltipData && tooltipOpen && (
          <TooltipWithBounds
            key={Math.random()}
            style={{
              top: -axisHeight - 16,
              left: tooltipLeft ?? 0,
              position: 'absolute',
              width: 'auto',
              fontSize: '14px',
              textAlign: 'center',
              borderRadius: '8px',
              padding: '8px',
              background: colors.gray700,
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
              hasContracts ? `${tooltipData.value} days` : 'No data yet'
            }`}
          </TooltipWithBounds>
        )}
      </div>
    </div>
  );
};

export default TimeToOnboardChart;
