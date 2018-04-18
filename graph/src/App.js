import React, { Component } from 'react';
import { ResponsiveChord } from 'nivo';
import request from 'superagent';
import chordBackendToChordData from './utils/transforms/chordBackendToChordData';
import HistorySelect from './appComponents/HistorySelect';
import PlayPause from './appComponents/PlayPause';
// import sample from './mock/sample.json';
import './App.css';

const TICK_PERIOD = 3000;

class App extends Component {
  constructor(props) {
    super(props);
    this.state = ({
      playing: true,
      history: [],
      initialFetching: true,
      data: [[]],
      labels:[],
    })
  }
  tick = () => {

    this.doRequest();

    clearTimeout(this.timeout);
    this.timeout = setTimeout(this.tick, TICK_PERIOD);
  }
  doRequest() {
    request
      .get('http://52.42.55.249:10000/chord')
      .end((err, res) => {
        if(err) {
          console.log('error', err);
          return;
        }
        const graph = JSON.parse(res.text);
        const newState = chordBackendToChordData(graph.graph)
        // const newState = chordBackendToChordData(sample.graph);
        const time = new Date();
        this.state.history.push({
          time,
          item: newState,
        });
        if (this.state.history.length > 10) {
          // cap history
          this.state.history.splice(0, 1);
        }
        if (this.state.playing) {
          this.setState({
            history: this.state.history.slice(),
            time,
            initialFetching: false,
            data: newState.matrix,
            labels: newState.labels
          })
        } else {
          this.setState({
            history: this.state.history.slice(),
          })
        }
      })
  }
  onHistoryChange = (historyItem) => {
    this.setState({
      playing: false,
      time: historyItem.time,
      data: historyItem.item.matrix,
      labels: historyItem.item.labels,
    })
  }
  componentDidMount() {
    clearTimeout(this.timeout);
    this.timeout = setTimeout(this.tick, TICK_PERIOD);
  }
  onPlayPauseChange = (newVal) => {
    this.setState({
      playing: newVal
    })
  }
  render() {
    if (this.state.initialFetching) {
      return(
        <div className="loading">
          loading...
        </div>
      );
    }
    return (
      <div className="App">
        <div className="controls">
          <PlayPause onChange={this.onPlayPauseChange} playing={this.state.playing}/>
        </div>
        <div className="controls">
          <HistorySelect history={this.state.history} onChange={this.onHistoryChange}/>
        </div>
        <h3>
          connections:
          {this.state.time.toLocaleTimeString('en-US')}
        </h3>
        <ResponsiveChord
        matrix={this.state.data}
        keys={this.state.labels}
        margin={{
            "top": 60,
            "right": 60,
            "bottom": 90,
            "left": 60
        }}
        padAngle={0.02}
        innerRadiusRatio={0.96}
        innerRadiusOffset={0.02}
        arcOpacity={1}
        arcBorderWidth={1}
        arcBorderColor="inherit:darker(0.4)"
        ribbonOpacity={0.5}
        ribbonBorderWidth={1}
        ribbonBorderColor="inherit:darker(0.4)"
        enableLabel={true}
        label="id"
        labelOffset={12}
        labelRotation={-90}
        labelTextColor="inherit:darker(1)"
        colors="set1"
        isInteractive={true}
        arcHoverOpacity={1}
        arcHoverOthersOpacity={0.25}
        ribbonHoverOpacity={0.75}
        ribbonHoverOthersOpacity={0.25}
        animate={true}
        motionStiffness={90}
        motionDamping={17}
        legends={[
            {
                "anchor": "bottom",
                "direction": "row",
                "translateY": 70,
                "itemWidth": 80,
                "itemHeight": 14,
                "symbolSize": 14,
                "symbolShape": "circle"
            }
        ]}
    />
      </div>
    );
  }
}

export default App;
