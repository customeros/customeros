'use client';
import { useRef, useEffect } from 'react';

import * as d3 from 'd3';

import { dodgeVariable } from './dodge';

interface CircleData {
  x: number;
  y: number;
  radius: number;
}
interface UseBubblesOptions {
  width: number;
  height: number;
  data: CircleData[];
  transform?: string;
}

const LOW_R = 4; // min radius value
const TOP_R = 30; // max radius value
const PADDING = 1.5; // padding between circles

export const Bubbles3 = ({
  data,
  width,
  height,
  transform,
}: UseBubblesOptions) => {
  const ref = useRef<SVGGElement>(null);

  useEffect(() => {
    const target = d3.select(ref.current);
    const minMaxX = d3.extent(data, (d) => d.x) as [number, number];
    const minMaxR = d3.extent(data, (d) => d.radius) as [number, number];

    const xScale = d3.scaleLinear().domain(minMaxX).range([50, width]);
    const rScale = d3.scaleLinear().domain(minMaxR).range([LOW_R, TOP_R]);

    const _data = dodgeVariable(data, {
      x: (d: CircleData) => xScale(d.x),
      r: (d: CircleData) => rScale(d.radius),
      padding: PADDING,
    }) as (CircleData & { r: number })[];

    target
      .selectAll<SVGCircleElement, CircleData>('circle')
      .data(_data)
      .enter()
      .append('circle')
      .attr('cx', (d) => d.x)
      .attr('cy', (d) => height - PADDING - d.y)
      .attr('r', (d) => d.r)
      .style('fill', 'red')
      .style('opacity', 0.5);
  }, [data, width, height]);

  return <g ref={ref} width={width} height={height} transform={transform}></g>;
};
