// @flow
import React, { Component } from 'react';
import * as d3 from 'd3';
import sortedIndexBy from 'lodash/sortedIndexBy';
const pauseCursorColor = 'blue';

type MouseMoveEvent = {
  clientX: number,
  clientY: number,
};

type XYPoint = {
  /** X coordinate (always a timestamp) */
  x: number,
  /** Y coordinate */
  y: number
};

type Line = {
  /** Color of the line */
  color: string,
  /** Coordinates for the data */
  data: Array<XYPoint>,
  key: string,
  label: string,
};

type Props = {
  /** An array of lines containing the data (see below for more info) */
  lines: Array<Line>,
  /** Width of the chart */
  width: number, // eslint-disable-line react/no-unused-prop-types
  /** Height of the chart */
  height: number, // eslint-disable-line react/no-unused-prop-types
  /** Theme (loaded with Glamorous) */
  theme: Object,
};

type State = {
  px: number,
  py: number,
  matches: Array<any>,
};

const defaultProps = {
  lines: [],
  width: 1000,
  height: 300,
  theme: {
    colors: {
      gray2: "#b7bbbe",
      gray3: "#70787D",
    },
  },
};

const margin = {top: 5, right: 40, bottom: 24, left: 40 };

// const ChartTooltip = glamorous.div({
//   position: 'absolute',
//   background: 'white',
//   pointerEvents: 'none',
//   border: 'solid 1px #70787D',
//   padding: '10px',
//   zIndex: 99,
//   whiteSpace: 'nowrap',
//   fontFamily: 'OpenSans',
//   fontSize: '14px',
// });

class ChartTooltip extends Component {
  render () {
    return (
      <div {...this.props}>
      </div>
    );
  }
}
const getPoint = (node, event) => {
  const svg = node.ownerSVGElement || node;
  if (svg.createSVGPoint) {
    let point = svg.createSVGPoint();
    point.x = event.clientX;
    point.y = event.clientY;
    point = point.matrixTransform(node.getScreenCTM().inverse());
    return point;
  }

  const rect = node.getBoundingClientRect();
  return {
    x: event.clientX - rect.left - node.clientLeft,
    y: event.clientY - rect.top - node.clientTop,
  };
};

const MouseEventCaptureRectangle =
  (props: { height: number, width: number }) => {
  const height = props.height - margin.top - margin.bottom > 0 ?
    props.height - margin.top - margin.bottom : 0;
  const width = props.width - margin.left - margin.right > 0 ?
    props.width - margin.left - margin.right : 0;

  return (
    <g>
      {/* full background rectangle to capture mousemove ev */}
      <rect width={width} height={height} opacity="0" />
    </g>
  )
}

/**
 * A line chart for showing data through time
 */
class TimeSeriesChart extends Component<Props, State> {
  // eslint-disable-next-line react/sort-comp
  rootNode: any
  yGridLines: any
  d3Root: any
  d3Chart: any
  graphDiv: HTMLDivElement
  xScale: any
  yScale: any
  line: any
  width: number
  height: number

  static defaultProps = defaultProps
  state = {
    pauseCursorTime: 0,
    px: 0,
    py: 0,
    matches: [],
  }
  getClosestX(point: XYPoint) {

    const chartX = this.xScale.invert(point.x);
    let closestX = Infinity;
    this.props.lines.forEach((curLine) => {
      const data = curLine.data;
      let i = sortedIndexBy(data, {x: chartX}, (p) => p.x);
      i = data.length === i ? i - 1 : i;
      const closestDiff = Math.abs(closestX - chartX);
      const pointDiff = Math.abs(data[i].x - chartX);
      if (closestDiff > pointDiff) {
        closestX = data[i].x;
      }
      if ( i > 0 && data.length > 2) {
        const curDiff = Math.abs(closestX - chartX);
        const nextDiff = Math.abs(data[i - 1].x - chartX);
        if (curDiff > nextDiff) {
          closestX = data[i - 1].x;
        }
      }
    });
    return closestX;
  }
  focusMatchingData(point: XYPoint) {
    const closestX = this.getClosestX(point); 
    const cursorX = closestX;

    const foundMatch = [];
    this.props.lines.forEach(curLine => {
      let i = sortedIndexBy(curLine.data, {x: closestX}, (p) => p.x)
      i = i === curLine.data.length ? i - 1 : i;
      if (curLine.data[i].x === closestX) {
        foundMatch.push({
          x: curLine.data[i].x,
          y: curLine.data[i].y,
          color: curLine.color,
          key: curLine.key,
          label: curLine.label,
        })
      }
    });
    const focusRef= this.d3Chart.select('.focus')
      .selectAll('circle')
      .data(foundMatch)

    // set attr on enter to avoid glitch where a black dot
    // shows up on the upper left corner when adding a new point
    focusRef
      .enter()
      .append("circle")
      .attr("r", 3.5)
      .attr('fill', (d) => {
        return d.color;
      })
      .attr(
        "transform", (p) => {
          return [
            `translate(`,
              this.xScale(p.x),
              ',',
              this.yScale(p.y),
            `)`,
          ].join('');
        }
      );
    focusRef
      .attr(
        "transform", (p) => {
          return [
            `translate(`,
              this.xScale(p.x),
              ',',
              this.yScale(p.y),
            `)`,
          ].join('');
        })
      .attr('fill', (p) => {
        return p.color;
      })
    focusRef
      .exit()
      .remove()
    this.setState({
      matches: foundMatch,
    })


    this.d3Chart.select('.cursor')
      .select('line')
      .attr('x1', this.xScale(cursorX))
      .attr('x2', this.xScale(cursorX))
      .attr('y1', this.yScale.range()[0])
      .attr('y2', this.yScale.range()[1])
      .attr('stroke-width', 0.5)
      .attr('stroke', 'white');
  }

  focusInterpolatedData(point: XYPoint) {
    const foundMatch = [];
    const lines = this.d3Chart.selectAll('.lines path').nodes();
    const focusRef = this.d3Chart.select(".mouse-per-line")
      .selectAll('circle')
      .data(this.props.lines)
    focusRef.enter()
      .append('circle')
      .attr('r', 3.5)
      .attr('fill', 'transparent')
      .attr('stroke', (line) => {
        return line.color;
      });
    focusRef
      .exit()
      .remove();
    focusRef.attr("transform", (d, i) => {
        let beginning = 0;
        let end = lines[i].getTotalLength();
        let target = null;
        let pos = null;
        while (true){ // eslint-disable-line no-constant-condition
          target = Math.floor((beginning + end) / 2);
          pos = lines[i].getPointAtLength(target);
          if ((target === end || target === beginning) && pos.x !== point.x) {
              break;
          }
          if (pos.x > point.x) {
            end = target;
          } else if (pos.x < point.x){
            beginning = target;
          } else {
            break; // position found
          }
        }
        if (pos) {
          foundMatch.push({
            x: this.xScale.invert(point.x),
            y: this.yScale.invert(pos.y),
            color: d.color,
            key: d.key,
            label: d.label,
          })
          return `translate(${point.x},${pos.y})`;
        }
        return null;
      });
    this.setState({
      matches: foundMatch,
    })
    this.d3Chart.select('.cursor')
      .select('line')
      .attr('x1', point.x)
      .attr('x2', point.x)
      .attr('y1', this.yScale.range()[0])
      .attr('y2', this.yScale.range()[1])
      .attr('stroke-width', 0.5)
      .attr('stroke', 'white');
    this.d3Chart.select('.pause-cursor')
      .select('line')
      .attr('x1', point.x)
      .attr('x2', point.x)
      .attr('y1', this.yScale.range()[0])
      .attr('y2', this.yScale.range()[1])
      .attr('stroke-width', 0.5)
      .attr('stroke', 'white');
  }

  componentDidMount() {
    this.d3Root = d3.select(this.rootNode);
    this.d3Chart = this.d3Root
      .select('.chart')
      .attr('transform', `translate(${margin.left}, ${margin.top})`);
    this.setInstanceVariables(this.props);
    this.setGraphRefs(this.props);
    this.renderChart();

    const thisTimeSeries = this;
    this.d3Chart
      // need to use function instead of =>, need path bound to this
      // eslint-disable-next-line func-names
      .on('mousemove.chart', function () {
        const point = getPoint(this, d3.event);
        // thisTimeSeries.focusInterpolatedData(point);
        thisTimeSeries.focusMatchingData(point);

      })
      .on('click', function() {
        const point = getPoint(this, d3.event);
        thisTimeSeries.onClick(point)
      })
  }

  setInstanceVariables(propsToUse: Props) {
    const width = propsToUse.width - margin.left - margin.right;
    const height = propsToUse.height - margin.top - margin.bottom;
    this.width = width;
    this.height = height;
  }

  setGraphRefs(propsToUse: Props) {
    const { lines } = propsToUse;
    this.line = d3.line()
      .x(d => this.xScale(d.x))
      .y(d => this.yScale(d.y))
      .curve(d3.curveLinear);
   this.xScale = d3
      .scaleTime()
      .domain([
        d3.min(lines, line => d3.min(line.data, d => d.x)),
        d3.max(lines, line => d3.max(line.data, d => d.x)),
      ])
      .range([0, this.width]);
    this.yScale = d3
      .scaleLinear()
      .domain([
        0, // d3.min(lines, line => d3.min(line.data, d => d.y)),
        d3.max(lines, line => d3.max(line.data, d => d.y)) * 1.75,
      ])
      .range([this.height, 0]);
  }

  componentWillReceiveProps(nextProps: Props) {
    this.setInstanceVariables(nextProps);
    this.setGraphRefs(nextProps);
    if (this.props.hasPauseCursor !== nextProps.hasPauseCursor) {
      if (nextProps.hasPauseCursor) {
        this.d3Chart.select('.pause-cursor')
          .select('line')
          // .attr('x1', this.xScale(cursorX))
          // .attr('x2', this.xScale(cursorX))
          // .attr('y1', this.yScale.range()[0])
          // .attr('y2', this.yScale.range()[1])
          .attr('stroke-width', 0.5)
          .attr('stroke', pauseCursorColor);
      } else {
        this.d3Chart.select('.pause-cursor')
          .select('line')
          .attr('stroke-width', 0)
          .attr('stroke', 'transparent');
      }
    }
  }

  componentDidUpdate() {
    this.renderChart();
  }

  mouseout = () => {
    const focusRef = this.d3Chart.select('.focus')
      .selectAll('circle')
      .data([])
    focusRef
      .exit()
      .remove();

    const focusRef2 = this.d3Chart.select('.mouse-per-line')
      .selectAll('circle')
      .data([]);
    focusRef2
      .exit()
      .remove();

    this.d3Chart.select('.cursor')
      .select('line')
      .attr('stroke-width', 0)
      .attr('stroke', 'transparent');
    this.setState({
      matches: [],
    })
  }

  onMouseMove = (ev: MouseMoveEvent) => {
    const rect = this.graphDiv.getBoundingClientRect();
    this.setState({
      px: ev.clientX - rect.left,
      py: ev.clientY - rect.top,
    })
  }

  getTheme() {
    if (this.props.theme && this.props.theme.colors) {
      return this.props.theme;
    }
    return defaultProps.theme;
  }

  renderChart() {
    const {lines} = this.props;
    const tickPadding = 10;
    this.d3Root
      .attr('height', this.height + margin.left + margin.right)
      .attr('width', this.width + margin.top + margin.bottom)

    // at least 80 pixels between ticks, with a minimum of 1
    let xAxisTickCount = Math.max(Math.floor(this.width/100), 2);

    const domain = this.xScale.domain();
    const seconds = d3.timeSecond.count(domain[0],domain[1]);

    const interval = Math.floor(seconds / xAxisTickCount);

    // by default, D3 will pick the closest smaller interval in a list
    // of time periods. We end up with at least xAxisTickCount, but could
    // be way more. We want xAxisTickCount to be the maximum amount of ticks.
    // So we manually pick the closest bigger interval in the list and
    // use it to define the tick count.

    // time periods used by D3 between ticks (in seconds)
    const d3Ticks = [ 1, 5, 15, 30, // 1,5,15,30 seconds
                      40, 50,
                      1*60, 5*60, 15*60, 30*60, // 1,5,15,30 minutes
                      1*3600, 3*3600, 6*3600, 12*3600, // 1,3,6,12 hours
                      1*3600*24, 1*3600*24*2, 1*3600*24*7 ]; // 1,2,7 days

    // look for closest bigger time interval
    let optimizedInterval = d3Ticks.find(timePeriod => interval <= timePeriod)
    if (!optimizedInterval) optimizedInterval = d3Ticks[d3Ticks.length - 1]

    // re-evaluate xAxisTickCount
    xAxisTickCount = Math.floor(seconds / optimizedInterval)

    // time format, display seconds under a minute
    let timeFormat = d3.timeFormat('%I:%M %p');
    if (optimizedInterval < 60) {
      timeFormat = d3.timeFormat('%I:%M:%S %p');
    }

    const xAxis = d3
      .axisBottom(this.xScale)
      .tickFormat(timeFormat)
      .tickSizeOuter(0) // start & end ticks
      .tickSizeInner(0) // ticks in between
      .ticks(xAxisTickCount)
      .tickPadding(tickPadding);

    const yAxis = d3
      .axisLeft(this.yScale)
      .tickFormat(d => `${d}`)
      .tickPadding(tickPadding)
      .ticks(4)
      .tickSizeOuter(0)
      .tickSizeInner(0);

    // ------------------------------------------------------------------------
    // Render axis
    // ------------------------------------------------------------------------
    this.d3Chart
      .select('.x')
      .attr('transform', `translate(0, ${this.height})`)
      .call(xAxis);

    this.d3Chart
      .select('.y')
      .call(yAxis);

    const { colors } = this.getTheme();

    this.d3Chart.select('.y').select('.domain').attr('stroke', 'transparent');
    this.d3Chart.select('.y').selectAll('text').attr('fill', colors.gray3);
    this.d3Chart.select('.y').selectAll('text')
      .attr('style', 'font-size: 12px; font-family: "OpenSans"');
    this.d3Chart.select('.x').select('.domain').attr('stroke', colors.gray3);
    this.d3Chart.select('.x').selectAll('text').attr('fill', colors.gray3);
    this.d3Chart.select('.x').selectAll('text')
      .attr('style', 'font-size: 12px; font-family: "OpenSans"');

    // ------------------------------------------------------------------------
    // Render lines
    // ------------------------------------------------------------------------
    const d3Lines = this.d3Chart.select('.lines');
    const paths = d3Lines.selectAll('path').data(lines);

    const enter = paths.enter()
      .append('path')
      .attr('fill', 'none')
      .attr('stroke-linejoin', 'round')
      .attr('stroke-linecap', 'round')
      .attr('stroke-width', 1.5)

      enter.merge(paths)
        .attr('d', ({ data:d }) => this.line(d))
        .attr('stroke', ({ color }) => color)
      enter.each(function animate() {
        const element = d3.select(this);
        const totalLength = element.node().getTotalLength();
          element
            .attr('stroke-dasharray', `${totalLength} ${totalLength}`)
            .attr('stroke-dashoffset', totalLength)
            .transition()
            .ease(d3.easeCubic)
            .duration(2000)
            .attr('stroke-dashoffset', 0);
      });
      paths
        .attr('stroke-dasharray', null)
        .exit()
        .remove();
      if (this.props.hasPauseCursor) {
        this.d3Chart.select('.pause-cursor')
          .select('line')
          .attr('x1', this.xScale(this.state.pauseCursorTime))
          .attr('x2', this.xScale(this.state.pauseCursorTime))
          .attr('y1', this.yScale.range()[0])
          .attr('y2', this.yScale.range()[1])
          .attr('stroke-width', 0.5)
          .attr('stroke', pauseCursorColor);
      }
  }
  onClick = (point) => {
    console.log('mouse event')
    console.log(point)
    // const rect = this.graphDiv.getBoundingClientRect();
    const clickX = point.x;
    const clickY = point.y;
    this.setState({
      clickX,
      clickY,
    })
    const closestX = this.getClosestX({x: clickX, y: clickY});
    this.setState({
      pauseCursorTime: closestX,
    })

    console.log(console.log(new Date(closestX)))
    if (this.props.onClick) {
      this.props.onClick(closestX)
    }
  }
  render() {
    const matches = this.state.matches || [];
    return (
      // relative position wrapper to get
      // the tooltip placed properly w/ absolute positioning
      // eslint-disable-next-line no-return-assign
      <div
          ref={
            // $FlowFixMe
            (node) => {this.graphDiv = node}
          }
          style={{position: 'relative', cursor: 'default'}}
          onMouseMove={this.onMouseMove}
          onMouseLeave={this.mouseout}>
        <svg
            style={{
              cursor: 'crosshair',
            }}
            ref={node => {this.rootNode = node}}>
          <g className="chart">
            <MouseEventCaptureRectangle {...this.props} />
            <g className="axis x"/>
            <g className="axis y"/>
            <g className="lines" />
            <g className="pause-cursor">
              <line/>
            </g>
            <g className="cursor">
              <line/>
            </g>
            <g className="focus" />
            <g className="mouse-per-line"/>
          </g>
        </svg>
        <ChartTooltip style={{
            display: matches.length ?  'block' : 'none',
            top: this.state.py,
            left: this.state.px,
            color: 'black',
            position: 'absolute',
            background: 'grey',
            pointerEvents: 'none',
            border: 'solid 1px #70787D',
            padding: '10px',
            zIndex: 99,
            whiteSpace: 'nowrap',
            fontFamily: 'OpenSans',
            fontSize: '14px',
          }}
        >
          { matches.map((match) => {
            return (
              <div key={match.key}>
                <div style={{color: match.color}}>
                  {`${match.label}: ${match.y}` }
                </div>
              </div>
            )
          })}
          { matches.length ?
            <div style={{ color: 'white', paddingTop: '8px' }}>
              <div>
                {(new Date(matches[0].x)).toLocaleTimeString()}
              </div>
            </div>:
            null
          }
        </ChartTooltip>
      </div>
    )
  }
}

export default TimeSeriesChart;
