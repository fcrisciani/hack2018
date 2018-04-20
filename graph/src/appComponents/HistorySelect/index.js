import React, { Component } from 'react';
import TimeSeriesChart from '../TimeSeriesChart';
import mockdata from './mockdata.json';

class HistoryButton extends Component {
  onClick = () => {
    this.props.onClick(this.props.item)
  }
  render() {
    return (
      <button onClick={this.onClick}>
        {this.props.item.time.toLocaleTimeString('en-US')}
      </button>
    )
  }
}
const transformHistoryDataToHistoryChartData = (history) => {
  return history.map(item => {
    return {
      x: item.time,
      y: item.totalCount,
    }
  });
}
export default class HistorySelect extends Component {
  onChartClick = (time) => {
    const found = this.props.history.filter((item) => {
      return (+item.time) === (+time)
    })[0]
    if (found) {
      this.props.onChange(found);
    } else {
      alert('could not find match')
    }
  }
  onClickHistoryButton = (item) => {
    this.props.onChange(item);
  }
  render() {
    const now = new Date();
    const historyChartLines = transformHistoryDataToHistoryChartData(this.props.history);
    return (
      <div>
        <TimeSeriesChart
          hasPauseCursor={this.props.hasPauseCursor}
          /** An array of lines containing the data (see below for more info) */
          lines={
            [{  
              color: 'white',
              data: historyChartLines,
              key: 'h',
              label: 'connection total count',
            }]
          }
          width={500}
          height={200}
          onClick={this.onChartClick}
        />
        <div style={{
          display: 'flex',
          flex: 1,
          justifyContent: 'center',
        }}>
          {
            /*
            this.props.history.map((instance) =>{
              return (
                <HistoryButton onClick={this.onClickHistoryButton} key={instance.time} item={instance}/>
              )
            })
            */
          }
        </div>
      </div>
    );
  }
}