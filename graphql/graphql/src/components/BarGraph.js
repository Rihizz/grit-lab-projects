import React from "react";
import { useRef, useEffect } from "react";
import { scaleLinear, scaleBand, axisBottom, axisLeft, select } from "d3";
import * as d3 from "d3";

function BarGraph({ data, width, height }) {
  const margin = { top: 20, right: 20, bottom: 80, left: 40 };
  const innerWidth = width - margin.left - margin.right;
  const innerHeight = height - margin.top - margin.bottom;

  // Create x and y scales
  const xScale = scaleBand()
    .domain(data.map((d) => d.label))
    .range([0, innerWidth])
    .padding(0.1);

  const yScale = scaleLinear()
    .domain([0, Math.max(...data.map((d) => d.value))])
    .nice()
    .range([innerHeight, 0]);

  // Create a custom tick format function
  const tickFormat = (n) => {
    const absN = Math.abs(n);
    if (absN >= 1e6) {
      return `${n / 1e6}M`;
    } else if (absN >= 1e3) {
      return `${n / 1e3}K`;
    } else {
      return n;
    }
  };

  // Create x and y axes
  const xAxis = axisBottom(xScale).tickSizeOuter(0);
  const yAxis = axisLeft(yScale).tickFormat(tickFormat);

  const xAxisRef = (element) => {
    if (element) {
      select(element)
        .call(xAxis)
        .selectAll("text")
        .attr("transform", "rotate(-45)")
        .style("text-anchor", "end");
    }
  };

  const yAxisRef = (element) => {
    if (element) {
      select(element).call(yAxis);
    }
  };

  const svgRef = useRef();
  const tooltipRef = useRef();

  useEffect(() => {
    if (tooltipRef.current) {
      d3.select(tooltipRef.current)
        .style("opacity", 0)
        .style("position", "absolute")
        .style("background-color", "white")
        .style("padding", "5px")
        .style("border", "1px solid black")
        .style("border-radius", "5px");
    }
  }, []);

  return (
    <div>
      <div ref={tooltipRef} className="tooltip" />
      <svg ref={svgRef} width={width} height={height}>
        <g transform={`translate(${margin.left}, ${margin.top})`}>
          {data.map((d, i) => {
            const barHeight = innerHeight - yScale(d.value);
            const x = xScale(d.label);
            const y = yScale(d.value);

            return (
              <g key={i} transform={`translate(${x}, ${y})`}>
                <rect
                  width={xScale.bandwidth()}
                  height={barHeight}
                  fill="steelblue"
                  onMouseOver={() => {
                    d3.select(tooltipRef.current)
                      .style("opacity", 1)
                      .html(`XP: ${d.value}`);
                  }}
                  onMouseMove={(event) => {
                    d3.select(tooltipRef.current)
                      .style("left", `${event.pageX + 10}px`)
                      .style("top", `${event.pageY}px`);
                  }}
                  onMouseOut={() => {
                    d3.select(tooltipRef.current).style("opacity", 0);
                  }}
                />
              </g>
            );
          })}
          <g ref={xAxisRef} transform={`translate(0, ${innerHeight})`} />
          <g ref={yAxisRef} />
        </g>
      </svg>
    </div>
  );
}

export default BarGraph;
