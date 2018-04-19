import React, { Component } from 'react';

export default class PlayPause extends Component {
  onClick = () => {
    this.props.onChange(!this.props.playing);
  }
  renderImg() {
    const imgStyle = {
      // width: '100%',
      height: '100%',
    }
    if (this.props.playing) {
      return (
        <img style={imgStyle} src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADIAAAAyCAYAAAAeP4ixAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAACLSURBVGhD7dkxCoAwAATB/L+x9yuW/sAP6aYXjrMQkR3YLoFcnSFJv7HQ9qCVknnm7m5qvqm20/mgg5J55u5uar6p5hBKHEI1h1DiEKo5hBKHUM0hlDiEag6hxCFUcwglDqGaQyhxCNUcQolDqOYQShxCNYdQ4hCqOYQSh1DNIZS8OuQ3n6GS9DVjXBcWQR16aeAQAAAAAElFTkSuQmCC"/>
      );
    }
    return (
      <img style={imgStyle} src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADIAAAAyCAYAAAAeP4ixAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAGnSURBVGhD7dk/KIVRHMbxiygxGAwmAxajUWS2kpLVKsPNaJaS1Solq1WKDAazYjIqi4XEIvF96Fdv6nLv+55/b52nPnXv+Cz3PuecRk5OjpP0ouvnY32ziFfcYxXdqGXO8Vlwg3nULpcoFjEqOIXapFUR+cAhRpF8/ipi3rCDISSbdoqYRzTRh+TSSRFzh2Uk9ZNdpoi5whySSJUi5hiTiBoXReQdexhBlLgqYl6wiQEEjesi5gGaPD0IEl9FzC2CTB7fRcwFvE6eUEVEk+cIXiZPyCJGk2cXTidPjCJGk2cDTiZPzCJGk2cFlSZPCkXMGQZRKikVER29SyWlIk8YQ6mkUkTDcxylE7uIjgIzqJxYRfRLtQRnh7PQRfTfsQ7nx+VQRfRvvg1vFxi+i2hfHcD7lZLPIqcIdsnno8g1gl+7uiwS9SLcRZFn6Jzej2ipUsRuToYRPWWLaFJMIJl0WkSTYhrJpd0imhSa2Mk+0f1XRJNiDUnewBfTqogmxRaSfhMp5ncRTYp91OKVqhi9RFmJE9Tq3bAY3c0uYPb7W05OTk68NBpfT/sOMP4jl8QAAAAASUVORK5CYII="/>
    );
  }
  render() {
    return (
      <button style={{width: 34, height: 24, background: 'transparent', margin: '0 5px'}} onClick={this.onClick}>
        {

          this.renderImg()
        }
      </button>
    )
  }
}
