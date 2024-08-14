import React, { useEffect, useRef, useMemo } from "react";
import * as d3 from "d3";

function PieChart({ data, width, height }) {
  const svgRef = useRef();
  const margin = useMemo(
    () => ({ top: 40, right: 20, bottom: 40, left: 20 }),
    []
  );
  const innerWidth = width - margin.left - margin.right;
  const innerHeight = height - margin.top - margin.bottom;

  useEffect(() => {
    const pie = d3.pie().value((d) => d.value);
    const data_ready = pie([
      { name: "Up", value: data.up.aggregate.sum.amount },
      { name: "Down", value: data.down.aggregate.sum.amount },
    ]);

    const arcGenerator = d3
      .arc()
      .innerRadius(0)
      .outerRadius(Math.min(innerWidth, innerHeight) / 2);

    const colorScale = d3.scaleOrdinal(d3.schemeCategory10);

    const svg = d3.select(svgRef.current);

    const tooltip = d3
      .select("body")
      .append("div")
      .style("opacity", 0)
      .attr("class", "tooltip")
      .style("position", "absolute")
      .style("background-color", "white")
      .style("padding", "5px")
      .style("border", "1px solid black")
      .style("border-radius", "5px");

    const arcs = svg.selectAll(".arc").data(data_ready);

    arcs
      .enter()
      .append("path")
      .attr("class", "arc")
      .attr("fill", (d) => colorScale(d.data.name))
      .attr("stroke", "white")
      .attr("d", d3.arc().innerRadius(0).outerRadius(0))
      .attr(
        "transform",
        `translate(${innerWidth / 2 + margin.left}, ${
          innerHeight / 2 + margin.top
        })`
      )
      .on("mouseover", (event, d) => {
        tooltip
          .style("opacity", 1)
          .html(`${d.data.name}: ${d.data.value.toFixed(2)}`);
      })
      .on("mousemove", (event, d) => {
        tooltip
          .style("left", `${event.pageX + 10}px`)
          .style("top", `${event.pageY}px`);
      })
      .on("mouseout", () => {
        tooltip.style("opacity", 0);
      })
      .transition()
      .duration(1000)
      .attrTween("d", function (d) {
        const i = d3.interpolate({ startAngle: 0, endAngle: 0 }, d);
        return (t) => arcGenerator(i(t));
      });

    // Create a legend
    const legend = svg
      .append("g")
      .attr("transform", `translate(${margin.left}, ${innerHeight + 50})`);

    // Legend items
    const legendItems = [
      { name: "Up", color: colorScale("Up") },
      { name: "Down", color: colorScale("Down") },
    ];

    // Add legend items
    legendItems.forEach((item, index) => {
      // Colored rectangle
      legend
        .append("rect")
        .attr("x", innerWidth / 2 - 50 + index * 100)
        .attr("y", 0)
        .attr("width", 20)
        .attr("height", 20)
        .attr("fill", item.color);

      // Text label
      legend
        .append("text")
        .attr("x", innerWidth / 2 - 25 + index * 100)
        .attr("y", 15)
        .text(item.name)
        .attr("font-size", 14)
        .attr("font-weight", "bold");
    });
  }, [data, innerWidth, innerHeight, margin.left, margin.top]);

  return (
    <svg ref={svgRef} width={width} height={height}>
      <g transform={`translate(${margin.left}, ${margin.top})`}></g>
    </svg>
  );
}

export default PieChart;
