import React, { Component } from 'react';

export default class PlayPause extends Component {
  onClick = () => {
    this.props.onChange(!this.props.playing);
  }
  render() {
    return (
      <button onClick={this.onClick}>
        {
          this.props.playing ?
          'pause' : 'play'
        }
      </button>
    )
  }
}
