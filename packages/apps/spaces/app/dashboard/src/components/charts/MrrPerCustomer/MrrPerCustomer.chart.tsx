'use client';

import { LinearGradient } from '@visx/gradient';
import {
  XYChart,
  Tooltip,
  AnimatedAxis,
  AnimatedAreaSeries,
  AnimatedGlyphSeries,
} from '@visx/xychart';

import { useToken } from '@ui/utils';
import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';

import { getMonthLabel } from '../util';

export type MrrPerCustomerDatum = {
  month: number;
  value: number;
};

interface MrrPerCustomerProps {
  width: number;
  height?: number;
  data?: MrrPerCustomerDatum[];
}

const getX = (d: MrrPerCustomerDatum) => getMonthLabel(d.month);
const getY = (d: MrrPerCustomerDatum) => d.value;

const MrrPerCustomerChart = ({ data = [], width }: MrrPerCustomerProps) => {
  const [primary600, gray300, gray500, gray700] = useToken('colors', [
    'primary.600',
    'gray.300',
    'gray.500',
    'gray.700',
  ]);

  return (
    <XYChart
      height={200}
      width={width || 500}
      margin={{ top: 12, right: 0, bottom: 20, left: 0 }}
      xScale={{ type: 'band', padding: 0 }}
      yScale={{ type: 'linear' }}
    >
      <LinearGradient
        fromOpacity={0}
        toOpacity={0.3}
        to={'white'}
        from={primary600}
        id='visx-area-gradient'
      />
      <AnimatedAreaSeries
        dataKey=''
        data={data}
        fill={'url(#visx-area-gradient)'}
        lineProps={{ stroke: primary600 }}
        xAccessor={(d) => getX(d)}
        yAccessor={(d) => getY(d)}
      />
      <AnimatedGlyphSeries
        dataKey=''
        data={data}
        renderGlyph={({ x, y, datum }) => {
          const isLast = datum.month === data[data.length - 1].month;

          return (
            <>
              {isLast && <circle cx={x} cy={y} r={7} fill={gray300} />}
              <circle
                r={4}
                cx={x}
                cy={y}
                stroke='white'
                strokeWidth={2}
                fill={isLast ? gray500 : primary600}
              />
            </>
          );
        }}
        xAccessor={(d) => getX(d)}
        yAccessor={(d) => getY(d)}
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
        showVerticalCrosshair
        verticalCrosshairStyle={{
          stroke: gray300,
          strokeDasharray: 4,
          strokeWidth: 1.5,
        }}
        showSeriesGlyphs
        glyphStyle={{
          fill: primary600,
          r: 6,
          stroke: 'white',
          strokeWidth: 2,
        }}
        style={{
          position: 'absolute',
          padding: '8px',
          background: gray700,
          borderRadius: '8px',
        }}
        renderTooltip={({ tooltipData }) => {
          const xLabel = getX(
            tooltipData?.nearestDatum?.datum as MrrPerCustomerDatum,
          );
          const yLabel = getY(
            tooltipData?.nearestDatum?.datum as MrrPerCustomerDatum,
          );

          return (
            <Flex>
              <Text color='white' fontWeight='normal'>
                {xLabel}
                {': '}
                {yLabel}
              </Text>
            </Flex>
          );
        }}
      />
    </XYChart>
  );
};

export default MrrPerCustomerChart;
