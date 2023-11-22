'use client';
import React, { useMemo, RefObject } from 'react';

import { Zoom } from '@visx/zoom';
import sample from 'lodash/sample';
import { localPoint } from '@visx/event';
import { LinearGradient } from '@visx/gradient';
import { AnimatedAxis } from '@visx/react-spring';
import { scaleUtc, coerceNumber } from '@visx/scale';
import { timeFormat } from '@visx/vendor/d3-time-format';
import { Axis, AxisScale, Orientation, SharedAxisProps } from '@visx/axis';
import genRandomNormalPoints, {
  PointsRange,
} from '@visx/mock-data/lib/generators/genRandomNormalPoints';

// import { Bubbles } from './Bubbles2';
import { Bubbles3 } from './Bubble3';

const points: PointsRange[] = genRandomNormalPoints(600, 0.5).filter(
  (_, i) => i < 600,
);

const pointsWithRadii = points.map(([x, y]) => ({
  x,
  y: 0,
  radius: sample([2, 15, 6, 12, 22, 5, 10, 15, 25, 7, 120]),
}));
const data = pointsWithRadii;

const backgroundColor = '#FCFCFD';
const axisColor = '#344054';
const tickLabelColor = '#344054';
const margin = {
  top: 40,
  right: 50,
  bottom: 20,
  left: 50,
};

const tickLabelProps = {
  fill: tickLabelColor,
  fontSize: 12,
  // fontFamily: 'sans-serif',
  textAnchor: 'middle',
} as const;

const getMinMax = (vals: (number | { valueOf(): number })[]) => {
  const numericVals = vals.map(coerceNumber);

  return [Math.min(...numericVals), Math.max(...numericVals)];
};

export type AxisProps = {
  width: number;
  height: number;
  showControls?: boolean;
};

type AnimationTrajectory = 'outside' | 'center' | 'min' | 'max' | undefined;

type AxisComponentType = React.FC<
  SharedAxisProps<AxisScale> & {
    animationTrajectory: AnimationTrajectory;
  }
>;

const CustomerMap = ({
  width: outerWidth = 800,
  height: outerHeight = 800,
  showControls = true,
}: AxisProps) => {
  const width = outerWidth - margin.left - margin.right;
  const height = outerHeight - margin.top - margin.bottom;

  // use non-animated components if prefers-reduced-motion is set
  const prefersReducedMotionQuery =
    typeof window === 'undefined'
      ? false
      : window.matchMedia('(prefers-reduced-motion: reduce)');
  const prefersReducedMotion =
    !prefersReducedMotionQuery || !!prefersReducedMotionQuery.matches;

  const AxisComponent: AxisComponentType = !prefersReducedMotion
    ? AnimatedAxis
    : Axis;

  const { scale, tickFormat, values } = useMemo(() => {
    const timeValues = [
      new Date('2020-01-15'),
      new Date('2020-02-15'),
      new Date('2020-03-15'),
      new Date('2020-04-15'),
      new Date('2020-05-15'),
      new Date('2020-06-15'),
      new Date('2020-07-15'),
      new Date('2020-08-15'),
      new Date('2020-09-15'),
      new Date('2020-10-15'),
      new Date('2020-11-15'),
      new Date('2020-12-15'),
    ];

    return {
      scale: scaleUtc({
        domain: getMinMax(timeValues),
        range: [0, width],
      }),
      values: timeValues,
      tickFormat: (v: Date, i: number) => timeFormat('%b')(v),
    };
  }, [width]);
  if (width < 10) return null;

  const scalePadding = 20;
  const scaleHeight = height / 2 - scalePadding;

  return (
    <Zoom
      width={outerWidth}
      height={outerHeight}
      scaleXMin={1 / 2}
      scaleXMax={4}
      scaleYMin={1 / 2}
      scaleYMax={4}
      // constrain={(transformMatrix, prevTransformMatrix) => {
      //   const min = applyMatrixToPoint(transformMatrix, { x: width, y: -100 });
      //   const max = applyMatrixToPoint(transformMatrix, {
      //     x: width,
      //     y: height,
      //   });
      //   if (max.x < width || max.y < height) {
      //     return prevTransformMatrix;
      //   }
      //   if (min.x > 0 || min.y > 0) {
      //     return prevTransformMatrix;
      //   }

      //   return transformMatrix;
      // }}
    >
      {(zoom) => (
        <div style={{ position: 'relative' }}>
          <svg
            width={outerWidth}
            height={outerHeight}
            style={{
              cursor: zoom.isDragging ? 'grabbing' : 'grab',
              touchAction: 'none',
            }}
            ref={zoom.containerRef as RefObject<SVGSVGElement>}
          >
            <LinearGradient
              fromOpacity={0}
              to={backgroundColor}
              from={backgroundColor}
              id='visx-axis-gradient'
            />

            <rect
              x={0}
              y={0}
              rx={14}
              width={outerWidth}
              height={outerHeight}
              fill={'url(#visx-axis-gradient)'}
            />

            <Bubbles3
              data={data}
              width={width}
              height={height}
              transform={zoom.toString()}
            />

            <rect
              width={outerWidth}
              height={outerHeight}
              rx={14}
              fill='transparent'
              onTouchStart={zoom.dragStart}
              onTouchMove={zoom.dragMove}
              onTouchEnd={zoom.dragEnd}
              onMouseDown={zoom.dragStart}
              onMouseMove={zoom.dragMove}
              onMouseUp={zoom.dragEnd}
              onMouseLeave={() => {
                if (zoom.isDragging) zoom.dragEnd();
              }}
              onDoubleClick={(event) => {
                const point = localPoint(event) || { x: 0, y: 0 };
                zoom.scale({ scaleX: 1.1, scaleY: 1.1, point });
              }}
            />

            {/* <g
              clipPath='url(#zoom-clip)'
              transform={`
                    scale(0.25)
                    translate(${width * 4 - width - 60}, ${
                height * 4 - height - 60
              })
                  `}
            >
              <rect
                width={width}
                height={height}
                fill='white'
                fillOpacity={0.2}
                stroke='white'
                strokeWidth={4}
                transform={zoom.toStringInvert()}
              />
            </g> */}

            <g transform={`translate(${margin.left},${margin.top})`}>
              <g transform={`translate(0, ${scaleHeight + scalePadding})`}>
                <AxisComponent
                  scale={scale}
                  top={scaleHeight}
                  stroke={axisColor}
                  tickValues={values}
                  tickStroke={axisColor}
                  tickFormat={tickFormat}
                  animationTrajectory={'center'}
                  tickLabelProps={tickLabelProps}
                  orientation={Orientation.bottom}
                />
              </g>
            </g>
          </svg>
        </div>
      )}
    </Zoom>
  );
};

export default CustomerMap;
