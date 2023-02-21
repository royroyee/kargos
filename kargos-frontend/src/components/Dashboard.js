import { Col, Row, Card, Dropdown, Button, ButtonGroup, ProgressBar } from '@themesberg/react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faServer, faQuestionCircle, faExclamationCircle, faCheckCircle } from '@fortawesome/free-solid-svg-icons';
import { Link } from 'react-router-dom';

import React, { useState, useEffect } from 'react';
import Chartist from "react-chartist";
import ChartistTooltip from 'chartist-plugin-tooltips-updated';

/**
 * Generate Dropdown item for a specific name of the element.
 * @param {string} element: The name of the element.
 * @returns The dropdown item that has hyperlink to pod's specific information.
 */
function controllersDropdown (element, baseurl) {
    return (
    <Dropdown.Item>
      <Link to={baseurl + element}> {element} </Link>
    </Dropdown.Item>
    );
  };

export const DashboardPage = () => {
    return (
        <>
            <Row className="justify-content-md-center">
                <StatusWidgets></StatusWidgets>
            </Row>
                <LastStatus></LastStatus>
            <Row className="justify-content-md-center">
                <IntensiveNodeWidgets></IntensiveNodeWidgets>
                <IntensivePodWidgets></IntensivePodWidgets>
            </Row>
        </>
    );
}

/**
 * Generate top CPU and RAM intensive node widgets.
 * @param {props} props The props to use.
 * @returns JSX Component that represents Top CPU and RAM intensive node widgets.
 */
const IntensiveNodeWidgets = (props) => {
    const [usages, setUsages] = useState({ ram: [], cpu: [] });

    // Retrieve data from REST API.
    const getUsage = () => {
      const url = "/overview/nodes/top";
      var requestOptions = {
          method: 'GET',
          redirect: 'follow'
      };
  
      fetch(url, requestOptions)
          .then(response => response.text())
          .then(result => {
            setUsages(JSON.parse(result));
          })
          .catch(error => console.log('error', error));
    }

    useEffect(() => {
        getUsage();
    }, []);

    // Generate BarWidget using data retrieved.
    return ( 
        <>
        <Col xs={12} xl={3} className="mb-4">
            <BarWidget title="Top CPU Intensive Nodes" data={usages.cpu} metric={"%"}></BarWidget>
        </Col>
        <Col xs={12} xl={3} className="mb-4">
            <BarWidget title="Top RAM Intensive Nodes" data={usages.ram} metric={"%"}></BarWidget>
        </Col>
        </>
    );
}

/**
 * Generate top CPU and RAM intensive pod widgets.
 * @param {props} props The props to use.
 * @returns JSX Component that represents Top CPU and RAM intensive pod widgets.
 */
const IntensivePodWidgets = (props) => {
    const [usages, setUsages] = useState({ ram: [], cpu: [] });

    // Retrieve data from REST API.
    const getUsage = () => {
      const url = "/overview/pods/top";
      var requestOptions = {
          method: 'GET',
          redirect: 'follow'
      };
  
      fetch(url, requestOptions)
          .then(response => response.text())
          .then(result => {
            setUsages(JSON.parse(result));
          })
          .catch(error => console.log('error', error));
    }

    useEffect(() => {
        getUsage();
    }, []);

    // Generate BarWidget using data retrieved.
    return ( 
        <>
        <Col xs={12} xl={3} className="mb-4">
            <BarWidget title="Top CPU Intensive Pods" data={usages.cpu} metric={"Mi"}></BarWidget>
        </Col>
        <Col xs={12} xl={3} className="mb-4">
            <BarWidget title="Top RAM Intensive Pods" data={usages.ram} metric={"MB"}></BarWidget>
        </Col>
        </>
    );
}

/**
 * Generate bar widget using given data.
 * @param {props} props The props to use.
 * @returns JSX Component that represents a single bar widget
 */
const BarWidget = (props) => {
    const { title, data, metric } = props;
    
    /**
     * Genrate progress bar according to the given value.
     * @param {{}} element The element. First key must be name and second key must be usage.
     * @param {string} metric The metric unit.
     * @returns JSX component that represents a row of progress bars.
     */
    function percentageBars (element, metric) {    
        const keys = Object.keys(element);
        const name = element[keys[0]];
        const usage = element[keys[1]];
        
        var colorValue;
        if (usage < 35) {
            colorValue = "green";
        } else if (usage < 70) {
            colorValue = "orange";
        } else {
                colorValue = "red";
        }

        const Progress = (props) => {
        const { title, percentage, color, last = false } = props;
        const extraClassName = last ? "" : "mb-2";
    
        return (
            <Row className={`align-items-center ${extraClassName}`}>
            <Col>
                <div className="progress-wrapper">
                <div className="progress-info">
                    <h6 className="mb-0">{title}</h6>
                    <small className="fw-bold text-dark">
                    <span>{percentage} {metric}</span>
                    </small>
                </div>
                <ProgressBar variant={color} now={percentage} min={0} max={100} />
                </div>
            </Col>
            </Row>
        );
        };
    
        return <Progress title={name} percentage={usage} color={colorValue}></Progress>;
    };

    return (
        <Card border="light" className="shadow-sm">
          <Card.Header className="border-bottom border-light">
            <h5 className="mb-0">{title}</h5>
          </Card.Header>
          <Card.Body>
            {data.map(elements => (percentageBars(elements, metric)))}
          </Card.Body>
        </Card>
      );
}

/**
 * Generate status widgets.
 * @returns A JSX Component that represents both pod and node resource usage.
 */
const StatusWidgets = () => {
    const [status, setStatus] = useState({ node_status: { not_ready: [], ready: [], no_connection: []}, pod_status: {error: [], pending: [], running:1 }});

    // Retrieve data from REST API.
    const getStatus = () => {
      const url = "/overview/status";
      var requestOptions = {
          method: 'GET',
          redirect: 'follow'
      };
  
      fetch(url, requestOptions)
          .then(response => response.text())
          .then(result => {
            setStatus(JSON.parse(result));
          })
          .catch(error => console.log('error', error));
    }

    useEffect(() => {
        getStatus();
    }, []);

    // Generate BarWidget using data retrieved.
    return ( 
        <>
        <Col xs={12} xl={6} className="mb-4">
            <NodeStatusWidget data={status.node_status}></NodeStatusWidget>
        </Col>
        <Col xs={12} xl={6} className="mb-4">
            <PodStatusWidget data={status.pod_status}></PodStatusWidget>
        </Col>
        </>
    );
}

/**
 * Generate status widget for node.
 * @param {props} props The props
 * @returns A Widget that is for node status.
 */
const NodeStatusWidget = (props) => {
    const { data } = props;

    return (
        <Card border="light" className="shadow-sm">
        <Card.Body>
          <Row className="d-block d-xl-flex align-items-center">
            <Col xs={12} xl={5} className="text-xl-center d-flex align-items-center justify-content-xl-center mb-3 mb-xl-0">
            <div className={`icon icon-shape icon-md rounded me-4 me-sm-0`}>
                <FontAwesomeIcon icon={faServer} />
              </div>
            </Col>
            <Col xs={12} xl={7} className="px-xl-0">
              <h5 className="mb-3">Nodes</h5>
              <h6 className="fw-normal text-gray">
                <Dropdown as={ButtonGroup}>
                  <Dropdown.Toggle as={Button} split variant="link" className="text-dark m-0 p-0">
                    <FontAwesomeIcon icon={faExclamationCircle} className={`icon icon-xs text-danger w-20 me-1`} />
                  </Dropdown.Toggle>
                <Dropdown.Menu>     
                    {    
                        data === undefined || data.no_connection === undefined || data.no_connection === null ? (
                        <></>
                    ) : (
                        <>
                            {data.no_connection.map(node => (controllersDropdown(node, "/nodes/detail/")))}
                        </>
                    )
                    }
                </Dropdown.Menu>
                </Dropdown>
                  No Connection: 
                    {
                        data.no_connection === undefined || data.no_connection === null ? (
                            0
                        ) : (
                            <>
                                {data.no_connection.length}
                            </>
                        )
                    }               
              </h6>
              <h6 className="fw-normal text-gray">
                <Dropdown as={ButtonGroup}>
                  <Dropdown.Toggle as={Button} split variant="link" className="text-dark m-0 p-0">
                    <FontAwesomeIcon icon={faQuestionCircle} className={`icon icon-xs text-warning w-20 me-1`} />
                    </Dropdown.Toggle>
                <Dropdown.Menu>     
                    {    
                        data === undefined || data.not_ready === undefined || data.not_ready === null ? (
                        <></>
                    ) : (
                        <>
                            {data.not_ready.map(node => (controllersDropdown(node, "/nodes/detail/")))}
                        </>
                    )
                    }
                </Dropdown.Menu>
                </Dropdown>
                  Not Ready: 
                    {
                        data.not_ready === undefined || data.not_ready === null ? (
                            0
                        ) : (
                            <>
                                {data.not_ready.length}
                            </>
                        )
                    }               
              </h6>
              <h6 className="fw-normal text-gray">
                <Dropdown as={ButtonGroup}>
                  <Dropdown.Toggle as={Button} split variant="link" className="text-dark m-0 p-0">
                  <FontAwesomeIcon icon={faCheckCircle} className={`icon icon-xs text-success w-20 me-1`} />
                  </Dropdown.Toggle>
                <Dropdown.Menu>     
                    {    
                        data === undefined || data.ready === undefined || data.ready === null ? (
                        <></>
                    ) : (
                        <>
                            {data.ready.map(node => (controllersDropdown(node, "/nodes/detail/")))}
                        </>
                    )
                    }
                </Dropdown.Menu>
                </Dropdown>
                  Ready: 
                    {
                        data.ready === undefined || data.ready === null ? (
                            0
                        ) : (
                            <>
                                {data.ready.length}
                            </>
                        )
                    }               
              </h6>
            </Col>
          </Row>
        </Card.Body>
        </Card>
    );
}

/**
 * Generate status widget for pod.
 * @param {props} props The props
 * @returns A Widget that is for pod status.
 */
const PodStatusWidget = (props) => {
    const { data } = props;

    return (
        <Card border="light" className="shadow-sm">
        <Card.Body>
          <Row className="d-block d-xl-flex align-items-center">
            <Col xs={12} xl={5} className="text-xl-center d-flex align-items-center justify-content-xl-center mb-3 mb-xl-0">
            <div className={`icon icon-shape icon-md rounded me-4 me-sm-0`}>
                <FontAwesomeIcon icon={faServer} />
              </div>
            </Col>
            <Col xs={12} xl={7} className="px-xl-0">
              <h5 className="mb-3">Pods</h5>
              <h6 className="fw-normal text-gray">
                <Dropdown as={ButtonGroup}>
                  <Dropdown.Toggle as={Button} split variant="link" className="text-dark m-0 p-0">
                    <FontAwesomeIcon icon={faExclamationCircle} className={`icon icon-xs text-danger w-20 me-1`} />
                  </Dropdown.Toggle>
                <Dropdown.Menu>     
                    {    
                        data === undefined || data.error === undefined || data.error === null ? (
                        <></>
                    ) : (
                        <>
                            {data.error.map(node => (controllersDropdown(node, "/nodes/detail/")))}
                        </>
                    )
                    }
                </Dropdown.Menu>
                </Dropdown>
                  Error: 
                    {
                        data.error === undefined || data.error === null ? (
                            0
                        ) : (
                            <>
                                {data.error.length}
                            </>
                        )
                    }               
              </h6>
              <h6 className="fw-normal text-gray">
                <Dropdown as={ButtonGroup}>
                  <Dropdown.Toggle as={Button} split variant="link" className="text-dark m-0 p-0">
                    <FontAwesomeIcon icon={faQuestionCircle} className={`icon icon-xs text-warning w-20 me-1`} />
                    </Dropdown.Toggle>
                <Dropdown.Menu>     
                    {    
                        data === undefined || data.pending === undefined || data.pending === null ? (
                        <></>
                    ) : (
                        <>
                            {data.pending.map(node => (controllersDropdown(node, "/nodes/detail/")))}
                        </>
                    )
                    }
                </Dropdown.Menu>
                </Dropdown>
                  Pending: 
                    {
                        data.pending === undefined || data.pending === null ? (
                            0
                        ) : (
                            <>
                                {data.pending.length}
                            </>
                        )
                    }               
              </h6>
              <h6 className="fw-normal text-gray">
              <FontAwesomeIcon icon={faCheckCircle} className={`icon icon-xs text-success w-20 me-1`} />
                  Running: 
                    {
                        data.running === undefined || data.running === null ? (
                            0
                        ) : (
                            data.running
                        )
                    }               
              </h6>
            </Col>
          </Row>
        </Card.Body>
        </Card>
    );
}


/**
 * Generate last 24 hours graph for cpu, ram, disk and network.
 * @returns A JSX Component that represents last 24 hours graph.
 */
const LastStatus = () => {
    const [status, setStatus] = useState({ "cpu_usage": [], "ram_usage": [], "network_usage": [], "disk_usage": [] });

    // Retrieve data from REST API.
    const getStatus = () => {
      const url = "/overview/nodes/usage";
      var requestOptions = {
          method: 'GET',
          redirect: 'follow'
      };
  
      fetch(url, requestOptions)
          .then(response => response.text())
          .then(result => {
            setStatus(JSON.parse(result));
          })
          .catch(error => console.log('error', error));
    }

    useEffect(() => {
        getStatus();
    }, []);

    // Generate BarWidget using data retrieved.
    return ( 
        <>
        <Row className="justify-content-md-center">
            <Col xs={12} xl={6} className="mb-4">
                <LastStatusGraph data={status.cpu_usage} title={"CPU Usage"}></LastStatusGraph>
            </Col>
            <Col xs={12} xl={6} className="mb-4">
                <LastStatusGraph data={status.ram_usage} title={"RAM Usage"}></LastStatusGraph>
            </Col>
        </Row>
        <Row className="justify-content-md-center">
            <Col xs={12} xl={6} className="mb-4">
                <LastStatusGraph data={status.network_usage} title={"Network Usage"}></LastStatusGraph>
            </Col>
            <Col xs={12} xl={6} className="mb-4">
                <LastStatusGraph data={status.disk_usage} title={"Disk Usage"}></LastStatusGraph>
            </Col>
        </Row>
        </>
    );
}

/**
 * Generate last 24 hours graph.
 * @param {props} props The props
 * @returns A JSX Component that shows last status graph for 24 hours.
 */
const LastStatusGraph = (props) => {
    const { data, title } = props;
    return (
        <Card className="shadow-sm">
          <Card.Header className="d-flex flex-row align-items-center flex-0">
            <div className="d-block">
              <h5 className="fw-normal mb-2">
                {title}
              </h5>
              <h3>{data === null || data === undefined ? "UNKNOWN" : data[data.length - 1] + "%"}</h3>
              <small>Last 24 Hours</small>
            </div>
          </Card.Header>
          <Card.Body className="p-2">
            <Valuechart data={data === null || data === undefined ? [] : data}></Valuechart>
          </Card.Body>
        </Card>
      );
}

/**
 * Generate value chart from data.
 * @param {props} props The props
 * @returns A JSX Component that shows graph in series.
 */
export const Valuechart = (props) => {
    const { data } = props; 
    var graphValues = data;

    // If data was null or undefined, set it 0.
    if (data === null || data === undefined) {
        data = new Array(24).fill(0);
    } else if (data.length != 24) { // When data was not enough for last 24 hours, fill rest with 0.
        var original = data;
        var newData = new Array(24).fill(0);

        for (let i = 0; i < original.length; i++) {
            newData[i] = original[i];
        }
        graphValues = newData;
    }

    // Add labels
    var label = []
    for (var i = 0; i < data.length; i++) {
        label.push(24 - i)
    }
  
    const graphData = {
      labels: label,
      series: [graphValues]
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