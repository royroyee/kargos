import React from "react";
import Chartist from "react-chartist";
import ChartistTooltip from 'chartist-plugin-tooltips-updated';

/**
 * Generate Resource Value chart according to the data.
 * @todo automate generating resource chart using REST API.
 * @returns Resource Value Chart
 */
export const ResourceChart = (props) => {

  const { data } = props;

  var label = []
  for (var i = 0; i < data.length; i++) {
    label.push(24 - i)
  }

  const graphData = {
    labels: label,
    series: [data]
  }

  const options = {
    low: 0,
    showArea: true,
    fullWidth: true,
    axisX: {
      position: 'end',
      showGrid: true
    },
    axisY: {
      // On the y-axis start means left and end means right
      showGrid: false,
      showLabel: false,
    }
  };

  const plugins = [
    ChartistTooltip()
  ]

  return (
    <Chartist data={graphData} options={{...options, plugins}} type="Line" className="ct-series-g ct-double-octave" />
  );
}